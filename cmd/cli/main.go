package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	fs "github.com/Olyxz16/kunzru/internal/filesystem/infrastructure"
	ia "github.com/Olyxz16/kunzru/internal/ia/infrastructure"
	context "github.com/Olyxz16/kunzru/internal/context/application"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kunzru",
		Short: "Kunzru CLI",
		Long:  "Kunzru CLI is a command-line interface for managing contexts for IAs.",
	}
	rootCmd.AddCommand(generateCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var generateCommand = &cobra.Command{
	Use: "generate",
	Short: "Generate the context tree",
	Run: func(cmd *cobra.Command, args []string) {
		fsRepository := fs.NewFileRepository()
		aiRepository := ia.NewGeminiService()
		contextService := context.NewContextService(fsRepository, aiRepository)
		_, err := contextService.GenerateContextTree(".")
		if err != nil {
			panic(err)
		}
	},
}
