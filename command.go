package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// TODO Introduce a command struct similar to cli.go
// and use that to implement help and tab completion
type CommandFunc func(*Project, []string) error

type Command struct {
	// The name of the command
	Name string
	// A short description of the usage of this command
	Usage string
	// The function to call when checking for bash command completions
	Completer func(query, ctx string) []string
	// The function to call when this command is invoked
	Action func(*Project, []string) error
}

func show(project *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing argument. Usage: show config|server-groups|hosts|source|docker")
	}
	if len(args) > 1 {
		return errors.New("Too many arguments. Usage: show config|server-groups|hosts|source|docker")
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
		return errors.New(
			fmt.Sprintf(`Unsupported argument "%s". Usage: show config|server-groups|hosts|source|docker`, args[0]))
	}
	return nil
}

func cd(_ *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing argument. Usage: cd path")
	}
	if len(args) > 1 {
		return errors.New("Too many arguments. Usage: cd path")
	}

	// TODO validate path
	workingDir = args[0]
	return nil
}

func add(_ *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing arguments. Usage: add server-group|host|server|deployment|user value,... [--times=n]")
	}
	return nil
}

func set(_ *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing arguments. Usage: set path value,...")
	}
	if len(args) > 2 {
		return errors.New("Too many arguments. Usage: set path value,...")
	}
	return nil
}

func rm(_ *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing argument. Usage: rm path")
	}
	if len(args) > 1 {
		return errors.New("Too many arguments. Usage: rm path")
	}
	return nil
}

func validate(_ *Project, args []string) error {
	if len(args) != 0 {
		return errors.New("Illegal argument. Usage: validate")
	}
	return nil
}

func docker(_ *Project, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing argument. Usage: docker create|start")
	}
	if len(args) > 1 {
		return errors.New("Too many arguments. Usage: docker create|start")
	}
	switch args[0] {
	case "create":
		fmt.Printf("\nDocker create...\n")
	case "start":
		fmt.Printf("\nDocker start...\n")
	default:
		return errors.New(
			fmt.Sprintf(`Unsupported argument "%s". Usage: docker create|start`, args[0]))
	}
	return nil
}

func topLevelCompleter(query, ctx string) []string {
	var commands []string
	for key, _ := range commandRegistry {
		if strings.HasPrefix(key, query) {
			commands = append(commands, key)
		}
	}
	return commands
}

var commandRegistry = map[string]CommandFunc{
	"help":     help,
	"show":     show,
	"cd":       cd,
	"add":      add,
	"set":      set,
	"rm":       rm,
	"validate": validate,
	"docker":   docker,
	"exit": func(_ *Project, _ []string) error {
		saveHistory()
		fmt.Println("\nHaere rā")
		os.Exit(0)
		return nil
	},
}
