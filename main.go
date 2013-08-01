// TODO: - help string for positional arguments
// - streams
// - split percentage of the data into train and test files

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

var logger *log.Logger
var SAMPLE map[int]string
var count int

func init() {
	logger = log.New(os.Stderr, "[SNL] ", log.LstdFlags|log.Lshortfile)
	SAMPLE = make(map[int]string)

}

func main () {
	rand.Seed(time.Now().UTC().UnixNano())
	sample_size_help := "The number of lines to sample from a file"
	var sampleSize = flag.Int("sample_size", 0, sample_size_help)

	filename_help := "The file which to sample from"
	var filename = flag.String("file", "", filename_help)

	flag.Parse()

	// try to get the positional argument for count argument
	if *sampleSize == 0 {
		stringSize := flag.Arg(0)
		s, err := strconv.Atoi(stringSize)
		sampleSize = &s
		if err != nil {
			logger.Fatal("[Error] error converting %s to integer: %s", stringSize, err)
		}
	}

	// try to get the positional argument for filename argument
	if *filename == "" {
		f := flag.Arg(1)
		if f == "" {
			logger.Fatal("[Error] missing filename")
		}
		filename = &f
	}

	file, err := os.Open(*filename)

	if err != nil {
		logger.Fatal("[Error] error opening %s: %s", *filename, err)
	}
	defer file.Close()

	// TODO: - for streams,
	// we need to store the first `count` items from the file in a hash or circular buffer
	// then we replace lines at random in something that ends up being a random
	// distribution

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// for files
		// randomly store this line after first collecting the first `sampleSize` items
		SAMPLE[count] = fmt.Sprint(scanner.Text())
		count++
	}

	// a log of which line numbers we have seen
	seen := make(map[int]bool)

	var candidate int
	var done int
	for {
		candidate = rand.Intn(count)

		// if we haven't seen this before, print to stdout
		if seen[candidate] != true {
			fmt.Println(SAMPLE[candidate])
			seen[candidate] = true
			done++
		}
		if done == *sampleSize {
			goto DONE
		}
	}
DONE:
	os.Exit(0)
}
