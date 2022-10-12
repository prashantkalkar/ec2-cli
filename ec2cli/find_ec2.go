package ec2cli

func FindEC2InstanceIds(ip string, id bool, client AWSClient) ([]string, error) {
	privateIPResult := client.GetInstancesWithPrivateIP(ip, "Name")
	if len(privateIPResult) > 0 {
		return []string{privateIPResult[0].Id}, nil
	}
	publicIPResult := client.GetInstancesWithPublicIP(ip)
	if len(publicIPResult) > 0 {
		return []string{publicIPResult[0].Id}, nil
	}
	return []string{}, nil
}
