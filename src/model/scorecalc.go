package model

import (
	"math"
)

// score gives a cosine similarity between two vecs.
// https://en.wikipedia.org/wiki/Cosine_similarity
func score(vec1, vec2 []float64) float64 {
	norm1, norm2 := norm(vec1), norm(vec2)

	// # Avoid div by 0
	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	dot := 0.0
	for i := 0; i < len(vec1); i++ {
		dot += vec1[i] * vec2[i]
	}

	return dot / norm1 / norm2
}

// Norm: mathworld.wolfram.com/Norm.html
func norm(vec []float64) float64 {
	x := 0.0
	l := len(vec)
	for i := 0; i < l; i++ {
		x += vec[i] * vec[i]
	}
	return math.Sqrt(x)
}
