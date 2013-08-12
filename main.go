package main

import (
	"fmt"
	"flag"
	"bufio"
	"log"
	"os"
	"os/signal"
	"math/rand"
	"strconv"
	"time"
)


var SAMPLE_TYPE int
const (
	INTEGER = iota
	PERCENTAGE
	)

var COUNT int // a count of how many lines have been collected
var SAMPLE_INTEGER []string
var PERCENT_SAMPLE *PercentageSample
var SAMPLE_VALUE int // either a percentage or a sum to keep

var command = os.Args[0]
var invocation = fmt.Sprintf("%s [[sample size]%%] [file path]\n", command)

var logger *log.Logger

// flag.Usage help message override
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s", invocation)
}

func init() {
	logger = log.New(os.Stderr, "[SNL] ", log.LstdFlags|log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Usage = Usage
	flag.Parse()
}

// parseValue determines if the value is a percentage or
// an integer, and sets the global value `SAMPLE_VALUE`,
func parseValue(s string) {
	var value string
	if string(s[len(s) - 1]) == "%" {
		SAMPLE_TYPE = PERCENTAGE
		value = s[:len(s)-1]
	} else {
		SAMPLE_TYPE = INTEGER
		value = s
	}

	// convert value to integer
	intValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Printf("[Error] error converting sample_size: %s to integer: %s", value, err)
		fmt.Printf("Usage: %s", invocation)
		os.Exit(1)
	}
	SAMPLE_VALUE = intValue
}

// parseFile validates a string and returns an *os.File
func parseFile(s string) (file *os.File) {
	if s == "" {
		logger.Print("[Error] missing filename argument")
		fmt.Printf("Usage: %s", invocation)
		os.Exit(1)
	}

	file, err := os.Open(s)
	if err != nil {
		logger.Fatalf("[Error] error opening %s: %s", s, err)
	}

	return file
}

// forgetOrReplace will choose a number, N, between
// 0 and count, if N is >= threshold we return the sample;
// if N is < len(`sample`) we replace it with `value`
func forgetOrReplace(sample []string, count, threshold int, value string) {
	var candidate = rand.Intn(count)

	if candidate < threshold {
		sample[candidate] = value
	}
}

type PercentageSample struct {
	sample []string // actual sample from all lines seen
	percentageKeep int // the percentage of all samples to keep
	well []string // the maximum size of the elements to take samples from
	wellSize int
	wellSeen int // the total number of new lines in well
	keep int
}

// implements the  "Algorithm 235: Random permutation" by Richard Durstenfeld.
// http://en.wikipedia.org/wiki/Fisher-Yates_shuffle#The_modern_algorithm
func (percentSample *PercentageSample) shuffleAlgorithm235() {
	var choice int
	var old string
	for i := percentSample.wellSeen - 1; i > 1; i-- {
		choice = rand.Intn(i)
		old = percentSample.well[i]
		percentSample.well[i] = percentSample.well[choice]
		percentSample.well[choice] = old
	}
}

// add number of shuffled samples from the well to the sample.
func (sample *PercentageSample) addPercentageToTotal() {
	sample.shuffleAlgorithm235()
	sample.keep = int((float64(sample.percentageKeep) / 100.0) * float64(sample.wellSeen))
	for i := 0; i < sample.keep; i++ {
		sample.sample = append(sample.sample, sample.well[i])
	}
}

// sampleLine is a method that incrementally collects a percentage of all
// samples seen.
func (sample *PercentageSample) sampleLine(line string, count int) {
	if count > 0 && count % sample.wellSize == 0 {
		// add samples from well
		sample.addPercentageToTotal()

		// restart sampling
		sample.wellSeen = 0
		sample.well = make([]string, 0)
	}
	sample.well = append(sample.well, line)
	sample.wellSeen++
}

func handleSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<- sigChannel

	printSample()
}

func printSample() {
	if SAMPLE_TYPE == PERCENTAGE {
		PERCENT_SAMPLE.addPercentageToTotal()
		for _, line := range PERCENT_SAMPLE.sample {
			fmt.Println(line)
		}
	} else {
		for _, line := range SAMPLE_INTEGER {
			fmt.Println(line)
		}
	}
}



func main () {
	var file *os.File

	sampleSize := flag.Arg(0)
	parseValue(sampleSize)

	if SAMPLE_TYPE == INTEGER {
		SAMPLE_INTEGER = make([]string, SAMPLE_VALUE)
	} else if SAMPLE_TYPE == PERCENTAGE {
		PERCENT_SAMPLE = &PercentageSample{
			sample: make([]string, 0),
			percentageKeep: SAMPLE_VALUE,
			well:  make([]string, 0),
			wellSize: 100,
			wellSeen: 0,
		}
	}

	fileName := flag.Arg(1)
	if fileName == "-" {
		file = os.Stdin
	} else {
		file = parseFile(fileName)
		defer file.Close()
	}

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprint(scanner.Text())

		if SAMPLE_TYPE == PERCENTAGE {
			PERCENT_SAMPLE.sampleLine(line, COUNT)
		} else {
			if COUNT < SAMPLE_VALUE {
				SAMPLE_INTEGER[COUNT] = line
			} else {
				forgetOrReplace(SAMPLE_INTEGER, COUNT, SAMPLE_VALUE, line)
			}
		}
		COUNT++
	}
	printSample()
}
