package application

import (
	"github.com/Olyxz16/kunzru/internal/context/domain"
)

type pathModulePair struct {
	dir		string
	module  *domain.ContextModule
}

func (s ContextService) BuildContextTree(inputDir string) (*domain.ContextModule, error) {
	result := &domain.ContextModule{}
	queue := []pathModulePair{{inputDir, result}}
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		dir := elem.dir
		module := elem.module

		entries, err := s.fsPort.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		contextFile, _, dirs, _ := s.parseDirectory(entries)

		if contextFile != nil {
			content, err := contextFile.Content()
			if err != nil {
				return nil, err
			}
			newModule, err := domain.ModuleFromMarkdown(dir, string(content))
			if err != nil {
				return nil, err
			}
			module.AddModule(newModule)
			module = newModule
		}

		for _, dir := range dirs {
			newDir := dir.Path()
			queue = append(queue, pathModulePair{dir: newDir, module: module})
		}
	}

	if len(result.GetModules()) == 0 {
		return nil, nil
	}
	return result.GetModules()[0], nil
}
