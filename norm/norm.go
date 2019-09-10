// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"math"

	"github.com/chewxy/math32"
	"github.com/goki/ki/kit"
)

///////////////////////////////////////////
//  DivNorm

// DivNorm32 does divisive normalization by given norm function
// i.e., it divides each element by the norm value computed from nfunc.
func DivNorm32(a []float32, nfunc Func32) {
	nv := nfunc(a)
	if nv != 0 {
		MultVec32(a, 1/nv)
	}
}

// DivNorm64 does divisive normalization by given norm function
// i.e., it divides each element by the norm value computed from nfunc.
func DivNorm64(a []float64, nfunc Func64) {
	nv := nfunc(a)
	if nv != 0 {
		MultVec64(a, 1/nv)
	}
}

///////////////////////////////////////////
//  SubNorm

// SubNorm32 does subtractive normalization by given norm function
// i.e., it subtracts norm computed by given function from each element.
func SubNorm32(a []float32, nfunc Func32) {
	nv := nfunc(a)
	AddVec32(a, -nv)
}

// SubNorm64 does subtractive normalization by given norm function
// i.e., it subtracts norm computed by given function from each element.
func SubNorm64(a []float64, nfunc Func64) {
	nv := nfunc(a)
	AddVec64(a, -nv)
}

///////////////////////////////////////////
//  ZScore

// ZScore32 subtracts the mean and divides by the standard deviation
func ZScore32(a []float32) {
	DivNorm32(a, Mean32)
	SubNorm32(a, Std32)
}

// ZScore64 subtracts the mean and divides by the standard deviation
func ZScore64(a []float64) {
	DivNorm64(a, Mean64)
	SubNorm64(a, Std64)
}

///////////////////////////////////////////
//  MultVec

// MultVec32 multiplies vector elements by scalar
func MultVec32(a []float32, val float32) {
	for i, av := range a {
		if math32.IsNaN(av) {
			continue
		}
		a[i] *= val
	}
}

// MultVec64 multiplies vector elements by scalar
func MultVec64(a []float64, val float64) {
	for i, av := range a {
		if math.IsNaN(av) {
			continue
		}
		a[i] *= val
	}
}

///////////////////////////////////////////
//  AddVec

// AddVec32 adds scalar to vector
func AddVec32(a []float32, val float32) {
	for i, av := range a {
		if math32.IsNaN(av) {
			continue
		}
		a[i] += val
	}
}

// AddVec64 adds scalar to vector
func AddVec64(a []float64, val float64) {
	for i, av := range a {
		if math.IsNaN(av) {
			continue
		}
		a[i] += val
	}
}

// Func32 is a norm function operating on slice of float32 numbers
type Func32 func(a []float32) float32

// Func64 is a norm function operating on slices of float64 numbers
type Func64 func(a []float64) float64

// StdNorms are standard norm functions, including stats
type StdNorms int

const (
	L1 StdNorms = iota
	L2
	SumSquares
	N
	Sum
	Mean
	Var
	Std
	Max
	MaxAbs

	StdNormsN
)

//go:generate stringer -type=StdNorms

var KiT_StdNorms = kit.Enums.AddEnum(StdNormsN, false, nil)

func (ev StdNorms) MarshalJSON() ([]byte, error)  { return kit.EnumMarshalJSON(ev) }
func (ev *StdNorms) UnmarshalJSON(b []byte) error { return kit.EnumUnmarshalJSON(ev, b) }

// StdFunc32 returns a standard norm function as specified
func StdFunc32(std StdNorms) Func32 {
	switch std {
	case L1:
		return L132
	case L2:
		return L232
	case SumSquares:
		return SumSquares32
	case N:
		return N32
	case Sum:
		return Sum32
	case Mean:
		return Mean32
	case Var:
		return Var32
	case Std:
		return Std32
	case Max:
		return Max32
	case MaxAbs:
		return MaxAbs32
	}
	return nil
}

// StdFunc64 returns a standard norm function as specified
func StdFunc64(std StdNorms) Func64 {
	switch std {
	case L1:
		return L164
	case L2:
		return L264
	case SumSquares:
		return SumSquares64
	case N:
		return N64
	case Sum:
		return Sum64
	case Mean:
		return Mean64
	case Var:
		return Var64
	case Std:
		return Std64
	case Max:
		return Max64
	case MaxAbs:
		return MaxAbs64
	}
	return nil
}
