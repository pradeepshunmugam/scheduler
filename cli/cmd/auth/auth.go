package cmd

import (
	"fmt"
	"regexp"
	"scheduler/cli/db"
	"scheduler/cli/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var username string
var password string
var email string
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
		//isLoggedIn = true
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

func init() {
	rootCmd.AddCommand(loginCmd, createUser)
	loginCmd.Flags().StringVarP(&username, "user", "u", "", "user name")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	createUser.Flags().StringVarP(&username, "user", "u", "", "user name")
	createUser.Flags().StringVarP(&password, "password", "p", "", "password")
	createUser.Flags().StringVarP(&email, "email", "m", "", "email")

}
