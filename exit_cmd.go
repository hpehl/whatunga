package main

import (
	"fmt"
	"os"
)

var exit = Command{
	"exit",
	"exit",
	"Get out of here.",
	// tab completer
	func(_, _ string) []string {
		return nil
	},
	// action
	func(_ *Project, _ []string) error {
		saveHistory()
		fmt.Println("\nHaere rā")
		os.Exit(0)
		return nil
	},
}
