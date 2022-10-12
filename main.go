package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/prashantkalkar/ec2-cli/ec2cli"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var ec2IP string
	var getId bool
	var cmdFindEC2One = &cobra.Command{
		Use:   "find-ec2-one [flags]",
		Short: "Find single ec2 instance, if multiple found first will be provided",
		Long:  `Find single ec2 instance, if multiple found first will be provided. If non found empty result is shown.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := ec2cli.FindEC2InstanceId(ec2IP, getId, ec2cli.AWSClientImpl{Client: getAWSEC2Client()})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id)
		},
	}

	cmdFindEC2One.Flags().StringVarP(&ec2IP, "ip", "p", "", "Provide public or private ip")
	cmdFindEC2One.Flags().BoolVarP(&getId, "get-id", "i", false, "Provide only instance Id (default false)")

	var rootCmd = &cobra.Command{Use: "ec2-cli - lookup ec2 instance by ip or tags"}
	rootCmd.AddCommand(cmdFindEC2One)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getAWSEC2Client() *ec2.Client {
	client := ec2.NewFromConfig(defaultAWSConfig())
	return client
}

func defaultAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
