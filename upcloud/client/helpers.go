package client

import (
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	log "github.com/UpCloudLtd/upcloud-go-api/logger"

	"golang.org/x/term"
)

func dumpHTTPRequest(req *http.Request) {
	originalHeader := req.Header.Clone()
	anonymizeRequest(req)
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Warnf("cannot dump http request: %s", err)
	}

	terminalWidth := GetTerminalWidth()
	horizontalDelimiterWidth := terminalWidth - 6 // sub the word "debug"
	if horizontalDelimiterWidth < 0 {
		horizontalDelimiterWidth = 0
	}
	dumpString := strings.Repeat("#", horizontalDelimiterWidth)
	dumpString += "\n\t### UPCLOUD SDK HTTP REQUEST ###\n"
	dumpString += "%s\n"
	dumpString += strings.Repeat("#", terminalWidth)

	log.Debugf(dumpString, requestDump)

	//Restore original headers
	req.Header = originalHeader
}

func dumpHTTPResponse(res *http.Response) {
	responseDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Warnf("cannot dump http response: %s", err)
	}

	terminalWidth := GetTerminalWidth()
	horizontalDelimiterWidth := terminalWidth - 6 // sub the word "debug "
	if horizontalDelimiterWidth < 0 {
		horizontalDelimiterWidth = 0
	}
	dumpString := strings.Repeat("#", horizontalDelimiterWidth)
	dumpString += "\n\t### UPCLOUD SDK HTTP RESPONSE ###\n"
	dumpString += "%s\n"
	dumpString += strings.Repeat("#", terminalWidth)

	log.Debugf(dumpString, responseDump)
}

// anonymizeRequest hides critical and sensitive data in a request for dumping purpose
func anonymizeRequest(req *http.Request) {
	// Anonymize basic auth hdr
	req.Header.Set("Authorization", "XXXXXXXXXXXXXX")
}

// GetTerminalWidth tries to figure out the width of the terminal and returns it
// returns 0 if there are problems in getting the width.
func GetTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return w
}
