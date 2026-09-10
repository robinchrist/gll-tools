package sofaexport

import (
	"math"
	"testing"

	"github.com/cwbudde/gll-tools/pkg/gll"
)

func TestComplexFromLevelPhase(t *testing.T) {
	tests := []struct {
		levelDB, phaseRad  float64
		wantReal, wantImag float64
	}{
		{0, 0, 1, 0},           // 0 dB, 0 rad → 1+0j
		{0, math.Pi / 2, 0, 1}, // 0 dB, 90° → 0+1j
		{20, 0, 10, 0},         // +20 dB → 10+0j
		{-20, math.Pi, -0.1, 1.224646799147353e-17}, // -20 dB, 180° ≈ -0.1
	}
	for _, tt := range tests {
		gotR, gotI := complexFromLevelPhase(tt.levelDB, tt.phaseRad)
		if math.Abs(gotR-tt.wantReal) > 1e-12 || math.Abs(gotI-tt.wantImag) > 1e-12 {
			t.Errorf("complexFromLevelPhase(%g,%g) = (%g,%g), want (%g,%g)",
				tt.levelDB, tt.phaseRad, gotR, gotI, tt.wantReal, tt.wantImag)
		}
	}
}

func TestCombineResponse_Relative(t *testing.T) {
	def := gll.LogSpectrumDefinition{BandsPerOctave: 1, StartFreq: 100, PointCount: 3}
	resp := &gll.TransferFunction{
		Definition: def,
		Level:      []float64{0, -3, -6},
		Phase:      []float64{0, 0, 0},
	}
	src := &gll.SourceDefinition{
		OnAxisSpectrum: &gll.TransferFunction{
			Definition: def,
			Level:      []float64{96, 95, 94},
			Phase:      []float64{0, 0, 0},
		},
		OnAxisLevel: 0,
	}

	reArr, imArr, err := combineResponse(resp, src, true)
	if err != nil {
		t.Fatalf("combineResponse: %v", err)
	}
	// Expected: just 10^(level/20) for each band, imag=0
	want := []float64{1.0, math.Pow(10, -3.0/20), math.Pow(10, -6.0/20)}
	for i := range want {
		if math.Abs(reArr[i]-want[i]) > 1e-9 {
			t.Errorf("real[%d] = %g, want %g", i, reArr[i], want[i])
		}
		if math.Abs(imArr[i]) > 1e-12 {
			t.Errorf("imag[%d] = %g, want 0", i, imArr[i])
		}
	}
}

func TestCombineResponse_Combined(t *testing.T) {
	def := gll.LogSpectrumDefinition{BandsPerOctave: 1, StartFreq: 100, PointCount: 2}
	resp := &gll.TransferFunction{
		Definition: def,
		Level:      []float64{-3, -6},
		Phase:      []float64{0, 0},
	}
	src := &gll.SourceDefinition{
		OnAxisSpectrum: &gll.TransferFunction{
			Definition: def,
			Level:      []float64{90, 92},
			Phase:      []float64{0, 0},
		},
		OnAxisLevel: 0,
	}

	reArr, imArr, err := combineResponse(resp, src, false)
	if err != nil {
		t.Fatalf("combineResponse: %v", err)
	}
	// Combined level = resp + onAxis: 87, 86. Real = 10^(L/20).
	want := []float64{math.Pow(10, 87.0/20), math.Pow(10, 86.0/20)}
	for i := range want {
		if math.Abs(reArr[i]-want[i]) > 1e-3 {
			t.Errorf("real[%d] = %g, want %g", i, reArr[i], want[i])
		}
		if math.Abs(imArr[i]) > 1e-9 {
			t.Errorf("imag[%d] = %g, want 0", i, imArr[i])
		}
	}
}

func TestCombineResponse_GridMismatch(t *testing.T) {
	resp := &gll.TransferFunction{
		Definition: gll.LogSpectrumDefinition{BandsPerOctave: 1, StartFreq: 100, PointCount: 2},
		Level:      []float64{0, 0},
		Phase:      []float64{0, 0},
	}
	src := &gll.SourceDefinition{
		OnAxisSpectrum: &gll.TransferFunction{
			Definition: gll.LogSpectrumDefinition{BandsPerOctave: 3, StartFreq: 100, PointCount: 2},
			Level:      []float64{0, 0},
			Phase:      []float64{0, 0},
		},
	}
	if _, _, err := combineResponse(resp, src, false); err == nil {
		t.Fatal("expected grid-mismatch error, got nil")
	}
}

// A 5 ms differential delay makes equal 100 Hz sources cancel. Losing the
// separate delay fields changes destructive interference to reinforcement.
func TestCombineResponse_StoredDelaysAndReferenceLevel(t *testing.T) {
	def := gll.LogSpectrumDefinition{BandsPerOctave: 1, StartFreq: 100, PointCount: 1}
	resp := &gll.TransferFunction{Definition: def, Level: []float64{0}, Phase: []float64{0}, Delay: .002}
	src := &gll.SourceDefinition{OnAxisLevel: 94, OnAxisSpectrum: &gll.TransferFunction{
		Definition: def, Level: []float64{94}, Phase: []float64{0}, Delay: .003,
	}}
	re, im, err := combineResponse(resp, src, false)
	if err != nil {
		t.Fatal(err)
	}
	amplitude := math.Pow(10, 94.0/20)
	if math.Abs(re[0]+amplitude) > 1e-8 || math.Abs(im[0]) > 1e-8 {
		t.Fatalf("delayed sum does not cancel reference: %g + %gi", re[0]+amplitude, im[0])
	}
	re, im, err = combineResponse(resp, src, true)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(re[0]-math.Cos(-.4*math.Pi)) > 1e-12 || math.Abs(im[0]-math.Sin(-.4*math.Pi)) > 1e-12 {
		t.Fatalf("relative mode lost balloon delay: %g + %gi", re[0], im[0])
	}
}
