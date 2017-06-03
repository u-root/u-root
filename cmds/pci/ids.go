package main
var vendor = map[VID]Vendor {
0x1a22: Vendor{Name: "Ambric Inc.", }, // vendor
0x1a29: Vendor{Name: "Fortinet, Inc.", Devs: map[DID]Device {

	0x4338: Device{Name: " CP8 Content Processor ASIC",	}, // device 

	0x4e36: Device{Name: " NP6 Network Processor",	}, // device 
	}, // devices
}, // vendor
0x1a2b: Vendor{Name: "Ascom AG", Devs: map[DID]Device {

	0x0000: Device{Name: " GESP v1.2",	}, // device 

	0x0001: Device{Name: " GESP v1.3",	}, // device 

	0x0002: Device{Name: " ECOMP v1.3",	}, // device 

	0x0005: Device{Name: " ETP v1.4",	}, // device 

	0x000a: Device{Name: " ETP-104 v1.1",	}, // device 

	0x000e: Device{Name: " DSLP-104 v1.1",	}, // device 
	}, // devices
}, // vendor
} // table
