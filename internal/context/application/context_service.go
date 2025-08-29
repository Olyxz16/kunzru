package application

import (
	fs "github.com/Olyxz16/kunzru/internal/filesystem/application"
	fsd "github.com/Olyxz16/kunzru/internal/filesystem/domain"
	ia "github.com/Olyxz16/kunzru/internal/ia/application"
)

type ContextService struct {
	fsPort			fs.FileSystemPort
	iaPort			ia.IAPort
}

func NewContextService(fsPort fs.FileSystemPort, iaPort ia.IAPort) *ContextService {
	return &ContextService{fsPort, iaPort}
}

func (s *ContextService) parseDirectory(dir []*fsd.RawFile) (contextFile *fsd.RawFile, dirs []*fsd.RawFile, files []*fsd.RawFile, ignore *IgnoreService) {
	return s.parseDirectoryIgnore(dir, NewIgnoreService())
}

func (s *ContextService) parseDirectoryIgnore(dir []*fsd.RawFile, ignoreService *IgnoreService) (contextFile *fsd.RawFile, dirs []*fsd.RawFile, files []*fsd.RawFile, ignore *IgnoreService) {
	ignore = ignoreService
	dirs = make([]*fsd.RawFile, 0, len(dir))
	files = make([]*fsd.RawFile, 0, len(dir))
	for _, file := range dir {
		if !file.IsDir() && file.Name() == ".gitignore" {
			var err error
			ignore, err = ignoreService.AddFile(file)
			if err != nil {
				return	
			}
			continue
		}
	}
	for _, file := range dir {
		name := file.Name()
		isDir := file.IsDir()
		if ignore.IsIgnored(file.Path()) {
			continue
		}
		if !isDir && name == CONTEXT_FILE_NAME {
			contextFile = file
			continue
		}
		if !isDir {
			files = append(files, file)
		}
		if isDir {
			dirs = append(dirs, file)
		}
	}
	return
}
