package ec2cli

import "fmt"

func PrintInstancesByTags(getId bool, searchTags []string) {
	// find instance by tags
	instances := FindEC2InstancesByTagValues(searchTags, AWSClientImpl{Client: getAWSEC2Client()})

	// print ids if required
	if getId {
		printInstanceIdsByTag(instances)
	}

	// print tabs if required
	tabPrintInstances(instances)
}

func printInstanceIdsByTag(instances []InstanceResult) {
	for _, instance := range instances {
		fmt.Println(instance.Id)
	}
}
