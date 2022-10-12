package ec2cli

import "testing"
import "github.com/stretchr/testify/assert"

func TestShouldListEC2InstancesIdsByPrivateIP(t *testing.T) {
	// Given
	expectedInstanceId := []string{"InstanceId"}
	// When
	instanceIds := FindEC2InstanceIds(testIP(), testAWSClientWithPrivateIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceIds, "they should be equal")
}

func TestShouldListEC2InstanceIdsByPublicIP(t *testing.T) {
	// Given
	expectedInstanceId := []string{"instanceId"}
	// When
	instanceIds := FindEC2InstanceIds(testIP(), testAWSClientWithPublicIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceIds, "they should be equal")
}

func TestShouldReturnEmptyResultInstanceNotFound(t *testing.T) {
	// Given
	// When
	instanceIds := FindEC2InstanceIds(testIP(), TestAWSClient{})
	// then
	assert.Equal(t, 0, len(instanceIds))
}

func TestShouldReturnMultipleInstancesWhenFoundForPrivateIPs(t *testing.T) {
	// unlikely scenario but implementation does not need to make distinction between one and many.
	// Given
	expectedInstanceId := []string{"InstanceId1", "InstanceId2"}
	// When
	instanceIds := FindEC2InstanceIds(testIP(), testAWSClientWithPrivateIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceIds, "they should be equal")
}

func TestShouldReturnMultipleInstancesWhenFoundForPublicIPs(t *testing.T) {
	// unlikely scenario but implementation does not need to make distinction between one and many.
	// Given
	expectedInstanceId := []string{"InstanceId1", "InstanceId2"}
	// When
	instanceIds := FindEC2InstanceIds(testIP(), testAWSClientWithPublicIP(testIP(), expectedInstanceId))
	// then
	assert.Equal(t, expectedInstanceId, instanceIds, "they should be equal")
}

func testIP() string {
	return "172.16.102.74"
}

func testAWSClientWithPrivateIP(ip string, instanceIds []string) TestAWSClient {
	return TestAWSClient{ExpectedPrivateResult: getInstanceResult(ip, instanceIds)}
}

func getInstanceResult(ip string, instanceIds []string) []InstanceResult {
	var result []InstanceResult
	for _, instanceId := range instanceIds {
		result = append(result, InstanceResult{Id: instanceId, Name: "", IP: ip})
	}
	return result
}

func testAWSClientWithPublicIP(ip string, expectedInstanceIds []string) TestAWSClient {
	return TestAWSClient{ExpectedPublicResult: getInstanceResult(ip, expectedInstanceIds)}
}

type TestAWSClient struct {
	ExpectedPrivateResult []InstanceResult
	ExpectedPublicResult  []InstanceResult
}

func (t TestAWSClient) GetInstancesWithPrivateIP(_ string, _ string) []InstanceResult {
	return t.ExpectedPrivateResult
}

func (t TestAWSClient) GetInstancesWithPublicIP(_ string, _ string) []InstanceResult {
	return t.ExpectedPublicResult
}

func (t TestAWSClient) GetInstancesWithTags(_ []string) []InstanceResult {
	//TODO implement me
	panic("implement me")
}
