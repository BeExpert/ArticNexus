package domain

import "time"

// Response is the standard envelope for every API response.
//
//	{
//	  "success": true,
//	  "data":    {},
//	  "message": "...",
//	  "code":    "auth.invalid_credentials",  // only on errors
//	  "errors":  null,
//	  "meta":    { "timestamp": "..." }
//	}
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Errors  interface{} `json:"errors"`
	Meta    Meta        `json:"meta"`
}

// Meta carries metadata such as pagination and timestamp.
type Meta struct {
	Timestamp  string      `json:"timestamp"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination is included in list responses.
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// NewResponse builds a successful response envelope.
func NewResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
		Errors:  nil,
		Meta:    Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	}
}

// NewResponseWithPagination builds a successful paginated list response.
func NewResponseWithPagination(data interface{}, message string, p Pagination) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
		Errors:  nil,
		Meta: Meta{
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Pagination: &p,
		},
	}
}

// NewErrorResponse builds an error response envelope without a code.
func NewErrorResponse(message string, errors interface{}) Response {
	return Response{
		Success: false,
		Data:    nil,
		Message: message,
		Errors:  errors,
		Meta:    Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	}
}

// NewErrorResponseCode builds an error response envelope with a stable error code.
func NewErrorResponseCode(message, code string, errors interface{}) Response {
	return Response{
		Success: false,
		Data:    nil,
		Message: message,
		Code:    code,
		Errors:  errors,
		Meta:    Meta{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	}
}

// PaginationParams holds query parameters extracted from requests.
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
}

// Offset calculates the SQL OFFSET value.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}
