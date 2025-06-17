package service

import (
	"FM_techincaltest/helpers"
	"FM_techincaltest/models"
	"FM_techincaltest/repository"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserService interface {
	CreateUserAndEmployee(data models.UserandEmployeeRegisterPayload) error
	LoginUser(data models.LoginUserPayload) (models.UserReadData, error)
}

type UserServiceImplement struct {
	userRepository repository.UserRepository
}

func NewUserServiceImplement(userRepository repository.UserRepository) UserService {
	return &UserServiceImplement{userRepository: userRepository}
}

func (u UserServiceImplement) CreateUserAndEmployee(data models.UserandEmployeeRegisterPayload) error {
	parsedDateOfBirth, err := time.Parse("2006-01-02", data.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid DateOfBirth format: %w", err)
	}
	parsedJoinedDate, err := time.Parse("2006-01-02", data.JoinedDate)
	if err != nil {
		return fmt.Errorf("invalid JoinedDate format: %w", err)
	}
	hashPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	registeredData := models.UserandEmployeeRegister{
		Email:             data.Email,
		PasswordHash:      hashPassword,
		IsActive:          true,
		CompanyID:         data.CompanyID,
		FullName:          data.FullName,
		NIK:               data.NIK,
		PhoneNumber:       data.PhoneNumber,
		Address:           data.Address,
		DateOfBirth:       parsedDateOfBirth,
		PositionID:        data.PositionID,
		DepartmentID:      data.DepartmentID,
		ProfilePictureURL: data.ProfilePictureURL,
		JoinedDate:        parsedJoinedDate,
		LateTolerance:     data.LateTolerance,
	}
	err = u.userRepository.CreateUserAndEmployee(registeredData)
	if err != nil {
		return fmt.Errorf("failed to create a new user: %w", err)
	}
	return nil
}
func (u UserServiceImplement) LoginUser(data models.LoginUserPayload) (models.UserReadData, error) {
	checkDataExist, err := u.userRepository.GetUserByEmail(data.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.UserReadData{}, fmt.Errorf("user with email %s not found: %w", data.Email, err)
		}
		return models.UserReadData{}, fmt.Errorf("failed to retrieve user data: %w", err)
	}
	// Verify the password given
	isPasswordCorrect := helpers.CheckPasswordHash(data.Password, checkDataExist.Password)
	if !isPasswordCorrect {
		return models.UserReadData{}, errors.New("invalid email or password")
	}
	// Generate jwt Token
	tokenGenerate, err := helpers.GenerateToken(int64(checkDataExist.ID), checkDataExist.Email)
	if err != nil {
		return models.UserReadData{}, fmt.Errorf("could not generate authentication token: %w", err)
	}
	return models.UserReadData{
		ID:                checkDataExist.ID,
		Email:             checkDataExist.Email,
		IsActive:          checkDataExist.IsActive,
		CompanyID:         checkDataExist.CompanyID,
		FullName:          checkDataExist.FullName,
		PhoneNumber:       checkDataExist.PhoneNumber,
		PositionID:        checkDataExist.PositionID,
		DepartmentID:      checkDataExist.DepartmentID,
		ProfilePictureURL: checkDataExist.ProfilePictureURL,
		JoinedDate:        checkDataExist.JoinedDate,
		Token:             tokenGenerate,
	}, err
}
