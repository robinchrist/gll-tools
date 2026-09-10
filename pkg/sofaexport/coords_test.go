package sofaexport

import (
	"math"
	"testing"
)

func TestDirectionToCartesian(t *testing.T) {
	tests := []struct {
		name      string
		azDeg     float64
		elDeg     float64
		r         float64
		x, y, z   float64
		tolerance float64
	}{
		{"on-axis +X", 0, 0, 1, 1, 0, 0, 1e-12},
		{"east +Y at azimuth 90", 90, 0, 1, 0, 1, 0, 1e-12},
		{"behind -X at azimuth 180", 180, 0, 1, -1, 0, 0, 1e-12},
		{"south pole -Z at elevation -90", 0, -90, 1, 0, 0, -1, 1e-12},
		{"north pole +Z at elevation +90", 0, 90, 1, 0, 0, 1, 1e-12},
		{"45° NE at distance 2 m", 45, 0, 2, math.Sqrt(2), math.Sqrt(2), 0, 1e-12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := directionToCartesian(tt.azDeg, tt.elDeg, tt.r)
			if math.Abs(got.X-tt.x) > tt.tolerance ||
				math.Abs(got.Y-tt.y) > tt.tolerance ||
				math.Abs(got.Z-tt.z) > tt.tolerance {
				t.Errorf("got (%g,%g,%g), want (%g,%g,%g)",
					got.X, got.Y, got.Z, tt.x, tt.y, tt.z)
			}
		})
	}
}

func TestGridAngles(t *testing.T) {
	const merStep, parStep = 5.0, 5.0
	tests := []struct {
		merIdx, parIdx int
		azDeg, elDeg   float64
	}{
		{0, 0, 0, 0},     // front, independent of meridian
		{0, 18, 0, 90},   // top
		{0, 36, 180, 0},  // rear
		{18, 18, 90, 0},  // east equator
		{36, 18, 0, -90}, // bottom
	}
	for _, tt := range tests {
		az, el := gridAngles(tt.merIdx, tt.parIdx, merStep, parStep)
		if math.Abs(el-tt.elDeg) > 1e-9 || (math.Abs(tt.elDeg) != 90 && math.Abs(az-tt.azDeg) > 1e-9) {
			t.Errorf("gridAngles(%d,%d) = (%g,%g), want (%g,%g)",
				tt.merIdx, tt.parIdx, az, el, tt.azDeg, tt.elDeg)
		}
	}
}
