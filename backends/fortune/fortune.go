package fortune

import (
	"os/exec"
)

// FIXME: Alpine Linux's fortune package returns offensive fortunes
// even when the "-o" flag is omitted.

// Return executes the locally installed fortune command (thereby
// violating 12-factor principles around external dependencies)
// and returns the output.
func Return(mode string) string {
	var opt string
	if mode == "obscene" {
		opt = "-o"
	} else {
		opt = ""
	}
	fortuneOut, fortuneErr := exec.Command("fortune", opt).Output()
	if fortuneErr != nil {
		panic(fortuneErr)
	}
	output := string(fortuneOut)
	return output
}
