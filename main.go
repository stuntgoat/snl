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

var CANDIDATE int
var MUST_KEEP int
var SAMPLE_MAP map[int]string
var SAMPLE_VALUE int
var SAMPLE_PERCENTAGE float64
var CURRENT_COUNT int

var logger *log.Logger

var command = os.Args[0]
var invocation = fmt.Sprintf("%s [sample percentage] -\n", command)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s", invocation)
}

func parseValue(s string) int {
	var value string

	// convert value to integer
	intValue, err := strconv.Atoi(s)
	if err != nil {
		logger.Printf("[Error] error converting sample_size: %s to integer: %s", value, err)
		fmt.Printf("Usage: %s", invocation)
		os.Exit(1)
	}
	return intValue
}

func init() {
	logger = log.New(os.Stderr, "[snls] ", log.LstdFlags|log.Lshortfile)

	SAMPLE_MAP = make(map[int]string)

	rand.Seed(time.Now().UTC().UnixNano())

	flag.Usage = Usage
	flag.Parse()
}

func keepPercentage() {
	MUST_KEEP = int(SAMPLE_PERCENTAGE * float64(CURRENT_COUNT))
	if MUST_KEEP == 0 {
		os.Exit(0)
	}

	// calculate a safe threshold to randomly remove from the map
	mustDelete := len(SAMPLE_MAP) - MUST_KEEP
	if mustDelete == 0 {
		return
	}
	var candidate int
	for  {
		candidate = rand.Intn(CURRENT_COUNT)
		_, exists := SAMPLE_MAP[candidate]

		if exists {
			delete(SAMPLE_MAP, candidate)
			mustDelete--
		}
		if mustDelete == 0 {
			break
		}
	}
}

func printFinal() {
	for _, line := range SAMPLE_MAP {
		fmt.Println(line)
	}

	os.Exit(0)
}

func handleSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<- sigChannel

	keepPercentage()

	printFinal()
}

func main() {
	sizeRequest := flag.Arg(0)

	SAMPLE_VALUE := parseValue(sizeRequest)
	SAMPLE_PERCENTAGE = float64(SAMPLE_VALUE) / 100.0

	SAMPLE_MAP = make(map[int]string)

	go handleSignal()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		SAMPLE_MAP[CURRENT_COUNT] = fmt.Sprint(scanner.Text())
		CURRENT_COUNT++

		if CURRENT_COUNT % 300 == 0 {
			keepPercentage()
		}
	}
	keepPercentage()

	printFinal()
}
