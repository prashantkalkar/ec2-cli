package ec2cli

type AWSClient interface {
	GetInstancesWithPrivateIP(ip string, nameKey string) []InstanceResult
	GetInstancesWithPublicIP(ip string) []InstanceResult
	GetInstancesWithTags(tagValues []string) []InstanceResult
}
