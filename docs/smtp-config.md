# Configuración SMTP — ArticNexus

ArticNexus usa **dos cuentas SMTP separadas** para distintos tipos de correo:

| Cuenta | Variables | Propósito |
|---|---|---|
| Soporte técnico | `SUPPORT_SMTP_*` | Restablecimiento de contraseña, tickets técnicos |
| Empresa / Negocio | `BUSINESS_SMTP_*` | Formulario de contacto, propuestas comerciales |

Ambas cuentas usan Gmail con **App Passwords** (contraseñas de aplicación), no la contraseña real de la cuenta de Google.

## 1. Crear un App Password de Gmail

1. Ir a [myaccount.google.com/security](https://myaccount.google.com/security)
2. Activar **Verificación en 2 pasos** si no está activada
3. Buscar **Contraseñas de aplicaciones** (en la sección "Cómo inicias sesión en Google")
4. Nombre de la app: `ArticNexus Backend` (libre, solo es descriptivo)
5. Copiar la contraseña de 16 caracteres generada (ej. `abcd efgh ijkl mnop`)
6. Eliminar los espacios al usarla en las variables de entorno

> Las App Passwords solo aparecen una vez. Guardarlas en un gestor de contraseñas.

## 2. Variables de entorno

```env
# Servidor SMTP (Gmail — el mismo para ambas cuentas)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Cuenta de soporte técnico (password reset, tickets)
SUPPORT_SMTP_USER=soporte@tudominio.com
SUPPORT_SMTP_PASSWORD=abcdefghijklmnop
SUPPORT_SMTP_FROM=ArticNexus Soporte <soporte@tudominio.com>

# Cuenta de empresa/negocio (formulario de contacto)
BUSINESS_SMTP_USER=contacto@tudominio.com
BUSINESS_SMTP_PASSWORD=abcdefghijklmnop
BUSINESS_SMTP_FROM=ArticDev S.A. <contacto@tudominio.com>

# Email de destino para el formulario de contacto
CONTACT_EMAIL=info@tudominio.com

# URL del frontend (para construir el link de reset)
FRONTEND_URL=https://nexus.articdev.com
```

## 3. Comportamiento si SMTP no está configurado

Si las variables `SUPPORT_SMTP_*` no están configuradas:
- El restablecimiento de contraseña **no enviará correo** pero no fallará — el token se genera igual
- Se registrará un warning en `storage/logs/security.log`

En desarrollo local es normal no tener SMTP configurado. Para probar el flujo de reset, revisar el log para obtener el token.

## 4. TTL del token de reset

```env
PASSWORD_RESET_EXP_MIN=30    # minutos de validez del token (default: 30)
```

## 5. Testing del SMTP

Desde el sistema, ir a la pantalla de login → "¿Olvidaste tu contraseña?" e ingresar un usuario que tenga email configurado. Si el correo llega, el SMTP está operativo.

Para debug avanzado, revisar `storage/logs/security.log` que registrará errores de envío con el mensaje exacto del servidor SMTP.

## 6. Usar solo una cuenta

Si se prefiere una sola cuenta para todo:

```env
SUPPORT_SMTP_USER=cuenta@gmail.com
SUPPORT_SMTP_PASSWORD=apppassword
SUPPORT_SMTP_FROM=ArticNexus <cuenta@gmail.com>

BUSINESS_SMTP_USER=cuenta@gmail.com
BUSINESS_SMTP_PASSWORD=apppassword
BUSINESS_SMTP_FROM=ArticDev S.A. <cuenta@gmail.com>
```
