package infrastructure

import (
	"os"
	"path"

	"github.com/Olyxz16/kunzru/internal/filesystem/domain"
)

type FileRepository struct {}

func ReadFile(filePath string) (*domain.RawFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return domain.NewFile(filePath, string(content)), nil
}

func ReadDir(dirPath string) ([]*domain.RawFile, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	files := []*domain.RawFile{}
	for _, e := range entries {
		var file *domain.RawFile
		p := path.Join(dirPath, e.Name())
		if e.IsDir() {
			file = domain.NewDir(p)
		} else {
			content, err := os.ReadFile(p)
			if err != nil {
				return nil, err
			}
			file = domain.NewFile(p, string(content))
		}
		files = append(files, file)
	}
	return files, nil
}

func SaveFile(file *domain.RawFile) error {
	content, err := file.Content()
	if err != nil {
		return err
	}
	err = os.WriteFile(file.Path(), []byte(content), os.FileMode(os.O_WRONLY))
	return err
}
