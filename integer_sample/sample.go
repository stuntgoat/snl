package integer_sample

import (
	"fmt"
	"math/rand"
)

type Sample struct {
	Sample []string
	Size int
	candidate int
	count int
}

// Print prints the values of the sample to STDOUT.
func (sample *Sample) Print() {
	for _, line := range sample.Sample {
		fmt.Println(line)
	}
}


// SampleLine performs a classic 'resevoir sample`.
// Choose a number, N, between 0 and , if N is >=  we return;
// if N is < len(`sample`) we replace it with `value`
func (sample *Sample) SampleLine(line string) {

	if sample.count < sample.Size {
		goto ADD_SAMPLE
	}

	sample.candidate = rand.Intn(sample.count)
	if sample.candidate < sample.Size {

		sample.Sample[sample.candidate] = line
		goto INCREMENT_COUNT
	}

	goto INCREMENT_COUNT

ADD_SAMPLE:
	sample.Sample[sample.count] = line
	goto INCREMENT_COUNT

INCREMENT_COUNT:
	sample.count++

	return
}
