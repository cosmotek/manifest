package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"

	// "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourceexplorer2"
	// "github.com/aws/aws-sdk-go/service/sns"

	"github.com/robfig/cron/v3"
)

var regions = []string{
	"global",
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-south-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-southeast-3",
	"ap-southeast-4",
	"ca-central-1",
	"cn-norht-1",
	"eu-central-1",
	"eu-central-2",
	"eu-north-1",
	"eu-south-1",
	"eu-south-2",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"me-central-1",
	"me-south-1",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-gov-east-1",
	"us-gov-west-1",
	"us-west-1",
	"us-west-2",
}

type ResourceList []*resourceexplorer2.Resource

var currentResources = ResourceList{}
var compareOpts = []cmp.Option{
	cmpopts.IgnoreUnexported(
		resourceexplorer2.Resource{},
	),
	cmpopts.IgnoreFields(
		resourceexplorer2.Resource{},
		"LastReportedAt",
	),
}

func (r ResourceList) Len() int {
	return len(r)
}
func (r ResourceList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ResourceList) Less(i, j int) bool {
	return len(*r[i].Arn) < len(*r[j].Arn)
}

type Config struct {
	ScanCronSchedule string   `split_words:"true" required:"false" default:"@daily"`
	ScanRegions      []string `split_words:"true" required:"false"`

	ResourceExplorerAWSRegion string `split_words:"true" required:"false" default:"us-east-1"`
}

func main() {
	conf := Config{}
	err := envconfig.Process("manifest", &conf)
	if err != nil {
		panic(err)
	}

	if len(conf.ScanRegions) == 0 {
		conf.ScanRegions = regions
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(conf.ResourceExplorerAWSRegion),
	})
	if err != nil {
		panic(err)
	}
	rsrc := resourceexplorer2.New(session)

	types, err := GetResourceTypes(rsrc)
	if err != nil {
		panic(err)
	}

	schedule := cron.New()
	schedule.AddFunc(conf.ScanCronSchedule, func() {
		resources := ResourceList{}

		for _, region := range regions {

			// TODO convert to a bunch of goroutines
			for _, resourceType := range types {
				fmt.Println("getting resources for query", fmt.Sprintf("arn region:%s resourcetype:%s", region, *resourceType.ResourceType))

				resourcesReturned, err := GetResources(rsrc, region, *resourceType.ResourceType)
				if err != nil {
					panic(err)
				}

				fmt.Printf("adding %d resources\n", len(resourcesReturned))
				resources = append(resources, resourcesReturned...)
			}
		}

		fmt.Printf("found %d resources\n", len(resources))
		sort.Sort(resources)

		if len(currentResources) > 0 && !cmp.Equal(currentResources, resources) {
			diff := cmp.Diff(currentResources, resources, compareOpts...)
			fmt.Println("diff detected", diff)
			err = ioutil.WriteFile("./diff.txt", []byte(diff), os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("no diff detected", len(currentResources))
		}

		currentResources = resources
		err = ioutil.WriteFile("./resources.txt", []byte(fmt.Sprint(resources)), os.ModePerm)
		if err != nil {
			panic(err)
		}
	})

	fmt.Printf("starting manifest scheduler with cron '%s'...\n", conf.ScanCronSchedule)
	schedule.Run()
}

func NotifyAssetChange() {
	// svc := sns.New(sess)

	// result, err := svc.Publish(&sns.PublishInput{
	// 	Message:  msgPtr,
	// 	TopicArn: topicPtr,
	// })
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
}

func GetResourceTypes(rsrc *resourceexplorer2.ResourceExplorer2) ([]*resourceexplorer2.SupportedResourceType, error) {
	types := []*resourceexplorer2.SupportedResourceType{}

	var nextToken *string
	for {
		output, err := rsrc.ListSupportedResourceTypes(&resourceexplorer2.ListSupportedResourceTypesInput{
			MaxResults: aws.Int64(1000),
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, err
		}

		types = append(types, output.ResourceTypes...)
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}

		nextToken = output.NextToken
	}

	return types, nil
}

func GetResources(rsrc *resourceexplorer2.ResourceExplorer2, region, resourceType string) (ResourceList, error) {
	query := fmt.Sprintf("arn region:%s resourcetype:%s", region, resourceType)

	resources := ResourceList{}
	var token *string

	for {
		output, err := rsrc.Search(&resourceexplorer2.SearchInput{
			MaxResults:  aws.Int64(1000),
			NextToken:   token,
			QueryString: aws.String(query),
		})
		if err != nil {
			return nil, err

			// TODO check if UnauthorizedException: Unauthorized because
			// they may need to enable the resource explorer service
		}

		resources = append(resources, output.Resources...)
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}

		token = output.NextToken
	}

	return resources, nil
}

type WaitGroupCount struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroupCount) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupCount) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}
