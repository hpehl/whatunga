package path

import (
	"github.com/hpehl/whatunga/model"
	. "gopkg.in/check.v1"
)

// ------------------------------------------------------ setup

type PathResolveSuite struct {
	project *model.Project
}

func (s *PathResolveSuite) SetUpSuite(_ *C) {
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
				Name:          "serverGroup0",
				Profile:       "profile0",
				SocketBinding: "socketBinding0",
				Jvm: model.Jvm{
					Name: "serverGroup0-jvm",
					Heap: model.BoundedMemory{
						Initial: "1GB",
						Max:     "2GB",
					},
				},
				Deployments: []model.Deployment{
					model.Deployment{"deployment0", "deployment0-rt", "/path/to/deployment0"},
					model.Deployment{"deployment1", "deployment1-rt", "/path/to/deployment1"},
				},
			},
			model.ServerGroup{
				Name:          "serverGroup0",
				Profile:       "profile0",
				SocketBinding: "socketBinding0",
			},
		},
		Hosts: []model.Host{
			model.Host{
				Name: "host0",
				DC:   true,
				Servers: []model.Server{
					model.Server{
						Name:        "host0-server0",
						ServerGroup: "serverGroup0",
					},
					model.Server{
						Name:        "host0-server1",
						ServerGroup: "serverGroup0",
						PortOffset:  50,
						AutoStart:   true,
					},
					model.Server{
						Name:        "host0-server2",
						ServerGroup: "serverGroup0",
						PortOffset:  100,
						Jvm: model.Jvm{
							Name: "host0-server2-jvm",
							Heap: model.BoundedMemory{
								Initial: "1GB",
								Max:     "2GB",
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
						Name:        "host1-server1",
						ServerGroup: "serverGroup1",
					},
					model.Server{
						Name:        "host1-server1",
						ServerGroup: "serverGroup1",
						PortOffset:  50,
					},
					model.Server{
						Name:        "host1-server2",
						ServerGroup: "serverGroup1",
						PortOffset:  100,
					},
				},
				Jvm: model.Jvm{
					Name: "host1-jvm",
					Heap: model.BoundedMemory{
						Initial: "1GB",
						Max:     "2GB",
					},
					PermGem: "128MB",
					Stack:   "256MB",
					Options: []string{"-server"},
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

func (s *PathResolveSuite) TestResolveName(c *C) {
	path, _ := Parse("name")
	value, err := path.Resolve(s.project)

	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, "test")
}
