package handler

import (
	"FM_techincaltest/models"
	"FM_techincaltest/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EmployeeHandler interface {
	ClockInRequest(c echo.Context) error
	ClockOutRequest(c echo.Context) error
	AbsenceHistory(c echo.Context) error
}

type EmployeeHandlerImplement struct {
	employeService service.EmployeeService
}

func NewEmployeeHandlerImplement(employeService service.EmployeeService) EmployeeHandler {
	return &EmployeeHandlerImplement{employeService: employeService}
}

// ClockInRequest godoc
// @Summary Clock In Request
// @Description Melakukan clock-in (absensi masuk) oleh karyawan berdasarkan lokasi dan waktu
// @Tags absence
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param payload body models.AbsenceClockIn true "Data absensi masuk (lokasi, waktu, dll.)"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /employee/clock_in [post]
func (e EmployeeHandlerImplement) ClockInRequest(c echo.Context) error {
	var payload models.AbsenceClockIn
	if err := c.Bind(&payload); err != nil {
		c.Logger().Errorf("Error binding absence payload: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format. Please check your input.",
		})
	}
	emailUser, ok := c.Get("email").(string)
	if !ok {
		c.Logger().Error("Email not found in context after JWT middleware. Invalid setup.")
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Authentication context missing.",
		})
	}
	absenceID, err := e.employeService.ClockInRequest(emailUser, payload)
	if err != nil {
		c.Logger().Errorf("Service error during employee check-in: %v", err)
		if err.Error() == "location outside designated work area" {
			return c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Check-in failed: You are outside the designated work area.",
			})
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Check-in failed: Something went wrong.",
		})
	}
	return c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Data Successfully Created",
		Data:    absenceID,
	})
}

// ClockOutRequest godoc
// @Summary Clock Out Request
// @Description Melakukan clock-out (absensi pulang) oleh karyawan berdasarkan lokasi
// @Tags absence
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param payload body models.AbsenceClockOut true "Data absensi pulang (lokasi, waktu, dll.)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /employee/clock_out [post]
func (e EmployeeHandlerImplement) ClockOutRequest(c echo.Context) error {
	var payload models.AbsenceClockOut
	if err := c.Bind(&payload); err != nil {
		c.Logger().Errorf("Error binding absence payload: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format. Please check your input.",
		})
	}
	emailUser, ok := c.Get("email").(string)
	if !ok {
		c.Logger().Error("Email not found in context after JWT middleware. Invalid setup.")
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Authentication context missing.",
		})
	}
	err := e.employeService.ClockOutRequest(emailUser, payload)
	if err != nil {
		c.Logger().Errorf("Service error during employee check-out: %v", err)
		if err.Error() == "location outside designated work area" {
			return c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Check-out failed: You are outside the designated work area.",
			})
		} else if err.Error() == "not yet clock-in" {
			return c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Check-out failed: You not yet clock-in.",
			})
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Check-out failed: Something went wrong.",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Data Successfully Created",
	})
}

// AbsenceHistory godoc
// @Summary Get Absence History
// @Description Mengambil riwayat absensi berdasarkan email dari JWT
// @Tags absence
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /employee/absence/log [get]
func (e EmployeeHandlerImplement) AbsenceHistory(c echo.Context) error {
	emailUser, ok := c.Get("email").(string)
	if !ok {
		c.Logger().Error("Email not found in context after JWT middleware. Invalid setup.")
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Authentication context missing.",
		})
	}
	absenceLog, err := e.employeService.AbsenceHistory(emailUser)
	if err != nil {
		c.Logger().Errorf("Service error during employee check-out: %v", err)
		if err.Error() == "record not found" {
			return c.JSON(http.StatusOK, models.Response{
				Success: true,
				Message: "Absence-History Success: Record not Found.",
			})
		}
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Absence-History failed: Something went wrong.",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Successfully Get Data",
		Data:    absenceLog,
	})
}
