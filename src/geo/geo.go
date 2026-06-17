package geo

import "math"

const (
	earthRadiusKm = 6371.0

	// kmPerDeg is the approximate number of km per degree of latitude.
	// Earth's mean circumference is ~40 075 km → 40075 / 360 ≈ 111.32 km/°.
	kmPerDeg = 111.32
)

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// BoundingBox returns a lat/lon box guaranteed to contain all points within
// distKm of (lat, lon). Filter results with Haversine to remove box corners.
func BoundingBox(lat, lon, distKm float64) (minLat, maxLat, minLon, maxLon float64) {
	deltaLat := distKm / kmPerDeg
	deltaLon := distKm / (kmPerDeg * math.Cos(lat*math.Pi/180))
	return lat - deltaLat, lat + deltaLat, lon - deltaLon, lon + deltaLon
}
