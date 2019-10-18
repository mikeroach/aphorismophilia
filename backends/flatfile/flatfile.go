package flatfile

/* The flatfile backend reads a quote file from the local filesystem and
   can return either the file's entire contents or a randomly selected
   single line (default behavior) based on the 'mode' parameter. */

import (
	//"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ReadFile reads a hardcoded file and returns it as a string.
func ReadFile() string {
	const filename string = "/wisdom.txt"
	pwd, _ := os.Getwd()
	data, readErr := ioutil.ReadFile(pwd + filename)
	if readErr != nil {
		panic(readErr)
	}
	content := string(data)
	return content
}

// RandomLine selects a random one-line quote to return from its input,
// which we expect to receive from the ReadFile function.
// TODO: Exclude comment-prefixed and blank lines.
func RandomLine(content string) string {
	// Split input string into slice elements delimited by newline
	quotes := strings.SplitAfter(content, "\n")
	// For debugging: fmt.Println(len(quotes))
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Return a random slice element
	return quotes[rand.Intn(len(quotes))]
}

// Return outputs a quote(s) based on requested mode.
func Return(mode string) string {
	var output string
	switch mode {
	case "all":
		output = ReadFile() // Return entire quotes file if requested...
	case "random":
		output = RandomLine(ReadFile()) // ...otherwise return a single random quote.
	default:
		output = RandomLine(ReadFile()) // Default to a single randomly selected line.
	}
	return output
}
