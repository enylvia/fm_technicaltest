package models

import "time"

type AbsenceClockIn struct {
	Latitude        float64 `json:"latitude" example:"-6.200000"`
	Longitude       float64 `json:"longitude" example:"106.816666"`
	ClockInPhotoURL string  `json:"clock_in_photo_url" form:"clock_in_photo_url" example:"https://example.com/absensi_foto/checkin_john_doe_20250616.jpg"`
	Notes           string  `json:"notes" form:"notes" example:"optional"`
}
type AbsenceClockOut struct {
	ID               int64   `json:"id" form:"id"`
	Latitude         float64 `json:"latitude" example:"-6.200000"`
	Longitude        float64 `json:"longitude" example:"106.816666"`
	ClockOutPhotoURL string  `json:"clock_in_photo_url" form:"clock_in_photo_url" example:"https://example.com/absensi_foto/checkin_john_doe_20250616.jpg"`
}
type Absence struct {
	ID               int64     `json:"id" form:"id"`
	EmployeeID       int       `json:"employee_id" form:"employee_id"`
	ClockInTime      time.Time `json:"clock_in_time" form:"clock_in_time"`
	ClockOutTime     time.Time `json:"clock_out_time" form:"clock_out_time"`
	ClockInPhotoURL  string    `json:"clock_in_photo_url" form:"clock_in_photo_url"`
	ClockOutPhotoURL string    `json:"clock_out_photo_url" form:"clock_out_photo_url"`
	Status           string    `json:"status" form:"status"`
	Notes            string    `json:"notes" form:"notes"`
	Type             string    `json:"type" form:"type"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
}
type AbsencePayload struct {
	Email   string  `json:"email"`
	Absence Absence `json:"absence"`
}
type ResponseAbsence struct {
	AbsenceID int `json:"absence_id"`
}

type EmployeeAbsenceLog struct {
	ID           int64      `json:"id"`
	ClockInTime  *time.Time `json:"clock_in_time"`
	ClockOutTime *time.Time `json:"clock_out_time"`
	DateAbsence  string     `json:"date_absence"`
}
