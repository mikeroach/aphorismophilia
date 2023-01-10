package main

import (
	"aphorismophilia/backends/flatfile"
	"aphorismophilia/backends/fortune"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

const version string = "0.2.0"

// Build string is overridden at compile time by linker flags in Dockerfile (locally or via Jenkins).
var build = "local" //nolint since we want part of our version identifier to be a global variable.

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(response http.ResponseWriter, request *http.Request) {

	// Capture client request details to display debug information.
	dump, dumpErr := httputil.DumpRequest(request, true)
	if dumpErr != nil {
		panic(dumpErr)
	}

	/* Determine whether to show debug information by default based
	   on GET parameters from client HTTP request. 					*/

	var debug = false
	getDebug := strings.Join(request.URL.Query()["debug"], " ")

	switch getDebug {
	case "yes", "1", "true":
		debug = true
	case "no", "0", "false":
		debug = false
	}

	// Determine backend from GET parameters of client HTTP request.
	// TODO: Implement this elegantly.
	// Note: Multiple GET parameters with duplicate key result in dummy response.
	backend := strings.Join(request.URL.Query()["backend"], " ")

	// Mode passes client requested options to some backends to affect their operation.
	mode := strings.Join(request.URL.Query()["mode"], " ")

	/*
	   TODO: Document and automatically expose backend-specific mode options. Meanwhile:
	   Flatfile: 'all' = returns entire quotes file, 'random' (default) = returns single random quote
	   Fortune: 'obscene' = returns potentially offensive fortunes,
				'all' = returns both offensive and SFW fortunes,
				'<null>' (default) = returns SFW fortunes only
	   Dummy: No options. This backend isn't included in the HTML template for end-user display.
	*/

	// Call the desired backend function to return a quote. Default to fortune.
	// TODO: Expose configuration parameters to define defaults at runtime.
	var wisdom string
	switch backend {
	case "flatfile":
		wisdom = flatfile.Return(mode)
	case "fortune":
		wisdom = fortune.Return(mode)
	case "dummy":
		wisdom = "My mind is going... I can feel it. Dummy response returned."
	default:
		wisdom = fortune.Return(mode)
	}

	// Define a struct type with all input needed to render the HTML template.
	type HTMLOutput struct {
		WisdomOutput  string // Quote to display
		DebugOutput   string // Debug output to display
		VersionOutput string // Version and build info to display
		ShowDebug     bool   // Whether to expand debug output by default
	}

	// Declare an HTMLOutput struct with specific values for rendering this request's HTML.
	templateIn := HTMLOutput{
		WisdomOutput:  wisdom,
		DebugOutput:   string(dump),
		VersionOutput: "Version: " + version + " Build: " + build,
		ShowDebug:     debug,
	}

	// FIXME: This needs test coverage.

	// Parse HTML template (declared over in html.go).
	template := template.Must(template.New("output").Parse(html))
	// Render HTML template based on request data and write it to the HTTP response.
	templateExecErr := template.Execute(response, templateIn)
	if templateExecErr != nil {
		panic(templateExecErr)
	}

	/* TODO: Throw an HTTP error and custom failure message if template
	rendering fails. Since Go's template execute function writes directly
	to an IO writer, we should probably use "defer" or capture the output
	into a string with something like:

	var renderedTemplate string
	templateOut := new(bytes.Buffer)
	template.Execute(templateOut, templateIn)
	renderedTemplate = templateOut.String()

	// Write rendered template string in HTTP response
	fmt.Fprintf(response, "%s\n\n", renderedTemplate)

	*/
}
