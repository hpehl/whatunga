package path

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

// ------------------------------------------------------ setup

func TestPath(t *testing.T) { TestingT(t) }

type PathSuite struct{}

var emptyIndex = Index{}
var emptyRange = Range{Undefined, Undefined}
var _ = Suite(&PathSuite{})

// ------------------------------------------------------ parse tests

func (s *PathSuite) TestParseEmptyPath(c *C) {
	path, err := Parse("")
	assertPath(c, path, err, 0)
}

func (s *PathSuite) TestParseSegment(c *C) {
	path, err := Parse("foo")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", PlainSegment, emptyIndex, emptyRange)
}

func (s *PathSuite) TestParseSegments(c *C) {
	path, err := Parse("foo.bar")
	assertPath(c, path, err, 2)
	assertSegment(c, path[0], "foo", PlainSegment, emptyIndex, emptyRange)
	assertSegment(c, path[1], "bar", PlainSegment, emptyIndex, emptyRange)
}

func (s *PathSuite) TestParseNumericIndex(c *C) {
	path, err := Parse("foo[42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", IndexSegment, Index{NumericIndex, 42}, emptyRange)
}

func (s *PathSuite) TestParseAlphaNumericIndex(c *C) {
	path, err := Parse("foo[bar]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", IndexSegment, Index{AlphaNumericIndex, "bar"}, emptyRange)
}

func (s *PathSuite) TestParseSliceRangeFrom(c *C) {
	path, err := Parse("foo[42:]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, emptyIndex, Range{42, Undefined})
}

func (s *PathSuite) TestParseSliceRangeTo(c *C) {
	path, err := Parse("foo[:42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, emptyIndex, Range{Undefined, 42})
}

func (s *PathSuite) TestParseSliceRangeFromTo(c *C) {
	path, err := Parse("foo[23:42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, emptyIndex, Range{23, 42})
}

func (s *PathSuite) TestParseSliceRangeAll(c *C) {
	path, err := Parse("foo[:]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, emptyIndex, Range{Undefined, Undefined})
}

func (s *PathSuite) TestParseMixed(c *C) {
	path, err := Parse("a[0].b[z].c[1:].d[:2].e[3:4].f[:].g")
	assertPath(c, path, err, 7)
	assertSegment(c, path[0], "a", IndexSegment, Index{NumericIndex, 0}, emptyRange)
	assertSegment(c, path[1], "b", IndexSegment, Index{AlphaNumericIndex, "z"}, emptyRange)
	assertSegment(c, path[2], "c", RangeSegment, emptyIndex, Range{1, Undefined})
	assertSegment(c, path[3], "d", RangeSegment, emptyIndex, Range{Undefined, 2})
	assertSegment(c, path[4], "e", RangeSegment, emptyIndex, Range{3, 4})
	assertSegment(c, path[5], "f", RangeSegment, emptyIndex, Range{Undefined, Undefined})
	assertSegment(c, path[6], "g", PlainSegment, emptyIndex, emptyRange)
}

// ------------------------------------------------------ error tests

func (s *PathSuite) TestParseMalformed(c *C) {
	path, err := Parse("foo[bar.")
	c.Assert(path, IsNil)
	c.Assert(err, NotNil)
}

// ------------------------------------------------------ misc tests

func (s *PathSuite) TestString(c *C) {
	in := "a[0].b[z].c[1:].d[:2].e[3:4].f[:].g"
	path, _ := Parse(in)
	out := fmt.Sprint(path)
	c.Assert(out, Equals, in)
}

// ------------------------------------------------------ helper functions

func assertPath(c *C, path Path, err error, length int) {
	if err != nil {
		c.Error(err)
	}

	c.Assert(path, NotNil)
	c.Assert(err, IsNil)
	c.Assert(len(path), Equals, length)
}

func assertSegment(c *C, segment Segment, name string, segmentKind SegmentKind, index Index, rng Range) {
	c.Assert(segment.Name, Equals, name)
	c.Assert(segment.Kind, Equals, segmentKind)
	c.Assert(segment.Index, Equals, index)
	c.Assert(segment.Range, Equals, rng)
}
