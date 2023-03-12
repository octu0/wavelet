package wavelet

import (
	"math"
)

type Float interface {
	float32 | float64
}

func Wavelet[T Float](signal []T) (high, low []T) {
	N := len(signal)
	high, low = make([]T, N/2), make([]T, N/2)

	for i := 0; i < N/2; i += 1 {
		high[i] = (signal[2*i] + signal[2*i+1]) / math.Sqrt2
		low[i] = (signal[2*i] - signal[2*i+1]) / math.Sqrt2
	}
	return high, low
}

func WaveletClamp[T Float](signal []T, min, max T) ([]T, []T) {
	high, low := Wavelet(signal)
	for i, v := range high {
		high[i] = Clamp(v, min, max)
	}
	return high, low
}

func Compare[T Float](a, b []T) (low []T) {
	N := len(a)
	low = make([]T, N/2)

	for i := 0; i < N/2; i += 1 {
		d := (a[2*i] - b[2*i+1]) / 2
		if 0 < d {
			d = 0
		}
		low[i] = d
	}
	return low
}

func Inverse[T Float](high, low []T) []T {
	N := len(high) * 2
	out := make([]T, N)

	for i := 0; i < N/2; i += 1 {
		out[2*i] = (high[i] + low[i]) / math.Sqrt2
		out[2*i+1] = (high[i] - low[i]) / math.Sqrt2
	}
	return out
}

func Clamp[T Float](data, min, max T) T {
	if data < min {
		return min
	}
	if max < data {
		return max
	}
	return data
}
