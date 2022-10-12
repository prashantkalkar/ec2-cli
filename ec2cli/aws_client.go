package ec2cli

type AWSClient interface {
	GetInstancesWithPrivateIP(ip string, nameKey string) []InstanceResult
	GetInstancesWithPublicIP(ip string, nameKey string) []InstanceResult
	GetInstancesWithTags(tagValues []string) []InstanceResult
}
