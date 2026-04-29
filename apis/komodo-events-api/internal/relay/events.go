package relay

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-events-api/internal/models"
)

// PublishEvent handles POST /events.
// Requires a valid internal service JWT (svc: scope). Returns 202 on success.
func (p *Publisher) PublishEvent(wtr http.ResponseWriter, req *http.Request) {
	var env EventEnvelope
	if err := json.NewDecoder(req.Body).Decode(&env); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail(err.Error()))
		return
	}

	if env.ID == "" || string(env.Type) == "" || env.Version == "" ||
		string(env.Source) == "" || env.OccurredAt.IsZero() || env.Payload == nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest,
			httpErr.WithDetail("id, type, version, source, occurred_at, and payload are required"))
		return
	}

	if _, ok := KnownEventTypes[env.Type]; !ok {
		httpErr.SendError(wtr, req, models.Err.UnknownType,
			httpErr.WithDetail("unrecognised event type: "+string(env.Type)))
		return
	}

	if p.transport == "sns" {
		msgID, err := p.Publish(req.Context(), env)
		if err != nil {
			logger.Error("failed to publish event", err,
				logger.Attr("event_id", env.ID),
				logger.Attr("event_type", string(env.Type)),
			)
			httpErr.SendError(wtr, req, models.Err.PublishFailed)
			return
		}

		logger.Info("event accepted",
			logger.Attr("event_id", env.ID),
			logger.Attr("event_type", string(env.Type)),
			logger.Attr("message_id", msgID),
		)

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusAccepted)
		json.NewEncoder(wtr).Encode(map[string]string{
			"event_id":   env.ID,
			"message_id": msgID,
		})
		return
	}

	// dynamo transport (default)
	if p.repo != nil {
		if err := p.repo.SaveEvent(req.Context(), env); err != nil {
			logger.Error("failed to persist event", err, logger.Attr("event_id", env.ID))
		}
	}
	if p.dispatcher != nil {
		if err := p.dispatcher.Dispatch(req.Context(), env); err != nil {
			logger.Error("dispatch failed", err, logger.Attr("event_id", env.ID))
		}
	}

	logger.Info("event accepted",
		logger.Attr("event_id", env.ID),
		logger.Attr("event_type", string(env.Type)),
	)

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusAccepted)
	json.NewEncoder(wtr).Encode(map[string]string{
		"event_id": env.ID,
	})
}
