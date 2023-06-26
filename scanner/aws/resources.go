package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourceexplorer2"
)

func getResources(rsrc *resourceexplorer2.ResourceExplorer2, region, resourceType string) ([]*resourceexplorer2.Resource, error) {
	query := fmt.Sprintf("arn region:%s resourcetype:%s", region, resourceType)

	resources := []*resourceexplorer2.Resource{}
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
