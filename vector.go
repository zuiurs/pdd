package main

import (
	"fmt"
	"math"
)

type Vector []float64

func HistToVector(h Hist) Vector {
	vector := make(Vector, len(h))
	for i, v := range h {
		vector[i] = float64(v)
	}

	return vector
}

func Normalize(vector Vector) Vector {
	retv := make(Vector, len(vector))

	var norm float64
	for _, v := range vector {
		norm += float64(v * v)
	}
	norm = math.Sqrt(norm)

	for i, v := range vector {
		retv[i] = v / norm
	}

	return retv
}

func Distance(v1, v2 Vector) (float64, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("Both vector's dimension is not match.")
	}

	d := len(v1)
	var diff float64
	for i := 0; i < d; i++ {
		diff += math.Pow(v2[i]-v1[i], 2.0)
	}

	distance := math.Sqrt(diff)
	return distance, nil
}
