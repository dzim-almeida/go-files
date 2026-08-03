package upload

import (
	"path/filepath"

	"github.com/google/uuid"
)

type service struct {
	storage    UploadStorage
	validator  UploadValidator
	repository UploadRepository
}

type UploadService interface {
	Upload(input UploadInput) (*UploadOutput, error)
	GetUploads() ([]UploadOutput, error)
	FindUploadById(uploadId int) (*UploadOutput, error)
	FindUploadByOrgininalFilename(originalFilename string) ([]UploadOutput, error)
	DeleteUpload(uploadId int) error
}

func NewService(storage UploadStorage, validator UploadValidator, repository UploadRepository) UploadService {
	return &service{
		storage:    storage,
		validator:  validator,
		repository: repository,
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
	err := s.validator.Validate(input)

	if err != nil {
		return nil, err
	}

	uniqueFilename, err := generateUniqueFilename(input.FileName)

	if err != nil {
		return nil, err
	}

	storagePath, err := s.storage.Save(input.Reader, uniqueFilename)

	if err != nil {
		return nil, err
	}

	uploadModel := UploadModel{
		OriginalFileName: input.FileName,
		StoredFileName:   uniqueFilename,
		StoragePath:      storagePath,
		ContentType:      input.ContentType,
		Size:             input.Size,
	}

	uploadId, err := s.repository.CreateUpload(uploadModel)

	if err != nil {
		return nil, err
	}

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
	uploads, err := s.repository.GetUploads()

	if err != nil {
		return nil, err
	}

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
	upload, err := s.repository.FindUploadById(uploadId)

	if err != nil {
		return nil, err
	}

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
	uploads, err := s.repository.FindUploadByOriginalFilename(originalFilename)

	if err != nil {
		return nil, err
	}

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
	upload, err := s.repository.FindUploadById(uploadId)

	err = s.storage.Delete(upload.StoragePath)

	if err != nil {
		return err
	}

	err = s.repository.DeleteUploadById(upload.ID)

	if err != nil {
		return err
	}

	return nil
}