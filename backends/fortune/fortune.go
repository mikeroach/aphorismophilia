package fortune

import (
	"os/exec"
)

// Return executes the locally installed fortune command (thereby
// violating 12-factor principles around external dependencies)
// and returns the output.
func Return(mode string) string {
	var opt string

	switch mode {
	case "obscene":
		opt = "-o"
	case "all":
		opt = "-a"
	default:
		opt = ""
	}

	fortuneOut, fortuneErr := exec.Command("fortune", opt).Output()
	if fortuneErr != nil {
		panic(fortuneErr)
	}
	output := string(fortuneOut)
	return output
}
