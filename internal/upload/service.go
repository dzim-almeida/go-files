package upload

import (
	"path/filepath"

	"github.com/google/uuid"
)

type service struct {
	storage UploadStorage
	validator UploadValidator
	repository UploadRepository
}

type UploadService interface {
	Upload(input UploadInput) (*UploadOutput, error)
}

func NewService(storage UploadStorage, validator UploadValidator, repository UploadRepository) UploadService {
	return &service{
		storage: storage,
		validator: validator,
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
		OrinialFileName: input.FileName,
		StoredFileName: uniqueFilename,
		StoragePath: storagePath,
		ContentType: input.ContentType,
		Size: input.Size,
	}

	uploadId, err := s.repository.CreateUpload(uploadModel)

	if err != nil {
		return nil, NewSaveError()
	}

	output := &UploadOutput{
		Id: uploadId,
		OriginalFileName: input.FileName,
		StoredFileName: uniqueFilename,
		StoragePath: storagePath,
		ContetType: input.ContentType,
		Size: input.Size,
	}

	return output, nil
}