package config

import "github.com/cosmotek/manifest/scanner/aws"

type Config struct {
	ScanCronSchedule string `split_words:"true" required:"false" default:"@daily"`
	MaxAsyncRoutines uint64 `split_words:"true" required:"false" default:"2"`

	AwsScanner aws.Config `split_words:"true" required:"true"`
}
