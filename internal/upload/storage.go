package upload

import "io"

type UploadStorage interface {
	Save(reader io.Reader, filename string) (string, error)
	Delete(filepath string) error
}