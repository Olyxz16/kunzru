package application

import (
	"strings"

	"github.com/Olyxz16/kunzru/internal/filesystem/domain"
	gitignore "github.com/sabhiram/go-gitignore"
)

type IgnoreService struct {
	rules []string
}

func NewIgnoreService() *IgnoreService {
	return &IgnoreService{
		[]string{".git/"},
	}
}

func (s *IgnoreService) AddFile(file *domain.RawFile) (*IgnoreService, error) {
	content, err := file.Content()
	if err != nil {
		return nil, err
	}
	rules := strings.Split(content, "\n")
	return &IgnoreService{
		append(s.rules, rules...),
	}, nil
}

func (s *IgnoreService) IsIgnored(filepath string) bool {
	gi := gitignore.CompileIgnoreLines(s.rules...)
	return gi.MatchesPath(filepath)
}
