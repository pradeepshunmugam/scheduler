package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"scheduler/cli/db"
	"scheduler/cli/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var username string
var password string
var url string
var name string
var cron string
var sample string
var email string
var csvFile string
var isLoggedIn bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login with username and password",
	Long:  "to run/execute any command login first",
	Run: func(cmd *cobra.Command, args []string) {
		ValidUser()
	},
}

func ValidUser() {
	db := db.GetDB()
	userQuery := `SELECT username, password from users where username = $1`
	row := db.QueryRow(userQuery, username)
	var user, pwd string
	err := row.Scan(&user, &pwd)
	if err != nil {
		logger.Log.Error("unable to query the user detail", zap.Error(err))
	}
	if user == "" && pwd == "" {
		fmt.Println("user is not available in system")
		return
	}
	if username == user && password == pwd {
		isLoggedIn = true
		fmt.Println("Logged in succesfully")
	} else {
		fmt.Println("Enter valid credential")
	}

}

var createUser = &cobra.Command{
	Use:   "user",
	Short: "create user",
	Run: func(cmd *cobra.Command, args []string) {
		matched, err := regexp.MatchString(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`, email)
		if err != nil {
			logger.Log.Error("Email pattern is not valid", zap.Error(err))
		}
		if username != "" && password != "" && matched {
			hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println(err)
			}
			hashedPwdStr := string(hashedPwd)
			userCreation := db.InsertUser(username, hashedPwdStr, email)
			fmt.Println(userCreation)
		} else {
			fmt.Println("Given details are not valid.")
		}

	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add the url",
	Long:  "To monitor or initialize the URL, add it using 'sch add -flag <> '. Add url in format - 'https://website.com'. ",
	Run: func(cmd *cobra.Command, args []string) {
		//Get single line input and validate mandatory field and assign default value
		matched, err := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, url)
		if err != nil {
			fmt.Println("url is not in valid format", err)
		}
		if url != "" && name != "" && matched {
			if cron == "" || sample == "" || email == "" {
				cron = "*/5 * * * *"
				sample = "1"
				email = ""
			}
			//result := query.DBInsert(name, url, cron, sample, email)
			db.DBInsert(name, url, cron, sample, email)
			//fmt.Println(result)
		} else {
			fmt.Println("Either url or name field is not provided. Provided : ", url, name)
		}
		//Get fields by ready csv to bulk import it and validate mandatory field and assign default value
		if csvFile != "" {
			logger.Log.Info("Reading csv file")
			file, err := os.Open(csvFile)
			if err != nil {
				logger.Log.Error("Unable to read", zap.Error(err))
				return
			}
			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				logger.Log.Error("Unable to read", zap.Error(err))
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
				matched, err := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, url)
				if err != nil {
					fmt.Println("url is not in valid format")
				}
				if name != "" && url != "" && matched {
					//result := db.DBInsert(name, url, cron, sample, email) //calling dbpackage to insert the files
					db.DBInsert(name, url, cron, sample, email) //calling dbpackage to insert the files
					// fmt.Println(result)
				} else {
					logger.Log.Warn("Missing Fields.")
				}
			}

		}
	},
}

var deletecmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the url from monitoring",
	Long:  "Delete the url from monitoring/scheduling which is no longer needed.",
	Run: func(cmd *cobra.Command, args []string) {
		//db := db.GetDB()
		if name != "" {
			fmt.Printf("Got the name - %v.\n", name)

			db.DeleteRow(name)
			//fmt.Println(result)
		} else {
			fmt.Println("Provide the url job name.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd, deletecmd, loginCmd, createUser)
	addCmd.Flags().StringVarP(&url, "url", "l", "", "url to monitor")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "name of the url/job.")
	addCmd.Flags().StringVarP(&cron, "cron", "c", "", "cron definition to check the url")
	addCmd.Flags().StringVarP(&sample, "sample", "s", "", "number of samples to check")
	addCmd.Flags().StringVarP(&email, "email", "m", "", "email id to send notification")
	addCmd.Flags().StringVarP(&csvFile, "file", "g", "", "csv file to add bulk data.")
	deletecmd.Flags().StringVarP(&name, "name", "n", "", "name of the url to delete.")
	loginCmd.Flags().StringVarP(&username, "user", "u", "", "user name")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	createUser.Flags().StringVarP(&username, "user", "u", "", "user name")
	createUser.Flags().StringVarP(&password, "password", "p", "", "password")
	createUser.Flags().StringVarP(&email, "email", "m", "", "email")

}
