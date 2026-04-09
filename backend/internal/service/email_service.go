package service

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"articnexus/backend/internal/config"
)

// EmailService sends transactional emails.
type EmailService interface {
	SendPasswordReset(toEmail, resetURL string) error
	SendContactForm(nombre, email, tipo, descripcion string) error
	// SendDemoInvitation sends a demo access email containing the link URL and
	// a temporary plaintext password that has already been applied to the demo
	// user account. Sent from the Business SMTP account (articdevsa@gmail.com).
	SendDemoInvitation(guestName, toEmail, appName, demoURL, tempPass, demoUsername string) error
}

type emailService struct {
	cfg *config.Config
}

// NewEmailService returns an EmailService. If SMTP is not configured, emails
// are logged to stdout instead of being sent.
func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendPasswordReset(toEmail, resetURL string) error {
	subject := "Restablecimiento de contraseña — ArticNexus"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f8fafc;padding:40px 0">
<tr><td align="center">
<table role="presentation" width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;border:1px solid #e2e8f0;overflow:hidden">

  <!-- Header -->
  <tr>
    <td style="background-color:#0f172a;padding:24px 32px">
      <h1 style="margin:0;color:#ffffff;font-size:20px;font-weight:700;letter-spacing:0.5px">ArticNexus</h1>
    </td>
  </tr>

  <!-- Body -->
  <tr>
    <td style="padding:32px">
      <h2 style="margin:0 0 16px;color:#0f172a;font-size:18px;font-weight:600">Solicitud de restablecimiento de contraseña</h2>
      <p style="margin:0 0 12px;color:#334155;font-size:14px;line-height:1.6">
        Estimado usuario,
      </p>
      <p style="margin:0 0 12px;color:#334155;font-size:14px;line-height:1.6">
        Hemos recibido una solicitud para restablecer la contraseña asociada a tu cuenta en la plataforma ArticNexus.
      </p>
      <p style="margin:0 0 24px;color:#334155;font-size:14px;line-height:1.6">
        Para continuar con el proceso, haz clic en el siguiente boton:
      </p>

      <!-- CTA Button -->
      <table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 24px">
        <tr>
          <td style="background-color:#0f172a;border-radius:6px;padding:12px 32px">
            <a href="%s" target="_blank" style="color:#ffffff;font-size:14px;font-weight:600;text-decoration:none;display:inline-block">Restablecer contraseña</a>
          </td>
        </tr>
      </table>

      <p style="margin:0 0 8px;color:#64748b;font-size:13px;line-height:1.5">
        Este enlace tiene una validez de <strong>%d minutos</strong> y solo puede utilizarse una vez.
      </p>
      <p style="margin:0 0 0;color:#64748b;font-size:13px;line-height:1.5">
        Si no realizaste esta solicitud, puedes ignorar este mensaje. Tu cuenta permanecera segura.
      </p>
    </td>
  </tr>

  <!-- Divider -->
  <tr>
    <td style="padding:0 32px">
      <hr style="border:none;border-top:1px solid #e2e8f0;margin:0">
    </td>
  </tr>

  <!-- Footer -->
  <tr>
    <td style="padding:20px 32px">
      <p style="margin:0 0 4px;color:#94a3b8;font-size:12px;line-height:1.5">
        Este es un correo automatico enviado por ArticNexus.
      </p>
      <p style="margin:0;color:#94a3b8;font-size:12px;line-height:1.5">
        Por favor, no respondas a este mensaje.
      </p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, resetURL, s.cfg.PasswordResetExpMin)

	// If SMTP is not configured, log the link to stdout (development mode).
	if s.cfg.SupportSMTPUser == "" || s.cfg.SupportSMTPPassword == "" {
		log.Printf("[EMAIL DEV] Password reset for %s → %s", toEmail, resetURL)
		return nil
	}

	msg := strings.Join([]string{
		"From: " + s.cfg.SupportSMTPFrom,
		"To: " + toEmail,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.cfg.SupportSMTPUser, s.cfg.SupportSMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, s.cfg.SupportSMTPUser, []string{toEmail}, []byte(msg))
}

func (s *emailService) SendContactForm(nombre, email, tipo, descripcion string) error {
	to := s.cfg.ContactEmail
	if to == "" {
		to = s.cfg.BusinessSMTPUser
	}

	subject := fmt.Sprintf("[ArticDev S.A.] Solicitud: %s", tipo)
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f8fafc;padding:40px 0">
<tr><td align="center">
<table role="presentation" width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;border:1px solid #e2e8f0;overflow:hidden">
  <tr>
    <td style="background-color:#0f172a;padding:24px 32px">
      <h1 style="margin:0;color:#ffffff;font-size:20px;font-weight:700">ArticDev S.A. — Nueva Solicitud</h1>
    </td>
  </tr>
  <tr>
    <td style="padding:32px">
      <h2 style="margin:0 0 16px;color:#0f172a;font-size:18px;font-weight:600">%s</h2>
      <table style="width:100%%;border-collapse:collapse;margin-bottom:24px">
        <tr><td style="padding:8px 0;color:#64748b;font-size:13px;width:120px;vertical-align:top">Nombre:</td>
            <td style="padding:8px 0;color:#0f172a;font-size:14px"><strong>%s</strong></td></tr>
        <tr><td style="padding:8px 0;color:#64748b;font-size:13px;vertical-align:top">Email:</td>
            <td style="padding:8px 0;color:#0f172a;font-size:14px"><a href="mailto:%s" style="color:#00D4FF">%s</a></td></tr>
        <tr><td style="padding:8px 0;color:#64748b;font-size:13px;vertical-align:top">Tipo:</td>
            <td style="padding:8px 0;color:#0f172a;font-size:14px">%s</td></tr>
      </table>
      <h3 style="margin:0 0 8px;color:#0f172a;font-size:15px;font-weight:600">Descripción:</h3>
      <p style="margin:0;color:#334155;font-size:14px;line-height:1.7;white-space:pre-wrap">%s</p>
    </td>
  </tr>
  <tr>
    <td style="padding:16px 32px;background-color:#f8fafc">
      <p style="margin:0;color:#94a3b8;font-size:12px">Correo generado automáticamente por ArticNexus.</p>
    </td>
  </tr>
</table>
</td></tr>
</table>
</body>
</html>`, tipo, nombre, email, email, tipo, descripcion)

	if s.cfg.BusinessSMTPUser == "" || s.cfg.BusinessSMTPPassword == "" {
		log.Printf("[EMAIL DEV] Contact form from %s (%s) — tipo: %s", nombre, email, tipo)
		return nil
	}

	msg := strings.Join([]string{
		"From: " + s.cfg.BusinessSMTPFrom,
		"To: " + to,
		"Reply-To: " + email,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.cfg.BusinessSMTPUser, s.cfg.BusinessSMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, s.cfg.BusinessSMTPUser, []string{to}, []byte(msg))
}

func (s *emailService) SendDemoInvitation(guestName, toEmail, appName, demoURL, tempPass, demoUsername string) error {
	subject := fmt.Sprintf("Acceso demo a %s — ArticDev", appName)
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif">
<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f8fafc;padding:40px 0">
<tr><td align="center">
<table role="presentation" width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;border:1px solid #e2e8f0;overflow:hidden">

  <!-- Header -->
  <tr>
    <td style="background-color:#0f172a;padding:24px 32px">
      <h1 style="margin:0;color:#ffffff;font-size:20px;font-weight:700;letter-spacing:0.5px">ArticDev &#183; Demo %s</h1>
    </td>
  </tr>

  <!-- Body -->
  <tr>
    <td style="padding:32px">
      <h2 style="margin:0 0 16px;color:#0f172a;font-size:18px;font-weight:600">Hola, %s</h2>
      <p style="margin:0 0 12px;color:#334155;font-size:14px;line-height:1.6">
        Has sido invitado a explorar una demo completa de <strong>%s</strong>.
        La demo te permite navegar el sistema con datos gen&eacute;ricos en modo de solo lectura.
      </p>
      <p style="margin:0 0 24px;color:#334155;font-size:14px;line-height:1.6">
        Accede haciendo clic en el bot&oacute;n o usando las credenciales de abajo:
      </p>

      <!-- CTA Button -->
      <table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 28px">
        <tr>
          <td style="background-color:#0f172a;border-radius:6px;padding:12px 32px">
            <a href="%s" target="_blank" style="color:#ffffff;font-size:14px;font-weight:600;text-decoration:none;display:inline-block">Acceder a la demo</a>
          </td>
        </tr>
      </table>

      <!-- Credentials box -->
      <table role="presentation" width="100%%" cellpadding="0" cellspacing="0"
             style="background-color:#f1f5f9;border-radius:6px;margin-bottom:24px;border:1px solid #e2e8f0">
        <tr>
          <td style="padding:16px 20px">
            <p style="margin:0 0 6px;color:#64748b;font-size:12px;font-weight:600;text-transform:uppercase;letter-spacing:0.8px">Acceso directo (alternativo al link)</p>
            <p style="margin:0 0 4px;color:#0f172a;font-size:13px"><strong>Usuario:</strong>&nbsp; <code style="background:#e2e8f0;padding:2px 6px;border-radius:4px;font-size:13px">%s</code></p>
            <p style="margin:0;color:#0f172a;font-size:13px"><strong>Contrase&ntilde;a temporal:</strong>&nbsp; <code style="background:#e2e8f0;padding:2px 6px;border-radius:4px;font-size:13px">%s</code></p>
          </td>
        </tr>
      </table>

      <p style="margin:0 0 8px;color:#64748b;font-size:13px;line-height:1.5">
        Este acceso expira pronto y es de solo lectura. No se requiere registro.
      </p>
      <p style="margin:0;color:#64748b;font-size:13px;line-height:1.5">
        Si recibiste este correo por error, puedes ignorarlo.
      </p>
    </td>
  </tr>

  <tr><td style="padding:0 32px"><hr style="border:none;border-top:1px solid #e2e8f0;margin:0"></td></tr>

  <tr>
    <td style="padding:20px 32px">
      <p style="margin:0;color:#94a3b8;font-size:12px">Correo enviado por ArticDev S.A. &mdash; no responder.</p>
    </td>
  </tr>

</table>
</td></tr>
</table>
</body>
</html>`, appName, guestName, appName, demoURL, demoUsername, tempPass)

	if s.cfg.BusinessSMTPUser == "" || s.cfg.BusinessSMTPPassword == "" {
		log.Printf("[EMAIL DEV] Demo invitation for %s → %s  user:%s / pass:%s", toEmail, demoURL, demoUsername, tempPass)
		return nil
	}

	msg := strings.Join([]string{
		"From: " + s.cfg.BusinessSMTPFrom,
		"To: " + toEmail,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		body,
	}, "\r\n")

	smtpAuth := smtp.PlainAuth("", s.cfg.BusinessSMTPUser, s.cfg.BusinessSMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, smtpAuth, s.cfg.BusinessSMTPUser, []string{toEmail}, []byte(msg))
}
