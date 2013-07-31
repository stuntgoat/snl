// TODO: - help string for positional arguments
// -

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
var pseudoRandom *rand.Rand


func init () {
	pseudoRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
	logger = log.New(os.Stderr, "[SNL] ", log.LstdFlags|log.Lshortfile)
}

func main () {
	count_help := "The number of lines to sample from a file"
	var sampleSize = flag.Int("count", 0, count_help)

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

	// we need to store the first `count` items from the file in a hash or circular buffer
	// then we replace lines at random in something that ends up being a random
	// distribution

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if count == *sampleSize {
			os.Exit(0)
		}
		// randomly store this line after first collecting the first `sampleSize` items
		fmt.Println(scanner.Text())
		count++
	}
}
