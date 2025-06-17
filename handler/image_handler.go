package handler

import (
	"FM_techincaltest/models"
	"FM_techincaltest/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ImageHandler interface {
	UploadImage(c echo.Context) error
}

type ImageHandlerImplement struct {
	imageService service.ImageService
}

func NewImageHandler(imageService service.ImageService) ImageHandler {
	return &ImageHandlerImplement{imageService: imageService}
}

// UploadImage godoc
// @Summary Upload Image
// @Description Mengunggah gambar berdasarkan jenis dan file
// @Tags image
// @Accept mpfd
// @Security BearerAuth
// @Produce json
// @Param jenis formData string true "Jenis file (misalnya: avatar, profile, etc)"
// @Param file formData file true "File gambar yang akan diunggah"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /image/save [post]
func (i ImageHandlerImplement) UploadImage(c echo.Context) error {
	jenis := c.FormValue("jenis")
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		c.Logger().Errorf("Failed to get file from request: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Failed to get file from request: " + err.Error(),
		})
	}

	filename, err := i.imageService.UploadImage(jenis, file, fileHeader)
	if err != nil {
		c.Logger().Errorf("Failed to upload image: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Failed to upload image: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success",
		Data:    filename,
	})
}
