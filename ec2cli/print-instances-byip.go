package ec2cli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
	"os"
	"text/tabwriter"
)

func PrintInstancesByIP(getId bool, ec2IP string) {
	if getId {
		printInstanceIds(func() []string { return FindEC2InstanceIds(ec2IP, AWSClientImpl{Client: getAWSEC2Client()}) })
		return
	}
	tabPrintInstances(FindEC2Instances(ec2IP, AWSClientImpl{Client: getAWSEC2Client()}))
}

func tabPrintInstances(instances []InstanceResult) {
	writer := getTabWriter()
	_, _ = fmt.Fprintln(writer, "INSTANCE_ID\tNAME\tIP_ADDRESS\t")
	for _, instance := range instances {
		_, _ = fmt.Fprintln(writer, fmt.Sprintf("%s\t%s\t%s\t", instance.Id, instance.Name, instance.IP))
	}
	_ = writer.Flush()
}

func getTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
}

func printInstanceIds(fetchIds func() []string) {
	for _, id := range fetchIds() {
		fmt.Println(id)
	}
}

func getAWSEC2Client() *ec2.Client {
	return ec2.NewFromConfig(defaultAWSConfig())
}

func defaultAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
