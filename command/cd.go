package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"github.com/oleiade/reflections"
	"reflect"
	"strings"
)

var cdUsage = "cd <path>"

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
			pth, reminder := split(tokens[1])
			if pth == "" {
				// is the reminder a valid path?
				children := ls(path.CurrentPath, project)
				if contains(keys(children), reminder) {
					// yep. check the type to find the right CompletionAppendChar and return that match
					if children[reminder] == reflect.Struct {
						return []string{reminder}, '.'
					} else {
						return []string{reminder}, '['
					}
				} else {
					// no. return a slice of matches based on the reminder
					var matches []string
					for _, name := range keys(children) {
						if strings.HasPrefix(name, reminder) {
							matches = append(matches, name)
						}
					}
					if len(matches) == 1 {
						if children[matches[0]] == reflect.Struct {
							return matches, '.'
						} else {
							return matches, '['
						}
					} else {
						return matches, 0
					}
				}
			}
		}

		return nil, 0

		//		readline.CompletionAppendChar = 0
		//		var matches []string
		//		var completePath path.Path
		//		var reminder string
		//
		//		lastDot := strings.LastIndex(query, ".")
		//		if lastDot != -1 {
		//			q, err := path.Parse(query[0:lastDot])
		//			if err != nil {
		//				return nil
		//			}
		//			completePath = path.CurrentPath.Append(q)
		//			reminder = query[lastDot+1:]
		//		} else {
		//			completePath = path.CurrentPath
		//			reminder = query
		//		}
		//
		//		context, err := completePath.Resolve(project)
		//		if err == nil {
		//			tags, err := reflections.Tags(context, "json")
		//			if err != nil {
		//				return nil
		//			}
		//			for field, tag := range tags {
		//				kind, _ := reflections.GetFieldKind(context, field)
		//				if kind == reflect.Struct || kind == reflect.Slice {
		//					if strings.HasPrefix(tag, reminder) {
		//						matches = append(matches, tag)
		//					}
		//				}
		//			}
		//		}
		//		return matches
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", cdUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", cdUsage)
		}
		// TODO not yet implemented
		fmt.Println("Not yet implemented!")
		return nil
	},
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
