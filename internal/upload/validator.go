package upload

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dutchcoders/go-clamd"
)

type UploadValidator interface {
	Validate(input UploadInput) error
}

type ValidatorConfig struct {
	MaxFileSize      int64
	AllowedFileTypes map[string]string
}

type validator struct {
	config ValidatorConfig
}

func NewValidator(config ValidatorConfig) UploadValidator {
	return &validator{
		config: config,
	}
}

func _restartReader(reader io.Reader) error {
	if seeker, ok := reader.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return NewProcessingError()
		}
	}

	return nil
}

func validContentType(reader io.Reader, expectedContentType string) error {
	const (
		bufferSize = 512
	)
	
	head := make([]byte, bufferSize)
	magicBytes, err := reader.Read(head)

	if err != nil && err != io.EOF {
		return NewProcessingError()
	}

	realContentType := http.DetectContentType(head[:magicBytes])
	
	if !strings.HasPrefix(realContentType, expectedContentType) {
		return NewBadRequestError(fmt.Sprintf("o tipo de conteúdo real %s não corresponde ao tipo de conteúdo esperado %s", realContentType, expectedContentType))
	}

	err = _restartReader(reader)
	if err != nil {
		return err
	}

	return nil
}

func scanVirus(reader io.Reader) error {
	c := clamd.NewClamd("tcp://localhost:3310")

	response, err := c.ScanStream(reader, make(chan bool))

	if err != nil {
		return NewProcessingError()
	}

	for s := range response {
		if s.Status == clamd.RES_FOUND {
			return NewVirusDetectedError(s.Description)
		}
	}

	err = _restartReader(reader)
	if err != nil {
		return err
	}

	return nil
}

func (v *validator) Validate(input UploadInput) error {

	if input.FileName == "" {
		return NewInvalidFileName()
	}

	if input.Size == 0 {
		return NewEmptyFileError()
	}

	if input.Size > v.config.MaxFileSize {
		return NewMaxSizeExceededError(v.config.MaxFileSize)
	}

	ext := filepath.Ext(input.FileName)

	if ext == "" {
		return NewInvalidExtensionError(ext)
	}

	if _, ok := v.config.AllowedFileTypes[ext]; !ok {
		return NewInvalidExtensionError(ext)
	}

	if input.ContentType != v.config.AllowedFileTypes[ext] {
		return NewInvalidContentTypeError(input.ContentType)
	}

	if err := validContentType(input.Reader, input.ContentType); err != nil {
		return err
	}

	if err := scanVirus(input.Reader); err != nil {
		return err
	}

	return nil
}
