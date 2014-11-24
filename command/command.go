package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"github.com/oleiade/reflections"
	"reflect"
	"strconv"
	"strings"
)

var cmdOrder []string

type registry map[string]Command

func (r registry) Add(cmd Command) {
	r[cmd.Name] = cmd
	cmdOrder = append(cmdOrder, cmd.Name)
}

// use internal iteration to control the iteration order
func (r registry) forEach(fn func(cmd Command)) {
	for _, name := range cmdOrder {
		fn(r[name])
	}
}

var Registry = make(registry)

func init() {
	Registry.Add(show)
	Registry.Add(cd)
	Registry.Add(ls)
	Registry.Add(add)
	Registry.Add(set)
	Registry.Add(rm)
	Registry.Add(validate)
	Registry.Add(docker)
	Registry.Add(exit)
	Registry.Add(help)
}

type Command struct {
	// The name of the command
	Name string
	// A short description of this command
	Description string
	// How to use this command
	Usage string
	// A more detailed description for this command
	UsageDescription string
	// The function to call when checking for tab completion
	Completer func(project *model.Project, query, cmdline string) ([]string, int)
	// The function to call when this command is invoked
	Action func(project *model.Project, args []string) error
}

// ------------------------------------------------------ functions used by several commands

// TODO add parameter []reflect.Kind to select which kind of objects / attributes to return
func completion(project *model.Project, query, cmdline string) ([]string, int) {
	tokens := strings.Fields(cmdline)
	if len(tokens) == 1 && query == "" {
		// just the command was given, return matches based on the current path
		return keys(children(project, path.CurrentPath)), 0

	} else if len(tokens) > 1 {
		// check the stuff after "ls"
		possiblePath, segment := path.SplitLastSegment(tokens[1])

		if possiblePath == "" {
			// no full path. return matches based on the segment
			return matchesFor(project, path.CurrentPath, segment)

		} else {
			// some kind of path given. append to path.CurrentPath
			// and return matches for the full path
			pth, err := path.Parse(possiblePath)
			if err != nil {
				return nil, 0
			}
			fullPath := path.CurrentPath.Append(pth)
			return matchesFor(project, fullPath, segment)
		}
	}
	return nil, 0
}

func matchesFor(project *model.Project, context path.Path, segment string) ([]string, int) {
	children := children(project, context)
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
				matches = append(matches, indices(project, context, beforeIndex, index)...)
				matches = append(matches, names(project, context, beforeIndex, index)...)
			} else {
				_, err := strconv.ParseUint(index, 10, 32)
				if err == nil {
					// unclosed numeric index
					matches = append(matches, indices(project, context, beforeIndex, index)...)
				} else {
					// unclosed alphanumeric index
					matches = append(matches, names(project, context, beforeIndex, index)...)
				}
			}

			//			log.Printf(`
			//-+-+-+-+-
			//# matches: %d
			//matches:   %v
			//-+-+-+-+-
			//`, len(matches), matches)

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

func children(project *model.Project, context path.Path) map[string]reflect.Kind {
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

func indices(project *model.Project, context path.Path, name string, index string) []string {
	slice := getSlice(project, context, name)
	if !slice.IsValid() {
		return nil
	}
	var matches []string
	for i := 0; i < slice.Len(); i++ {
		strIndex := fmt.Sprintf("%d", i)
		if strings.HasPrefix(strIndex, index) {
			matches = append(matches, strIndex)
		}
	}
	return matches
}

func names(project *model.Project, context path.Path, name string, index string) []string {
	slice := getSlice(project, context, name)
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

func getSlice(project *model.Project, context path.Path, name string) reflect.Value {
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
