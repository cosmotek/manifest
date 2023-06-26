package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourceexplorer2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/resourceexplorer2"
	// "github.com/aws/aws-sdk-go/service/sns"
)

type ResourceList []*resourceexplorer2.Resource

var scanInterval = time.Second * 30
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

func main() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	// // Create S3 service client
	// svc := s3.New(sess)

	ticker := time.NewTicker(scanInterval)
	rsrc := resourceexplorer2.New(session)

	for {
		resources := ResourceList{}
		var token *string

		for {
			fmt.Println("fetching page of results")
			output, err := rsrc.Search(&resourceexplorer2.SearchInput{
				QueryString: aws.String("arn"),
				MaxResults:  aws.Int64(1000),
				NextToken:   token,
			})
			if err != nil {
				panic(err)
			}

			resources = append(resources, output.Resources...)
			if output.NextToken == nil {
				break
			} else {
				token = output.NextToken
			}
		}

		sort.Sort(resources)
		fmt.Println(resources, len(resources))

		if len(currentResources) > 0 && !cmp.Equal(currentResources, resources) {
			diff := cmp.Diff(resources, currentResources, compareOpts...)
			fmt.Println("diff found", diff)
		} else {
			fmt.Println("no diff found", len(currentResources))
		}

		currentResources = resources
		<-ticker.C
	}
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
