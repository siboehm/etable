// Copyright (c) 2019, The eTable Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eplot

import (
	"fmt"
	"log"
	"math"

	"github.com/emer/etable/etable"
	"github.com/emer/etable/minmax"
	"github.com/emer/etable/split"
	"github.com/goki/gi/gist"
	"github.com/goki/ki/ints"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// bar plot is on integer positions, with different Y values and / or
// legend values interleaved

// GenPlotBar generates a Bar plot, setting GPlot variable
func (pl *Plot2D) GenPlotBar() {
	plt := plot.New() // note: not clear how to re-use, due to newtablexynames
	plt.Title.Text = pl.Params.Title
	plt.X.Label.Text = pl.XLabel()
	plt.Y.Label.Text = pl.YLabel()
	plt.BackgroundColor = nil

	if pl.Params.BarWidth > 1 {
		pl.Params.BarWidth = .8
	}

	// process xaxis first
	xi, xview, _, err := pl.PlotXAxis(plt, pl.Table)
	if err != nil {
		return
	}
	xp := pl.Cols[xi]

	var lsplit *etable.Splits
	nleg := 1
	if pl.Params.LegendCol != "" {
		_, err = pl.Table.Table.ColIdxTry(pl.Params.LegendCol)
		if err != nil {
			log.Println("eplot.LegendCol: " + err.Error())
		} else {
			xview.SortColNames([]string{pl.Params.LegendCol, xp.Col}, etable.Ascending) // make it fit!
			lsplit = split.GroupBy(xview, []string{pl.Params.LegendCol})
			nleg = ints.MaxInt(lsplit.Len(), 1)
		}
	}

	var firstXY *TableXY
	var strCols []*ColParams
	nys := 0
	for _, cp := range pl.Cols {
		cp.UpdateVals()
		if !cp.On {
			continue
		}
		if cp.IsString {
			strCols = append(strCols, cp)
			continue
		}
		if cp.TensorIdx < 0 {
			yc := pl.Table.Table.ColByName(cp.Col)
			_, sz := yc.RowCellSize()
			nys += sz
		} else {
			nys++
		}
		if cp.Range.FixMin {
			plt.Y.Min = math.Min(plt.Y.Min, cp.Range.Min)
		}
		if cp.Range.FixMax {
			plt.Y.Max = math.Max(plt.Y.Max, cp.Range.Max)
		}
	}

	if nys == 0 {
		return
	}

	stride := nys * nleg
	if stride > 1 {
		stride += 1 // extra gap
	}

	yoff := 0
	yidx := 0
	maxx := 0 // max number of x values
	for _, cp := range pl.Cols {
		if !cp.On || cp == xp {
			continue
		}
		if cp.IsString {
			continue
		}
		start := yoff
		for li := 0; li < nleg; li++ {
			lview := xview
			leg := ""
			if lsplit != nil && len(lsplit.Values) > li {
				leg = lsplit.Values[li][0]
				lview = lsplit.Splits[li]
			}
			nidx := 1
			stidx := cp.TensorIdx
			if cp.TensorIdx < 0 { // do all
				yc := pl.Table.Table.ColByName(cp.Col)
				_, sz := yc.RowCellSize()
				nidx = sz
				stidx = 0
			}
			for ii := 0; ii < nidx; ii++ {
				idx := stidx + ii
				xy, _ := NewTableXYName(lview, xi, xp.TensorIdx, cp.Col, idx, cp.Range)
				if xy == nil {
					continue
				}
				maxx = ints.MaxInt(maxx, lview.Len())
				if firstXY == nil {
					firstXY = xy
				}
				lbl := cp.Label()
				clr := cp.Color
				if leg != "" {
					lbl = leg + " " + lbl
				}
				if nleg > 1 {
					cidx := yidx*nleg + li
					clr, _ = gist.ColorFromString(PlotColorNames[cidx%len(PlotColorNames)], nil)
				}
				if nidx > 1 {
					clr, _ = gist.ColorFromString(PlotColorNames[idx%len(PlotColorNames)], nil)
					lbl = fmt.Sprintf("%s_%02d", lbl, idx)
				}
				ec := -1
				if cp.ErrCol != "" {
					ec = pl.Table.Table.ColIdx(cp.ErrCol)
				}
				var bar *ErrBarChart
				if ec >= 0 {
					exy, _ := NewTableXY(lview, ec, 0, ec, 0, minmax.Range64{})
					bar, err = NewErrBarChart(xy, exy)
					if err != nil {
						log.Println(err)
						continue
					}
				} else {
					bar, err = NewErrBarChart(xy, nil)
					if err != nil {
						log.Println(err)
						continue
					}
				}
				bar.Color = clr
				bar.Stride = float64(stride)
				bar.Start = float64(start)
				bar.Width = pl.Params.BarWidth
				plt.Add(bar)
				plt.Legend.Add(lbl, bar)
				start++
			}
		}
		yidx++
		yoff += nleg
	}
	mid := (stride - 1) / 2
	if stride > 1 {
		mid = (stride - 2) / 2
	}
	if firstXY != nil && len(strCols) > 0 {
		firstXY.Table = xview
		n := xview.Len()
		for _, cp := range strCols {
			xy, _ := NewTableXYName(xview, xi, xp.TensorIdx, cp.Col, cp.TensorIdx, firstXY.YRange)
			xy.LblCol = xy.YCol
			xy.YCol = firstXY.YCol
			xy.YIdx = firstXY.YIdx

			xyl := plotter.XYLabels{}
			xyl.XYs = make(plotter.XYs, n)
			xyl.Labels = make([]string, n)

			for i := range xview.Idxs {
				y := firstXY.Value(i)
				x := float64(mid + (i%maxx)*stride)
				xyl.XYs[i] = plotter.XY{x, y}
				xyl.Labels[i] = xy.Label(i)
			}
			lbls, _ := plotter.NewLabels(xyl)
			if lbls != nil {
				plt.Add(lbls)
			}
		}
	}

	netn := pl.Table.Len() * stride
	xc := pl.Table.Table.Cols[xi]
	vals := make([]string, netn)
	for i, dx := range pl.Table.Idxs {
		pi := mid + i*stride
		if pi < netn && dx < xc.Len() {
			vals[pi] = xc.StringVal1D(dx)
		}
	}
	plt.NominalX(vals...)

	plt.Legend.Top = true
	plt.X.Tick.Label.Rotation = math.Pi * (pl.Params.XAxisRot / 180)
	if pl.Params.XAxisRot > 10 {
		plt.X.Tick.Label.YAlign = draw.YCenter
		plt.X.Tick.Label.XAlign = draw.XRight
	}
	pl.GPlot = plt
}
