package main

var Predefined = map[ProductVersion]Templates{
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.0              */
	/*                           */
	/* ------------------------- */
	ProductVersion{WildFly, "8.0"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  WildFly 8.1              */
	/*                           */
	/* ------------------------- */
	ProductVersion{WildFly, "8.1"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
	/* ------------------------- */
	/*                           */
	/*  EAP 6.3                  */
	/*                           */
	/* ------------------------- */
	ProductVersion{EAP, "6.3"}: Templates{
		Domain:     `domain.xml...`,
		HostMaster: `host-master.xml...`,
		HostSlave:  `host-slave.xml...`,
	},
}
