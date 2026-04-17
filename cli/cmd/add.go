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
			if cron == "" || sample == "" || email == "" {
				cron = "*/5 * * * *"
				sample = "1"
				email = ""
			}
			//fmt.Printf("got the url : %v and name %v to monitor.", url, name)
			//dbconn.DBInsert(name, url) //Once get the url and respective details from user. call the function to uopdate in  DB - scheduler
			dbconn.DBInsert(name, url, cron, sample, email)
			fmt.Println("Data inserted", name, url, cron, sample, email)
		}
		if csvFile != "" {
			fmt.Println("Reading csv file")
			file, err := os.Open(csvFile)
			if err != nil {
				fmt.Println("Unable to read file.", err)
				return
			}
			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			for i := range records {
				name, url, cron, sample, email := records[i][0], records[i][1], records[i][2], records[i][3], records[i][4]
				for _, row := range records {
					if len(row) < 2 {
						fmt.Println("Skipping invalid value..")
						continue
					}
					name := row[0]
					url := row[1]
					cron = "*/5 * * * *"
					sample = "1"
					email = ""
					if len(row) > 2 && row[2] != "" {
						cron = row[2]
					}
					if len(row) > 3 && row[2] != "" {
						sample = row[3]
					}

					dbconn.DBInsert(name, url, cron, sample, email)
				}
				fmt.Println("Data inserted", name, url, cron, sample, email)
			}

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
