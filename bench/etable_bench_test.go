package benchtable

import (
	"testing"

	"github.com/emer/etable/etensor"
)

func sum(idx int, agg float64, val float64) float64 {
	return agg + val
}

func BenchmarkAggSmall(b *testing.B) {
	AggBenchmark(b, SmallMat)
}

func BenchmarkAggMed(b *testing.B) {
	AggBenchmark(b, MediumMat)
}

func BenchmarkAggLarge(b *testing.B) {
	AggBenchmark(b, LargeMat)
}

func AggBenchmark(b *testing.B, size int) {
	x := etensor.NewFloat32([]int{size}, nil, nil)
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		x.Agg(0.0, sum)
	}
}