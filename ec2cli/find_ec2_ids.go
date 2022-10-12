package ec2cli

func FindEC2InstanceIds(ip string, client AWSClient) []string {
	instances := FindEC2Instances(ip, client)
	var instanceIds []string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.Id)
	}
	return instanceIds
}
