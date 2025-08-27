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
	cmd := exec.Command("gemini", "-m", "gemini-2.5-flash", "-p")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(prompt))
	}()

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	output := string(outputBytes)
	slices := strings.Split(output, "\n")
	slices = slices[1:]
	output = strings.Join(slices, "\n")

	return output, nil
}
