package ec2cli

import "errors"

func FindEC2InstanceId(ip string, id bool, client AWSClient) (string, error) {
	privateIPResult := client.GetInstancesWithPrivateIP(ip, "Name")
	if len(privateIPResult) > 0 {
		return privateIPResult[0].Id, nil
	}
	publicIPResult := client.GetInstancesWithPublicIP(ip)
	if len(publicIPResult) > 0 {
		return publicIPResult[0].Id, nil
	}
	return "", errors.New("IP not found")
}
