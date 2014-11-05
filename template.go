package main

var Predefined = map[Target]Templates{
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.0              */
	/*                           */
	/* ------------------------- */
	Target{WildFly, "8.0"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.1              */
	/*                           */
	/* ------------------------- */
	Target{WildFly, "8.1"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  EAP 6.3                  */
	/*                           */
	/* ------------------------- */
	Target{EAP, "6.3"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
}
