package application

import "github.com/Olyxz16/kunzru/internal/filesystem/domain"


type FileSystemPort interface {
	ReadFile(path string) (*domain.RawFile, error)
	ReadDir(path string) ([]*domain.RawFile, error)
	SaveFile(file *domain.RawFile) (error)
}
