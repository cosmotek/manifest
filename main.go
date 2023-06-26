package main

import (
	"fmt"
	"sort"

	"github.com/cosmotek/manifest/config"
	"github.com/cosmotek/manifest/notifications"
	"github.com/cosmotek/manifest/scanner"
	"github.com/cosmotek/manifest/state"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron/v3"
	// "github.com/aws/aws-sdk-go/service/sns"
)

var currentAssets = scanner.AssetList{}

type ManifestService struct {
	notifications.NotifiationProvider
	state.StorageProvider
	scanner.EnvironmentScanner
}

func main() {
	conf := config.Config{}
	err := envconfig.Process("manifest", &conf)
	if err != nil {
		panic(err)
	}

	manifest := ManifestService{}

	schedule := cron.New()
	schedule.AddFunc(conf.ScanCronSchedule, func() {
		assets, err := manifest.EnvironmentScanner.RunScan()
		if err != nil {
			panic(err)
		}

		sort.Sort(assets)

		diff := scanner.ComputeDiff(currentAssets, assets)
		if diff != nil {
			manifest.NotifiationProvider.Notify(fmt.Sprintf("diff detected:\n\n%s\n\n", *diff))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("no diff detected", len(currentAssets))
		}

		currentAssets = assets
		manifest.StorageProvider.Put("assets", currentAssets)
		if err != nil {
			panic(err)
		}
	})

	fmt.Printf("starting manifest scheduler with cron '%s'...\n", conf.ScanCronSchedule)
	schedule.Run()
}
