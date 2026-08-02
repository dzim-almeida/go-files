package upload

import "io"

type UploadInput struct {
	Reader      io.Reader
	FileName    string
	Size        int64
	ContentType string
}

type UploadOutput struct {
	Id               int    `json:"id"`
	OriginalFileName string `json:"original_filename"`
	StoredFileName   string `json:"stored_filename"`
	StoragePath      string `json:"storage_path"`
	ContetType       string `json:"content_type"`
	Size             int64  `json:"size"`
}
