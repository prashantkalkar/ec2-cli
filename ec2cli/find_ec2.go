package ec2cli

func FindEC2Instances(ip string, client AWSClient) []InstanceResult {
	privateIPResult := client.GetInstancesWithPrivateIP(ip, "Name")
	if len(privateIPResult) > 0 {
		return privateIPResult
	}

	publicIPResult := client.GetInstancesWithPublicIP(ip, "Name")
	return publicIPResult
}
