package ec2cli

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyResultWhenNoInstancesFound(t *testing.T) {
	// Given
	// When
	instances := FindEC2Instances(testIP(), TestAWSClient{})
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestEmptyResultWhenInstanceWithPrivateIPFound(t *testing.T) {
	// Given
	// When
	instances := FindEC2Instances(testIP(), testAWSClientWithPrivateIP(testIP(), []string{"instanceId"}))
	// Then
	assert.Equal(t, InstanceResult{Id: "instanceId", Name: "", IP: testIP()}, instances[0])
}

func TestEmptyResultWhenInstanceWithPublicIPFound(t *testing.T) {
	// Given
	// When
	instances := FindEC2Instances(testIP(), testAWSClientWithPublicIP(testIP(), []string{"instanceId"}))
	// Then
	assert.Equal(t, InstanceResult{Id: "instanceId", Name: "", IP: testIP()}, instances[0])
}

func TestEmptyResultWhenNoTagsAreProvided(t *testing.T) {
	// When
	client := TestAWSClient{AllInstancesResult: []InstanceResult{}}
	instances := FindEC2InstancesByTagValues([]string{}, client)
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestEmptyResultWhenNoTagsAreProvidedEvenWhenInstancesExists(t *testing.T) {
	// When
	someInstances := getInstanceResultsWithTagValues(1, "instanceId", 0, []string{})
	client := TestAWSClient{AllInstancesResult: someInstances}
	instances := FindEC2InstancesByTagValues([]string{}, client)
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestEmptyResultWhenNoInstancesExists(t *testing.T) {
	client := TestAWSClient{AllInstancesResult: []InstanceResult{}}
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1"}, client)
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestEmptyResultWhenInstancesDontHaveTags(t *testing.T) {
	noTagInstances := getInstanceResultsWithTagValues(1, "instanceId", 0, []string{})
	client := TestAWSClient{AllInstancesResult: noTagInstances}
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1"}, client)
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestEmptyResultWhenInstancesDontMatchTags(t *testing.T) {
	nonMatchingInstances := getInstanceResultsWithTagValues(1, "instanceId", 0, []string{"tagV2"})
	client := TestAWSClient{AllInstancesResult: nonMatchingInstances}
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1"}, client)
	// Then
	assert.Equal(t, 0, len(instances))
}

func TestShouldProvideInstancesWhenInstancesMatchTagSingleSearchTag(t *testing.T) {
	matchingTagInstances := getInstanceResultsWithTagValues(2, "instanceId", 0, []string{"tagV1"})
	client := TestAWSClient{AllInstancesResult: matchingTagInstances}
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1"}, client)
	// Then
	assert.Equal(t, 2, len(instances))
	assert.Equal(t, "instanceId-0", instances[0].Id)
	assert.Equal(t, "instanceId-1", instances[1].Id)

	assert.Equal(t, "", instances[0].Name)
	assert.Equal(t, "", instances[1].Name)

	assert.Equal(t, "172.168.20.0", instances[0].IP)
	assert.Equal(t, "172.168.20.1", instances[1].IP)
}

func TestShouldProvideInstancesWhenInstancesMatchTagMultiSearchTag(t *testing.T) {
	matchingTagInstances := getInstanceResultsWithTagValues(2, "instanceId",
		0, []string{"tagV1", "tagV2", "tagV3"})
	client := TestAWSClient{AllInstancesResult: matchingTagInstances}
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1", "tagV2"}, client)
	// Then
	assert.Equal(t, 2, len(instances))
	assert.Equal(t, "instanceId-0", instances[0].Id)
	assert.Equal(t, "instanceId-1", instances[1].Id)

	assert.Equal(t, "", instances[0].Name)
	assert.Equal(t, "", instances[1].Name)

	assert.Equal(t, "172.168.20.0", instances[0].IP)
	assert.Equal(t, "172.168.20.1", instances[1].IP)
}

func TestShouldProvideOnlyInstancesWhenInstancesMatchTagMultiSearchTag(t *testing.T) {
	nonMatchingInstances := getInstanceResultsWithTagValues(2, "NonInstanceId",
		10, []string{"otherTag1", "otherTag2"})
	matchingTagInstances := getInstanceResultsWithTagValues(2, "instanceId",
		0, []string{"tagV1", "tagV2", "tagV3"})
	client := testAWSClientWithAllInstances(nonMatchingInstances, matchingTagInstances)
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1", "tagV2"}, client)
	// Then
	assert.Equal(t, 2, len(instances))
	assert.Equal(t, "instanceId-0", instances[0].Id)
	assert.Equal(t, "instanceId-1", instances[1].Id)

	assert.Equal(t, "", instances[0].Name)
	assert.Equal(t, "", instances[1].Name)

	assert.Equal(t, "172.168.20.0", instances[0].IP)
	assert.Equal(t, "172.168.20.1", instances[1].IP)
}

func TestShouldMatchTagCaseInsensitively(t *testing.T) {
	nonMatchingInstances := getInstanceResultsWithTagValues(2, "NonInstanceId",
		10, []string{"otherTag1", "otherTag2"})
	matchingTagInstances := getInstanceResultsWithTagValues(2, "instanceId",
		0, []string{"TagV1", "TAGV2", "tAgV3"})
	client := testAWSClientWithAllInstances(nonMatchingInstances, matchingTagInstances)
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1", "tagV2"}, client)
	// Then
	assert.Equal(t, 2, len(instances))
	assert.Equal(t, "instanceId-0", instances[0].Id)
	assert.Equal(t, "instanceId-1", instances[1].Id)

	assert.Equal(t, "", instances[0].Name)
	assert.Equal(t, "", instances[1].Name)

	assert.Equal(t, "172.168.20.0", instances[0].IP)
	assert.Equal(t, "172.168.20.1", instances[1].IP)
}

func TestShouldMatchInstanceTagPartially(t *testing.T) {
	nonMatchingInstances := getInstanceResultsWithTagValues(2, "NonInstanceId",
		10, []string{"otherTag1", "otherTag2"})
	matchingTagInstances := getInstanceResultsWithTagValues(2, "instanceId",
		0, []string{"SomeTagV1thing", "sTAGV2omething", "tAgV3"})
	client := testAWSClientWithAllInstances(nonMatchingInstances, matchingTagInstances)
	// When
	instances := FindEC2InstancesByTagValues([]string{"tagV1", "tagV2"}, client)
	// Then
	assert.Equal(t, 2, len(instances))
	assert.Equal(t, "instanceId-0", instances[0].Id)
	assert.Equal(t, "instanceId-1", instances[1].Id)

	assert.Equal(t, "", instances[0].Name)
	assert.Equal(t, "", instances[1].Name)

	assert.Equal(t, "172.168.20.0", instances[0].IP)
	assert.Equal(t, "172.168.20.1", instances[1].IP)
}

func testAWSClientWithAllInstances(instanceGroup1 []InstanceResult, instanceGroup2 []InstanceResult) TestAWSClient {
	return TestAWSClient{AllInstancesResult: append(instanceGroup1, instanceGroup2...)}
}

func getInstanceResultsWithTagValues(numOfInstances int, instanceIPPrefix string, incrementalIPLastDigit int, tagValues []string) []InstanceResult {
	var allInstances []InstanceResult
	for i := 0; i < numOfInstances; i++ {
		allInstances = append(allInstances, InstanceResult{
			Id:        fmt.Sprintf("%s-%d", instanceIPPrefix, i),
			Name:      "",
			IP:        fmt.Sprintf("172.168.20.%d", (incrementalIPLastDigit+i)%256),
			TagValues: tagValues,
		})
	}
	return allInstances
}
