package percent_sample

import (
	"fmt"
	"math/rand"
)

// Sample is an object that is used for sampling a
// percentage of lines.
type Sample struct {
	Sample []string // actual sample from all lines seen
	PercentageKeep int // the percentage of all samples to keep
	Well []string // the maximum size of the elements to take samples from
	WellSize uint64
	WellSeen int64 // the total number of new lines in well
	keep uint64
	count uint64
}

func (sample *Sample) Print() {
	for _, line := range sample.Sample {
		fmt.Println(line)
	}
}

// Shuffle235 randomly shuffles an array in place using:
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle#The_modern_algorithm
func Shuffle235(well []string, count int64) []string {
	var choice int64
	var old string

	for i := count - 1; i > 1; i-- {
		choice = rand.Int63n(i)
		old = well[i]
		well[i] = well[choice]
		well[choice] = old
	}
	return well
}

// implements the  "Algorithm 235: Random permutation" by Richard Durstenfeld.
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle#The_modern_algorithm
func (sample *Sample) shuffleAlgorithm235() {
	sample.Well = Shuffle235(sample.Well, sample.WellSeen)
}

// add number of shuffled samples from the well to the sample.
func (sample *Sample) AddPercentageToTotal() {
	sample.shuffleAlgorithm235()
	sample.keep = uint64((float64(sample.PercentageKeep) / 100.0) * float64(sample.WellSeen))
	var i uint64 = 0
	for i = 0; i < sample.keep; i++ {
		sample.Sample = append(sample.Sample, sample.Well[i])
	}
}

// sampleLine is a method that incrementally collects a percentage of all
// samples seen.
func (sample *Sample) SampleLine(line string) {
	if sample.count > 0 && sample.count % sample.WellSize == 0 {
		// add samples from well
		sample.AddPercentageToTotal()

		// restart sampling
		sample.WellSeen = 0
		sample.Well = make([]string, 0)
	}
	sample.Well = append(sample.Well, line)
	sample.WellSeen++
	sample.count++
}
