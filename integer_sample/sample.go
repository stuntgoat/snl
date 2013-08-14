package integer_sample

import (
	"fmt"
	"math/rand"
)


type IntegerSample struct {
	Sample []string
	Size int
	candidate int
	count int
}

// Print prints the values of the sample to STDOUT.
func (sample *IntegerSample) Print() {
	for _, line := range sample.Sample {
		fmt.Println(line)
	}
}

// SampleLine performs a classic 'resevoir sample'.
// Choose a number, N, between 0 and , if N is >= sample.Size
// we return; if N is < sample.Size we replace it with `line` at
// the Nth index of sample.Sample.
func (sample *IntegerSample) SampleLine(line string) {

	if sample.count < sample.Size {
		goto ADD_SAMPLE
	}

	sample.candidate = rand.Intn(sample.count)
	if sample.candidate < sample.Size {
		sample.Sample[sample.candidate] = line
		goto INCREMENT_COUNT
	}

	return

ADD_SAMPLE:
	sample.Sample[sample.count] = line
	goto INCREMENT_COUNT

INCREMENT_COUNT:
	sample.count++

	return
}
