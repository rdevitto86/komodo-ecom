package httperrors

import "net/http"

type ErrorCode struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// 10xxx errors
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
	MockNotFound        ErrorCode
}

var Global = GlobalErrors{
	BadRequest: 					ErrorCode{ID: "10001", Status: http.StatusBadRequest, Message: "Bad request"},
	Unauthorized: 				ErrorCode{ID: "10002", Status: http.StatusUnauthorized, Message: "Unauthorized"},
	PaymentRequired: 			ErrorCode{ID: "10003", Status: http.StatusPaymentRequired, Message: "Payment required"},
	Forbidden: 						ErrorCode{ID: "10004", Status: http.StatusForbidden, Message: "Forbidden"},
	NotFound: 						ErrorCode{ID: "10005", Status: http.StatusNotFound, Message: "Not found"},
	MethodNotAllowed: 		ErrorCode{ID: "10006", Status: http.StatusMethodNotAllowed, Message: "Method not allowed"},
	Conflict: 						ErrorCode{ID: "10007", Status: http.StatusConflict, Message: "Conflict"},
	UnprocessableEntity: 	ErrorCode{ID: "10008", Status: http.StatusUnprocessableEntity, Message: "Unprocessable entity"},
	TooManyRequests: 			ErrorCode{ID: "10009", Status: http.StatusTooManyRequests, Message: "Too many requests"},
	Internal: 						ErrorCode{ID: "10010", Status: http.StatusInternalServerError, Message: "Internal server error"},
	NotImplemented: 			ErrorCode{ID: "10011", Status: http.StatusNotImplemented, Message: "Not implemented"},
	BadGateway: 					ErrorCode{ID: "10012", Status: http.StatusBadGateway, Message: "Bad gateway"},
	ServiceUnavailable: 	ErrorCode{ID: "10013", Status: http.StatusServiceUnavailable, Message: "Service unavailable"},
	GatewayTimeout: 			ErrorCode{ID: "10014", Status: http.StatusGatewayTimeout, Message: "Gateway timeout"},
}

// 11xxx errors
type DBErrors struct {
	ConnectionFailed ErrorCode
	QueryFailed      ErrorCode
	TransactionFailed ErrorCode
	RecordNotFound   ErrorCode
	DuplicateEntry   ErrorCode
}

var DB = DBErrors{
	ConnectionFailed: ErrorCode{ID: "110001", Status: http.StatusInternalServerError, Message: "Database connection failed"},
	QueryFailed:      ErrorCode{ID: "110002", Status: http.StatusInternalServerError, Message: "Database query failed"},
	TransactionFailed: ErrorCode{ID: "110003", Status: http.StatusInternalServerError, Message: "Database transaction failed"},
	RecordNotFound:   ErrorCode{ID: "110004", Status: http.StatusNotFound, Message: "Database record not found"},
	DuplicateEntry:   ErrorCode{ID: "110005", Status: http.StatusConflict, Message: "Database duplicate entry"},
}

// 20xxx errors
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
	InvalidKey:								ErrorCode{ID: "20005", Status: http.StatusUnauthorized, Message: "Invalid auth key"},
	ExpiredToken:             ErrorCode{ID: "20006", Status: http.StatusUnauthorized, Message: "Token expired"},
	UnauthorizedClient:       ErrorCode{ID: "20007", Status: http.StatusUnauthorized, Message: "Unauthorized client"},
	UnsupportedGrantType:     ErrorCode{ID: "20008", Status: http.StatusBadRequest, Message: "Unsupported grant type"},
	UnsupportedResponseType:  ErrorCode{ID: "20009", Status: http.StatusBadRequest, Message: "Unsupported response type"},
	InvalidRedirectURI:       ErrorCode{ID: "20010", Status: http.StatusBadRequest, Message: "Invalid redirect URI"},
	AccessDenied:             ErrorCode{ID: "20011", Status: http.StatusForbidden, Message: "Access denied"},
	InsufficientScope:        ErrorCode{ID: "20012", Status: http.StatusForbidden, Message: "Insufficient scope"},
}

// 30xxx errors
type UserErrors struct {
	NotFound         		ErrorCode
	AlreadyExists    		ErrorCode
	AccountLocked    		ErrorCode
	AccountSuspended 		ErrorCode
	EmailNotVerified 		ErrorCode
	PhoneNotVerified 		ErrorCode
	InvalidCredentials 	ErrorCode
	PasswordExpired 		ErrorCode
	WeakPassword 				ErrorCode
	MFARequired 				ErrorCode
	InvalidMFACode 			ErrorCode
}

var User = UserErrors{
	NotFound:         	ErrorCode{ID: "100001", Status: http.StatusNotFound, Message: "User not found"},
	AlreadyExists:    	ErrorCode{ID: "100002", Status: http.StatusConflict, Message: "User already exists"},
	AccountLocked:    	ErrorCode{ID: "100003", Status: http.StatusForbidden, Message: "Account locked"},
	AccountSuspended: 	ErrorCode{ID: "100004", Status: http.StatusForbidden, Message: "Account suspended"},
	EmailNotVerified: 	ErrorCode{ID: "100005", Status: http.StatusForbidden, Message: "Email not verified"},
	PhoneNotVerified: 	ErrorCode{ID: "100006", Status: http.StatusForbidden, Message: "Phone not verified"},
	InvalidCredentials: ErrorCode{ID: "100007", Status: http.StatusUnauthorized, Message: "Invalid credentials"},
	PasswordExpired:  	ErrorCode{ID: "100008", Status: http.StatusForbidden, Message: "Password expired"},
	WeakPassword:     	ErrorCode{ID: "100009", Status: http.StatusBadRequest, Message: "Weak password"},
	MFARequired:      	ErrorCode{ID: "100010", Status: http.StatusForbidden, Message: "Multi-factor authentication required"},
	InvalidMFACode:   	ErrorCode{ID: "100011", Status: http.StatusUnauthorized, Message: "Invalid MFA code"},
}

// 40xxx errors
type OrderErrors struct {
	NotFound ErrorCode
}

var Order = OrderErrors{
	NotFound: ErrorCode{ID: "100001", Status: http.StatusNotFound, Message: "Order not found"},
}

// 41xxx errors
type OrderItemErrors struct {
	NotFound ErrorCode
}

var OrderItem = OrderItemErrors{
	NotFound: ErrorCode{ID: "100001", Status: http.StatusNotFound, Message: "Order item not found"},
}

// 50xxx errors
type PaymentErrors struct {
	InsufficientFunds	ErrorCode
	Declined					ErrorCode
	MethodInvalid			ErrorCode
	TransactionFailed	ErrorCode
	RefundFailed			ErrorCode
	ProviderError			ErrorCode
}

var Payment = PaymentErrors{
	InsufficientFunds: ErrorCode{ID: "120001", Status: http.StatusPaymentRequired, Message: "Insufficient funds"},
	Declined:          ErrorCode{ID: "120002", Status: http.StatusPaymentRequired, Message: "Payment declined"},
	MethodInvalid:     ErrorCode{ID: "120003", Status: http.StatusBadRequest, Message: "Payment method invalid"},
	TransactionFailed: ErrorCode{ID: "120004", Status: http.StatusInternalServerError, Message: "Transaction failed"},
	RefundFailed:      ErrorCode{ID: "120005", Status: http.StatusInternalServerError, Message: "Refund failed"},
	ProviderError:     ErrorCode{ID: "120006", Status: http.StatusBadGateway, Message: "Payment provider error"},
}

// 60xxx errors
type ShopItemErrors struct {
	ItemNotFound         ErrorCode
	InventoryUnavailable ErrorCode
	InvalidSKU           ErrorCode
	SuggestionFailed     ErrorCode
	StorageError         ErrorCode
}

var ShopItem = ShopItemErrors{
	ItemNotFound:         ErrorCode{ID: "60001", Status: http.StatusNotFound, Message: "Shop item not found"},
	InventoryUnavailable: ErrorCode{ID: "60002", Status: http.StatusServiceUnavailable, Message: "Inventory data unavailable"},
	InvalidSKU:           ErrorCode{ID: "60003", Status: http.StatusBadRequest, Message: "Invalid SKU format"},
	SuggestionFailed:     ErrorCode{ID: "60004", Status: http.StatusInternalServerError, Message: "Failed to generate suggestions"},
	StorageError:         ErrorCode{ID: "60005", Status: http.StatusInternalServerError, Message: "Storage retrieval error"},
}
