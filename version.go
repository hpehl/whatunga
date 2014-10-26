package main

const (
	WILDFLY = "WildFly"
	EAP     = "EAP"
)

type Version struct {
	Product string
	Version string
}

func (v Version) String() string {
	return v.Product + ":" + v.Version
}

var SupportedVersions = []Version{
	{WILDFLY, "8.0"},
	{WILDFLY, "8.1"},
	{EAP, "6.3"},
}

// the xmlns version of the config files
var ModelVersions = map[string]Version{
	"2.0": {WILDFLY, "8.0"},
	"2.1": {WILDFLY, "8.1"},
	"1.6": {EAP, "6.3"},
}
