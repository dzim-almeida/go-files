package upload

import (
	"errors"
	"net/http"
	"strconv"

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

var uploadError *UploadError

func (h *Handler) Upload(c *gin.Context) {
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

func (h *Handler) GetUploads(c *gin.Context) {
	uploads, err := h.service.GetUploads()

	if errors.As(err, &uploadError) {
		c.JSON(uploadError.HTTPStatus, gin.H{
			"code": uploadError.Code,
			"message": uploadError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, uploads)
}

func (h *Handler) FindUploadById(c *gin.Context) {
	id := c.Param("id")

	uploadId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "INVALID_ID",
			"message": "o id fornecido é inválido",
		})
		return
	}

	upload, err := h.service.FindUploadById(uploadId)

	if errors.As(err, &uploadError) {
		c.JSON(uploadError.HTTPStatus, gin.H{
			"code": uploadError.Code,
			"message": uploadError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, upload)
}

func (h *Handler) FindUploadByOriginalFilename(c *gin.Context) {
	filename := c.Param("filename")

	uploads, err := h.service.FindUploadByOrgininalFilename(filename)

	if errors.As(err, &uploadError) {
		c.JSON(uploadError.HTTPStatus, gin.H{
			"code": uploadError.Code,
			"message": uploadError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, uploads)
}

func (h *Handler) DeleteUpload(c *gin.Context) {
	id := c.Param("id")

	uploadId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "INVALID_ID",
			"message": "o id fornecido é inválido",
		})
		return
	}

	err = h.service.DeleteUpload(uploadId)

	if errors.As(err, &uploadError) {
		c.JSON(uploadError.HTTPStatus, gin.H{
			"code": uploadError.Code,
			"message": uploadError.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "UPLOAD_DELETED",
		"message": "upload deletado com sucesso",
	})
}