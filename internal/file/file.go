package file

import (
	"os"
)

type FileManager struct {
	File *os.File
}

func NewFileManager(filePath string) (*FileManager, error) {
	//file, err := os.Open(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &FileManager{File: file}, nil
}
