package repository

import (
	"FM_techincaltest/app/database"
	"FM_techincaltest/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	CreateUserAndEmployee(data models.UserandEmployeeRegister) error
	GetUserByEmail(email string) (models.FullResponseUser, error)
}

type UserRepositoryImplement struct {
	DBClient *database.DBClient
}

func NewUserRepository(db *database.DBClient) UserRepository {
	return &UserRepositoryImplement{DBClient: db}
}
func (u UserRepositoryImplement) CreateUserAndEmployee(data models.UserandEmployeeRegister) error {
	tx, err := u.DBClient.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(
		`CALL register_new_user ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		data.Email,
		data.PasswordHash,
		data.IsActive,
		data.CompanyID,
		data.PositionID,
		data.DepartmentID,
		data.FullName,
		data.NIK,
		data.PhoneNumber,
		data.Address,
		data.DateOfBirth,
		data.ProfilePictureURL,
		data.JoinedDate,
		data.LateTolerance)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err) // Wrap error
	}
	return nil
}
func (u UserRepositoryImplement) GetUserByEmail(email string) (models.FullResponseUser, error) {
	var result models.FullResponseUser

	row := u.DBClient.DB.QueryRow(`SELECT * FROM find_user_by_email($1)`, email)

	var profilePicture sql.NullString
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.Password,
		&result.IsActive,
		&result.CompanyID,
		&result.FullName,
		&result.PhoneNumber,
		&result.PositionID,
		&result.DepartmentID,
		&profilePicture,
		&result.JoinedDate,
		&result.LateTolerance,
		&result.EmployeeID,
	)
	if err != nil {
		return models.FullResponseUser{}, fmt.Errorf("error scanning user details from SP result: %w", err)
	}
	if profilePicture.Valid {
		result.ProfilePictureURL = profilePicture.String
	}
	return result, nil
}
