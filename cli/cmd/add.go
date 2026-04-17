package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"scheduler/cli/cmd/dbconn"

	"github.com/spf13/cobra"
)

var url string
var name string
var cron string
var sample string
var email string
var csvFile string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add the url",
	Long:  "To monitor or initialize the URL, add it using 'sch add -flag <> '. Add url in format - 'https://website.com'. ",
	Run: func(cmd *cobra.Command, args []string) {
		if url != "" && name != "" {
			fmt.Printf("got the url : %v and name %v to monitor.", url, name)
			dbconn.DBInsert(name, url) //Once get the url and respective details from user. call the function to uopdate in  DB - scheduler
		}
		if csvFile != "" {
			file, err := os.Open(csvFile)
			if err != nil {
				fmt.Println("Unable to read file.")
				return
			}
			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			for i := range records {
				name, url := records[i][0], records[i][1]
				dbconn.DBInsert(name, url)
			}
			fmt.Println("Data inserted")

		}
		fmt.Println("Please provide url and name. Provided is empty")

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&url, "url", "u", "", "url to monitor")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "name of the url/job.")
	addCmd.Flags().StringVarP(&cron, "cron", "c", "", "cron definition to check the url")
	addCmd.Flags().StringVarP(&sample, "sample", "s", "", "number of samples to check")
	addCmd.Flags().StringVarP(&email, "email", "m", "", "email id to send notification")
	addCmd.Flags().StringVarP(&csvFile, "file", "g", "", "csv file to add bulk data.")
}
