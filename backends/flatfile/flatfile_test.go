package flatfile

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	run := ReadFile()
	fmt.Printf("%s\n", run)
}

func TestRandom(t *testing.T) {
	content := ReadFile()
	run := RandomLine(content)
	fmt.Printf("%s\n", run)
}

func TestReturnAll(t *testing.T) {
	run := Return("all")
	fmt.Printf("%s\n", run)
}

func TestReturnRandom(t *testing.T) {
	run := Return("random")
	fmt.Printf("%s\n", run)
}
