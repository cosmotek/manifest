package notifications

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type NotifiationProvider interface {
	Notify(message string) error
}

type SNSNotificationConfig struct {
	SNSTopic string `split_words:"true" required:"true"`
}

func NewSNSNotifier(conf SNSNotificationConfig) (*SNSNotifier, error) {
	session, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}

	return &SNSNotifier{
		notificationTopic: conf.SNSTopic,
		svc:               sns.New(session),
	}, nil
}

type SNSNotifier struct {
	svc               *sns.SNS
	notificationTopic string
}

func (s *SNSNotifier) Notify(message string) error {
	result, err := s.svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(s.notificationTopic),
	})
	if err != nil {
		return err
	}

	log.Println("notified SNS topic", s.notificationTopic, "messageID:", result.MessageId)
	return nil
}
