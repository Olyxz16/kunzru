package domain

import (
	"fmt"
	"strings"
)

type ContextModule struct {
	path		string
	subject 	string
	description string
	files 		[]*ContextFile
	modules		[]*ContextModule
}

func NewContextModule(path, subject, description string, files []*ContextFile, modules []*ContextModule) *ContextModule {
	return &ContextModule{path, subject, description, files, modules}
}

func EmptyModule(path string) *ContextModule {
	return &ContextModule{path, "", "", []*ContextFile{}, []*ContextModule{}}
}

func ModuleFromMarkdown(path, markdown string) (*ContextModule, error) {
	module := &ContextModule{path: path}
	lines := splitAndTrim(markdown)
	if len(lines) < 2 {
		return EmptyModule(path), nil
	}

	var subject string
	_, err := fmt.Sscanf(lines[0], "# %s", &subject)
	if err != nil {
		return nil, fmt.Errorf("Error when parsing module subject : %s\n", err.Error())
	}
	var description = lines[2]

	module.subject = subject
	module.description = description
	
	index := 3
	var equals bool

	/* Modules */
	equals = index < len(lines) && lines[index] == "### Modules"
	if equals {
		index++

		for index < len(lines) && !strings.HasPrefix(lines[index], "#") && len(lines[index]) > 0 {
			submoduleFile, err := FileFromMarkdown(lines[index])
			if err != nil {
				fmt.Printf("%s\n", lines[index])
				return nil, fmt.Errorf("Error when parsing submodule : %s\n", err.Error())
			}
			submodule := EmptyModule(submoduleFile.path)
			submodule.description = submoduleFile.description
			module.AddModule(submodule)
			index++
		}
	}

	/* Files */
	equals = index < len(lines) && lines[index] == "### Files"
	if equals {
		index++

		for index < len(lines) && !strings.HasPrefix(lines[index], "#") && len(lines[index]) > 0 {
			contextFile, err := FileFromMarkdown(lines[index])
			if err != nil {
				return nil, fmt.Errorf("Error when parsing file : %s\n", err.Error())
			}
			module.AddFile(contextFile)
			index++
		}
	}

	return module, nil
}

func (m *ContextModule) GetPath() string {
	return m.path
}

func (m *ContextModule) GetSubject() string {
	return m.subject
}

func (m *ContextModule) GetDescription() string {
	return m.description
}

func (m *ContextModule) GetFiles() []*ContextFile {
	return m.files
}

func (m *ContextModule) GetModules() []*ContextModule {
	return m.modules
}

func (m *ContextModule) IsEmpty() bool {
	return m.subject == ""
}

func (m *ContextModule) AddFile(file *ContextFile) {
	m.files = append(m.files, file)	
}

func (m *ContextModule) AddModule(module *ContextModule) {
	m.modules = append(m.modules, module)	
}

func (m *ContextModule) ToMarkdown() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# %s\n\n", m.subject))
	builder.WriteString("### Description\n\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", m.description))
	
	if len(m.modules) > 0 {
		builder.WriteString("### Modules\n\n")
	}
	for _, f := range m.modules {
		builder.WriteString(fmt.Sprintf("%s : %s\n\n", f.path, f.description))
	}
	
	if len(m.files) > 0 {
		builder.WriteString("### Files\n\n")
	}
	for _, f := range m.files {
		builder.WriteString(f.ToMarkdown())
	}
	return builder.String()
}

func splitAndTrim(markdown string) []string {
	lines := strings.Split(markdown, "\n")
	result := make([]string, 0, len(lines))
	
	for _, l := range lines {
		trimmed := strings.Trim(l, " \t")	
		if len(trimmed) > 0 {
			result = append(result, l)
		}
	}

	return result
}
