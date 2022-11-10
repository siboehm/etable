package benchtable

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
)

func BenchmarkMatmulSmall(b *testing.B) {
	MatmulBenchmark(b, SmallMat)
}

func BenchmarkMatmulMed(b *testing.B) {
	MatmulBenchmark(b, MediumMat)
}

func BenchmarkMatmulLarge(b *testing.B) {
	MatmulBenchmark(b, LargeMat)
}

func BenchmarkMatmulHuge(b *testing.B) {
	MatmulBenchmark(b, HugeMat)
}

// Convert this function to Python:
func MatmulBenchmark(b *testing.B, dim int) {
	A := mat.NewDense(dim, dim, nil)
	B := mat.NewDense(dim, dim, nil)
	C := mat.NewDense(dim, dim, nil)


	// Fill A and B with random numbers.
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			A.Set(i, j, rand.Float64())
			B.Set(i, j, rand.Float64())
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C.Mul(A, B)
	}
}



