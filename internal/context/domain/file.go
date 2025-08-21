package domain

import (
	"fmt"
)

type ContextFile struct {
	path		string
	description	string
}

func NewContextFile(path, description string) *ContextFile {
	return &ContextFile{path, description}
}

func FileFromMarkdown(markdown string) (*ContextFile, error) {
	var path string
	var description string
	
	_, err := fmt.Sscanf(markdown, "%s : %s\n", &path, &description)
	if err != nil {
		return nil, err	
	}

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
	return fmt.Sprintf("%s : %s\n", f.path, f.description)
}
