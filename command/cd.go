package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"reflect"
)

var cdUsage = "cd <path> | cd .. | cd / | cd -"

var cd = Command{
	"cd",
	"Changes the current context to the given path.",
	cdUsage,
	`Changes the current context to the given path. The path is composed of
the names of the objects in the project model seperated with dots:

    config.templates.domain

If the object is part of a collection you can also use an index on the object's
type. Both numeric and name based indizes are supported:

    hosts[master].servers[4]

Addresses the fifth server of host "master".`,
	// tab completer
	func(project *model.Project, query, cmdline string) ([]string, int) {
		return completion(project, query, cmdline, []reflect.Kind{reflect.Struct, reflect.Slice})
	},
	// action
	func(project *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", cdUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", cdUsage)
		}

		if args[0] == ".." {
			if path.CurrentPath.IsEmpty() {
				return fmt.Errorf("Cannot go up one level: Already at root")
			}
			internalCd(path.CurrentPath[0 : len(path.CurrentPath)-1])

		} else if args[0] == "/" {
			if path.CurrentPath.IsEmpty() {
				return fmt.Errorf("Already at root")
			}
			internalCd(path.CurrentPath[:0])

		} else if args[0] == "-" {
			if path.LastPath == nil {
				return fmt.Errorf("No previous path")
			}
			internalCd(path.LastPath)

		} else {
			changeTo, err := path.Parse(args[0])
			if err != nil {
				return err
			}
			full := path.CurrentPath.Append(changeTo)
			if _, err := full.Resolve(project); err != nil {
				return err
			}
			internalCd(full)
		}
		return nil
	},
}

func internalCd(p path.Path) {
	path.LastPath = path.CurrentPath
	path.CurrentPath = p
}
