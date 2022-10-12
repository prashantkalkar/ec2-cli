package ec2cli

import (
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
