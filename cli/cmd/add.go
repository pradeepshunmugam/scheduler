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
		//Get single line input and validate mandatory field and assign default value
		if url != "" && name != "" {
			if cron == "" || sample == "" || email == "" {
				cron = "*/5 * * * *"
				sample = "1"
				email = ""
			}
			dbconn.DBInsert(name, url, cron, sample, email)
		}
		//Get fields by ready csv to bulk import it and validate mandatory field and assign default value
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
			for _, row := range records {
				var name, url, cron, sample, email string
				name = row[0]
				url = row[1]
				cron = "*/5 * * * *"
				sample = "1"
				email = ""
				if len(row) > 2 && row[2] != "" {
					cron = row[2]
				}
				if len(row) > 3 && row[3] != "" {
					sample = row[3]
				}
				if len(row) > 4 && row[4] != "" {
					email = row[4]
				}
				if name != "" && url != "" {
					dbconn.DBInsert(name, url, cron, sample, email) //calling dbpackage to insert the files
				} else {
					fmt.Println("Missing Fields.")
				}
			}

		}
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
