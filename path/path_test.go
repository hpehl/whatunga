package path

import (
	. "gopkg.in/check.v1"
	"testing"
)

// ------------------------------------------------------ setup

func TestPath(t *testing.T) { TestingT(t) }

type PathSuite struct{}

var _ = Suite(&PathSuite{})

// ------------------------------------------------------ parse tests

func (s *PathSuite) TestParseEmptyPath(c *C) {
	path, err := Parse("")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "", EmptyIndex, EmptyRange)
}

func (s *PathSuite) TestParseSegment(c *C) {
	path, err := Parse("foo")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", EmptyIndex, EmptyRange)
}

func (s *PathSuite) TestParseSegments(c *C) {
	path, err := Parse("foo.bar")
	assertValidPath(c, path, err, 2)
	assertSegment(c, path[0], "foo", EmptyIndex, EmptyRange)
	assertSegment(c, path[1], "bar", EmptyIndex, EmptyRange)
}

func (s *PathSuite) TestParseNumericIndex(c *C) {
	path, err := Parse("foo[42]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", Index{Numeric, 42}, EmptyRange)
}

func (s *PathSuite) TestParseAlphaNumericIndex(c *C) {
	path, err := Parse("foo[bar]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", Index{AlphaNumeric, "bar"}, EmptyRange)
}

func (s *PathSuite) TestParseSliceRangeFrom(c *C) {
	path, err := Parse("foo[42:]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", EmptyIndex, Range{42, Undefined})
}

func (s *PathSuite) TestParseSliceRangeTo(c *C) {
	path, err := Parse("foo[:42]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", EmptyIndex, Range{Undefined, 42})
}

func (s *PathSuite) TestParseSliceRangeFromTo(c *C) {
	path, err := Parse("foo[23:42]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", EmptyIndex, Range{23, 42})
}

func (s *PathSuite) TestParseSliceRangeAll(c *C) {
	path, err := Parse("foo[:]")
	assertValidPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", EmptyIndex, Range{Undefined, Undefined})
}

func (s *PathSuite) TestParseMixed(c *C) {
	path, err := Parse("a[0].b[z].c[1:].d[:2].e[3:4].f[:].g")
	assertValidPath(c, path, err, 7)
	assertSegment(c, path[0], "a", Index{Numeric, 0}, EmptyRange)
	assertSegment(c, path[1], "b", Index{AlphaNumeric, "z"}, EmptyRange)
	assertSegment(c, path[2], "c", EmptyIndex, Range{1, Undefined})
	assertSegment(c, path[3], "d", EmptyIndex, Range{Undefined, 2})
	assertSegment(c, path[4], "e", EmptyIndex, Range{3, 4})
	assertSegment(c, path[5], "f", EmptyIndex, Range{Undefined, Undefined})
	assertSegment(c, path[6], "g", EmptyIndex, EmptyRange)
}

// ------------------------------------------------------ error tests

func (s *PathSuite) TestParseMalformed(c *C) {
	path, err := Parse("foo[bar.")
	c.Logf("path: %v, err: %v", path, err)
}

// ------------------------------------------------------ helper functions

func assertValidPath(c *C, path Path, err error, length int) {
	if err != nil {
		c.Error(err)
	}

	c.Assert(path, NotNil)
	c.Assert(err, IsNil)
	c.Assert(len(path), Equals, length)
}

func assertSegment(c *C, segment Segment, name string, index Index, rng Range) {
	c.Assert(segment.Name, Equals, name)
	c.Assert(segment.Index, Equals, index)
	c.Assert(segment.Range, Equals, rng)
}
