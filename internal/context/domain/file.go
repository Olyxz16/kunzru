package domain

import (
	"errors"
	"fmt"
	"regexp"
)

type ContextFile struct {
	path		string
	description	string
}

func NewContextFile(path, description string) *ContextFile {
	return &ContextFile{path, description}
}

func FileFromMarkdown(markdown string) (*ContextFile, error) {
	re := regexp.MustCompile(`^-\s?["'` + "\\`" + `]?(\S*)["'` + "`" + `]?\s?:\s*(.*)$`)
	var path string
	var description string
	
	m := re.FindStringSubmatch(markdown)
	if len(m) == 0 {
		return nil, errors.New("no match found")
	}

	path = m[1]
	description = m[2]

	return &ContextFile{
		path: path,
		description: description,
	}, nil
}

func (f *ContextFile) GetPath() string {
	return f.path
}

func (f *ContextFile) GetDescription() string {
	return f.description
}

func (f *ContextFile) ToMarkdown() string {
	return fmt.Sprintf("- %s : %s\n\n", f.path, f.description)
}
