package wavelet

import (
	"math"
)

type Float interface {
	float32 | float64
}

func Haar[T Float](signal []T) (high, low []T) {
	N := len(signal)
	low, high = make([]T, N/2), make([]T, N/2)

	for i := 0; i < N/2; i += 1 {
		low[i] = (signal[2*i] + signal[2*i+1]) / math.Sqrt2
		high[i] = (signal[2*i] - signal[2*i+1]) / math.Sqrt2
	}
	return low, high
}

func InverseHaar[T Float](low, high []T) []T {
	N := len(high) * 2
	out := make([]T, N)

	for i := 0; i < N/2; i += 1 {
		out[2*i] = (low[i] + high[i]) / math.Sqrt2
		out[2*i+1] = (low[i] - high[i]) / math.Sqrt2
	}
	return out
}

func Threshold[T Float](signal []T, ratio T) {
	min, max := 0*ratio, 0*ratio
	for i := 0; i < len(signal); i += 1 {
		if signal[i] < min {
			min = signal[i]
		}
		if max < signal[i] {
			max = signal[i]
		}
	}

	thMin, thMax := (min * ratio), (max * ratio)
	for i := 0; i < len(signal); i += 1 {
		if signal[i] < thMin {
			signal[i] = thMin
		}
		if thMax < signal[i] {
			signal[i] = thMax
		}
	}
}

func Compare[T Float](a, b []T) (high []T) {
	N := len(a)
	high = make([]T, N/2)

	for i := 0; i < N/2; i += 1 {
		d := (a[2*i] - b[2*i+1]) / 2
		if 0 < d {
			d = 0
		}
		high[i] = d
	}
	return high
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
