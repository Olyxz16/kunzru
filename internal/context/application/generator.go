package application

import (
	"errors"
	"path"

	"github.com/Olyxz16/kunzru/internal/context/domain"
	fsd "github.com/Olyxz16/kunzru/internal/filesystem/domain"
)

type contextEntry struct {
	dir				string
	ignoreService 	*IgnoreService
	container   	*moduleContainer
}

type moduleContainer struct {
	modulePath string
	filePaths   []string
	children	[]*moduleContainer
}

func emptyContainer(path string) *moduleContainer {
	return &moduleContainer{path, []string{}, []*moduleContainer{}}	
}

func (s ContextService) GenerateContextTree(inputDir string) (*domain.ContextModule, error) {
	tree, err := s.buildAbstractContextTree(inputDir)
	if err != nil {
		return nil, err
	}

	return s.buildContextTreeFromAbstractTree(tree)
}

func (s ContextService) buildAbstractContextTree(inputDir string) (*moduleContainer, error) {
	base := emptyContainer(inputDir)
	ignore := NewIgnoreService()
	queue := []contextEntry{{inputDir, ignore, base}}
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		dir := elem.dir
		container := elem.container

		entries, err := s.fsPort.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		contextFile, ignoreFile, dirs, files := s.parseDirectoryIgnore(entries, ignore)

		if contextFile != nil {
			newContainer := emptyContainer(dir)
			container.children = append(container.children, newContainer)
			container = newContainer
		}

		if ignoreFile != nil {
			ignore, err = ignore.AddFile(ignoreFile)
			if err != nil {
				return nil, err
			}
		}	

		for _, dir := range dirs {
			queue = append(queue, contextEntry{dir.Path(), ignore, container})
		}

		for _, file := range files {
			container.filePaths = append(container.filePaths, file.Path())
		}
	}

	l := len(base.children)
	if l == 0 {
		return nil, errors.New("No CONTEXT.mdc found.")
	}
	if l > 1 {
		return nil, errors.New("No root CONTEXT.mdc found.")	
	}
	
	return base.children[0], nil
}

func (s ContextService) buildContextTreeFromAbstractTree(container *moduleContainer) (*domain.ContextModule, error) {
	submodules := make([]*domain.ContextModule, 0, len(container.children))	
	for _, m := range container.children {
		submodule, err := s.buildContextTreeFromAbstractTree(m)
		if err != nil {
			return nil, err
		}
		submodules = append(submodules, submodule)
	}

	content, err := s.toPrompt(container)
	if err != nil {
		return nil, err
	}

	result, err := s.iaPort.Prompt(content)
	if err != nil {
		return nil, err
	}

	module, err := domain.ModuleFromMarkdown(container.modulePath, result)
	if err != nil {
		return nil, err
	}

	rawFile := fsd.NewFile(path.Join(module.GetPath(), "CONTEXT.mdc"), module.ToMarkdown())
	err = s.fsPort.SaveFile(rawFile)
	if err != nil {
		return nil, err
	}

	return module, nil
}
