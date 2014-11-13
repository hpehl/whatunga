package main

import (
	"fmt"
	"strings"
)

var dockerUsage = "docker create|start"

var docker = Command{
	"docker",
	dockerUsage,
	`Docker related commands
    - create: Creates docker images based on the current project model.
    - start: Starts the docker images.`,
	// tab completer
	func(query, _ string) []string {
		var results []string
		subCommands := [...]string{"create", "start"}
		for _, subCommand := range subCommands {
			if strings.HasPrefix(subCommand, query) {
				results = append(results, subCommand)
			}
		}
		return results
	},
	// action
	func(_ *Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", dockerUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", dockerUsage)
		}
		switch args[0] {
		case "create":
			fmt.Printf("\nDocker create...\n")
		case "start":
			fmt.Printf("\nDocker start...\n")
		default:
			return fmt.Errorf(`Unsupported argument "%s". Usage: %s`, args[0], dockerUsage)
		}
		return nil
	},
}
