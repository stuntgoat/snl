package main

import (
	"bufio"
	"fmt"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/stuntgoat/snl/percent_sample"
	"github.com/stuntgoat/snl/integer_sample"
)


var SAMPLE_VALUE int // either a percentage or a sum to keep
var SAMPLE_TYPE int

const (
	INTEGER = iota
	PERCENTAGE
	)

var INTEGER_SAMPLE *integer_sample.IntegerSample
var PERCENT_SAMPLE *percent_sample.PercentageSample

var command = os.Args[0]
var invocationFile = fmt.Sprintf("%s [[sample size]%%] [file path]\n", command)
var invocationStdin = fmt.Sprintf("%s [[sample size]%%] -\n", command)

var logger *log.Logger

// flag.Usage help message override
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n%s%s", invocationFile, invocationStdin)
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

	intValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Printf("[Error] error converting sample_size: %s to integer: %s", value, err)
		Usage()
		os.Exit(1)
	}
	SAMPLE_VALUE = intValue
}

// parseFile validates a string and returns an *os.File
func parseFile(s string) (file *os.File) {
	if s == "" {
		logger.Print("[Error] missing filename argument")
		Usage()
		os.Exit(1)
	}

	file, err := os.Open(s)
	if err != nil {
		Usage()
		logger.Fatalf("[Error] error opening %s: %s", s, err)
	}

	return file
}

// handleSignal handles a SIGINT (control-c) when the user
// might want to break from a stream while sampling a percentage.
func handleSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<- sigChannel

	printSample()
}

func printSample() {
	if SAMPLE_TYPE == PERCENTAGE {
		PERCENT_SAMPLE.AddPercentageToTotal()
		PERCENT_SAMPLE.Print()
	} else {
		INTEGER_SAMPLE.Print()
	}
}

func main () {
	var file *os.File
	var line string

	var sampleSize = flag.Arg(0)
	parseValue(sampleSize)

	if SAMPLE_TYPE == INTEGER {
		INTEGER_SAMPLE = &integer_sample.IntegerSample{
			Sample: make([]string, SAMPLE_VALUE),
			Size: SAMPLE_VALUE,
		}
	} else if SAMPLE_TYPE == PERCENTAGE {
		PERCENT_SAMPLE = &percent_sample.PercentageSample{
			Sample: make([]string, 0),
			PercentageKeep: SAMPLE_VALUE,
			Well:  make([]string, 0),
			WellSize: 100,
			WellSeen: 0,
		}
	}

	fileName := flag.Arg(1)
	if fileName == "-" {
		file = os.Stdin
	} else {
		file = parseFile(fileName)
		defer file.Close()
	}

	go handleSignal()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprint(scanner.Text())

		if SAMPLE_TYPE == PERCENTAGE {
			PERCENT_SAMPLE.SampleLine(line)
		} else {
			INTEGER_SAMPLE.SampleLine(line)
		}
	}
	printSample()
}
