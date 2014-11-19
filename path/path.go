package path

import (
	"bytes"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/oleiade/reflections"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	// SegmentKind
	PlainSegment SegmentKind = iota
	IndexSegment
	RangeSegment

	// IndexKind
	NumericIndex IndexKind = iota
	AlphaNumericIndex

	// Undefined range
	Undefined int = -1
)

type SegmentKind int

type IndexKind int

type Path []Segment

type Segment struct {
	Name  string
	Kind  SegmentKind
	Index Index
	Range Range
}

type Index struct {
	Kind  IndexKind
	Value interface{}
}

type Range struct {
	From, To int
}

// regular expression to distinguish between the different segments
var plainSegment = regexp.MustCompile(`^([\w-]+)$`)
var indexSegment = regexp.MustCompile(`^([\w-]+)\[((\d+)|([A-Za-z0-9_-]+))\]$`)
var rangeSegment = regexp.MustCompile(`^([\w-]+)\[((\d*)(:)(\d*))\]$`)

// the current path which is used by the commands and the shell
var CurrentPath Path = []Segment{}

// Turns a string into a path
func Parse(p string) (Path, error) {
	if p == "" {
		return make(Path, 0), nil
	}

	var path Path
	segments := strings.Split(p, ".")
	for _, s := range segments {
		segment := Segment{"", PlainSegment, Index{}, Range{Undefined, Undefined}}

		if s != "" {
			// check most specific re first!
			if rangeSegment.MatchString(s) {
				groups := rangeSegment.FindStringSubmatch(s)
				segment.Name = groups[1]
				segment.Kind = RangeSegment
				if groups[3] != "" {
					from, err := strconv.Atoi(groups[3])
					if err != nil {
						return nil, fmt.Errorf(`Unable to resolve path "%s": "%s:%s" is not a valid range`, p, groups[3], groups[5])
					}
					segment.Range.From = from
				}
				if groups[5] != "" {
					to, err := strconv.Atoi(groups[5])
					if err != nil {
						return nil, fmt.Errorf(`Unable to resolve path "%s": "%s:%s" is not a valid range`, p, groups[3], groups[5])
					}
					segment.Range.To = to
				}

			} else if indexSegment.MatchString(s) {
				groups := indexSegment.FindStringSubmatch(s)
				segment.Name = groups[1]
				segment.Kind = IndexSegment
				if groups[3] != "" {
					// numeric index
					index, err := strconv.Atoi(groups[3])
					if err != nil {
						return nil, fmt.Errorf(`Unable to resolve path "%s": "%s" is not a valid numeric index`, p, groups[3])
					}
					segment.Index.Kind = NumericIndex
					segment.Index.Value = index
				} else if groups[4] != "" {
					// alpha-numeric range
					segment.Index.Kind = AlphaNumericIndex
					segment.Index.Value = groups[4]
				}

			} else if plainSegment.MatchString(s) {
				groups := plainSegment.FindStringSubmatch(s)
				segment.Name = groups[1]
				segment.Kind = PlainSegment

			} else {
				return nil, fmt.Errorf(`Invalid segment "%s" in path "%s"`, s, p)
			}
		}
		path = append(path, segment)
	}
	return path, nil
}

// Get the value of the project model the given path points to. The path must be unambiguous,
// thus it must not contain ranges.
func (path Path) Resolve(project *model.Project) (interface{}, error) {
	var context interface{} = project

	for _, segment := range path {

		// Find field referenced by the tag <segment.Name>
		tags, err := reflections.Tags(context, "json")
		if err != nil {
			return nil, fmt.Errorf(`Unable to resolve path "%s": %s.`, path, err)
		}
		var fieldName = ""
		for name, tag := range tags {
			if tag == segment.Name {
				fieldName = name
				break
			}
		}
		if fieldName == "" {
			return nil, fmt.Errorf(`Unable to resolve path "%s": Segment "%s" not found.`, path, segment)
		}

		// Check the field type
		kind, err := reflections.GetFieldKind(context, fieldName)
		if err != nil {
			return nil, fmt.Errorf(`Unable to resolve path "%s": %s.`, path, err)
		}

		switch kind {
		case reflect.Struct:
			if segment.Kind == IndexSegment || segment.Kind == RangeSegment {
				return nil, fmt.Errorf(`Unable to resolve path "%s": Segment "%s" does not refer to a collection.`, path, segment)
			}
			nested, err := reflections.GetField(context, fieldName)
			if err != nil {
				return nil, fmt.Errorf(`Unable to resolve path "%s": %s.`, path, err)
			}
			context = nested

		case reflect.Slice:
			switch segment.Kind {

			case PlainSegment:
				return nil, fmt.Errorf(`Unable to resolve path "%s": Missig index given for collection "%s".`, path, segment)

			case IndexSegment:
				slice, err := reflections.GetField(context, fieldName)
				if err != nil {
					return nil, fmt.Errorf(`Unable to resolve path "%s": %s.`, path, err)
				}
				sliceValue := reflect.ValueOf(slice)

				if segment.Index.Kind == NumericIndex {
					var index = segment.Index.Value.(int)
					if index < 0 || index >= sliceValue.Len() {
						return nil, fmt.Errorf(`Unable to resolve path "%s": Index in segment "%s" is out of bounds.`, path, segment)
					}
					context = sliceValue.Index(index).Interface()

				} else if segment.Index.Kind == AlphaNumericIndex {
					var indexFound = false
					for i := 0; i < sliceValue.Len(); i++ {
						element := sliceValue.Index(i).Interface()
						if exists, _ := reflections.HasField(element, "Name"); exists {
							name, _ := reflections.GetField(element, "Name")
							if segment.Index.Value == name {
								indexFound = true
								context = element
								break
							}
						}
					}
					if !indexFound {
						return nil, fmt.Errorf(`Unable to resolve path "%s": Named index in segment "%s" not found.`, path, segment)
					}
				}

			case RangeSegment:
				return nil, fmt.Errorf(`Unable to resolve path "%s": Range in segment "%s" not supported.`, path, segment)
			}

		default:
			if segment.Kind != PlainSegment {
				return nil, fmt.Errorf(`Unable to resolve path "%s": Segment "%s" does refer to a collection.`, path, segment)
			}
			nested, err := reflections.GetField(context, fieldName)
			if err != nil {
				return nil, fmt.Errorf(`Unable to resolve path "%s": %s.`, path, err)
			}
			context = nested
		}

		if context == nil {
			return nil, fmt.Errorf(`Unable to resolve path "%s": Segment "%s" not found.`, path, segment)
		}
	}
	return context, nil
}

// Append the specified path to this path and return the result as a new path
func (self Path) Append(path Path) Path {
	var result Path
	if len(self) > 0 {
		copy(result, self)
	}
	for _, segment := range path {
		result = append(result, segment)
	}
	return result
}

func (path Path) IsEmpty() bool {
	return len(path) == 0
}

func (path Path) String() string {
	if len(path) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	for idx, segment := range path {
		buffer.WriteString(fmt.Sprint(segment))
		if idx < len(path)-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func (segment Segment) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(segment.Name)

	switch segment.Kind {
	case RangeSegment:
		buffer.WriteString("[")
		if segment.Range.From != Undefined {
			buffer.WriteString(fmt.Sprintf("%d", segment.Range.From))
		}
		buffer.WriteString(":")
		if segment.Range.To != Undefined {
			buffer.WriteString(fmt.Sprintf("%d", segment.Range.To))
		}
		buffer.WriteString("]")
	case IndexSegment:
		buffer.WriteString(fmt.Sprintf("[%v]", segment.Index.Value))
	}
	return buffer.String()
}
