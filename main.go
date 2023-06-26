package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	// "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourceexplorer2"
	// "github.com/aws/aws-sdk-go/service/sns"
)

var regions = []string{
	"global",
	// "af-south-1",
	// "ap-east-1",
	// "ap-northeast-1",
	// "ap-northeast-2",
	// "ap-northeast-3",
	// "ap-south-1",
	// "ap-south-2",
	// "ap-southeast-1",
	// "ap-southeast-2",
	// "ap-southeast-3",
	// "ap-southeast-4",
	// "ca-central-1",
	// "cn-norht-1",
	// "eu-central-1",
	// "eu-central-2",
	// "eu-north-1",
	// "eu-south-1",
	// "eu-south-2",
	// "eu-west-1",
	// "eu-west-2",
	// "eu-west-3",
	// "me-central-1",
	// "me-south-1",
	// "sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-gov-east-1",
	"us-gov-west-1",
	"us-west-1",
	"us-west-2",
}

var resourceTypes = []string{

	"cloudfront:cache-policy",

	"cloudfront:distribution",

	"cloudfront:function",

	"cloudfront:origin-access-identity",

	"cloudfront:origin-request-policy",

	"cloudfront:realtime-log-config",

	"cloudfront:response-headers-policy",

	"cloudwatch:alarm",

	"cloudwatch:dashboard",

	"cloudwatch:insight-rule",

	"cloudwatch:metric-stream",

	"logs:destination",

	"logs:log-group",
	"dynamodb:table",

	"elasticache:cluster",

	"elasticache:globalreplicationgroup",

	"elasticache:parametergroup",

	"elasticache:replicationgroup",

	"elasticache:reserved-instance",

	"elasticache:snapshot",

	"elasticache:subnetgroup",

	"elasticache:user",

	"elasticache:usergroup",

	"ec2:capacity-reservation",

	"ec2:capacity-reservation-fleet",

	"ec2:client-vpn-endpoint",

	"ec2:customer-gateway",

	"ec2:dedicated-host",

	"ec2:dhcp-options",

	"ec2:egress-only-internet-gateway",

	"ec2:elastic-gpu",

	"ec2:elastic-ip",

	"ec2:fleet",

	"ec2:fpga-image",

	"ec2:host-reservation",

	"ec2:image",

	"ec2:instance",

	"ec2:instance-event-window",

	"ec2:internet-gateway",

	"ec2:ipam",

	"ec2:ipam-pool",

	"ec2:ipam-scope",

	"ec2:ipv4pool-ec2",

	"ec2:key-pair",

	"ec2:launch-template",

	"ec2:natgateway",

	"ec2:network-acl",

	"ec2:network-insights-access-scope",

	"ec2:network-insights-access-scope-analysis",

	"ec2:network-insights-analysis",

	"ec2:network-insights-path",

	"ec2:network-interface",

	"ec2:placement-group",

	"ec2:prefix-list",

	"ec2:reserved-instances",

	"ec2:route-table",

	"ec2:security-group",

	"ec2:security-group-rule",

	"ec2:snapshot",

	"ec2:spot-fleet-request",

	"ec2:spot-instances-request",

	"ec2:subnet",

	"ec2:subnet-cidr-reservation",

	"ec2:traffic-mirror-filter",

	"ec2:traffic-mirror-filter-rule",

	"ec2:traffic-mirror-session",

	"ec2:traffic-mirror-target",

	"ec2:transit-gateway",

	"ec2:transit-gateway-attachment",

	"ec2:transit-gateway-connect-peer",

	"ec2:transit-gateway-multicast-domain",

	"ec2:transit-gateway-policy-table",

	"ec2:transit-gateway-route-table",

	"ec2:volume",

	"ec2:vpc",

	"ec2:vpc-endpoint",

	"ec2:vpc-flow-log",

	"ec2:vpc-peering-connection",

	"ec2:vpn-connection",

	"ec2:vpn-gateway",

	"ecs:cluster",

	"ecs:container-instance",

	"ecs:service",

	"ecs:task",

	"ecs:task-definition",

	"ecs:task-set",

	"elasticloadbalancing:listener",

	"elasticloadbalancing:listener-rule",

	"elasticloadbalancing:listener-rule/app",

	"elasticloadbalancing:listener/app",

	"elasticloadbalancing:listener/net",

	"elasticloadbalancing:loadbalancer",

	"elasticloadbalancing:loadbalancer/app",

	"elasticloadbalancing:loadbalancer/net",

	"elasticloadbalancing:targetgroup",

	"iam:group",

	"iam:instance-profile",

	"iam:oidc-provider",

	"iam:policy",

	"iam:role",

	"iam:saml-provider",

	"iam:server-certificate",

	"iam:user",

	"kinesis:stream",

	"lambda:code-signing-config",

	"lambda:event-source-mapping",

	"lambda:function",

	"es:domain",

	"redshift:cluster",

	"redshift:eventsubscription",

	"redshift:parametergroup",

	"redshift:snapshot",

	"redshift:snapshotcopygrant",

	"redshift:snapshotschedule",

	"redshift:subnetgroup",

	"redshift:usagelimit",

	"rds:auto-backup",

	"rds:cev",

	"rds:cluster",

	"rds:cluster-endpoint",

	"rds:cluster-pg",

	"rds:cluster-snapshot",

	"rds:db",

	"rds:db-proxy",

	"rds:db-proxy-endpoint",

	"rds:es",

	"rds:global-cluster",

	"rds:og",

	"rds:pg",

	"rds:ri",

	"rds:secgrp",

	"rds:snapshot",

	"rds:subgrp",

	"resource-explorer-2:index",

	"resource-explorer-2:view",

	"servicecatalog:applications",

	"servicecatalog:attribute-groups",

	"sns:topic",

	"sqs:queue",

	"s3:accesspoint",

	"s3:bucket",

	"s3:storage-lens",

	"ssm:association",

	"ssm:automation-execution",

	"ssm:document",

	"ssm:maintenancewindow",

	"ssm:managed-instance",

	"ssm:parameter",

	"ssm:patchbaseline",

	"ssm:windowtarget",

	"ssm:windowtask",
}

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
	// session, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-1"),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// // Create S3 service client
	// svc := s3.New(sess)

	ticker := time.NewTicker(scanInterval)
	// rsrc := resourceexplorer2.New(session)

	types, err := GetResourceTypes()
	if err != nil {
		panic(err)
	}

	for {
		resources := ResourceList{}

		for _, region := range regions {

			// TODO convert to a bunch of goroutines
			for _, resourceType := range types {
				resourcesReturned, err := GetResources(region, *resourceType.ResourceType)
				if err != nil {
					panic(err)
				}

				resources = append(resources, resourcesReturned...)
			}
		}

		fmt.Println("sorting resources")
		sort.Sort(resources)
		fmt.Println(resources, len(resources))

		if len(currentResources) > 0 && !cmp.Equal(currentResources, resources) {
			diff := cmp.Diff(currentResources, resources, compareOpts...)
			fmt.Println("diff found", diff)
			err = ioutil.WriteFile("./diff.txt", []byte(diff), os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("no diff found", len(currentResources))
		}

		currentResources = resources
		err = ioutil.WriteFile("./resources.txt", []byte(fmt.Sprint(resources)), os.ModePerm)
		if err != nil {
			panic(err)
		}

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

func GetResourceTypes() ([]*resourceexplorer2.SupportedResourceType, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}

	rsrc := resourceexplorer2.New(session)
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
		} else {
			nextToken = output.NextToken
		}
	}

	return types, nil
}

func GetResources(region, resourceType string) (ResourceList, error) {
	query := fmt.Sprintf("arn region:%s resourcetype:%s", region, resourceType)
	fmt.Println("getting resources for query", query)

	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, err
	}

	rsrc := resourceexplorer2.New(session)
	resources := ResourceList{}
	var token *string

	for {
		output, err := rsrc.Search(&resourceexplorer2.SearchInput{
			MaxResults:  aws.Int64(1000),
			NextToken:   token,
			QueryString: aws.String(query),
		})
		if err != nil {
			panic(err)

			// TODO check if UnauthorizedException: Unauthorized because
			// they may need to enable the resource explorer service
		}

		resources = append(resources, output.Resources...)
		if output.NextToken == nil || *output.NextToken == "" {
			break
		} else {
			token = output.NextToken
		}
	}

	return resources, nil
}
