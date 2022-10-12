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
	"os"
	"text/tabwriter"
)

func main() {
	var ec2IP string
	var getId bool
	var cmdFindEC2One = &cobra.Command{
		Use:   "find [flags]",
		Short: "Find ec2 instances, matching ip address or tags",
		Long:  `Find ec2 instances, matching ip address or tags. If non found empty result is shown.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if getId {
				ids := ec2cli.FindEC2InstanceIds(ec2IP, ec2cli.AWSClientImpl{Client: getAWSEC2Client()})
				for _, id := range ids {
					fmt.Println(id)
					return
				}
			}

			instances := ec2cli.FindEC2Instances(ec2IP, ec2cli.AWSClientImpl{Client: getAWSEC2Client()})
			writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			_, _ = fmt.Fprintln(writer, "INSTANCE_ID\tNAME\tIP_ADDRESS\t")
			for _, instance := range instances {
				instanceStr := fmt.Sprintf("%s\t%s\t%s\t\n", instance.Id, instance.Name, instance.IP)
				_, _ = fmt.Fprintln(writer, instanceStr)
			}
			_ = writer.Flush()
		},
	}

	cmdFindEC2One.Flags().StringVarP(&ec2IP, "ip", "p", "", "Provide public or private ip")
	cmdFindEC2One.Flags().BoolVarP(&getId, "id", "i", false, "Provide only instance Id (default false)")

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
