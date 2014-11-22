package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"github.com/oleiade/reflections"
	"reflect"
	"strings"
)

var cdUsage = "cd <path> | cd .. | cd /"

var cd = Command{
	"cd",
	"Changes the current context to the specified path.",
	cdUsage,
	`Changes the current context to the specified path. The path specifies
the name of an object in the project model like "config.templates.domain".

If the object is part of a collection you can also use index on the objects
type. Both numeric and name based indizes are supported:

	cd hosts[master].servers[4]

changes the current context to the fifth server of host "master".`,
	// tab completer
	func(project *model.Project, query, cmdline string) ([]string, int) {
		tokens := strings.Fields(cmdline)

		if len(tokens) == 1 && query == "" {
			// just the command was given, return matches based on the current path
			return keys(ls(path.CurrentPath, project)), 0

		} else if len(tokens) > 1 {
			// check the stuff after "cd"
			pathPart, reminder := split(tokens[1])

			if pathPart == "" {
				// no full path. return matches based on the reminder
				return matchesFor(path.CurrentPath, reminder, project)

			} else {
				// some kind of path given. append to path.CurrentPath
				p, err := path.Parse(pathPart)
				if err != nil {
					return nil, 0
				}
				fullPath := path.CurrentPath.Append(p)
				return matchesFor(fullPath, reminder, project)
			}
		}
		return nil, 0
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
			path.CurrentPath = path.CurrentPath[0 : len(path.CurrentPath)-1]
		} else if args[0] == "/" {
			if path.CurrentPath.IsEmpty() {
				return fmt.Errorf("Already at root")
			}
			path.CurrentPath = path.CurrentPath[:0]
		} else {
			changeTo, err := path.Parse(args[0])
			if err != nil {
				return err
			}
			full := path.CurrentPath.Append(changeTo)
			if _, err := full.Resolve(project); err != nil {
				return err
			}
			path.CurrentPath = full
		}
		return nil
	},
}

func matchesFor(context path.Path, reminder string, project *model.Project) ([]string, int) {
	children := ls(context, project)

	if contains(keys(children), reminder) {
		if children[reminder] == reflect.Struct {
			return []string{reminder}, 0
		} else {
			return []string{reminder}, '['
		}
	} else {
		var matches []string
		for _, name := range keys(children) {
			if strings.HasPrefix(name, reminder) {
				matches = append(matches, name)
			}
		}
		if len(matches) == 1 {
			if children[matches[0]] == reflect.Struct {
				return matches, 0
			} else {
				return matches, '['
			}
		} else {
			return matches, 0
		}
	}
}

func ls(context path.Path, project *model.Project) map[string]reflect.Kind {
	var matches = make(map[string]reflect.Kind)

	obj, err := context.Resolve(project)
	if err == nil {
		tags, err := reflections.Tags(obj, "json")
		if err == nil {
			for field, tag := range tags {
				kind, _ := reflections.GetFieldKind(obj, field)
				if kind == reflect.Struct || kind == reflect.Slice {
					matches[tag] = kind
				}
			}
		}
	}
	return matches
}

func keys(m map[string]reflect.Kind) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func split(arg string) (string, string) {
	var path, reminder string
	lastDot := strings.LastIndex(arg, ".")
	if lastDot != -1 {
		path = arg[0:lastDot]
		reminder = arg[lastDot+1:]
	} else {
		path = ""
		reminder = arg
	}
	return path, reminder
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
