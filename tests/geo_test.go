package tests

import (
	"math"
	"testing"

	"github.com/Zoltan-Balazs/wot/src/geo"
)

func TestHaversineSamePoint(t *testing.T) {
	if d := geo.Haversine(47.5, 19.0, 47.5, 19.0); d != 0 {
		t.Fatalf("expected 0 for same point, got %f", d)
	}
}

func TestHaversineBudapestVienna(t *testing.T) {
	// Budapest → Vienna is ~214 km great-circle (not road distance)
	d := geo.Haversine(47.4979, 19.0402, 48.2082, 16.3738)
	if math.Abs(d-214) > 5 {
		t.Fatalf("expected ~214 km Budapest→Vienna, got %.1f km", d)
	}
}

func TestBoundingBoxEquator(t *testing.T) {
	// At the equator, kmPerDeg (111.32 km) ≈ 1° in both lat and lon
	minLat, maxLat, minLon, maxLon := geo.BoundingBox(0, 0, 111.32)
	if math.Abs((maxLat-minLat)-2.0) > 0.01 {
		t.Fatalf("expected ~2° lat span at equator, got %.4f", maxLat-minLat)
	}
	if math.Abs((maxLon-minLon)-2.0) > 0.01 {
		t.Fatalf("expected ~2° lon span at equator, got %.4f", maxLon-minLon)
	}
}

func TestBoundingBoxEdgeMidpointsOnCircle(t *testing.T) {
	// The midpoint of each bounding box edge should lie ~distKm from the centre
	lat, lon, dist := 47.5, 19.0, 100.0
	minLat, maxLat, minLon, maxLon := geo.BoundingBox(lat, lon, dist)
	type pt struct{ lat, lon float64 }
	mids := []pt{
		{lat, minLon}, {lat, maxLon},
		{minLat, lon}, {maxLat, lon},
	}
	for _, p := range mids {
		d := geo.Haversine(lat, lon, p.lat, p.lon)
		if math.Abs(d-dist) > 1.0 {
			t.Fatalf("edge midpoint %.4f,%.4f: got %.2f km, want ~%.2f km", p.lat, p.lon, d, dist)
		}
	}
}
