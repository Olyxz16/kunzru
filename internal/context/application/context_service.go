package application

import (
	"os"
	"path"

	"github.com/Olyxz16/kunzru/internal/context/domain"
	fs "github.com/Olyxz16/kunzru/internal/filesystem/application"
)

type ContextService struct{
	contextPort	ContextPort
	fsPort		fs.FileSystemPort
}

func NewContextService() *ContextService {
	return &ContextService{}
}


type pathModulePair struct {
	dir		string
	module  *domain.ContextModule
}

func (s *ContextService) BuildContextTree(inputDir string) (*domain.ContextModule, error) {
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
	return result.GetModules()[0], nil
}

func (s *ContextService) GenerateContextTree(module *domain.ContextModule) (error) {
	return nil	
} 
