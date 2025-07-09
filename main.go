package main

import (
	"fmt"
	"github.com/nunseik/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}
	config.CurrentUserName = "Fabio"
	err = config.SetUser()
	if err != nil {
		fmt.Println("Error setting user:", err)
	} else {
		fmt.Println("User set successfully:", config.CurrentUserName)
		fmt.Println("Database URL:", config.DbURL)
	}
}