package models

import (
	"database/sql"
)

type CompanyResponse struct {
	ID           int64        `json:"id"`
	CompanyName  string       `json:"company_name"`
	Address      string       `json:"address"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	RadiusMeters int          `json:"radius_meters"`
	CheckInTime  sql.NullTime `json:"check_in_time"`
}
