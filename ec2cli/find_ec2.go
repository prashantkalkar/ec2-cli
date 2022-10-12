package ec2cli

import "golang.org/x/exp/slices"

func FindEC2Instances(ip string, client AWSClient) []InstanceResult {
	privateIPResult := client.GetInstancesWithPrivateIP(ip, "Name")
	if len(privateIPResult) > 0 {
		return privateIPResult
	}

	publicIPResult := client.GetInstancesWithPublicIP(ip, "Name")
	return publicIPResult
}

func FindEC2InstancesByTagValues(tagValues []string, client AWSClient) []InstanceResult {
	if len(tagValues) > 0 {
		return findInstancesForNonEmptySearchTags(tagValues, client)
	}

	return []InstanceResult{}
}

func findInstancesForNonEmptySearchTags(tagValues []string, client AWSClient) []InstanceResult {
	var result []InstanceResult
	instances := client.GetAllInstances("Name")
	for _, instance := range instances {
		if hasTagValues(instance, tagValues) {
			result = append(result, instance)
		}
	}
	return result
}

func hasTagValues(instance InstanceResult, values []string) bool {
	for _, value := range values {
		if !slices.Contains(instance.TagValues, value) {
			return false
		}
	}
	return true
}
