package path

import (
	. "gopkg.in/check.v1"
	"testing"
)

// ------------------------------------------------------ setup

// triggers all tests in this package
func TestPath(t *testing.T) { TestingT(t) }

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
