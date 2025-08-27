package application

import (
	"fmt"
	"path"
	"strings"
)

const system_prompt = `
	I need you to create context module, which has a simple goal : create context in markdown, meant to be read by developpers and ai agents. Each module takes between 0 and n files, and summarize each one, extracting a one liner describing the behaviour of the file. Each module can also contain submodules, and your role is to summarize each submodules, so that the reader is guided through the context tree, refining its search on each step. I will give you the files with format 

path: "path/to/file"
content: "content of the file"
---

I want you to answer only with the markdown content describing the module, with format:
# Subject
### Description
"Description of the module"

If there are submodules, then
### Modules
for each module,
	- <submodule_subject> : "submodule content"

If there are files, then
### Files
for each file,
	- <file_subject> : "file content"

DO NOT save this file, just return the raw content
DO NOT confirm or say anything else, I just want the raw markdown text\n\n
`

func (c ContextService) toPrompt(m *moduleContainer) (string, error) {
	builder := strings.Builder{}

	filePrompts, err := c.aggregate(m)
	if err != nil {
		return "", nil
	}
	builder.WriteString(system_prompt)
	builder.WriteString(filePrompts)

	return builder.String(), nil
}

func (c ContextService) aggregate(m *moduleContainer) (string, error) {
	builder := strings.Builder{}
	
	builder.WriteString("Files : \n")
	for _, p := range m.filePaths {
		rawFile, err := c.fsPort.ReadFile(p)	
		if err != nil {
			return "", nil
		}
		content, err := rawFile.Content()
		if err != nil {
			return "", nil
		}
		builder.WriteString(fmt.Sprintf("path : %s\n\n", p))
		builder.WriteString(fmt.Sprintf("content : %s\n\n", content))
		builder.WriteString("---\n\n")
	}

	builder.WriteString("Modules : \n")
	for _, p := range m.children {
		rawFile, err := c.fsPort.ReadFile(path.Join(p.modulePath, "CONTEXT.mdc"))
		if err != nil {
			return "", err
		}
		content, err := rawFile.Content()
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("path : %s\n\n", p))
		builder.WriteString(fmt.Sprintf("content : %s\n\n", content))
		builder.WriteString("---\n\n")
	}

	return builder.String(), nil
}
