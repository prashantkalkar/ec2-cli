package main

import (
	"github.com/prashantkalkar/ec2-cli/ec2cli"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var ec2IP string
	var getId bool
	var verbose bool
	var searchTags []string
	var cmdFindEC2 = &cobra.Command{
		Use:   "ec2-cli [flags]",
		Short: "Find ec2 instances, matching ip address or tags",
		Long:  `Find ec2 instances, matching ip address or tags. If non found empty result is shown.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(ec2IP) > 0 {
				ec2cli.PrintInstancesByIP(getId, ec2IP, verbose)
				return
			}
			if len(searchTags) > 0 {
				ec2cli.PrintInstancesByTags(getId, searchTags, verbose)
			}
		},
	}

	cmdFindEC2.Flags().StringVarP(&ec2IP, "ip", "p", "", "Provide public or private ip")
	cmdFindEC2.Flags().BoolVarP(&getId, "id", "i", false, "Provide only instance Id (default false)")
	cmdFindEC2.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logs (default false)")
	cmdFindEC2.Flags().StringSliceVarP(&searchTags, "tags", "t", []string{}, "Provide tags to be matched for instances")
	cmdFindEC2.MarkFlagsMutuallyExclusive("ip", "tags")

	//var rootCmd = &cobra.Command{Use: "ec2-cli - lookup ec2 instance by ip or tags"}
	//rootCmd.AddCommand(cmdFindEC2)
	err := cmdFindEC2.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
