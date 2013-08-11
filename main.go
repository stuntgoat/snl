package main

import (
	"fmt"
	"flag"
	"bufio"
	"log"
	"os"
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
var SAMPLE = make([]string, 0)
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


// keepPercentage returns a function that closes over
// an argument that represents a percentage. The function it returns
// accepts a count and it returns.
func keepPercentage(percentage float64) (fn func(int) int) {
	return func(count int) int {
		return int(percentage * float64(count))
	}
}

// forgetOrReplace will choose a number, N, between
// 0 and count, if N is >= threshold we return the sample;
// if N is < len(`sample`) we replace it with `value`
func forgetOrReplace(sample []string, count, threshold int, value string){
	var candidate = rand.Intn(count)

	if candidate < threshold {
		sample[candidate] = value
	}
}


func printSample() {
	for _, line := range SAMPLE {
		fmt.Println(line)
	}
}


type PercentageSample struct {
	total []string
	percentage float64
}

// add saved values to array
func (container *PercentageSample) addPercentageToTotal(sample []string) {

}


func handleSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<- sigChannel

	// check if this is a percentage or integer sample


	keepPercentage()

	printSample()
}

func main () {

	var file *os.File

	sampleSize := flag.Arg(0)
	parseValue(sampleSize)

	if SAMPLE_TYPE == INTEGER {
		SAMPLE = make([]string, SAMPLE_VALUE)
	} else if SAMPLE_TYPE == PERCENTAGE {
		// create a PercentageSample object and set it's percentage

		// make sample a defualt size of 100


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
			// collect the number of values

			// call collect at each threshold to siphon samples to objec
			logger.Println("percentage sampling not implemented yet")
			os.Exit(0)
		}

		if COUNT < SAMPLE_VALUE {
			SAMPLE[COUNT] = line
		} else {
			forgetOrReplace(SAMPLE, COUNT, SAMPLE_VALUE, line)
		}
		COUNT++
	}
	printSample()
}
