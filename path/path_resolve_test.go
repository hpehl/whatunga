package path

import (
	"github.com/hpehl/whatunga/model"
	. "gopkg.in/check.v1"
	"reflect"
)

// ------------------------------------------------------ setup

type PathResolveSuite struct {
	jvm     *model.Jvm
	project *model.Project
}

func (s *PathResolveSuite) SetUpSuite(_ *C) {
	s.jvm = &model.Jvm{
		Name: "server-group0-jvm",
		Heap: model.BoundedMemory{
			Initial: "1GB",
			Max:     "2GB",
		},
	}

	s.project = &model.Project{
		Name:    "test",
		Version: "1.0",
		Config: model.Config{
			Templates: model.Templates{
				Domain:     "templates/domain.xml",
				HostMaster: "templates/host-master.xml",
				HostSlave:  "templates/host-slave.xml",
			},
			ConsoleUser: model.User{
				Name:     "admin",
				Password: "passw0rd_",
			},
			DomainUser: model.User{
				Name:     "dc",
				Password: "passw0rd_",
			},
			DockerRemoteAPI: "unix:///var/run/docker.sock",
		},
		ServerGroups: []model.ServerGroup{
			model.ServerGroup{
				Name:          "server-group0",
				Profile:       "profile0",
				SocketBinding: "socket-binding0",
				Jvm:           s.jvm,
				Deployments: []model.Deployment{
					model.Deployment{"deployment0", "deployment0-rt", "/path/to/deployment0"},
					model.Deployment{"deployment1", "deployment1-rt", "/path/to/deployment1"},
				},
			},
			model.ServerGroup{
				Name:          "server-group1",
				Profile:       "profile1",
				SocketBinding: "socket-binding1",
			},
		},
		Hosts: []model.Host{
			model.Host{
				Name: "host0",
				DC:   true,
				Servers: []model.Server{
					model.Server{
						Name:        "host0-server0",
						ServerGroup: "server-group0",
					},
					model.Server{
						Name:        "host0-server1",
						ServerGroup: "server-group0",
						PortOffset:  50,
						AutoStart:   true,
					},
					model.Server{
						Name:        "host0-server2",
						ServerGroup: "server-group0",
						PortOffset:  100,
						Jvm: &model.Jvm{
							Name: "host0-server2-jvm",
							Heap: model.BoundedMemory{
								Initial: "3GB",
								Max:     "4GB",
							},
						},
					},
				},
			},
			model.Host{
				Name: "host1",
				DC:   false,
				Servers: []model.Server{
					model.Server{
						Name:        "host1-server0",
						ServerGroup: "server-group1",
					},
					model.Server{
						Name:        "host1-server1",
						ServerGroup: "server-group1",
						PortOffset:  50,
					},
					model.Server{
						Name:        "host1-server2",
						ServerGroup: "server-group1",
						PortOffset:  100,
					},
				},
				Jvm: &model.Jvm{
					Name: "host1-jvm",
					Heap: model.BoundedMemory{
						Initial: "5GB",
						Max:     "6GB",
					},
					PermGem: "128MB",
					Stack:   "256MB",
					Options: []string{"-server"},
				},
			},
			model.Host{
				Name: "host2",
				Servers: []model.Server{
					model.Server{Name: "foo", PortOffset: 1},
					model.Server{Name: "foo", PortOffset: 2},
					model.Server{Name: "foo", PortOffset: 3},
				},
			},
		},
		Users: []model.User{
			model.User{
				Name:     "user0",
				Password: "password0",
			},
			model.User{
				Name:     "user1",
				Password: "password1",
			},
			model.User{
				Name:     "user2",
				Password: "password2",
			},
		},
	}
}

var _ = Suite(&PathResolveSuite{})

// ------------------------------------------------------ resolve tests

func (s *PathResolveSuite) TestResolveSimple(c *C) {
	path, _ := Parse("name")
	value, err := path.Resolve(s.project)

	assertField(c, value, err, "test")
}

func (s *PathResolveSuite) TestResolveNested(c *C) {
	path, _ := Parse("config.domain-user.username")
	value, err := path.Resolve(s.project)

	assertField(c, value, err, "dc")
}

func (s *PathResolveSuite) TestResolveServerGroup0Jvm(c *C) {
	path, _ := Parse("server-groups[server-group0].jvm")
	value, err := path.Resolve(s.project)

	assertField(c, value, err, nil)
	c.Assert(reflect.DeepEqual(value, s.jvm), Equals, true)
}

func (s *PathResolveSuite) TestResolveNumericIndex(c *C) {
	path, _ := Parse("hosts[0].servers[2].jvm.heap.max")
	value, err := path.Resolve(s.project)

	assertField(c, value, err, "4GB")
}

func (s *PathResolveSuite) TestResolveAlphaNumericIndex(c *C) {
	path, _ := Parse("hosts[1].servers[host1-server2].port-offset")
	value, err := path.Resolve(s.project)

	assertField(c, value, err, 100)
}

// ------------------------------------------------------ error tests

func (s *PathResolveSuite) TestResolveUnknown(c *C) {
	emptyProject := &model.Project{}
	path, _ := Parse("foo")
	value, err := path.Resolve(emptyProject)

	expectError(c, value, err, `Unable to resolve path "foo": Segment "foo" not found.`)
}

func (s *PathResolveSuite) TestResolveWrongKind1(c *C) {
	emptyProject := &model.Project{}
	path, _ := Parse("config[0].templates.-domain")
	value, err := path.Resolve(emptyProject)

	expectError(c, value, err, `Unable to resolve path "config[0].templates.-domain": Segment "config[0]" does not refer to a collection.`)
}

func (s *PathResolveSuite) TestResolveWrongKind2(c *C) {
	emptyProject := &model.Project{}
	path, _ := Parse("server-groups.name")
	value, err := path.Resolve(emptyProject)

	expectError(c, value, err, `Unable to resolve path "server-groups.name": Missig index given for collection "server-groups".`)
}

func (s *PathResolveSuite) TestResolveRange(c *C) {
	emptyProject := &model.Project{}
	path, _ := Parse("server-groups[:].name")
	value, err := path.Resolve(emptyProject)

	expectError(c, value, err, `Unable to resolve path "server-groups[:].name": Range in segment "server-groups[:]" not supported.`)
}

func (s *PathResolveSuite) TestResolveInvalidIndex1(c *C) {
	emptyProject := &model.Project{}
	path, _ := Parse("hosts[0]")
	value, err := path.Resolve(emptyProject)

	expectError(c, value, err, `Unable to resolve path "hosts[0]": Index in segment "hosts[0]" is out of bounds.`)
}

func (s *PathResolveSuite) TestResolveInvalidIndex2(c *C) {
	path, _ := Parse("hosts[0].servers[100].name")
	value, err := path.Resolve(s.project)

	expectError(c, value, err, `Unable to resolve path "hosts[0].servers[100].name": Index in segment "servers[100]" is out of bounds.`)
}

func (s *PathResolveSuite) TestResolveInvalidIndex3(c *C) {
	path, _ := Parse("server-groups[foo].deployments[10]")
	value, err := path.Resolve(s.project)

	expectError(c, value, err, `Unable to resolve path "server-groups[foo].deployments[10]": Named index in segment "server-groups[foo]" not found.`)
}

// ------------------------------------------------------ helper functions

func assertField(c *C, value interface{}, err error, expected interface{}) {
	if err != nil {
		c.Error(err)
	}
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	if expected != nil {
		c.Assert(value, Equals, expected)
	}
}

func expectError(c *C, value interface{}, err error, why string) {
	c.Assert(value, IsNil)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, why)
}
