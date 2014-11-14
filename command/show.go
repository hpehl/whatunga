package command

import (
	"encoding/json"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"strings"
)

var showUsage = "show config|server-groups|hosts|source|docker"

var show = Command{
	"show",
	"Shows status information",
	showUsage,
	`Shows status information:
    - config: Shows the current configuration.
    - server-groups: Lists all server groups.
    - hosts: Lists all hosts.
    - source: Prints the complete project model.
    - docker: Provides information about the Docker status and version.`,
	// tab completer
	func(query, _ string) []string {
		var results []string
		subCommands := [...]string{"config", "server-groups", "hosts", "source", "docker"}
		for _, subCommand := range subCommands {
			if strings.HasPrefix(subCommand, query) {
				results = append(results, subCommand)
			}
		}
		return results
	},
	// action
	func(project *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", showUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", showUsage)
		}

		switch args[0] {
		case "config":
			fmt.Printf("\n%s\n", project.Config)
		case "server-groups":
			fmt.Printf("\n%v\n", project.ServerGroups)
		case "hosts":
			fmt.Printf("\n%v\n", project.Hosts)
		case "source":
			data, err := json.MarshalIndent(project, "", "  ")
			if err == nil {
				fmt.Printf("\n%s\n", string(data))
			} else {
				return err
			}
		case "docker":
			fmt.Printf("\nDocker not yet implemented\n")
		default:
			return fmt.Errorf(`Unsupported argument "%s". Usage: %s`, args[0], showUsage)
		}
		return nil
	},
}
