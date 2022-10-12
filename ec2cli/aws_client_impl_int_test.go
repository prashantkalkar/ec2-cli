package ec2cli

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInstancesWithPrivateIPWhenReservationsDoesNotExists(t *testing.T) {
	//Given
	var awsEc2Response = ec2.DescribeInstancesOutput{Reservations: []types.Reservation{}}

	impl := AWSClientImpl{Client: mockDescribeInstances{awsEc2Response}}

	//When
	instanceResult := impl.GetInstancesWithPrivateIP("172.16.102.42", "Name")

	//Then
	assert.Equal(t, 0, len(instanceResult))
}

func TestGetInstancesWithPrivateIPWhenInstancesDoesNotExists(t *testing.T) {
	//Given
	var reservations = []types.Reservation{{Instances: []types.Instance{}}}
	var awsEc2Response = ec2.DescribeInstancesOutput{Reservations: reservations}

	impl := AWSClientImpl{Client: mockDescribeInstances{awsEc2Response}}

	//When
	instanceResult := impl.GetInstancesWithPrivateIP("172.16.102.42", "Name")

	//Then
	assert.Equal(t, 0, len(instanceResult))
}

func TestGetInstancesWithPrivateIPWhenInstanceExistsWithoutNameTag(t *testing.T) {
	//Given
	instanceId := "MyInstanceId"
	instanceIP := "172.16.102.42"
	instances := []types.Instance{{InstanceId: aws.String(instanceId), Tags: []types.Tag{}}}
	var reservations = []types.Reservation{{Instances: instances}}
	var awsEc2Response = ec2.DescribeInstancesOutput{Reservations: reservations}

	impl := AWSClientImpl{Client: mockDescribeInstances{awsEc2Response}}

	//When
	instanceResult := impl.GetInstancesWithPrivateIP(instanceIP, "Name")

	//Then
	assert.Equal(t, 1, len(instanceResult))
	assert.Equal(t, instanceId, instanceResult[0].Id)
	assert.Equal(t, instanceIP, instanceResult[0].IP)
	assert.Equal(t, "", instanceResult[0].Name)
}

func TestGetInstancesWithPrivateIPWhenInstanceExistsWithNameTag(t *testing.T) {
	//Given
	instanceId := "MyInstanceId"
	instanceIP := "172.16.102.42"
	instanceName := "MyInstanceName"

	tags := []types.Tag{{Key: aws.String("Name"), Value: aws.String(instanceName)}}
	instances := []types.Instance{{InstanceId: aws.String(instanceId), Tags: tags}}
	var reservations = []types.Reservation{{Instances: instances}}
	var awsEc2Response = ec2.DescribeInstancesOutput{Reservations: reservations}

	impl := AWSClientImpl{Client: mockDescribeInstances{awsEc2Response}}

	//When
	instanceResult := impl.GetInstancesWithPrivateIP(instanceIP, "Name")

	//Then
	assert.Equal(t, 1, len(instanceResult))
	assert.Equal(t, instanceId, instanceResult[0].Id)
	assert.Equal(t, instanceIP, instanceResult[0].IP)
	assert.Equal(t, instanceName, instanceResult[0].Name)
}

type mockDescribeInstances struct {
	describeInstancesOutput ec2.DescribeInstancesOutput
}

func (m mockDescribeInstances) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return &m.describeInstancesOutput, nil
}
