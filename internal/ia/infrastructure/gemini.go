package infrastructure

import (
	"os/exec"
	"strings"
)

type GeminiService struct {}

func NewGeminiService() GeminiService {
	return GeminiService{}
}

func (g GeminiService) Prompt(prompt string) (string, error) {
	cmd := exec.Command("gemini", "-p")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		ioString := prompt
		stdin.Write([]byte(ioString))
	}()

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	output := string(outputBytes)
	slices := strings.Split(output, "\n")
	output = slices[2]

	return output, nil
}
