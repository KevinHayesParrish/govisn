package main

// Router is the structure representing a network router
type Router struct {
	System struct {
		RouterID    int
		Name        string
		Description string
		//		UpTime      string
		UpTime   uint32
		Contact  string
		Location string
		Services int
		GPS      struct {
			Latitude  string
			Longitude string
			Altitude  string
		}
		//		Coordinates struct {
		//			X float32
		//			Y float32
		//			Z float32
		//		}
	}
	Addresses struct {
		NetworkAddresses struct {
			IPAddress []string
		}
		MediaAddresses struct {
			MediaAddress string
		}
	}
	Neighbors struct {
		Neighbor []struct {
			DestinationAddress string
			NextHop            string
		}
	}
}

type ifTable struct {
	ifEntry struct {
		ifIndexOID    string
		ifIndexType   byte
		ifIndex       int
		ifIndexLogger string

		ifDescrOID    string
		ifDescrType   byte
		ifDescr       string
		ifDescrLogger string

		ifTypeOID    string
		ifTypeType   byte
		ifType       string
		ifTypeLogger string

		ifMtuOID    string
		ifMtuType   byte
		ifMtu       int
		ifMTULogger string

		ifSpeedOID    string
		ifSpeedType   byte
		ifSpeed       uint
		ifSpeedLogger string

		ifPhysAddressOID    string
		ifPhysAddressType   byte
		ifPhysAddress       string
		ifPhysAddressLogger string

		ifAdminStatusOID    string
		ifAdminStatusType   byte
		ifAdminStatus       string
		ifAdminStatusLogger string

		ifOperStatusOID    string
		ifOperStatusType   byte
		ifOperStatus       string
		ifOperStatusLogger string

		ifLastChangeOID    string
		ifLastChangeType   byte
		ifLastChange       uint32
		ifLastChangeLogger string

		ifInOctetsOID    string
		ifInOctetsType   byte
		ifInOctets       uint
		ifInOctetsLogger string

		ifInUcastPktsOID    string
		ifInUcastPktsType   byte
		ifInUcastPkts       uint
		ifInUcastPktsLogger string

		ifInNUcastPktsOID    string // deprecated
		ifInNUcastPktsType   byte   // deprecated
		ifInNUcastPkts       uint   // deprecated
		ifInNUcastPktsLogger string

		ifInDiscardsOID    string
		ifInDiscardsType   byte
		ifInDiscards       uint
		ifInDiscardsLogger string

		ifInErrorsOID    string
		ifInErrorsType   byte
		ifInErrors       uint
		ifInErrorsLogger string

		ifInUnknownProtosOID    string
		ifInUnknownProtosType   byte
		ifInUnknownProtos       uint
		ifInUnknownProtosLogger string

		ifOutOctetsOID    string
		ifOutOctetsType   byte
		ifOutOctets       uint
		ifOutOctetsLogger string

		ifOutUcastPktsOID    string
		ifOutUcastPktsType   byte
		ifOutUcastPkts       uint
		ifOutUcastPktsLogger string

		ifOutNUcastPktsOID    string // deprecated
		ifOutNUcastPktsType   byte   // deprecated
		ifOutNUcastPkts       uint   //deprecated
		ifOutNUcastPktsLogger string

		ifOutDiscardsOID    string
		ifOutDiscardsType   byte
		ifOutDiscards       uint
		ifOutDiscardsLogger string

		ifOutErrorsOID    string
		ifOutErrorsType   byte
		ifOutErrors       uint
		ifOutErrorsLogger string

		ifOutQLenOID    string
		ifOutQLenType   byte
		ifOutQLen       uint // deprecated
		ifOutQLenLogger string

		ifSpecificOID    string
		ifSpecificType   byte
		ifSpecific       string // deprecated
		ifSpecificLogger string
	}
}

type ipAddrTable struct {
	ipAddrEntry struct {
		ipAdEntAddr         string
		ipAdEntIfIndex      int
		ipAdEntNetMask      string
		ipAdEntBcastAddr    int
		ipAdEntReasmMaxSize int
	}
}

type ipRouteTable struct {
	ipRouteEntry struct {
		ipRouteDest    string
		ipRouteIfIndex int32
		ipRouteMetric1 int32
		ipRouteMetric2 int32
		ipRouteMetric3 int32
		ipRouteMetric4 int32
		ipRouteNextHop string
		ipRouteType    string
		ipRouteProto   string
		ipRouteAge     int32
		ipRouteMask    string
		ipRouteMetric5 int32
		ipRouteInfo    string
	}
}

// Link is the structure representing a network link between two routers
type Link struct {
	LinkID int
	//FromRouter     string
	//ToRouter       string
	FromRouterName string
	FromRouterIP   string
	ToRouterName   string
	ToRouterIP     string
	//	RouterName      string
	//	DestinationName string
	//	DestinationIP   string
	//	NextHopName     string
	//	NextHopIP       string
}
