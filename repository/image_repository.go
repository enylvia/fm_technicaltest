package repository

import (
	"fmt"
	"os"
	"time"
)

type ImageRepository interface {
	SaveImage(jenis string, fileData []byte, fileExtension string) (string, error)
}

type imageRepository struct {
	uploadDir string
}

func NewImageRepository(uploadDir string) ImageRepository {
	return &imageRepository{uploadDir: uploadDir}
}

func (r *imageRepository) SaveImage(jenis string, fileData []byte, fileExtension string) (string, error) {
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), fileExtension)
	folderPath := fmt.Sprintf("%s%s", r.uploadDir, jenis)
	filePath := fmt.Sprintf("%s/%s", folderPath, filename)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", err
		}
		if err := os.Chmod(folderPath, 0755); err != nil {
			return "", err
		}
	}

	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	if err := os.Chmod(filePath, 0644); err != nil {
		return "", err
	}

	return "/" + jenis + "/" + filename, nil
}
