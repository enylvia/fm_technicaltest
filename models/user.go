package models

import "time"

type UserandEmployeeRegisterPayload struct {
	Email             string `json:"email" form:"email" validate:"required" example:"john.doe@example.com"`
	Password          string `json:"password" form:"password" validate:"required" example:"StrongPassword123!"`
	IsActive          bool   `json:"is_active" form:"is_active" example:"true"`
	CompanyID         int    `json:"company_id" form:"company_id" validate:"required" example:"1"`
	FullName          string `json:"full_name" form:"full_name" validate:"required" example:"John Doe"`
	NIK               string `json:"nik" form:"nik" validate:"required" example:"123456789012345"`
	PhoneNumber       string `json:"phone_number" form:"phone_number" validate:"required" example:"081234567890"`
	Address           string `json:"address" form:"address" example:"Jl. Merdeka No. 10, Jakarta"`
	DateOfBirth       string `json:"date_of_birth" form:"date_of_birth" validate:"required" example:"1990-05-15"`
	PositionID        int    `json:"position_id" form:"position_id" validate:"required" example:"1"`
	DepartmentID      int    `json:"department_id" form:"department_id" validate:"required" example:"1"`
	ProfilePictureURL string `json:"profile_picture_url" form:"profile_picture_url" example:"https://example.com/images/john_doe_profile.jpg"`
	JoinedDate        string `json:"joined_date" form:"joined_date" example:"2023-01-20"`
	LateTolerance     int    `json:"late_tolerance" form:"late_tolerance" example:"15"`
}
type UserandEmployeeRegister struct {
	Email             string    `json:"email" form:"email"`
	PasswordHash      string    `json:"password_hash" form:"password_hash"`
	IsActive          bool      `json:"is_active" form:"is_active"`
	CompanyID         int       `json:"company_id" form:"company_id"`
	FullName          string    `json:"full_name" form:"full_name"`
	NIK               string    `json:"nik" form:"nik"`
	PhoneNumber       string    `json:"phone_number" form:"phone_number"`
	Address           string    `json:"address" form:"address"`
	DateOfBirth       time.Time `json:"date_of_birth" form:"date_of_birth"`
	PositionID        int       `json:"position_id" form:"position_id"`
	DepartmentID      int       `json:"department_id" form:"department_id"`
	ProfilePictureURL string    `json:"profile_picture_url" form:"profile_picture_url"`
	JoinedDate        time.Time `json:"joined_date" form:"joined_date"`
	LateTolerance     int       `json:"late_tolerance" form:"late_tolerance"`
}

type LoginUserPayload struct {
	Email    string `json:"email" form:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" form:"password" validate:"required" example:"StrongPassword123!"`
}

type UserReadData struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	IsActive          bool      `json:"is_active"`
	CompanyID         int       `json:"company_id"`
	FullName          string    `json:"full_name"`
	PhoneNumber       string    `json:"phone_number"`
	PositionID        int       `json:"position_id"`
	DepartmentID      int       `json:"department_id"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	JoinedDate        time.Time `json:"joined_date"`
	Token             string    `json:"token"`
}
type FullResponseUser struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	IsActive          bool      `json:"is_active"`
	CompanyID         int       `json:"company_id"`
	FullName          string    `json:"full_name"`
	PhoneNumber       string    `json:"phone_number"`
	PositionID        int       `json:"position_id"`
	DepartmentID      int       `json:"department_id"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	JoinedDate        time.Time `json:"joined_date"`
	LateTolerance     int       `json:"late_tolerance"`
	Token             string    `json:"token"`
	EmployeeID        int       `json:"employee_id"`
}
type ParamCheckUserByEmail struct {
	Email string `json:"email" form:"email" validate:"required"`
}
