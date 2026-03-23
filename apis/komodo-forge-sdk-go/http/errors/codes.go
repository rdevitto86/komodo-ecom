package httperrors

import "net/http"

// ErrorCode is the canonical error type for all Komodo services.
// Domain-specific codes are defined in each service's pkg/v1/models/errors.go
// using this type, within their registered range (see ranges.go).
type ErrorCode struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// 10xxx — generic HTTP-level errors, usable by any service
type GlobalErrors struct {
	BadRequest          ErrorCode
	Unauthorized        ErrorCode
	PaymentRequired     ErrorCode
	Forbidden           ErrorCode
	NotFound            ErrorCode
	MethodNotAllowed    ErrorCode
	Conflict            ErrorCode
	UnprocessableEntity ErrorCode
	TooManyRequests     ErrorCode
	Internal            ErrorCode
	NotImplemented      ErrorCode
	BadGateway          ErrorCode
	ServiceUnavailable  ErrorCode
	GatewayTimeout      ErrorCode
}

var Global = GlobalErrors{
	BadRequest:          ErrorCode{ID: "10001", Status: http.StatusBadRequest, Message: "Bad request"},
	Unauthorized:        ErrorCode{ID: "10002", Status: http.StatusUnauthorized, Message: "Unauthorized"},
	PaymentRequired:     ErrorCode{ID: "10003", Status: http.StatusPaymentRequired, Message: "Payment required"},
	Forbidden:           ErrorCode{ID: "10004", Status: http.StatusForbidden, Message: "Forbidden"},
	NotFound:            ErrorCode{ID: "10005", Status: http.StatusNotFound, Message: "Not found"},
	MethodNotAllowed:    ErrorCode{ID: "10006", Status: http.StatusMethodNotAllowed, Message: "Method not allowed"},
	Conflict:            ErrorCode{ID: "10007", Status: http.StatusConflict, Message: "Conflict"},
	UnprocessableEntity: ErrorCode{ID: "10008", Status: http.StatusUnprocessableEntity, Message: "Unprocessable entity"},
	TooManyRequests:     ErrorCode{ID: "10009", Status: http.StatusTooManyRequests, Message: "Too many requests"},
	Internal:            ErrorCode{ID: "10010", Status: http.StatusInternalServerError, Message: "Internal server error"},
	NotImplemented:      ErrorCode{ID: "10011", Status: http.StatusNotImplemented, Message: "Not implemented"},
	BadGateway:          ErrorCode{ID: "10012", Status: http.StatusBadGateway, Message: "Bad gateway"},
	ServiceUnavailable:  ErrorCode{ID: "10013", Status: http.StatusServiceUnavailable, Message: "Service unavailable"},
	GatewayTimeout:      ErrorCode{ID: "10014", Status: http.StatusGatewayTimeout, Message: "Gateway timeout"},
}

// 11xxx — database errors, owned by the forge-sdk DynamoDB/Aurora clients
type DBErrors struct {
	ConnectionFailed  ErrorCode
	QueryFailed       ErrorCode
	TransactionFailed ErrorCode
	RecordNotFound    ErrorCode
	DuplicateEntry    ErrorCode
}

var DB = DBErrors{
	ConnectionFailed:  ErrorCode{ID: "11001", Status: http.StatusInternalServerError, Message: "Database connection failed"},
	QueryFailed:       ErrorCode{ID: "11002", Status: http.StatusInternalServerError, Message: "Database query failed"},
	TransactionFailed: ErrorCode{ID: "11003", Status: http.StatusInternalServerError, Message: "Database transaction failed"},
	RecordNotFound:    ErrorCode{ID: "11004", Status: http.StatusNotFound, Message: "Record not found"},
	DuplicateEntry:    ErrorCode{ID: "11005", Status: http.StatusConflict, Message: "Duplicate entry"},
}

// 20xxx — auth/JWT errors, used by the forge-sdk auth middleware and komodo-auth-api
type AuthErrors struct {
	InvalidClientCredentials ErrorCode
	InvalidGrantType         ErrorCode
	InvalidScope             ErrorCode
	InvalidToken             ErrorCode
	InvalidKey               ErrorCode
	ExpiredToken             ErrorCode
	UnauthorizedClient       ErrorCode
	UnsupportedGrantType     ErrorCode
	UnsupportedResponseType  ErrorCode
	InvalidRedirectURI       ErrorCode
	AccessDenied             ErrorCode
	InsufficientScope        ErrorCode
}

var Auth = AuthErrors{
	InvalidClientCredentials: ErrorCode{ID: "20001", Status: http.StatusUnauthorized, Message: "Invalid client credentials"},
	InvalidGrantType:         ErrorCode{ID: "20002", Status: http.StatusBadRequest, Message: "Invalid grant type"},
	InvalidScope:             ErrorCode{ID: "20003", Status: http.StatusBadRequest, Message: "Invalid scope"},
	InvalidToken:             ErrorCode{ID: "20004", Status: http.StatusUnauthorized, Message: "Invalid token"},
	InvalidKey:               ErrorCode{ID: "20005", Status: http.StatusUnauthorized, Message: "Invalid auth key"},
	ExpiredToken:             ErrorCode{ID: "20006", Status: http.StatusUnauthorized, Message: "Token expired"},
	UnauthorizedClient:       ErrorCode{ID: "20007", Status: http.StatusUnauthorized, Message: "Unauthorized client"},
	UnsupportedGrantType:     ErrorCode{ID: "20008", Status: http.StatusBadRequest, Message: "Unsupported grant type"},
	UnsupportedResponseType:  ErrorCode{ID: "20009", Status: http.StatusBadRequest, Message: "Unsupported response type"},
	InvalidRedirectURI:       ErrorCode{ID: "20010", Status: http.StatusBadRequest, Message: "Invalid redirect URI"},
	AccessDenied:             ErrorCode{ID: "20011", Status: http.StatusForbidden, Message: "Access denied"},
	InsufficientScope:        ErrorCode{ID: "20012", Status: http.StatusForbidden, Message: "Insufficient scope"},
}
