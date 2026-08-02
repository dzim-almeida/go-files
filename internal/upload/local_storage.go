package upload

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	return &LocalStorage{
		uploadDir: uploadDir,
	}
}

func (l *LocalStorage) Save(reader io.Reader, filename string) (string, error) {
	err := os.MkdirAll(l.uploadDir, os.ModePerm)

	if err != nil {
		return "", NewSaveError()
	}

	path := filepath.Join(l.uploadDir, filename)

	file, err := os.Create(path)

	if err != nil {
		return "", NewSaveError()
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return "", NewSaveError()
	}

	return path, nil
}

func (l *LocalStorage) Delete(path string) error {
	return nil
}