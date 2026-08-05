package router

import (
	"github.com/dzim-almeida/go-files/internal/database/postgres"
	"github.com/dzim-almeida/go-files/internal/logger"
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

	logger := logger.NewLogger()

	service := upload.NewService(
		storage,
		validator,
		repository,
		logger,
	)
	
	handler := upload.NewHandler(service)

	router.POST("/uploads", handler.Upload)
	router.GET("/uploads", handler.GetUploads)
	router.GET("/uploads/:id", handler.FindUploadById)
	router.GET("/uploads/filename/:filename", handler.FindUploadByOriginalFilename)
	router.DELETE("/uploads/:id", handler.DeleteUpload)

	return router
}