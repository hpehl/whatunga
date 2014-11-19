package path

import (
	"fmt"
	. "gopkg.in/check.v1"
)

// ------------------------------------------------------ setup

type PathMiscSuite struct{}

var _ = Suite(&PathMiscSuite{})

// ------------------------------------------------------ misc tests

func (s *PathMiscSuite) TestString(c *C) {
	in := "a[0].b[z].c[1:].d[:2].e[3:4].f[:].g"
	path, _ := Parse(in)
	out := fmt.Sprint(path)
	c.Assert(out, Equals, in)
}
