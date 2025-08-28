package main

import (
	"fmt"

	"github.com/spf13/cobra"

	context "github.com/Olyxz16/kunzru/internal/context/application"
	fs "github.com/Olyxz16/kunzru/internal/filesystem/infrastructure"
	ia "github.com/Olyxz16/kunzru/internal/ia/infrastructure"
)

var generateCommand = &cobra.Command{
	Use: "generate",
	Short: "Generate the context tree",
	Run: generateCommandFunc,
}

func generateCommandFunc(cmd *cobra.Command, args []string) {
	fsRepository := fs.NewFileRepository()
	aiRepository := ia.NewGeminiService()
	contextService := context.NewContextService(fsRepository, aiRepository)
	_, err := contextService.GenerateContextTree(".")
	if err != nil {
		fmt.Printf("Error : %s\n", err.Error())
	}
}
