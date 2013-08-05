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

var SAMPLE_MAP map[int]string
var SAMPLE_VALUE int // either a percentage or an integer value

var command = os.Args[0]
var invocation = fmt.Sprintf("%s [[sample size]%%] [file path]\n", command)

var logger *log.Logger

// flag.Usage help message override
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s", invocation)
}

func init() {
	logger = log.New(os.Stderr, "[SNL] ", log.LstdFlags|log.Lshortfile)
	SAMPLE_MAP = make(map[int]string)
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
func parseFile(s string) (file *os.File){
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

func main () {
	var count int // a count of how many lines have been collected
	var candidate int // tmp variable for choosing a random number
	var done int // number of lines printed so far
	var totalOut int // number of total lines to print after calling parseValue

	sampleSize := flag.Arg(0)
	parseValue(sampleSize)

	fileName := flag.Arg(1)

	file := parseFile(fileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// store all lines in a map with a line number index
		SAMPLE_MAP[count] = fmt.Sprint(scanner.Text())
		count++
	}

	// a log of which line numbers we have seen
	seen := make(map[int]bool)

	// calulate the number of values we need to print to stdout
	if SAMPLE_TYPE == INTEGER {
		totalOut = SAMPLE_VALUE
	} else if SAMPLE_TYPE == PERCENTAGE {
		totalOut = int((float64(SAMPLE_VALUE) / 100.0) * float64(count))
	}

	for {
		candidate = rand.Intn(count)

		// if we haven't printed this line before, print to stdout
		if seen[candidate] != true {
			fmt.Println(SAMPLE_MAP[candidate])
			seen[candidate] = true
			done++
		}
		if done == totalOut {
			os.Exit(0)
		}
	}
}
