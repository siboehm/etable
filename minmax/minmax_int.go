// Copyright (c) 2019, The eTable Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package minmax

import "math"

// Int represents a min / max range for int values.
// Supports clipping, renormalizing, etc
type Int struct {
	Min int
	Max int
}

// Set sets the min and max values
func (mr *Int) Set(min, max int) {
	mr.Min, mr.Max = min, max
}

// SetInfinity sets the Min to +MaxFloat, Max to -MaxFloat -- suitable for
// iteratively calling Fit*InRange
func (mr *Int) SetInfinity() {
	mr.Min, mr.Max = math.MaxInt64, -math.MaxInt64
}

// IsValid returns true if Min <= Max
func (mr *Int) IsValid() bool {
	return mr.Min <= mr.Max
}

// InRange tests whether value is within the range (>= Min and <= Max)
func (mr *Int) InRange(val int) bool {
	return ((val >= mr.Min) && (val <= mr.Max))
}

// IsLow tests whether value is lower than the minimum
func (mr *Int) IsLow(val int) bool {
	return (val < mr.Min)
}

// IsHigh tests whether value is higher than the maximum
func (mr *Int) IsHigh(val int) bool {
	return (val > mr.Min)
}

// Range returns Max - Min
func (mr *Int) Range() int {
	return mr.Max - mr.Min
}

// Scale returns 1 / Range -- if Range = 0 then returns 0
func (mr *Int) Scale() float32 {
	r := mr.Range()
	if r != 0 {
		return 1 / float32(r)
	}
	return 0
}

// Midpoint returns point halfway between Min and Max
func (mr *Int) Midpoint() float32 {
	return 0.5 * float32(mr.Max+mr.Min)
}

// FitInRange adjusts our Min, Max to fit within those of other Int
// returns true if we had to adjust to fit.
func (mr *Int) FitInRange(oth Int) bool {
	adj := false
	if oth.Min < mr.Min {
		mr.Min = oth.Min
		adj = true
	}
	if oth.Max > mr.Max {
		mr.Max = oth.Max
		adj = true
	}
	return adj
}

// FitValInRange adjusts our Min, Max to fit given value within Min, Max range
// returns true if we had to adjust to fit.
func (mr *Int) FitValInRange(val int) bool {
	adj := false
	if val < mr.Min {
		mr.Min = val
		adj = true
	}
	if val > mr.Max {
		mr.Max = val
		adj = true
	}
	return adj
}

// NormVal normalizes value to 0-1 unit range relative to current Min / Max range
// Clips the value within Min-Max range first.
func (mr *Int) NormVal(val int) float32 {
	return float32(mr.ClipVal(val)-mr.Min) * mr.Scale()
}

// ProjVal projects a 0-1 normalized unit value into current Min / Max range (inverse of NormVal)
func (mr *Int) ProjVal(val float32) float32 {
	return float32(mr.Min) + (val * float32(mr.Range()))
}

// ClipVal clips given value within Min / Max rangee
func (mr *Int) ClipVal(val int) int {
	if val < mr.Min {
		return mr.Min
	}
	if val > mr.Max {
		return mr.Max
	}
	return val
}

// ClipNormVal clips then normalizes given value within 0-1
func (mr *Int) ClipNormVal(val int) float32 {
	if val < mr.Min {
		return 0
	}
	if val > mr.Max {
		return 1
	}
	return mr.NormVal(val)
}
