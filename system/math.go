package systemutil

import "math"

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Max[N1, N2 Number](n1 N1, n2 N2) N1 {
	return N1(math.Max(float64(n1), float64(n2)))
}

func Min[N1, N2 Number](n1 N1, n2 N2) N1 {
	return N1(math.Min(float64(n1), float64(n2)))
}

func Log[N Number](n N) N {
	return N(math.Log(float64(n)))
}

func Log10[N Number](n N) N {
	return N(math.Log10(float64(n)))
}
