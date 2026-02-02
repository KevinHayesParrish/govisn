// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
)

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
	}
	Addresses struct {
		NetworkAddresses struct {
			IPAddress []string
		}
		MediaAddresses struct {
			MediaAddress []string
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
		ipRouteIfIndex int
		ipRouteMetric1 int
		ipRouteMetric2 int
		ipRouteMetric3 int
		ipRouteMetric4 int
		ipRouteNextHop string
		ipRouteType    int
		ipRouteProto   string
		ipRouteAge     int
		ipRouteMask    string
		ipRouteMetric5 int
		ipRouteInfo    string
	}
}

// Link is the structure representing a network link between two routers
type Link struct {
	LinkID            int
	FromRouterID      int
	FromRouterName    string
	FromRouterIP      string
	FromRouterIfIndex string
	ToRouterID        int
	ToRouterName      string
	ToRouterIP        string
}

// ScannedRouter is the structure representing an SNMP capable router discovered on the network.
type ScannedRouter struct {
	Name      string
	IPAddress string
}

// FileSelect struct
type FileSelect struct {
	gui.Panel
	path *gui.Label
	list *gui.List
	bok  *gui.Button
	bcan *gui.Button
}

// ErrorDialog struct
type ErrorDialog struct {
	gui.Panel
	msg *gui.ImageLabel
	bok *gui.Button
}
type gvapp struct {
	*app.Application                // Embedded application object
	fs               *FileSelect    // File selection dialog
	ed               *ErrorDialog   // Error dialog
	axes             *helper.Axes   // Axis helper
	grid             *helper.Grid   // Grid helper
	viewAxes         bool           // Axis helper visible flag
	viewGrid         bool           // Grid helper visible flag
	camPos           math32.Vector3 // Initial camera position
	models           []*core.Node   // Models being shown
	scene            *core.Node
	cam              *camera.Camera
	orbit            *camera.OrbitControl
}

// LinkUpdate contains link visualization data to be applied in the main thread
type LinkUpdate struct {
	LinkID   int
	Color    math32.Color
	Width    float32
	FromName string
	ToName   string
}
