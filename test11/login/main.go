package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pwd := "hello"
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(hashedPwd))

	result := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	fmt.Println(result)
}
