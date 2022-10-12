package ec2cli

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

type Ec2api interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type AWSClientImpl struct {
	Client Ec2api
}

func (c AWSClientImpl) GetInstancesWithPrivateIP(ip string, nameKey string) []InstanceResult {
	awsResult := c.getInstancesForPrivateIP(ip)
	var instanceResults []InstanceResult
	for _, reservation := range awsResult.Reservations {
		for _, instance := range reservation.Instances {
			instanceResults = append(instanceResults, ExtractInstanceResult(instance, ip, nameKey))
		}
	}

	return instanceResults
}

func (c AWSClientImpl) GetInstancesWithPublicIP(ip string, nameKey string) []InstanceResult {
	return nil
}

func (c AWSClientImpl) GetInstancesWithTags(tagValues []string) []InstanceResult {
	return nil
}

func ExtractInstanceResult(instance types.Instance, ip string, nameKey string) InstanceResult {
	return InstanceResult{
		Id:   *instance.InstanceId,
		Name: GetInstanceName(instance, nameKey),
		IP:   ip,
	}
}

func GetInstanceName(instance types.Instance, nameKey string) string {
	for _, tag := range instance.Tags {
		if nameKey == *tag.Key {
			return *tag.Value
		}
	}
	return ""
}

func (c AWSClientImpl) getInstancesForPrivateIP(privateIP string) *ec2.DescribeInstancesOutput {
	result, err := c.Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("network-interface.addresses.private-ip-address"),
				Values: []string{privateIP},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	return result
}
