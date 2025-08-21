package domain

import (
	"errors"
	"path"
)

type RawFile struct {
	path 			string
	isDir			bool
	content			string
}

func NewFile(path, content string) (*RawFile) {
	return &RawFile{
		path: path,
		content: content,
		isDir: false,
	}
}

func NewDir(path string) (*RawFile) {
	return &RawFile{
		path: path,
		isDir: true,
		content: "",
	}
}

func (f *RawFile) Path() string {
	return f.path
}

func (f *RawFile) Name() string {
	return path.Base(f.path)
}

func (f *RawFile) IsDir() bool {
	return f.isDir
}

func (f *RawFile) Content() (string, error) {
	if f.isDir {
		return "", errors.New("Cannot read content of a directory")
	}
	return f.content, nil
}
