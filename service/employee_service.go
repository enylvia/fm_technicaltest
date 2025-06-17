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

type EmployeeService interface {
	ClockInRequest(emailUser string, payload models.AbsenceClockIn) (models.ResponseAbsence, error)
	ClockOutRequest(emailUser string, payload models.AbsenceClockOut) error
	AbsenceHistory(email string) ([]models.EmployeeAbsenceLog, error)
}
type EmployeeServiceImplement struct {
	userRepository     repository.UserRepository
	employeeRepository repository.EmployeeAbsenceRepository
}

func NewEmployeeServiceImplement(userRepository repository.UserRepository,
	employeeRepository repository.EmployeeAbsenceRepository) EmployeeService {
	return &EmployeeServiceImplement{
		userRepository:     userRepository,
		employeeRepository: employeeRepository,
	}
}
func (e EmployeeServiceImplement) ClockInRequest(emailUser string, payload models.AbsenceClockIn) (models.ResponseAbsence, error) {
	var absencePayload models.AbsencePayload
	userData, err := e.userRepository.GetUserByEmail(emailUser)
	if err != nil {
		return models.ResponseAbsence{}, fmt.Errorf("something went wrong when get user data: %w", err)
	}
	companyDetail, err := e.employeeRepository.GetCompanyByID(userData.CompanyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ResponseAbsence{}, fmt.Errorf("company not found: %w", err)
		}
		return models.ResponseAbsence{}, fmt.Errorf("something went wrong when get company data: %w", err)
	}
	if companyDetail.Latitude != 0 && companyDetail.Longitude != 0 && companyDetail.RadiusMeters > 0 {
		isInside := helpers.IsWithinRadius(payload.Latitude, payload.Longitude,
			companyDetail.Latitude, companyDetail.Longitude, companyDetail.RadiusMeters)
		if !isInside {
			return models.ResponseAbsence{}, errors.New("location outside designated work area")
		}
	} else {
		fmt.Printf("Warning: Company ID %d has incomplete location data. Geofencing skipped.\n", companyDetail.ID)
	}

	var standardCheckInHour, standardCheckInMinute int
	clockInTime := time.Now().In(time.Local)

	if companyDetail.CheckInTime.Valid {
		standardCheckInHour = companyDetail.CheckInTime.Time.Hour()
		standardCheckInMinute = companyDetail.CheckInTime.Time.Minute()
	} else {
		standardCheckInHour = 9
		standardCheckInMinute = 0
		fmt.Printf("Warning: Company ID %d has no standard check-in time. Using default 09:00.\n", companyDetail.ID)
	}
	referenceTime := time.Date(
		clockInTime.Year(), clockInTime.Month(), clockInTime.Day(),
		standardCheckInHour, standardCheckInMinute, 0, 0,
		clockInTime.Location(),
	)
	status := "present"
	lateToleranceSeconds := time.Duration(userData.LateTolerance) * time.Minute // Ambil toleransi dari data employee

	if clockInTime.After(referenceTime.Add(lateToleranceSeconds)) {
		status = "late"
	}
	absencePayload.Absence.ClockInTime = clockInTime
	absencePayload.Absence.Status = status
	absencePayload.Absence.Notes = payload.Notes
	absencePayload.Absence.ClockInPhotoURL = payload.ClockInPhotoURL
	absencePayload.Absence.EmployeeID = userData.EmployeeID

	absenceID, err := e.employeeRepository.ClockInRequest(absencePayload)

	if err != nil {
		return models.ResponseAbsence{}, fmt.Errorf("error when do absence : %w", err)
	}
	return models.ResponseAbsence{AbsenceID: absenceID}, nil
}
func (e EmployeeServiceImplement) ClockOutRequest(emailUser string, payload models.AbsenceClockOut) error {
	var absencePayload models.AbsencePayload
	userData, err := e.userRepository.GetUserByEmail(emailUser)
	if err != nil {
		return fmt.Errorf("something went wrong when get user data: %w", err)
	}
	companyDetail, err := e.employeeRepository.GetCompanyByID(userData.CompanyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("company not found: %w", err)
		}
		return fmt.Errorf("something went wrong when get company data: %w", err)
	}
	if companyDetail.Latitude != 0 && companyDetail.Longitude != 0 && companyDetail.RadiusMeters > 0 {
		isInside := helpers.IsWithinRadius(payload.Latitude, payload.Longitude,
			companyDetail.Latitude, companyDetail.Longitude, companyDetail.RadiusMeters)
		if !isInside {
			return errors.New("location outside designated work area")
		}
	} else {
		fmt.Printf("Warning: Company ID %d has incomplete location data. Geofencing skipped.\n", companyDetail.ID)
	}
	clockOutTime := time.Now().In(time.Local)
	absencePayload.Absence.ClockOutTime = clockOutTime
	absencePayload.Absence.ID = payload.ID
	absencePayload.Absence.ClockOutPhotoURL = payload.ClockOutPhotoURL
	rowAffected, err := e.employeeRepository.ClockOutRequest(absencePayload)
	if err != nil {
		return fmt.Errorf("failed to clock-out : %w", err)
	}

	if rowAffected > 0 {
		return nil
	}
	return errors.New("not yet clock-in")
}
func (e EmployeeServiceImplement) AbsenceHistory(email string) ([]models.EmployeeAbsenceLog, error) {
	userData, err := e.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("something went wrong when get user data: %w", err)
	}
	absenceHistory, err := e.employeeRepository.GetAbsenceLog(userData.EmployeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("something went wrong when get absence history: %w", err)
	}
	return absenceHistory, nil
}
