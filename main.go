// main package
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var password string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading environment variable")
	}

	password = os.Getenv("PASSWORD")
}

func main() {
	connectToDatabase(password)
	printMenu()
	result := handleUserInput()
	if result != nil {
		fmt.Println(result.Error())
	}
}
