package path

import . "gopkg.in/check.v1"

// ------------------------------------------------------ setup

type PathParseSuite struct {
	emptyIndex Index
	emptyRange Range
}

func (s *PathParseSuite) SetUpSuite(_ *C) {
	s.emptyIndex = Index{}
	s.emptyRange = Range{Undefined, Undefined}
}

var _ = Suite(&PathParseSuite{})

// ------------------------------------------------------ parse tests

func (s *PathParseSuite) TestParseEmptyPath(c *C) {
	path, err := Parse("")
	assertPath(c, path, err, 0)
}

func (s *PathParseSuite) TestParseSegment(c *C) {
	path, err := Parse("foo")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", PlainSegment, s.emptyIndex, s.emptyRange)
}

func (s *PathParseSuite) TestParseSegments(c *C) {
	path, err := Parse("foo.bar")
	assertPath(c, path, err, 2)
	assertSegment(c, path[0], "foo", PlainSegment, s.emptyIndex, s.emptyRange)
	assertSegment(c, path[1], "bar", PlainSegment, s.emptyIndex, s.emptyRange)
}

func (s *PathParseSuite) TestParseNumericIndex(c *C) {
	path, err := Parse("foo[42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", IndexSegment, Index{NumericIndex, 42}, s.emptyRange)
}

func (s *PathParseSuite) TestParseAlphaNumericIndex(c *C) {
	path, err := Parse("foo[bar]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", IndexSegment, Index{AlphaNumericIndex, "bar"}, s.emptyRange)
}

func (s *PathParseSuite) TestParseSliceRangeFrom(c *C) {
	path, err := Parse("foo[42:]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, s.emptyIndex, Range{42, Undefined})
}

func (s *PathParseSuite) TestParseSliceRangeTo(c *C) {
	path, err := Parse("foo[:42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, s.emptyIndex, Range{Undefined, 42})
}

func (s *PathParseSuite) TestParseSliceRangeFromTo(c *C) {
	path, err := Parse("foo[23:42]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, s.emptyIndex, Range{23, 42})
}

func (s *PathParseSuite) TestParseSliceRangeAll(c *C) {
	path, err := Parse("foo[:]")
	assertPath(c, path, err, 1)
	assertSegment(c, path[0], "foo", RangeSegment, s.emptyIndex, Range{Undefined, Undefined})
}

func (s *PathParseSuite) TestParseMixed(c *C) {
	path, err := Parse("a[0].b[z].c[1:].d[:2].e[3:4].f[:].g")
	assertPath(c, path, err, 7)
	assertSegment(c, path[0], "a", IndexSegment, Index{NumericIndex, 0}, s.emptyRange)
	assertSegment(c, path[1], "b", IndexSegment, Index{AlphaNumericIndex, "z"}, s.emptyRange)
	assertSegment(c, path[2], "c", RangeSegment, s.emptyIndex, Range{1, Undefined})
	assertSegment(c, path[3], "d", RangeSegment, s.emptyIndex, Range{Undefined, 2})
	assertSegment(c, path[4], "e", RangeSegment, s.emptyIndex, Range{3, 4})
	assertSegment(c, path[5], "f", RangeSegment, s.emptyIndex, Range{Undefined, Undefined})
	assertSegment(c, path[6], "g", PlainSegment, s.emptyIndex, s.emptyRange)
}

// ------------------------------------------------------ error tests

func (s *PathParseSuite) TestParseMalformed(c *C) {
	path, err := Parse("foo[bar.")
	c.Assert(path, IsNil)
	c.Assert(err, NotNil)
}
