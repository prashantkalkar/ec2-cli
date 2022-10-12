package ec2cli

func FindEC2InstanceIds(ip string, client AWSClient) []string {
	privateIPResult := client.GetInstancesWithPrivateIP(ip, "Name")
	if len(privateIPResult) > 0 {
		return getInstanceIds(privateIPResult)
	}
	publicIPResult := client.GetInstancesWithPublicIP(ip)
	if len(publicIPResult) > 0 {
		return getInstanceIds(publicIPResult)
	}
	return []string{}
}

func getInstanceIds(result []InstanceResult) []string {
	var instanceIds []string
	for _, instanceResult := range result {
		instanceIds = append(instanceIds, instanceResult.Id)
	}
	return instanceIds
}
