package handler

import (
	"FM_techincaltest/models"
	"FM_techincaltest/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler interface {
	RegisterUserAndEmployee(c echo.Context) error
	LoginUser(c echo.Context) error
}

type UserHandlerImplement struct {
	userService service.UserService
}

func NewUserHandlerImplement(userService service.UserService) UserHandler {
	return &UserHandlerImplement{userService: userService}
}

// RegisterUserAndEmployee godoc
// @Summary Register User and Employee
// @Description Register User and Employee
// @Tags authenticate
// @Accept json
// @Produce json
// @Param payload body models.UserandEmployeeRegisterPayload true "Request body untuk register"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/register [post]
func (u UserHandlerImplement) RegisterUserAndEmployee(c echo.Context) error {
	var payload models.UserandEmployeeRegisterPayload
	if err := c.Bind(&payload); err != nil {
		c.Logger().Errorf("Error binding register payload: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format. Please check your input.",
		})
	}
	if err := u.userService.CreateUserAndEmployee(payload); err != nil {
		c.Logger().Errorf("Validation error for register payload: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Validation failed. Please check the provided data.",
		})
	}
	return c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Data Successfully Created",
	})
}

// LoginUser godoc
// @Summary Login User
// @Description Autentikasi user dengan email dan password
// @Tags authenticate
// @Accept json
// @Produce json
// @Param payload body models.LoginUserPayload true "Email dan password untuk login"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/login [post]
func (u UserHandlerImplement) LoginUser(c echo.Context) error {
	var payload models.LoginUserPayload
	if err := c.Bind(&payload); err != nil {
		c.Logger().Errorf("Error binding login payload: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format. Please check your input.",
		})
	}
	userData, err := u.userService.LoginUser(payload)
	if err != nil {
		c.Logger().Errorf("Error for login err: %v", err)
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Login failed. Please check the provided data.",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Login Successfully",
		Data:    userData,
	})
}
