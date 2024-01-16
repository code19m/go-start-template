package main

import (

	// Importing the "automaxprocs" package from Uber's Gojuno enables automatic
	// tuning of the maximum number of OS threads (GOMAXPROCS) that Go can utilize,
	// improving performance particularly on containerized environments.

	"go-start-template/cmd/app"
	"log"

	_ "go.uber.org/automaxprocs"

	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use: "main",
	}

	rootCmd.AddCommand(app.AppCmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
