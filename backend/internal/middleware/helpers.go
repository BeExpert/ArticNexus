package middleware

import (
	"encoding/json"
	"net/http"
)

// encodeJSON writes v as JSON to w. Errors are silently swallowed because
// a write failure at this point is unrecoverable within the middleware chain.
func encodeJSON(w http.ResponseWriter, v interface{}) {
	_ = json.NewEncoder(w).Encode(v)
}
