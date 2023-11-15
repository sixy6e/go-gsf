package gsf

import (
	"math"
)

// Wgs84Coefficients contains the coefficients used to calculate longitude and latitude
// coordinates based on beam locations defined by across_track and along_track values.
// See https://en.wikipedia.org/wiki/Geographic_coordinate_system for more information.
// Also see https://math.stackexchange.com/questions/389942/why-is-it-necessary-to-use-sin-or-cos-to-determine-heading-dead-reckoning for complete calculations.
// These coeficients appear to be derived from an iterative process that is described here:
// https://gis.stackexchange.com/questions/75528/understanding-terms-in-length-of-degree-formula
type GeoCoefficients struct {
	A float64
	B float64
	C float64
	D float64
	E float64
	F float64
	G float64
}

// NewCoefWgs84 initialises a GeoCoefficients with coefficients set for WGS84.
// No thoughts, as of yet, to generate coefficients for other datums.
// Or is an alternative method more suited.
func NewCoefWgs84() *GeoCoefficients {
	g := new(GeoCoefficients)
	g.A = 111132.92
	g.B = 559.82
	g.C = 1.175
	g.D = 0.0023
	g.E = 111412.84
	g.F = 93.5
	g.G = 0.118

	return g
}

type LonLat struct {
	Longitude []float64
	Latitude  []float64
}

// BeamsLonLat calculates arrays of longitude and latitude of len(along_track).
// Most likely the func will change; potentially a method for ping data, header or GeoCoefficients.
// For formulae details: https://gis.stackexchange.com/questions/75528/understanding-terms-in-length-of-degree-formula
func BeamsLonLat(lon, lat float64, heading float32, along_track, across_track []float32, coef GeoCoefficients) LonLat {
	var (
		acr_trck float64
		aln_trck float64
		lonlat   LonLat
		deg2rad  float64
	)

	deg2rad = math.Pi / 180.0

	lat_rad := deg2rad * lat
	head_rad := deg2rad * float64(heading)

	// latitude metres scale factor
	lat_sf :=
		coef.A -
			coef.B*math.Cos(2.0*lat_rad) +
			coef.C*math.Cos(4.0*lat_rad) -
			coef.D*math.Cos(6.0*lat_rad)

	// longitude metres scale factor
	lon_sf :=
		coef.E*math.Cos(lat_rad) -
			coef.F*math.Cos(3.0*lat_rad) +
			coef.G*math.Cos(5.0*lat_rad)

	delta_x := math.Sin(head_rad)
	delta_y := math.Cos(head_rad)

	n := len(along_track)
	lon2 := make([]float64, len(along_track))
	lat2 := make([]float64, len(along_track))

	for i := 0; i < n; i++ {
		acr_trck = float64(across_track[i])
		aln_trck = float64(along_track[i])
		lon2[i] = lon + delta_y/lon_sf*acr_trck + delta_x/lon_sf*aln_trck
		lat2[i] = lat - delta_x/lat_sf*acr_trck + delta_y/lat_sf*aln_trck
	}

	lonlat.Longitude = lon2
	lonlat.Latitude = lat2

	return lonlat
}
