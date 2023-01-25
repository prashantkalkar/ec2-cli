package ec2cli

import "fmt"

func PrintInstancesByTags(getId bool, searchTags []string, verbose bool) {
	// find instance by tags
	instances := FindEC2InstancesByTagValues(searchTags, AWSClientImpl{Client: getAWSEC2Client(), Verbose: verbose})

	// print ids if required
	if getId {
		printInstanceIdsByTag(instances)
		return
	}

	// print tabs if required
	tabPrintInstances(instances)
}

func printInstanceIdsByTag(instances []InstanceResult) {
	for _, instance := range instances {
		fmt.Println(instance.Id)
	}
}
