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
	WellSize int
	WellSeen int // the total number of new lines in well
	keep int
	count int
}

func (sample *Sample) Print() {
	for _, line := range sample.Sample {
		fmt.Println(line)
	}
}

// Shuffle235 randomly shuffles an array in place using:
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle#The_modern_algorithm
func Shuffle235(well []string, count int) {
	var choice int
	var old string

	for i := count - 1; i > 1; i-- {
		choice = rand.Intn(i)
		old = well[i]
		well[i] = well[choice]
		well[choice] = old
	}
}

// implements the  "Algorithm 235: Random permutation" by Richard Durstenfeld.
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle#The_modern_algorithm
func (sample *Sample) shuffleAlgorithm235() {
	Shuffle235(sample.Well, sample.WellSeen)
}

// add number of shuffled samples from the well to the sample.
func (sample *Sample) AddPercentageToTotal() {
	sample.shuffleAlgorithm235()
	sample.keep = int((float64(sample.PercentageKeep) / 100.0) * float64(sample.WellSeen))
	for i := 0; i < sample.keep; i++ {
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
