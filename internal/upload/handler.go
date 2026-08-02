package upload

import (
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type Handler struct{
	service UploadService
}

func NewHandler(service UploadService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Upload(c *gin.Context) {
	var uploadError *UploadError

	fileHeader, err := c.FormFile("file")

	if err != nil {
		if err == http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "EMPTY_FILE",
				"message": "o arquivo enviado está vazio",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "PROCESSING_ERROR",
			"message": "falha ao processar o arquivo enviado",
		})
		return
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "PROCESSING_ERROR",	
			"message": "falha ao abrir o arquivo enviado",
		})
		return	
	}

	defer file.Close()

	input := UploadInput{
		Reader: file,
		FileName: fileHeader.Filename,
		Size: fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	output, err := h.service.Upload(input)

	if errors.As(err, &uploadError) {
		c.JSON(uploadError.HTTPStatus, gin.H{
			"code": uploadError.Code,
			"message": uploadError.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, output)
}