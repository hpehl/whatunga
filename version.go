package main

import (
	"fmt"
	"runtime"
)

const (
	AppName         = "whatunga"
	AppVersionMajor = 0
	AppVersionMinor = 1
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

	return fmt.Sprintf("%s %d.%d.%s (Go runtime %s).\nCopyright (c) 2010-2013, Jim Teeuwen.",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionRev, runtime.Version())
}

const (
	WildFly = "WildFly"
	EAP     = "EAP"
)

type ProductVersion struct {
	Product string
	Version string
}

func (v ProductVersion) String() string {
	return v.Product + ":" + v.Version
}

var SupportedVersions = []ProductVersion{
	{WildFly, "8.0"},
	{WildFly, "8.1"},
	{EAP, "6.3"},
}

// the xmlns version of the config files
var ModelVersions = map[string]ProductVersion{
	"2.0": {WildFly, "8.0"},
	"2.1": {WildFly, "8.1"},
	"1.6": {EAP, "6.3"},
}
