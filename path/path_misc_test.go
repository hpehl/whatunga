package path

import (
	"fmt"
	. "gopkg.in/check.v1"
)

// ------------------------------------------------------ setup

type PathMiscSuite struct{}

var _ = Suite(&PathMiscSuite{})

// ------------------------------------------------------ misc tests

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
