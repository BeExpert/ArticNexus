/**
 * i18n error helper.
 *
 * Usage:
 *   import { te, teError } from '@/i18n'
 *
 *   // Translate a known code:
 *   te('auth.invalid_credentials')   // → "Usuario o contraseña incorrectos..."
 *
 *   // Translate an Axios error automatically:
 *   catch (err) { errorMsg.value = teError(err) }
 *
 *   // Change language at runtime:
 *   import { setLang } from '@/i18n'
 *   setLang('en')
 */

import es from './errors.es.js'
import en from './errors.en.js'

const DICTS = { es, en }

// ─── Active language ──────────────────────────────────────────────────────────

/**
 * Returns the currently active language code ('es' | 'en').
 * Reads from localStorage so it persists across page loads.
 */
export function getLang() {
  const stored = localStorage.getItem('lang')
  return stored && DICTS[stored] ? stored : 'es'
}

/**
 * Changes the active language and persists the choice.
 * @param {'es'|'en'} lang
 */
export function setLang(lang) {
  if (!DICTS[lang]) return
  localStorage.setItem('lang', lang)
}

// ─── Core translation ─────────────────────────────────────────────────────────

/**
 * te — Translate Error code.
 * Returns the localised message for a stable error code.
 * Falls back to `generic.unknown` if the code is not in the active dictionary.
 *
 * @param {string} code - e.g. 'auth.invalid_credentials'
 * @returns {string}
 */
export function te(code) {
  const dict = DICTS[getLang()] || DICTS.es
  return dict[code] || dict['generic.unknown'] || 'Error desconocido.'
}

// ─── HTTP status fallback ─────────────────────────────────────────────────────

/** Maps HTTP status codes to generic error code keys. */
const STATUS_CODE_MAP = {
  400: 'http.bad_request',
  401: 'http.unauthorized',
  403: 'http.forbidden',
  404: 'http.not_found',
  409: 'http.conflict_duplicate',
  422: 'http.validation',
  429: 'http.too_many_requests',
}

/**
 * teStatus — Translate HTTP Status.
 * Returns a localised message for a given HTTP status code.
 *
 * @param {number} httpStatus
 * @returns {string}
 */
export function teStatus(httpStatus) {
  if (httpStatus >= 500) return te('http.internal')
  const code = STATUS_CODE_MAP[httpStatus]
  return code ? te(code) : te('generic.unknown')
}

// ─── Axios error helper ───────────────────────────────────────────────────────

/**
 * teError — Translate Axios Error.
 * Unified handler for all API call errors in the frontend.
 *
 * Resolution order:
 *  1. `error.response.data.code`  → te(code)          (backend emitted a stable code)
 *  2. `error.response.status`     → teStatus(status)  (no code but we have HTTP status)
 *  3. No response at all          → te('network.connection_error')
 *  4. Fallback                    → te('generic.unknown')
 *
 * @param {import('axios').AxiosError} error
 * @param {string} [fallbackCode] - Optional override code when all else fails
 * @returns {string}
 */
export function teError(error, fallbackCode) {
  if (!error) return te(fallbackCode || 'generic.unknown')

  const code = error?.response?.data?.code
  if (code) return te(code)

  const status = error?.response?.status
  if (status) return teStatus(status)

  // No response — network error or CORS
  if (!error.response) return te('network.connection_error')

  return te(fallbackCode || 'generic.unknown')
}
