package fortune

import (
	"fmt"
	"testing"
)

/*
TODO: Consider running this through an obscenity filter to increase chances of
preserving delicate sensibilities in production, though this will miss subtle
context and cleanliness is ultimately in the eye of the beholder.
*/
func TestClean(t *testing.T) {
	run := Return("clean")
	fmt.Printf("%s", run)
}

func TestObscene(t *testing.T) {
	run := Return("obscene")
	fmt.Printf("%s", run)
}
