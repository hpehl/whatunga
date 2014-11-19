package path

import (
	"fmt"
	. "gopkg.in/check.v1"
)

// ------------------------------------------------------ setup

type PathSuite struct{}

var _ = Suite(&PathSuite{})

// ------------------------------------------------------ misc tests

func (s *PathSuite) TestString(c *C) {
	in := "a[0].b[z].c[1:].d[:2].e[3:4].f[:].g"
	path, _ := Parse(in)
	out := fmt.Sprint(path)
	c.Assert(out, Equals, in)
}
