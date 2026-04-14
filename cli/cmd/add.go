package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var url string
var name string
var cron string
var sample string
var email string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add the url",
	Long:  "To monitor or initialize the URL, add it using 'sch add -flag <> '. Add url in format - 'https://website.com'. ",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" && name == "" {
			fmt.Println("Please provide url and name. Provided is empty")
			return
		}
		fmt.Printf("got the url : %v and name %v to monitor.", url, name)

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&url, "url", "u", "", "url to monitor")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "name of the url/job.")
	addCmd.Flags().StringVarP(&cron, "cron", "c", "", "cron definition to check the url")
	addCmd.Flags().StringVarP(&sample, "sample", "s", "", "number of samples to check")
	addCmd.Flags().StringVarP(&email, "email", "m", "", "email id to send notification")
}
