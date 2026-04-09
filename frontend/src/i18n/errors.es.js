/**
 * Diccionario de errores en español.
 * Las claves son los `code` estables que devuelve el backend.
 * Nunca renombres una clave existente — añade nuevas si las necesitas.
 */
export default {
  // ── HTTP genéricos ────────────────────────────────────────────────────────
  'http.bad_request':
    'La solicitud enviada no es válida. Revisa los datos e intenta de nuevo.',
  'http.unauthorized':
    'No tienes autorización para realizar esta acción. Inicia sesión nuevamente.',
  'http.forbidden':
    'No tienes permiso para acceder a este recurso. Contacta al administrador.',
  'http.not_found':
    'El recurso solicitado no fue encontrado.',
  'http.conflict_duplicate':
    'Ya existe un registro con esos datos. Verifica e intenta con valores distintos.',
  'http.conflict_dependency':
    'No se puede eliminar: existen registros que dependen de este.',
  'http.validation':
    'Algunos campos no son válidos. Revisa el formulario e intenta de nuevo.',
  'http.too_many_requests':
    'Demasiadas solicitudes en poco tiempo. Espera unos minutos e intenta de nuevo.',
  'http.internal':
    'Error interno del servidor. Intenta de nuevo más tarde o contacta al soporte.',

  // ── Autenticación — JWT y sesión ──────────────────────────────────────────
  'auth.invalid_credentials':
    'Usuario o contraseña incorrectos. Verifica tus datos e intenta de nuevo.',
  'auth.account_inactive':
    'Tu cuenta está inactiva. Contacta al administrador para rehabilitarla.',
  'auth.demo_expired':
    'Tu sesión de prueba ha expirado. Contacta al administrador para recuperar el acceso.',
  'auth.unauthenticated':
    'No se pudo verificar tu identidad. Por favor inicia sesión nuevamente.',
  'auth.missing_header':
    'Falta el encabezado de autorización. Por favor inicia sesión nuevamente.',
  'auth.invalid_header_format':
    'El formato del encabezado de autorización no es válido. Inicia sesión nuevamente.',
  'auth.invalid_token':
    'Tu token de sesión no es válido o ha expirado. Inicia sesión nuevamente.',
  'auth.session_expired':
    'Tu sesión ha expirado porque el servidor fue reiniciado. Inicia sesión nuevamente.',

  // ── Restablecimiento de contraseña ────────────────────────────────────────
  'auth.reset_link_invalid':
    'El enlace de restablecimiento no es válido. Solicita uno nuevo.',
  'auth.reset_link_used':
    'Este enlace ya fue utilizado. Solicita un nuevo enlace de restablecimiento.',
  'auth.reset_link_expired':
    'El enlace ha expirado. Solicita un nuevo enlace de restablecimiento.',

  // ── Dominio — usuarios ────────────────────────────────────────────────────
  'user.already_in_company':
    'Este usuario ya está asignado a esta empresa.',
  'user.not_in_company':
    'Este usuario no pertenece a esta empresa.',

  // ── Dominio — demo links ──────────────────────────────────────────────────
  'demolink.no_demo_user':
    'No hay un usuario demo configurado para esta aplicación. Contacta al soporte.',

  // ── Formulario de contacto ────────────────────────────────────────────────
  'contact.invalid_name':
    'El nombre es inválido. Debe tener al menos 2 caracteres.',
  'contact.invalid_email':
    'El correo electrónico no es válido.',
  'contact.invalid_type':
    'Selecciona un tipo de solicitud válido.',
  'contact.desc_too_short':
    'La descripción es demasiado corta. Mínimo 20 caracteres.',
  'contact.send_failed':
    'No se pudo enviar la solicitud. Intenta de nuevo más tarde.',

  // ── Validación de formularios (client-side) ───────────────────────────────
  'form.field_required':
    'Este campo es requerido.',
  'form.username_required':
    'El nombre de usuario es requerido.',
  'form.password_required':
    'La contraseña es requerida.',
  'form.passwords_mismatch':
    'Las contraseñas no coinciden.',
  'form.invalid_email':
    'El correo electrónico no es válido.',
  'form.min_2_chars':
    'Mínimo 2 caracteres.',
  'form.min_20_chars':
    'Mínimo 20 caracteres.',
  'form.select_option':
    'Debes seleccionar una opción.',
  'form.select_user':
    'Busca y selecciona un usuario del sistema.',
  'form.select_app':
    'Selecciona una aplicación.',
  'form.select_branch':
    'Selecciona al menos una sucursal.',
  'form.select_role':
    'Selecciona al menos un rol.',
  'form.name_required':
    'El nombre es requerido.',
  'form.first_name_required':
    'Nombre y primer apellido son obligatorios.',

  // ── Red y genéricos ───────────────────────────────────────────────────────
  'network.connection_error':
    'No se pudo conectar al servidor. Verifica tu conexión a internet.',
  'generic.unknown':
    'Ocurrió un error inesperado. Intenta de nuevo.',
}
