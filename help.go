package main

import (
	"errors"
	"fmt"
)

func help(_ *Project, args []string) error {
	if len(args) == 0 {

		// general help
		fmt.Printf(`
help [command]
  Shows the list of available commands or context sensitive help

show config|server-groups|hosts|source|docker
  Shows status information:
    - config: Shows the current configuration
    - server-groups: Lists all server groups
    - hosts: Lists all hosts
    - source: Prints the complete project model
    - docker: Provides information about the Docker status and version

cd path
  Changes the current context to the specified path.

add server-group|host|server|deployment|user value,... [--times=n]
  Adds one or several objects to the project model.

set path value,...
  Modifies an object / attribute of the project model.

rm path
  Removes an object from the project model.

path
  Describes the concepts and syntax behind paths and when and how they can be used.

value
  Give details about the different kind of values in %s

validate
  Checks whether the project model is valid.

docker create|start
  Docker related commands
    - create: Creates docker images based on the current project model.
    - start: Starts the docker images.

exit
  Get out of here.

!cmd
  Executes cmd as shell command.\n`, AppName)

	} else if len(args) == 1 {
		// context sensitive help
		switch args[0] {

		case "show":
			fmt.Println(`
Context sensitive help for show not yet implemented!`)

		case "cd":
			fmt.Println(`
Context sensitive help for cd not yet implemented!`)

		case "add":
			fmt.Println(`
Context sensitive help for add not yet implemented!`)

		case "set":
			fmt.Println(`
Context sensitive help for set not yet implemented!`)

		case "rm":
			fmt.Println(`
Context sensitive help for rm not yet implemented!`)

		case "path":
			fmt.Println(`
Context sensitive help for path not yet implemented!`)

		case "value":
			fmt.Println(`
Context sensitive help for value not yet implemented!`)

		case "validate":
			fmt.Println(`
Context sensitive help for validate not yet implemented!`)

		case "docker":
			fmt.Println(`
Context sensitive help for docker not yet implemented!`)

		case "exit":
			fmt.Printf(`
Exits %s\n`, AppName)

		default:
			return errors.New(fmt.Sprintf(`Unsupported command "%s"`, args[0]))
		}
	} else {
		return errors.New("Too many arguments. Usage: help [command]")
	}
	return nil
}
