package io

import (
	"path/filepath"
	"strings"
)

type File struct {
	FilePath string
	FileName string
}

func (file *File) FileNameWithoutExtension() string {
	return strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
}
