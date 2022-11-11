# Benchmarking etable

## BLAS

etable relies on `gonum`'s BLAS implementation, which is largely implemented in Go.
`gonum` also has some moderately optimized assembly, but only for amd64 ISA.

Instead, we use a proper BLAS library, like OpenBLAS.
For this, follow the instructions in [gonum/netlib](https://github.com/gonum/netlib).
Some more hints in [this thread](https://github.com/gonum/gonum/issues/511).

Then call:
```go
blas64.Use(netlib.Implementation{})
```

In my tests, using OpenBLAS is slightly faster for small (~100x100) matrix GEMM, and ~10x faster for large (10Kx10K) GEMM.
Relying on OpenBLAS makes gonum's GEMM performance comparable to Numpy.