package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kunzru",
		Short: "Kunzru CLI",
		Long:  "Kunzru CLI is a command-line interface for managing contexts for IAs.",
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
