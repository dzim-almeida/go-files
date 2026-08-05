package upload

import (
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
)

type service struct {
	storage    UploadStorage
	validator  UploadValidator
	repository UploadRepository
	logger     *slog.Logger
}

type UploadService interface {
	Upload(input UploadInput) (*UploadOutput, error)
	GetUploads() ([]UploadOutput, error)
	FindUploadById(uploadId int) (*UploadOutput, error)
	FindUploadByOrgininalFilename(originalFilename string) ([]UploadOutput, error)
	DeleteUpload(uploadId int) error
}

func NewService(storage UploadStorage, validator UploadValidator, repository UploadRepository, logger *slog.Logger) UploadService {
	return &service{
		storage:    storage,
		validator:  validator,
		repository: repository,
		logger:     logger,
	}
}

func generateUniqueFilename(originalFilename string) (string, error) {
	ext := filepath.Ext(originalFilename)

	uuid, err := uuid.NewRandom()

	if err != nil {
		return "", NewProcessingError()
	}

	uniqueFilename := uuid.String() + ext

	return uniqueFilename, nil
}

func (s *service) Upload(input UploadInput) (*UploadOutput, error) {
	logger := s.logger.With(
		"method", "Upload",
		"filename", input.FileName,
		"size", input.Size,
	)

	logger.Info(
		"iniciando upload",
		"filename", input.FileName,
		"size", input.Size,
		"content-type", input.ContentType,
	)

	logger.Debug("validando upload")
	err := s.validator.Validate(input)

	if err != nil {
		logger.Warn(
			"erro na validação",
			"error", err,
		)
		return nil, err
	}

	uniqueFilename, err := generateUniqueFilename(input.FileName)

	if err != nil {
		logger.Warn(
			"erro na geração de nome único",
			"generated filename", uniqueFilename,
			"error", err,
		)
		return nil, err
	}

	logger.Debug("salvando arquivo no storage")
	storagePath, err := s.storage.Save(input.Reader, uniqueFilename)

	if err != nil {
		logger.Error(
			"erro ao salvar arquivo",
			"error", err,
		)
		return nil, err
	}

	uploadModel := UploadModel{
		OriginalFileName: input.FileName,
		StoredFileName:   uniqueFilename,
		StoragePath:      storagePath,
		ContentType:      input.ContentType,
		Size:             input.Size,
	}

	logger.Debug("salvando metadados do upload no banco de dados")
	uploadId, err := s.repository.CreateUpload(uploadModel)

	if err != nil {
		logger.Error(
			"erro ao salvar metadados do upload",
			"error", err,
		)
		_ = s.storage.Delete(storagePath)
		return nil, err
	}

	logger.Info(
		"upload concluído com sucesso!",
		"upload id", uploadId,
		"path", storagePath,
		"stored filename", uniqueFilename,
	)
	output := &UploadOutput{
		Id:               uploadId,
		OriginalFileName: input.FileName,
		StoredFileName:   uniqueFilename,
		StoragePath:      storagePath,
		ContetType:       input.ContentType,
		Size:             input.Size,
	}

	return output, nil
}

func (s *service) GetUploads() ([]UploadOutput, error) {
	logger := s.logger.With(
		"method", "GetUploads",
	)

	logger.Debug(
		"buscando uploads no banco de dados",
	)

	uploads, err := s.repository.GetUploads()

	if err != nil {
		logger.Error(
			"erro na busca dos uploads",
			"error", err,
		)
		return nil, err
	}

	logger.Info(
		"uploads encontrados",
		"count", len(uploads),
	)

	outputs := make([]UploadOutput, 0, len(uploads))

	for _, upload := range uploads {
		outputs = append(outputs, UploadOutput{
			Id:               upload.ID,
			OriginalFileName: upload.OriginalFileName,
			StoredFileName:   upload.StoredFileName,
			StoragePath:      upload.StoragePath,
			ContetType:       upload.ContentType,
			Size:             upload.Size,
		})
	}

	return outputs, nil
}

func (s *service) FindUploadById(uploadId int) (*UploadOutput, error) {
	logger := s.logger.With(
		"method", "FindUploadById",
		"upload id", uploadId,
	)

	logger.Debug(
		"buscando upload no banco de dados",
	)

	upload, err := s.repository.FindUploadById(uploadId)

	if err != nil {
		logger.Error(
			"erro na busca do upload",
			"err", err,
		)
		return nil, err
	}

	logger.Info(
		"upload encontrado",
		"upload id", upload.ID,
		"path", upload.StoragePath,
		"stored filename", upload.StoredFileName,
	)

	output := &UploadOutput{
		Id:               upload.ID,
		OriginalFileName: upload.OriginalFileName,
		StoredFileName:   upload.OriginalFileName,
		StoragePath:      upload.StoragePath,
		ContetType:       upload.ContentType,
		Size:             upload.Size,
	}

	return output, nil
}

func (s *service) FindUploadByOrgininalFilename(originalFilename string) ([]UploadOutput, error) {
	logger := s.logger.With(
		"method", "FindUploadByOrgininalFilename",
		"search filename", originalFilename,
	)

	logger.Debug(
		"buscando upload no banco de dados",
	)
	uploads, err := s.repository.FindUploadByOriginalFilename(originalFilename)

	if err != nil {
		logger.Error(
			"erro na busca do upload",
			"error", err,
		)
		return nil, err
	}

	logger.Debug(
		"uploads encontrados",
		"count", len(uploads),
	)

	outputs := make([]UploadOutput, 0, len(uploads))

	for _, upload := range uploads {
		outputs = append(outputs, UploadOutput{
			Id:               upload.ID,
			OriginalFileName: upload.OriginalFileName,
			StoredFileName:   upload.StoredFileName,
			StoragePath:      upload.StoragePath,
			ContetType:       upload.ContentType,
			Size:             upload.Size,
		})
	}

	return outputs, nil
}

func (s *service) DeleteUpload(uploadId int) error {
	logger := s.logger.With(
		"method", "DeleteUpload",
		"upload id", uploadId,
	)

	logger.Debug(
		"verificando upload no banco de dados",
	)

	upload, err := s.repository.FindUploadById(uploadId)

	logger.Debug(
		"deletando upload do storage",
	)

	err = s.storage.Delete(upload.StoragePath)

	if err != nil {
		logger.Error(
			"erro em deletar no storage",
			"error", err,
		)
		return err
	}

	logger.Debug(
		"deletando upload do banco de dados",
	)

	err = s.repository.DeleteUploadById(upload.ID)

	if err != nil {
		logger.Error(
			"erro ao deletar registro no banco de dados",
			"error", err,
		)
		return err
	}

	logger.Info(
		"upload deletado com sucesso",
		"upload id", upload.ID,
		"path", upload.StoragePath,
		"stored filename", upload.StoredFileName,
	)

	return nil
}
