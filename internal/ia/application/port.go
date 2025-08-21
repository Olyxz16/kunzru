package application

type IAPort interface {
	Prompt(prompt string) (string, error)
}
