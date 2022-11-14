"""
This is the Python version of gonum_bench_test.go

Results will depend a lot on which BLAS backend you're using for Numpy (MKL / OpenBLAS / Accelerate etc.)
To install a given backend, use e.g `mamba install numpy "blas=*=openblas"`
"""

import numpy as np
import time

def test_benchmark_matmul(n):
    A = np.random.rand(n, n).astype(np.float64)
    B = np.random.rand(n, n).astype(np.float64)

    times = []
    for _ in range(10):
        start_time = time.time_ns()
        C = A @ B
        times.append(time.time_ns() - start_time)
        assert C.shape == A.shape

    print(f"Size {n}: {sum(times) / len(times)}ns")

np.show_config()
print()

for size in [10, 100, 1000, 10000]:
    test_benchmark_matmul(size)
