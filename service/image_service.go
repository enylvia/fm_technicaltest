package service

import (
	"FM_techincaltest/repository"
	"errors"
	"mime/multipart"
)

type ImageService interface {
	UploadImage(jenis string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type imageService struct {
	imageRepo repository.ImageRepository
}

func NewImageService(imageRepo repository.ImageRepository) ImageService {
	return &imageService{imageRepo: imageRepo}
}

func (s *imageService) UploadImage(jenis string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	// Read file data
	fileData := make([]byte, fileHeader.Size)
	if _, err := file.Read(fileData); err != nil {
		return "", err
	}
	allowedExtensions := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}
	fileExtension := allowedExtensions[fileHeader.Header.Get("Content-Type")]
	if fileExtension == "" {
		return "", errors.New("unsupported file type")
	}
	return s.imageRepo.SaveImage(jenis, fileData, fileExtension)
}
