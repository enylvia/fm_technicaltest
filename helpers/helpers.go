package helpers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math"
)

const EarthRadius = 6371000

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var (
		degLat = (lat2 - lat1) * (math.Pi / 180)
		degLon = (lon2 - lon1) * (math.Pi / 180)
		a      = math.Sin(degLat/2)*math.Sin(degLat/2) +
			math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
				math.Sin(degLon/2)*math.Sin(degLon/2)
		c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	)
	return EarthRadius * c
}

func IsWithinRadius(pointLat, pointLon, centerLat, centerLon float64, radiusMeters int) bool {
	distance := HaversineDistance(pointLat, pointLon, centerLat, centerLon)
	return distance <= float64(radiusMeters)
}
