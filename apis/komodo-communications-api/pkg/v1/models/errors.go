package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 70xxx — komodo-communications-api (see forge-sdk ranges.go)
type CommunicationsAPIErrors struct {
	DeliveryFailed   httpErr.ErrorCode
	InvalidRecipient httpErr.ErrorCode
	ProviderError    httpErr.ErrorCode
	TemplateNotFound httpErr.ErrorCode
}

var Err = CommunicationsAPIErrors{
	DeliveryFailed:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCommunications, 1), Status: http.StatusBadGateway, Message: "Message delivery failed"},
	InvalidRecipient: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCommunications, 2), Status: http.StatusUnprocessableEntity, Message: "Invalid recipient"},
	ProviderError:    httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCommunications, 3), Status: http.StatusBadGateway, Message: "Communications provider error"},
	TemplateNotFound: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCommunications, 4), Status: http.StatusNotFound, Message: "Message template not found"},
}
