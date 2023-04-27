package wavelet

import (
	"math"
	"testing"
)

func TestWavelet(t *testing.T) {
	t.Run("simple", func(tt *testing.T) {
		low, high := Haar([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 8.0})
		tt.Logf("log=%v", low)
		tt.Logf("high=%v", high)

		out := InverseHaar(low, high)
		tt.Logf("out=%v", out)

		for i, v := range out {
			out[i] = math.Ceil(v)
		}
		tt.Logf("out=%v", out)

		if out[0] != 1.0 {
			tt.Errorf("expect=1.0, actual=%v", out[0])
		}
		if out[1] != 2.0 {
			tt.Errorf("expect=2.0, actual=%v", out[1])
		}
		if out[2] != 3.0 {
			tt.Errorf("expect=3.0, actual=%v", out[2])
		}
		if out[3] != 4.0 {
			tt.Errorf("expect=3.0, actual=%v", out[3])
		}
		if out[4] != 5.0 {
			tt.Errorf("expect=5.0, actual=%v", out[4])
		}
		if out[5] != 8.0 {
			tt.Errorf("expect=8.0, actual=%v", out[5])
		}
	})
}

func TestCompare(t *testing.T) {
	out := Compare(
		[]float64{1, 1, 1, 0, 1, 1, 2, 2, 2, 0},
		[]float64{1, 1, 1, 8, 1, 1, 2, 2, 2, 0},
	)
	t.Logf("compare=%v", out)
	if out[0] != 0 {
		t.Errorf("no diff")
	}
	if out[1] == 0 {
		t.Errorf("no diff")
	}
	if out[2] != 0 {
		t.Errorf("no diff")
	}
	if out[3] != 0 {
		t.Errorf("has diff")
	}
	if out[4] != 0 {
		t.Errorf("no diff")
	}
}
