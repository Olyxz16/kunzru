package domain

import (
	"errors"
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
	lines := strings.Split(markdown, "\n")
	if len(lines) < 2 {
		return EmptyModule(path), nil
	}
	
	var subject string
	_, err := fmt.Sscanf(lines[0], "# %s\n", &subject)
	if err != nil {
		return nil, err
	}
	var description string
	_, err = fmt.Sscanf(lines[1], "### %s\n", &description)
	if err != nil {
		return nil, err
	}
	
	index := 2
	var equals bool

	/* Modules */
	equals = lines[index] == "### Modules\n"
	if !equals {
		return nil, errors.New("Wrong module delimiter")
	}
	index++

	for !strings.HasPrefix(lines[index], "#") {
		var path string
		var description string
		_, err := fmt.Sscanf(lines[index], "%s : %s\n", &path, &description)
		if err != nil {
			return nil, err
		}
		module := EmptyModule(path)
		module.description = description
		index++
	}

	/* Files */
	equals = lines[index] == "### Files\n"
	if !equals {
		return nil, errors.New("Wrong file delimiter")
	}
	index++

	for !strings.HasPrefix(lines[index], "#") {
		contextFile, err := FileFromMarkdown(lines[index])
		if err != nil {
			return nil, err
		}
		module.AddFile(contextFile)
		index++
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
	builder.WriteString(fmt.Sprintf("# %s\n", m.subject))
	builder.WriteString(fmt.Sprintf("%s\n", m.description))
	
	builder.WriteString("### Modules\n")
	for _, f := range m.modules {
		builder.WriteString(fmt.Sprintf("%s : %s\n", f.path, f.description))
	}

	builder.WriteString("### Files\n")
	for _, f := range m.files {
		builder.WriteString(f.ToMarkdown())
	}
	return builder.String()
}
