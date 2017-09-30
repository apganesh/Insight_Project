package main

import (
	"math"

	"github.com/golang/geo/s1"
)

const earthRadiusKm = 6371.01

func toRadians(num float64) float64 {
	return num * math.Pi / 180
}

func getHaversineDist(startLat, startLng, endLat, endLng float64) float64 {

	dLat := toRadians(endLat - startLat)
	dLon := toRadians(endLng - startLng)

	lat1 := toRadians(startLat)
	lat2 := toRadians(endLat)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*
			math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Convert distance on the Earth's surface to an angle.
func kmToAngle(km float64) s1.Angle {
	res := s1.Angle(km / earthRadiusKm)
	return res
}
