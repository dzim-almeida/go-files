package upload

import (
	"fmt"
	"net/http"
)

type CodeError string

const (
	ErrCodeInvalidFilename    CodeError = "INVALID_FILENAME"
	ErrCodeEmptyFile          CodeError = "EMPTY_FILE"
	ErrCodeInvalidExtension   CodeError = "INVALID_EXTENSION"
	ErrCodeMaxSizeExceeded    CodeError = "FILE_TOO_LARGE"
	ErrCodeInvalidContentType CodeError = "INVALID_CONTENT_TYPE"
	ErrCodeVirusDetected      CodeError = "VIRUS_DETECTED"
	ErrProcessingError        CodeError = "PROCESSING_ERROR"
	ErrSaveError              CodeError = "SAVE_ERROR"
	ErrDeleteError            CodeError = "DELETE_ERROR"
	ErrNotFound               CodeError = "NOT_FOUND"
	ErrBadRequest             CodeError = "BAD_REQUEST"
	ErrInternalServerError    CodeError = "INTERNAL_SERVER_ERROR"
)

type UploadError struct {
	Code       CodeError
	Message    string
	HTTPStatus int
}

func (e *UploadError) Error() string {
	return e.Message
}

func NewInvalidFileName() error {
	return &UploadError{
		Code:       ErrCodeInvalidFilename,
		Message:    "o nome do arquivo é inválido",
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewEmptyFileError() error {
	return &UploadError{
		Code:       ErrCodeEmptyFile,
		Message:    "o arquivo enviado está vazio",
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInvalidExtensionError(ext string) error {
	return &UploadError{
		Code:       ErrCodeInvalidExtension,
		Message:    fmt.Sprintf("a extensão %s não é permitida", ext),
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

func NewMaxSizeExceededError(maxSize int64) error {
	return &UploadError{
		Code:       ErrCodeMaxSizeExceeded,
		Message:    fmt.Sprintf("o tamanho do arquivo excede o limite de %d bytes", maxSize),
		HTTPStatus: http.StatusRequestEntityTooLarge,
	}
}

func NewInvalidContentTypeError(contentType string) error {
	return &UploadError{
		Code:       ErrCodeInvalidContentType,
		Message:    fmt.Sprintf("o tipo de conteúdo %s não é permitido", contentType),
		HTTPStatus: http.StatusUnsupportedMediaType,
	}
}

func NewVirusDetectedError(desc string) error {
	return &UploadError{
		Code:       ErrCodeVirusDetected,
		Message:    fmt.Sprintf("foi detectado um tipo de vírus no arquivo enviado: %s", desc),
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

func NewProcessingError() error {
	return &UploadError{
		Code:       ErrProcessingError,
		Message:    "falha ao processar o arquivo enviado",
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewSaveError() error {
	return &UploadError{
		Code:       ErrSaveError,
		Message:    "falha ao salvar o arquivo enviado",
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewDeleteError() error {
	return &UploadError{
		Code:       ErrDeleteError,
		Message:    "falha ao deletar o arquivo",
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewNotFoundError() error {
	return &UploadError{
		Code:       ErrNotFound,
		Message:    "arquivo não encontrado",
		HTTPStatus: http.StatusNotFound,
	}
}

func NewBadRequestError(message string) error {
	return &UploadError{
		Code:      	ErrBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInternalServerError() error {
	return &UploadError{
		Code:       ErrInternalServerError,
		Message:    "erro interno do servidor",
		HTTPStatus: http.StatusInternalServerError,
	}
}
