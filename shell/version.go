package shell

import (
	"fmt"
	"runtime"
)

const (
	AppName         = "whatunga"
	AppVersionMajor = 0
	AppVersionMinor = 3
)

// revision part of the program version.
// This will be set automatically at build time like so:
//
//     go build -ldflags "-X main.AppVersionRev `date -u +%s`"
var AppVersionRev string

func Version() string {
	if len(AppVersionRev) == 0 {
		AppVersionRev = "0"
	}

	return fmt.Sprintf("%s %d.%d.%s (Go runtime %s).",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionRev, runtime.Version())
}
