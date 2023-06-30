package config

import (
	"github.com/cosmotek/manifest/notifications"
	"github.com/cosmotek/manifest/scanner/aws"
	"github.com/cosmotek/manifest/state"
)

type Config struct {
	ScanCronSchedule string `split_words:"true" required:"false" default:"@daily"`
	MaxAsyncRoutines uint64 `split_words:"true" required:"false" default:"2"`

	AwsScanner  aws.Config                          `split_words:"true" required:"true"`
	SNSNotifier notifications.SNSNotificationConfig `envconfig:"SNS_NOTIFIER" required:"true"`
	S3Bucket    state.S3BucketConfig                `envconfig:"S3_BUCKET" required:"true"`
}
