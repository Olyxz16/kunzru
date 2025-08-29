package application

import (
	"fmt"
	"path"
	"strings"
)

const system_prompt = `
	I need you to create context module, which has a simple goal : create context in markdown, meant to be read by developpers and ai agents. Each module takes between 0 and n files, and summarize each one, extracting a one liner describing the behaviour of the file. Each module can also contain submodules, and your role is to summarize each submodules, so that the reader is guided through the context tree, refining its search on each step. I will give you the files with format :

path: "path/to/file"
content: "content of the file"

I will also give you the submodules description with format : 

path: "path/to/module/directory"
content: "the description of the submodule"

The files and modules will be divided into two sections :
The "## Files" title indicates the start of the file section. You will treat each entry as a file, and parse it with the file format defined previously
The "## Modules" title indidcates the start of the submodules section. You will treat each entry as a submodule, and parse it with the module format defined previously.

DO NOT assume something is a module or a file. Only parse what is given to you, one at a time. 

---

I want you to answer only with the markdown content describing the module, with format:
# Subject
### Description
"Description of the module"

If there are submodules, then
### Modules
for each module,
- <submodule_subject> : "submodule summary"

If there are files, then
### Files
for each file,
- <file_path> : "file summary"

FOR FILES ONLY : The <file_path> should be the full path given to you in the input. DO NOT diverge from it. 
FOR MODULES ONLY : DO NOT put the module path as the module name, always put a 1-2 words clear name stating the subject, with PascalCase naming convention.
DO NOT save this file, just return the raw content.
DO NOT confirm or say anything else, I just want the raw markdown text.

--- 

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
	
	if len(m.filePaths) > 0 {
		builder.WriteString("## Files : \n")
	}
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

	if len(m.children) > 0 {
		builder.WriteString("## Modules : \n")
	}
	for _, p := range m.children {
		rawFile, err := c.fsPort.ReadFile(path.Join(p.modulePath, "CONTEXT.mdc"))
		if err != nil {
			return "", err
		}
		content, err := rawFile.Content()
		if err != nil {
			return "", err
		}
		builder.WriteString(fmt.Sprintf("path : %s\n\n", p.modulePath))
		builder.WriteString(fmt.Sprintf("content : %s\n\n", content))
		builder.WriteString("---\n\n")
	}

	return builder.String(), nil
}
