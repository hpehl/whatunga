package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"strings"
)

var dockerUsage = "docker create|push|start"

var docker = Command{
	"docker",
	"Docker related commands",
	dockerUsage,
	`Docker related commands
    - create: Creates docker images based on the current project model.
    - start: Starts the docker images.`,
	// tab completer
	func(_ *model.Project, query, _ string) ([]string, int) {
		var results []string
		subCommands := [...]string{"create", "push", "start"}
		for _, subCommand := range subCommands {
			if strings.HasPrefix(subCommand, query) {
				results = append(results, subCommand)
			}
		}
		return results, ' '
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", dockerUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", dockerUsage)
		}
		// TODO not yet implemented
		switch args[0] {
		case "create":
			fmt.Printf("Docker create...\n")
		case "push":
			fmt.Printf("Docker push...\n")
		case "start":
			fmt.Printf("Docker start...\n")
		default:
			return fmt.Errorf(`Unsupported argument "%s". Usage: %s`, args[0], dockerUsage)
		}
		return nil
	},
}
