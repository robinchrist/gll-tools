package sofaexport

import (
	"math"

	sofa "github.com/cwbudde/go-sofa"
)

// gridAngles returns the (azimuthDeg, elevationDeg) pair for a given
// (merIdx, parIdx) pair in the GLL angular grid.
//
// GLL parallel is the polar angle from the firing axis (+X): 0 is front,
// 180 is rear. Meridian rotates around that axis, from +Z toward +Y.
// Convert that direction to SOFA azimuth/elevation; parallel is not elevation.
func gridAngles(merIdx, parIdx int, merStep, parStep float64) (azDeg, elDeg float64) {
	mer := float64(merIdx) * merStep * math.Pi / 180
	par := float64(parIdx) * parStep * math.Pi / 180
	x := math.Cos(par)
	y := math.Sin(par) * math.Sin(mer)
	z := math.Sin(par) * math.Cos(mer)
	azDeg = math.Atan2(y, x) * 180 / math.Pi
	elDeg = math.Atan2(z, math.Hypot(x, y)) * 180 / math.Pi
	return azDeg, elDeg
}

// directionToCartesian converts an (azimuth, elevation, radius) triple in
// degrees/metres to a SOFA Vector3 (right-handed cartesian, meters).
func directionToCartesian(azDeg, elDeg, radius float64) sofa.Vector3 {
	az := azDeg * math.Pi / 180.0
	el := elDeg * math.Pi / 180.0
	cosEl := math.Cos(el)
	return sofa.Vector3{
		X: radius * cosEl * math.Cos(az),
		Y: radius * cosEl * math.Sin(az),
		Z: radius * math.Sin(el),
	}
}
