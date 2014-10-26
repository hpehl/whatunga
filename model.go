package main

type Project struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Config       Config        `json:"config"`
	ServerGroups []ServerGroup `json:"server-groups"`
	Hosts        []Host        `json:"hosts"`
	Users        []User        `json:"users"`
}

type Config struct {
	Templates       Templates `json:"templates"`
	ConsoleUser     User      `json:"console-user"`
	DomainUser      User      `json:"domain-user"`
	DockerRemoteAPI string    `json:"docker-remote-api"`
}

type Templates struct {
	Domain     string `json:"domain"`
	HostMaster string `json:"host-master`
	HostSlave  string `json:host-slave`
}

type ServerGroup struct {
	Name          string       `json:"name"`
	Profile       string       `json:"profile"`
	SocketBinding string       `json:"socket-binding"`
	Jvm           Jvm          `json:"jvm"`
	Deployments   []Deployment `json:"deployments"`
}

type Deployment struct {
	Name        string `json:"name"`
	RuntimeName string `json:"runtime-name"`
	Path        string `json:"path"`
}

type Host struct {
	Name    string   `json:"name"`
	DC      bool     `json:"domain-controller"`
	Servers []Server `json:"servers"`
	Jvm     Jvm      `json:"jvm"`
}

type Server struct {
	Name        string `json:"name"`
	ServerGroup string `json:"server-group"`
	PortOffset  uint   `json:"port-offset"`
	AutoStart   bool   `json:"auto-start"`
	Jvm         Jvm    `json:"jvm"`
}

type Memory string

type BoundedMemory struct {
	Initial Memory `json:"initial"`
	Max     Memory `json:"max"`
}

type Jvm struct {
	Name    string        `json:"name"`
	Heap    BoundedMemory `json:"heap"`
	PermGem Memory        `json:"perm-gen"`
	Stack   Memory        `json:"stack"`
	Options []string      `json:"options"`
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
