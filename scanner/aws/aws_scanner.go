package aws

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourceexplorer2"
	"github.com/samber/lo"

	"github.com/cosmotek/manifest/async"
	"github.com/cosmotek/manifest/scanner"
)

type Scanner struct {
	resourceExplorerService *resourceexplorer2.ResourceExplorer2
	scanRegions             []string
	scanResourceTypes       []string

	conf Config
}

type Config struct {
	ScanRegions               []string `split_words:"true" required:"false"`
	ResourceExplorerAWSRegion string   `split_words:"true" required:"false" default:"us-east-1"`
	MaxConcurrency            uint64   `split_words:"true" required:"false" default:"2"`
}

func New(conf Config) (*Scanner, error) {
	if len(conf.ScanRegions) == 0 {
		conf.ScanRegions = awsRegions
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(conf.ResourceExplorerAWSRegion),
	})
	if err != nil {
		return nil, err
	}

	svc := resourceexplorer2.New(session)
	types, err := getResourceTypes(svc)
	if err != nil {
		return nil, err
	}

	return &Scanner{
		resourceExplorerService: svc,
		scanRegions:             conf.ScanRegions,
		scanResourceTypes: lo.Map(types, func(resource *resourceexplorer2.SupportedResourceType, _ int) string {
			return *resource.ResourceType
		}),
		conf: conf,
	}, nil
}

func (s *Scanner) RunScan() (scanner.AssetList, error) {
	writeLock := sync.Mutex{}
	assets := scanner.AssetList{}
	wg := async.NewSizedWaitGroup(s.conf.MaxConcurrency)

	for _, region := range s.scanRegions {
		go func() {
			wg.Add()
			defer wg.Done()

			for _, resourceType := range s.scanResourceTypes {
				fmt.Println("getting resources for query", fmt.Sprintf("arn region:%s resourcetype:%s", region, resourceType))

				resourcesReturned, err := getResources(s.resourceExplorerService, region, resourceType)
				if err != nil {
					panic(err)
				}

				assetsReturned := lo.Map(resourcesReturned, func(resource *resourceexplorer2.Resource, _ int) scanner.Asset {
					return scanner.Asset{
						Identifier: *resource.Arn,
						Metadata: map[string]any{
							"aws_account_id": resource.OwningAccountId,
							"region":         resource.Region,
							"props":          resource.Properties,
							"service":        resource.Service,
							"resource_type":  resource.ResourceType,
						},
					}
				})

				fmt.Printf("adding %d resources\n", len(resourcesReturned))
				writeLock.Lock()
				assets = append(assets, assetsReturned...)
				writeLock.Unlock()
			}
		}()
	}

	return assets, nil
}
