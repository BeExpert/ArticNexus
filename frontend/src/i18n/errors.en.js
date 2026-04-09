/**
 * Error dictionary in English.
 * Keys are the stable `code` values returned by the backend.
 * Never rename an existing key — add new ones if needed.
 */
export default {
  // ── HTTP generic ──────────────────────────────────────────────────────────
  'http.bad_request':
    'The request is invalid. Please check your data and try again.',
  'http.unauthorized':
    'You are not authorised to perform this action. Please log in again.',
  'http.forbidden':
    'You do not have permission to access this resource. Contact your administrator.',
  'http.not_found':
    'The requested resource was not found.',
  'http.conflict_duplicate':
    'A record with those details already exists. Please use different values.',
  'http.conflict_dependency':
    'Cannot delete: other records depend on this one.',
  'http.validation':
    'Some fields are invalid. Please check the form and try again.',
  'http.too_many_requests':
    'Too many requests in a short period. Please wait a few minutes and try again.',
  'http.internal':
    'Internal server error. Please try again later or contact support.',

  // ── Authentication — JWT & session ────────────────────────────────────────
  'auth.invalid_credentials':
    'Incorrect username or password. Please verify your details and try again.',
  'auth.account_inactive':
    'Your account is inactive. Contact your administrator to reactivate it.',
  'auth.demo_expired':
    'Your demo session has expired. Contact your administrator to regain access.',
  'auth.unauthenticated':
    'Your identity could not be verified. Please log in again.',
  'auth.missing_header':
    'The authorisation header is missing. Please log in again.',
  'auth.invalid_header_format':
    'The authorisation header format is invalid. Please log in again.',
  'auth.invalid_token':
    'Your session token is invalid or has expired. Please log in again.',
  'auth.session_expired':
    'Your session has expired because the server was restarted. Please log in again.',

  // ── Password reset ────────────────────────────────────────────────────────
  'auth.reset_link_invalid':
    'The reset link is not valid. Please request a new one.',
  'auth.reset_link_used':
    'This link has already been used. Please request a new reset link.',
  'auth.reset_link_expired':
    'The link has expired. Please request a new reset link.',

  // ── Domain — users ────────────────────────────────────────────────────────
  'user.already_in_company':
    'This user is already assigned to this company.',
  'user.not_in_company':
    'This user does not belong to this company.',

  // ── Domain — demo links ───────────────────────────────────────────────────
  'demolink.no_demo_user':
    'No demo user is configured for this application. Contact support.',

  // ── Contact form ──────────────────────────────────────────────────────────
  'contact.invalid_name':
    'The name is invalid. Must be at least 2 characters.',
  'contact.invalid_email':
    'The email address is not valid.',
  'contact.invalid_type':
    'Please select a valid request type.',
  'contact.desc_too_short':
    'The description is too short. Minimum 20 characters.',
  'contact.send_failed':
    'Could not send the request. Please try again later.',

  // ── Form validation (client-side) ─────────────────────────────────────────
  'form.field_required':
    'This field is required.',
  'form.username_required':
    'Username is required.',
  'form.password_required':
    'Password is required.',
  'form.passwords_mismatch':
    'Passwords do not match.',
  'form.invalid_email':
    'The email address is not valid.',
  'form.min_2_chars':
    'Minimum 2 characters.',
  'form.min_20_chars':
    'Minimum 20 characters.',
  'form.select_option':
    'You must select an option.',
  'form.select_user':
    'Search and select a system user.',
  'form.select_app':
    'Please select an application.',
  'form.select_branch':
    'Please select at least one branch.',
  'form.select_role':
    'Please select at least one role.',
  'form.name_required':
    'Name is required.',
  'form.first_name_required':
    'First name and first surname are required.',

  // ── Network & generic ────────────────────────────────────────────────────
  'network.connection_error':
    'Could not connect to the server. Please check your internet connection.',
  'generic.unknown':
    'An unexpected error occurred. Please try again.',
}
