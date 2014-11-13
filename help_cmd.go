package main

import (
	"bytes"
	"fmt"
	"strings"
)

var helpUsage = "help [command]"

var help = Command{
	"help",
	helpUsage,
	"Shows the list of available commands or context sensitive help.",
	// tab completer
	func(query, _ string) []string {
		var results []string
		for key, _ := range commandRegistry {
			if key == "help" {
				continue
			} else {
				if strings.HasPrefix(key, query) {
					results = append(results, key)
				}
			}
		}
		return results
	},
	// action
	func(_ *Project, args []string) error {
		if len(args) == 0 {
			// general help
			var buffer bytes.Buffer
			for _, cmd := range commandRegistry {
				buffer.WriteString(fmt.Sprintf("%s\n\t%s\n\n", cmd.Usage, cmd.Help))
			}
			fmt.Println(buffer.String())
		} else if len(args) == 1 {
			cmd, valid := commandRegistry[args[0]]
			if !valid {
				fmt.Printf("Unknown command: \"%s\"\n", args[0])
			} else {
				fmt.Println(cmd.Help)
			}
		} else {
			return fmt.Errorf("Too many arguments. Usage: %s", helpUsage)
		}
		return nil
	},
}
