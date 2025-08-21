package application

import "github.com/Olyxz16/kunzru/internal/context/domain"

const CONTEXT_FILE_NAME = "CONTEXT.mdc"

type ContextPort interface {
	BuildContextTree(baseDir string) (*domain.ContextModule, error)
	GenerateContextTree(module *domain.ContextModule) (error)
}
