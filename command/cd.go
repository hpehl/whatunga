package command

import (
	"bytes"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"github.com/oleiade/reflections"
	"reflect"
	"strconv"
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

		var buffer bytes.Buffer
		for i, t := range tokens {
			buffer.WriteString(strconv.Quote(t))
			if i < len(tokens)-1 {
				buffer.WriteString(", ")
			}
		}

		//		log.Printf(`
		//-+-+-+-+-
		//cmdline:     "%s"
		//query:       "%s"
		//# of tokens: %d
		//tokens:      [%s]
		//-+-+-+-+-
		//`, cmdline, query, len(tokens), buffer.String())

		if len(tokens) == 1 && query == "" {
			// just the command was given, return matches based on the current path
			return keys(ls(path.CurrentPath, project)), 0

		} else if len(tokens) > 1 {
			// check the stuff after "cd"
			possiblePath, segment := path.SplitLastSegment(tokens[1])

			if possiblePath == "" {
				// no full path. return matches based on the segment
				return matchesFor(path.CurrentPath, segment, project)

			} else {
				// some kind of path given. append to path.CurrentPath
				// and return matches for the full path
				pth, err := path.Parse(possiblePath)
				if err != nil {
					return nil, 0
				}
				fullPath := path.CurrentPath.Append(pth)
				return matchesFor(fullPath, segment, project)
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

func matchesFor(context path.Path, segment string, project *model.Project) ([]string, int) {
	children := ls(context, project)
	keys := keys(children)

	if contains(keys, segment) {
		if children[segment] == reflect.Struct {
			return []string{segment}, 0
		} else {
			return []string{segment}, '['
		}
	} else {
		var matches []string
		openBracket, beforeIndex, index := path.LastOpenSquareBracket(segment)

		//		log.Printf(`
		//-+-+-+-+-
		//openBracket: %v
		//beforeIndex: "%s"
		//index:       "%s"
		//-+-+-+-+-
		//`, openBracket, beforeIndex, index)

		if openBracket {
			// it's an unclosed index
			if index == "" {
				// return both numeric and alphanumeric matches
				matches = append(matches, indices(context, beforeIndex, index, project)...)
				matches = append(matches, names(context, beforeIndex, index, project)...)
			} else {
				_, err := strconv.ParseUint(index, 10, 32)
				if err == nil {
					// unclosed numeric index
					matches = append(matches, indices(context, beforeIndex, index, project)...)
				} else {
					// unclosed alphanumeric index
					matches = append(matches, names(context, beforeIndex, index, project)...)
				}
			}
			if len(matches) == 1 {
				return matches, ']'
			} else {
				return matches, 0
			}

		} else {
			// it's a normal segment
			for _, key := range keys {
				if strings.HasPrefix(key, segment) {
					matches = append(matches, key)
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

func indices(context path.Path, name string, index string, project *model.Project) []string {
	slice := getSlice(context, name, project)
	if !slice.IsValid() {
		return nil
	}
	matches := make([]string, slice.Len())
	for i, _ := range matches {
		strIndex := fmt.Sprintf("%d", i)
		if strings.HasPrefix(strIndex, index) {
			matches[i] = strIndex
		}
	}
	return matches
}

func names(context path.Path, name string, index string, project *model.Project) []string {
	slice := getSlice(context, name, project)
	if !slice.IsValid() {
		return nil
	}

	var matches []string
	for i := 0; i < slice.Len(); i++ {
		element := slice.Index(i)
		value, err := reflections.GetField(element.Interface(), "Name")
		if err == nil {
			strValue := value.(string)
			if strings.HasPrefix(strValue, index) {
				matches = append(matches, strValue)
			}
		}
	}
	return matches
}

func getSlice(context path.Path, name string, project *model.Project) reflect.Value {
	var slice reflect.Value
	var realName string // name is the json name!

	obj, err := context.Resolve(project)
	if err == nil {
		tags, err := reflections.Tags(obj, "json")
		if err == nil {
			for field, tag := range tags {
				if tag == name {
					realName = field
					break
				}
			}
		}

		if reflect.TypeOf(obj).Kind() == reflect.Ptr {
			slice = reflect.ValueOf(obj).Elem()
		} else {
			slice = reflect.ValueOf(obj)
		}
		check := slice.FieldByName(realName)
		if check.IsValid() && reflect.TypeOf(check).Kind() != reflect.Slice {
			slice = check
		}
	}
	return slice
}

func keys(m map[string]reflect.Kind) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
