package command

import (
	"fmt"
	"github.com/bobappleyard/readline"
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
	`Changes the current context to the specified path. The path addresses
the name of an object in the project model like "main-server-group".

If the object is part of a collection you can also use an index (zero based)
on the objects type. To avoid naming conflicts you have to prefix the relevant
path segment with ' in that case:

    cd 'hosts[2].'servers[4]

changes the current context to the fifth server of the third host.`,
	// tab completer
	func(project *model.Project, query, _ string) []string {
		if query != "" && query != "cd" {
			readline.CompletionAppendChar = 0
		}

		var results []string
		var completePath path.Path
		var reminder string

		lastDot := strings.LastIndex(query, ".")
		if lastDot != -1 {
			q, err := path.Parse(query[0:lastDot])
			if err != nil {
				return nil
			}
			completePath = path.CurrentPath.Append(q)
			reminder = query[lastDot+1:]
		} else {
			completePath = path.CurrentPath
			reminder = query
		}

		context, err := completePath.Resolve(project)
		if err == nil {
			tags, err := reflections.Tags(context, "json")
			if err != nil {
				return nil
			}
			for field, tag := range tags {
				kind, _ := reflections.GetFieldKind(context, field)
				if kind == reflect.Struct || kind == reflect.Slice {
					if strings.HasPrefix(tag, reminder) {
						if completePath.IsEmpty() {
							results = append(results, tag)
						} else {
							results = append(results, fmt.Sprintf("%s.%s", completePath, tag))
						}
					}
				}
			}
		}
		return results
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
