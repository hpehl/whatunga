package main

var Predefined = map[Version]Templates{
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.0              */
	/*                           */
	/* ------------------------- */
	Version{WILDFLY, "8.0"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.1              */
	/*                           */
	/* ------------------------- */
	Version{WILDFLY, "8.1"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  EAP 6.3                  */
	/*                           */
	/* ------------------------- */
	Version{EAP, "6.3"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
}
