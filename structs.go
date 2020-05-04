package main

import "math/big"

// Router is the structure representing a network router
type Router struct {
	System struct {
		RouterID    int
		Name        string
		Description string
		UpTime      string
		Contact     string
		Location    string
		Services    *big.Int
		GPS         struct {
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

type ipAddrTable struct {
	ipAddrEntry struct {
		ipAdEntAddr         string
		ipAdEntIfIndex      int32
		ipAdEntNetMask      string
		ipAdEntBcastAddr    int32
		ipAdEntReasmMaxSize int32
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
	//	FromRouter string
	//	ToRouter   string
	//	FromRouterName string
	//	FromRouterIP   string
	//	ToRouterName   string
	//	ToRouterIP     string
	RouterName      string
	DestinationName string
	DestinationIP   string
	NextHopName     string
	NextHopIP       string
}
