package path

import (
	"fmt"
	. "gopkg.in/check.v1"
)

// ------------------------------------------------------ setup

type PathMiscSuite struct{}

var _ = Suite(&PathMiscSuite{})

// ------------------------------------------------------ misc tests

func (s *PathMiscSuite) TestSplitLastSegment(c *C) {
	var pth, segment string

	pth, segment = SplitLastSegment("")
	c.Assert(pth, Equals, "")
	c.Assert(segment, Equals, "")

	pth, segment = SplitLastSegment(".")
	c.Assert(pth, Equals, "")
	c.Assert(segment, Equals, "")

	pth, segment = SplitLastSegment(".foo")
	c.Assert(pth, Equals, "")
	c.Assert(segment, Equals, "foo")

	pth, segment = SplitLastSegment("foo.")
	c.Assert(pth, Equals, "foo")
	c.Assert(segment, Equals, "")

	pth, segment = SplitLastSegment("foo.bar")
	c.Assert(pth, Equals, "foo")
	c.Assert(segment, Equals, "bar")
}

func (s *PathMiscSuite) TestLastOpenSquareBracket(c *C) {
	var hasIndex bool
	var pth, index string

	hasIndex, pth, index = LastOpenSquareBracket("")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo[")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo[1")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo")
	c.Assert(index, Equals, "1")

	hasIndex, pth, index = LastOpenSquareBracket("foo[bar")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo")
	c.Assert(index, Equals, "bar")

	hasIndex, pth, index = LastOpenSquareBracket("foo[1].bar")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo[1].bar[")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo[1].bar")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo[1].bar[2")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo[1].bar")
	c.Assert(index, Equals, "2")

	hasIndex, pth, index = LastOpenSquareBracket("foo[1].bar[meep")
	c.Assert(hasIndex, Equals, true)
	c.Assert(pth, Equals, "foo[1].bar")
	c.Assert(index, Equals, "meep")

	hasIndex, pth, index = LastOpenSquareBracket("[foo]")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("]")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")

	hasIndex, pth, index = LastOpenSquareBracket("foo]")
	c.Assert(hasIndex, Equals, false)
	c.Assert(pth, Equals, "")
	c.Assert(index, Equals, "")
}

func (s *PathMiscSuite) TestAppendToEmpty(c *C) {
	var empty Path
	foo, _ := Parse("foo")
	combined := empty.Append(foo)

	c.Assert(combined.String(), Equals, "foo")
}

func (s *PathMiscSuite) TestAppend(c *C) {
	foo, _ := Parse("foo")
	bar, _ := Parse("bar")
	foobar := foo.Append(bar)

	c.Assert(foobar.String(), Equals, "foo.bar")
}

func (s *PathMiscSuite) TestString(c *C) {
	in := "a[0].b[z].c[1:].d[:2].e[3:4].f[:].g"
	path, _ := Parse(in)
	out := fmt.Sprint(path)
	c.Assert(out, Equals, in)
}
