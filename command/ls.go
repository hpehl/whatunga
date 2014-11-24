package command

import (
	"encoding/json"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"reflect"
)

var lsUsage = "ls [path]"

var ls = Command{
	"ls",
	"Lists the model of the current context or specified path",
	lsUsage,
	`Lists the model of the current context or specified path`,
	// tab completer
	func(project *model.Project, query, cmdline string) ([]string, int) {
		return completion(project, query, cmdline, []reflect.Kind{})
	},
	// action
	func(project *model.Project, args []string) error {
		var context path.Path

		if len(args) == 0 {
			context = path.CurrentPath
		} else if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", lsUsage)
		} else {
			if args[0] == ".." {
				if path.CurrentPath.IsEmpty() {
					return fmt.Errorf("Cannot go up one level: Already at root")
				}
				context = path.CurrentPath[0 : len(path.CurrentPath)-1]
			} else if args[0] == "/" {
				if path.CurrentPath.IsEmpty() {
					return fmt.Errorf("Already at root")
				}
				context = path.CurrentPath[:0]
			} else {
				pth, err := path.Parse(args[0])
				if err != nil {
					return err
				}
				context = path.CurrentPath.Append(pth)
			}
		}

		obj, err := context.Resolve(project)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(obj, "", "  ")
		if err == nil {
			fmt.Printf("%s\n", string(data))
		} else {
			return err
		}
		return nil
	},
}
