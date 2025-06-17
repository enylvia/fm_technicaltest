package repository

import (
	"FM_techincaltest/app/database"
	"FM_techincaltest/models"
	"database/sql"
	"fmt"
)

type EmployeeAbsenceRepository interface {
	GetCompanyByID(id int) (models.CompanyResponse, error)
	ClockInRequest(payload models.AbsencePayload) (int, error)
	ClockOutRequest(payload models.AbsencePayload) (int, error)
	GetAbsenceLog(employee_id int) ([]models.EmployeeAbsenceLog, error)
}
type EmployeeAbsenceImplement struct {
	DB *database.DBClient
}

func NewEmployeeAbsenceImplement(db *database.DBClient) EmployeeAbsenceRepository {
	return &EmployeeAbsenceImplement{DB: db}
}

func (e EmployeeAbsenceImplement) GetCompanyByID(id int) (models.CompanyResponse, error) {
	var result models.CompanyResponse

	var sqlTime sql.NullTime
	row := e.DB.DB.QueryRow(`SELECT * FROM get_company_by_id($1)`, id)
	err := row.Scan(
		&result.ID,
		&result.CompanyName,
		&result.Address,
		&result.Latitude,
		&result.Longitude,
		&result.RadiusMeters,
		&sqlTime,
	)
	if err != nil {
		return models.CompanyResponse{}, fmt.Errorf("error scanning company data from SP result: %w", err)
	}
	if sqlTime.Valid {
		result.CheckInTime = sqlTime
	}

	return result, nil
}
func (e EmployeeAbsenceImplement) ClockInRequest(payload models.AbsencePayload) (int, error) {
	var absenceID int
	err := e.DB.DB.QueryRow(
		`SELECT absence_clock_in($1, $2, $3, $4, $5)`,
		payload.Absence.EmployeeID,
		payload.Absence.ClockInTime,
		payload.Absence.ClockInPhotoURL,
		payload.Absence.Status,
		payload.Absence.Notes,
	).Scan(&absenceID)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to call check_in_employee function: %w", err)
	}

	return absenceID, nil
}

func (e EmployeeAbsenceImplement) ClockOutRequest(payload models.AbsencePayload) (int, error) {
	var affectedRows int
	err := e.DB.DB.QueryRow(
		`SELECT absence_out($1, $2, $3, $4)`,
		payload.Absence.ID,
		payload.Absence.ClockOutTime,
		payload.Absence.ClockOutPhotoURL,
		payload.Absence.Notes,
	).Scan(&affectedRows)

	if err != nil {
		return 0, fmt.Errorf("repository: failed to call absence_out function: %w", err)
	}

	return affectedRows, nil
}

func (e EmployeeAbsenceImplement) GetAbsenceLog(employee_id int) ([]models.EmployeeAbsenceLog, error) {
	var results []models.EmployeeAbsenceLog

	rows, err := e.DB.DB.Query(`SELECT * from get_absence_history($1)`, employee_id)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to call get_absence_history function: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			clockIntNullSql sql.NullTime
			clockOutNullSql sql.NullTime
			dateabsence     sql.NullTime
		)
		log := models.EmployeeAbsenceLog{}
		rows.Scan(&log.ID, &clockIntNullSql, &clockOutNullSql, &dateabsence)
		if err != nil {
			return nil, err
		}
		if clockIntNullSql.Valid {
			log.ClockInTime = &clockIntNullSql.Time
		}
		if clockOutNullSql.Valid {
			log.ClockOutTime = &clockOutNullSql.Time
		}
		if dateabsence.Valid {
			log.DateAbsence = dateabsence.Time.Format("2006-01-02")
		}
		results = append(results, log)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
