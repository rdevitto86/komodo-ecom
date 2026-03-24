package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 50xxx — komodo-payments-api (see forge-sdk ranges.go)
type PaymentAPIErrors struct {
	InsufficientFunds httpErr.ErrorCode
	Declined          httpErr.ErrorCode
	MethodInvalid     httpErr.ErrorCode
	TransactionFailed httpErr.ErrorCode
	RefundFailed      httpErr.ErrorCode
	ProviderError     httpErr.ErrorCode
}

var Err = PaymentAPIErrors{
	InsufficientFunds: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 1), Status: http.StatusPaymentRequired, Message: "Insufficient funds"},
	Declined:          httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 2), Status: http.StatusPaymentRequired, Message: "Payment declined"},
	MethodInvalid:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 3), Status: http.StatusBadRequest, Message: "Payment method invalid"},
	TransactionFailed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 4), Status: http.StatusInternalServerError, Message: "Transaction failed"},
	RefundFailed:      httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 5), Status: http.StatusInternalServerError, Message: "Refund failed"},
	ProviderError:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangePayment, 6), Status: http.StatusBadGateway, Message: "Payment provider error"},
}
