package wavelet

import (
	"math"
	"testing"
)

func TestWavelet(t *testing.T) {
	t.Run("simple", func(tt *testing.T) {
		high, low := Wavelet([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 8.0})
		tt.Logf("high=%v", high)
		tt.Logf("log=%v", low)

		out := Inverse(high, low)
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
	t.Run("clamp", func(tt *testing.T) {
		high, low := WaveletClamp([]float64{1.0, 3.0, 11.0, -2.0, -4.0, -16.0}, -10, 5)
		tt.Logf("high=%v", high)
		tt.Logf("log=%v", low)

		out := Inverse(high, low)
		for i, v := range out {
			out[i] = math.Ceil(v)
		}
		if out[0] != 1.0 {
			tt.Errorf("expect=1.0, actual=%v", out[0])
		}
		if out[1] != 3.0 {
			tt.Errorf("expect=3.0, actual=%v", out[1])
		}
		if out[2] != 11.0 {
			tt.Errorf("expect=11.0, actual=%v", out[2])
		}
		if out[3] != -2.0 {
			tt.Errorf("expect=-2.0, actual=%v", out[3])
		}
		if out[4] != -1.0 {
			tt.Errorf("expect=-1.0, actual=%v", out[4])
		}
		if out[5] != -13.0 {
			tt.Errorf("expect=-13.0, actual=%v", out[5])
		}
		tt.Logf("out=%v", out)
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
