package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 71xxx — komodo-events-api (see forge-sdk ranges.go)
type EventsAPIErrors struct {
	UnknownType   httpErr.ErrorCode
	PublishFailed httpErr.ErrorCode
	FanOutFailed  httpErr.ErrorCode
}

var Err = EventsAPIErrors{
	UnknownType:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEvents, 1), Status: http.StatusBadRequest, Message: "Unknown event type"},
	PublishFailed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEvents, 2), Status: http.StatusInternalServerError, Message: "Event publish failed"},
	FanOutFailed:  httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEvents, 3), Status: http.StatusInternalServerError, Message: "Event fan-out failed"},
}
