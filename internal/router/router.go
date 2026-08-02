package router

import (
	"github.com/dzim-almeida/go-files/internal/database/postgres"
	"github.com/dzim-almeida/go-files/internal/upload"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	conn := postgres.ConnDB()
	repository := upload.NewRepository(conn)
	storage := upload.NewLocalStorage("C:/dev/go-files/uploads")

	validator := upload.NewValidator(
		upload.ValidatorConfig{
			MaxFileSize:      upload.MaxFileSize,
			AllowedFileTypes: upload.AllowedFileTypes,
		},
	)

	service := upload.NewService(
		storage,
		validator,
		repository,
	)
	
	handler := upload.NewHandler(service)

	router.POST("/upload", handler.Upload)

	return router
}