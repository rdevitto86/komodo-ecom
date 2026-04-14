package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

// Publisher routes validated event envelopes to SNS FIFO topics.
// One instance is shared for the lifetime of the process.
type Publisher struct {
	sns            *sns.Client
	topicARNPrefix string // e.g. "arn:aws:sns:us-east-1:123456789012:komodo-"
	env            string // e.g. "prod", "staging", "dev", "local"
}

func NewPublisher(snsClient *sns.Client, topicARNPrefix string) *Publisher {
	return &Publisher{
		sns:            snsClient,
		topicARNPrefix: topicARNPrefix,
		env:            os.Getenv("ENV"),
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
