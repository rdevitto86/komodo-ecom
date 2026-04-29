package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"komodo-events-api/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// Dispatcher is the minimal interface the Publisher requires from the dispatch package.
// Declared here to break the import cycle: relay ← dispatch ← relay.
type Dispatcher interface {
	Dispatch(ctx context.Context, env EventEnvelope) error
}

// DynamoRepo is the minimal interface the Publisher requires from the repo package.
type DynamoRepo interface {
	SaveEvent(ctx context.Context, env EventEnvelope) error
}

// Publisher routes validated event envelopes to the active transport.
// transport="dynamo" (default) persists to DynamoDB and HTTP-dispatches to subscribers.
// transport="sns" publishes to SNS FIFO topics (V2 path, not yet deployed).
type Publisher struct {
	sns            *sns.Client
	topicARNPrefix string
	env            string
	repo           DynamoRepo
	dispatcher     Dispatcher
	transport      string
}

func NewPublisher(snsClient *sns.Client, topicARNPrefix string, rep DynamoRepo, disp Dispatcher, transport string) *Publisher {
	return &Publisher{
		sns:            snsClient,
		topicARNPrefix: topicARNPrefix,
		env:            os.Getenv(config.ENV),
		repo:           rep,
		dispatcher:     disp,
		transport:      transport,
	}
}

// Publish sends the envelope to the correct SNS FIFO topic for its domain.
// Topic ARN is constructed as: <prefix><domain>-events-<env>.fifo
// MessageGroupId is the domain — preserves per-domain ordering while allowing
// cross-domain parallelism. MessageDeduplicationId is the event ID.
// Returns the SNS MessageId on success.
func (p *Publisher) Publish(ctx context.Context, env EventEnvelope) (string, error) {
	body, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}

	domain := domainFromType(string(env.Type))
	topicARN := fmt.Sprintf("%s%s-events-%s.fifo", p.topicARNPrefix, domain, p.env)

	out, err := p.sns.Publish(ctx, &sns.PublishInput{
		TopicArn:               aws.String(topicARN),
		Message:                aws.String(string(body)),
		MessageGroupId:         aws.String(domain),
		MessageDeduplicationId: aws.String(env.ID),
		MessageAttributes: map[string]snsTypes.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(env.Type)),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("sns publish to %s: %w", topicARN, err)
	}

	return aws.ToString(out.MessageId), nil
}
