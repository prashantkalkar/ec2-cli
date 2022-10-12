package ec2cli

import "testing"
import "github.com/stretchr/testify/assert"

func TestShouldListEC2InstancesIdsByPrivateIP(t *testing.T) {
	// Given
	expectedInstanceId := "InstanceId"
	// When
	instanceId, _ := FindEC2InstanceId(testIP(), true, testAWSClientWithPrivateIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceId, "they should be equal")
}

func TestShouldListEC2InstanceIdsByPublicIP(t *testing.T) {
	// Given
	expectedInstanceId := "instanceId"
	// When
	instanceId, _ := FindEC2InstanceId(testIP(), true, testAWSClientWithPublicIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceId, "they should be equal")
}

func TestShouldReturnErrorWhenInstanceNotFound(t *testing.T) {
	// Given
	// When
	_, err := FindEC2InstanceId(testIP(), true, TestAWSClient{})
	// then
	assert.NotNil(t, err, "Should generate error when IP not found")
}

func testIP() string {
	return "172.16.102.74"
}

func testAWSClientWithPrivateIP(ip string, instanceId string) TestAWSClient {
	return TestAWSClient{ExpectedPrivateResult: []InstanceResult{
		{
			Id: instanceId, Name: "", IP: ip,
		},
	}}
}

func testAWSClientWithPublicIP(ip string, expectedInstanceId string) TestAWSClient {
	return TestAWSClient{ExpectedPublicResult: []InstanceResult{
		{
			Id: expectedInstanceId, Name: "", IP: ip,
		},
	}}
}

type TestAWSClient struct {
	ExpectedPrivateResult []InstanceResult
	ExpectedPublicResult  []InstanceResult
}

func (t TestAWSClient) GetInstancesWithPrivateIP(_ string, _ string) []InstanceResult {
	return t.ExpectedPrivateResult
}

func (t TestAWSClient) GetInstancesWithPublicIP(_ string) []InstanceResult {
	return t.ExpectedPublicResult
}

func (t TestAWSClient) GetInstancesWithTags(_ []string) []InstanceResult {
	//TODO implement me
	panic("implement me")
}
