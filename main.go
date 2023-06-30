package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/cosmotek/manifest/config"
	"github.com/cosmotek/manifest/notifications"
	"github.com/cosmotek/manifest/scanner"
	"github.com/cosmotek/manifest/scanner/aws"
	"github.com/cosmotek/manifest/state"

	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron/v3"
)

const assetKey = "current_assets"

type ManifestService struct {
	notifications.NotifiationProvider
	state.StorageProvider
	scanner.EnvironmentScanner
}

func main() {
	conf := config.Config{}
	err := envconfig.Process("manifest", &conf)
	if err != nil {
		log.Fatalln(err)
	}

	notifier, err := notifications.NewSNSNotifier(conf.SNSNotifier)
	if err != nil {
		log.Fatalln(err)
	}

	envScanner, err := aws.New(conf.AwsScanner)
	if err != nil {
		log.Fatalln(err)
	}

	s3Bucket, err := state.NewS3Bucket(conf.S3Bucket)
	if err != nil {
		log.Fatalln(err)
	}

	manifest := ManifestService{
		NotifiationProvider: notifier,
		EnvironmentScanner:  envScanner,
		StorageProvider:     s3Bucket,
	}

	currentAssets := scanner.AssetList{}
	err = s3Bucket.Get(assetKey, &currentAssets)
	if err != nil {
		log.Fatalln(err)
	}

	schedule := cron.New()
	schedule.AddFunc(conf.ScanCronSchedule, func() {
		assets, err := manifest.EnvironmentScanner.RunScan()
		if err != nil {
			log.Fatalln(err)
		}

		// sort before diffing to prevent incidental diffing issues
		sort.Sort(assets)

		diff := scanner.ComputeDiff(currentAssets, assets)
		if diff != nil {
			manifest.NotifiationProvider.Notify(fmt.Sprintf("diff detected:\n\n%s\n\n", *diff))
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Println("no diff detected", len(currentAssets))
		}

		currentAssets = assets
		manifest.StorageProvider.Put(assetKey, currentAssets)
		if err != nil {
			log.Fatalln(err)
		}
	})

	log.Printf("starting manifest scheduler with cron '%s'...\n", conf.ScanCronSchedule)
	schedule.Run()
}
