package notifications

type NotifiationProvider interface {
	Notify(message string) error
}

// func NotifyAssetChange(svc *sns.SNS, topic, msg string) error {
// 	result, err := svc.Publish(&sns.PublishInput{
// 		Message:  aws.String(msg),
// 		TopicArn: aws.String(topic),
// 	})
// 	if err != nil {
// 		return nil
// 	}

// 	return nil
// }
