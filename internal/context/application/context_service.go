package application

import (
	"errors"
	"os"
	"path"

	"github.com/Olyxz16/kunzru/internal/context/domain"
	fs "github.com/Olyxz16/kunzru/internal/filesystem/application"
	fsd "github.com/Olyxz16/kunzru/internal/filesystem/domain"
	ia "github.com/Olyxz16/kunzru/internal/ia/application"
)

type ContextService struct{
	fsPort		fs.FileSystemPort
	iaPort		ia.IAPort
}

func NewContextService(fsPort fs.FileSystemPort, iaPort ia.IAPort) *ContextService {
	return &ContextService{fsPort, iaPort}
}


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

		for _, entry := range entries {
			if entry.Name() == CONTEXT_FILE_NAME {
				content, err := os.ReadFile(path.Join(dir, entry.Name()))
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
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				newDir := path.Join(dir, entry.Name())
				queue = append(queue, pathModulePair{dir: newDir, module: module})
			}
		}
	}

	if len(result.GetModules()) == 0 {
		return nil, nil
	}
	return result.GetModules()[0], nil
}


type contextDirPair struct {
	dir			string
	container   *moduleContainer
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
	queue := []contextDirPair{{inputDir, base}}
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		dir := elem.dir
		container := elem.container

		entries, err := s.fsPort.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.Name() == CONTEXT_FILE_NAME {
				newContainer := emptyContainer(dir)
				container.children = append(container.children, newContainer)
				container = newContainer
			}
		}

		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != CONTEXT_FILE_NAME {
				container.filePaths = append(container.filePaths, path.Join(dir, entry.Name()))
			}
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				newDir := path.Join(dir, entry.Name())
				queue = append(queue, contextDirPair{newDir, container})
			}
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
