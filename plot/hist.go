package plot

import (
	"fmt"
	"math"
	"os"
)

type Bound struct {
	Left  float64
	Right float64
}

type Bin struct {
	Bound
	Count int
}

// extend histogram flags from bar's flags
var histFlags = barFlags
var nBin int
var leftBound float64
var rightBound float64

func Histogram(args []string) error {
	// create histogram specific flag here because we reused bar plot flags for hist
	histFlags.IntVar(&nBin, "bin", 10, "number of bins in histogram")
	histFlags.Float64Var(&leftBound, "left", math.NaN(), "left bound of the histogram, default is min value")
	histFlags.Float64Var(&rightBound, "right", math.NaN(), "right bound of the histogram, default is max value")
	// replace the Usage method of histFlags since we cannot change its name from "bar"
	histFlags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of hist:")
		histFlags.PrintDefaults()
	}
	histFlags.Parse(args)

	scanner := inputScanner(barFlags.Args())
	values, err := readValues(scanner)
	if err != nil {
		return err
	}

	// left, right bounds is min, max value by default, which is indicated by NaN
	if math.IsNaN(leftBound) || math.IsNaN(rightBound) {
		bound := getBounds(values)
		if math.IsNaN(leftBound) {
			leftBound = bound.Left
		}
		if math.IsNaN(rightBound) {
			rightBound = bound.Right
		}
	}

	bins := groupValuesToBins(values, nBin, Bound{leftBound, rightBound})
	drawBins(bins)
	return nil
}

func getBounds(values []float64) Bound {
	// let's just return [0, 0] as default bound for an empty slice
	if len(values) == 0 {
		return Bound{0, 0}
	}

	min := values[0]
	max := values[0]
	for _, val := range values[1:] {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return Bound{min, max}
}

func groupValuesToBins(values []float64, nBin int, bound Bound) []Bin {
	binSize := (bound.Right - bound.Left) / float64(nBin)
	bins := make([]Bin, nBin)

	// label the bin by the upper/right bound
	for bin := 0; bin < nBin; bin++ {
		leftBound := bound.Left + float64(bin)*binSize
		rightBound := leftBound + binSize
		bins[bin].Bound = Bound{leftBound, rightBound}
	}

	for _, val := range values {
		switch {
		case val < bound.Left, val > bound.Right:
			continue
		case val == bound.Right:
			bins[nBin-1].Count++
		default:
			bins[int((val-bound.Left)/binSize)].Count++
		}
	}
	return bins
}

func drawBins(bins []Bin) {
	bars := make([]LabeledValue, len(bins))
	for i, bin := range bins {
		bars[i].Label = fmt.Sprintf("%.2f", bin.Right)
		bars[i].Value = float64(bin.Count)
	}
	if len(bins) > 1 {
		bars[0].Label = fmt.Sprintf("%.2f -> %.2f", bins[0].Left, bins[0].Right)
	}
	DrawBars(bars)
}
