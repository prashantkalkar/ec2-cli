package ec2cli

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

type Ec2api interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type AWSClientImpl struct {
	Client  Ec2api
	Verbose bool
}

func (c AWSClientImpl) GetInstancesWithPrivateIP(ip string, nameKey string) []InstanceResult {
	return c.ec2ResultToInstanceResult(func(instance types.Instance) string { return ip }, nameKey, c.getInstancesForPrivateIP(ip))
}

func (c AWSClientImpl) ec2ResultToInstanceResult(ipfunc func(instance types.Instance) string, nameKey string, awsResult *ec2.DescribeInstancesOutput) []InstanceResult {
	var instanceResults []InstanceResult
	for _, reservation := range awsResult.Reservations {
		for _, instance := range reservation.Instances {
			if c.Verbose {
				str, _ := PrettyStruct(instance)
				log.Printf("Converting instance details in instance result: %s", str)
			}
			instanceResults = append(instanceResults, ExtractInstanceResult(instance, ipfunc(instance), nameKey))
		}
	}

	return instanceResults
}

func (c AWSClientImpl) GetInstancesWithPublicIP(ip string, nameKey string) []InstanceResult {
	return c.ec2ResultToInstanceResult(func(instance types.Instance) string { return ip }, nameKey, c.getInstancesForPublicIP(ip))
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (c AWSClientImpl) GetAllInstances(nameKey string) []InstanceResult {
	instances := c.getAllInstances()
	if c.Verbose {
		str, _ := PrettyStruct(instances)
		log.Printf("Found instances, before filtering: %s", str)
	}
	return c.ec2ResultToInstanceResult(func(instance types.Instance) string { return *instance.PrivateIpAddress }, nameKey, instances)
}

func ExtractInstanceResult(instance types.Instance, ip string, nameKey string) InstanceResult {
	return InstanceResult{
		Id:        *instance.InstanceId,
		Name:      getInstanceName(instance, nameKey),
		IP:        ip,
		TagValues: getInstanceTagValues(instance),
	}
}

func getInstanceTagValues(instance types.Instance) []string {
	var tagValues []string
	for _, tag := range instance.Tags {
		tagValues = append(tagValues, *tag.Value)
	}
	return tagValues
}

func getInstanceName(instance types.Instance, nameKey string) string {
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
				Name:   aws.String("private-ip-address"),
				Values: []string{privateIP},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (c AWSClientImpl) getInstancesForPublicIP(publicIP string) *ec2.DescribeInstancesOutput {
	result, err := c.Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("ip-address"),
				Values: []string{publicIP},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (c AWSClientImpl) getAllInstances() *ec2.DescribeInstancesOutput {
	var runningInstancesFilter = []types.Filter{{Name: aws.String("instance-state-name"), Values: []string{"running"}}}
	result, err := c.Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{Filters: runningInstancesFilter})

	if err != nil {
		log.Fatal(err)
	}
	return result
}
