// Copyright 2020 Kevin Hayes Parrish. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"os"

	//"log"
	"strings"

	//"log"
	//"math"
	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"
)

// BUILD_LINKS_VERSION is the file version sequence number
const BUILD_LINKS_VERSION = "0.0.7"

/*
 * This function reads the database for Routers and their info, then build the links Table
 */
func buildLinks(log *logger.Logger, database *sql.DB) *sql.DB {

	log.Debug("func buildLinks version %s started", BUILD_LINKS_VERSION)

	routeTableRows, err := database.Query("SELECT RouterID, Name, DestAddr, IPRouteIfIndex, NextHop FROM Routers INNER JOIN RouteTable USING (RouterID)")
	if err != nil {
		log.Warn("buildLinks. databaseForRead JOIN error. %s", err.Error())
		log.Warn("No Routers Discovered. Check that all routers support SNMP with MIB II.")
		os.Exit(1)
		//		return database
	}
	defer routeTableRows.Close()

	var router Router
	var links []Link
	var link Link
	var RouterID int
	var Name string
	var DestAddr string
	var IPRouteIfIndex string
	var NextHop string
	for routeTableRows.Next() {
		routeTableRows.Scan(&RouterID, &Name, &DestAddr, &IPRouteIfIndex, &NextHop)
		router.System.RouterID = RouterID
		router.System.Name = Name

		//   Determine router interface using ipRouteIfIndex. This is the index of the interface. We can use this to get the interface IP address.

		link.FromRouterID = RouterID
		link.FromRouterName = Name // Current router
		link.FromRouterIfIndex = IPRouteIfIndex

		// Find FromRouterIP by DNS lookup by name
		fromIPs := getHostIP(Name)
		if len(fromIPs) < 1 {
			log.Warn("No Router IP Address the link from Router %s", Name)
			link.FromRouterIP = ""
		} else {
			link.FromRouterIP = fromIPs[0]
		}

		//! TODO: change from getRtrName function to getRtrIDByIP function,
		//! since some routers may not have a name or DNS entry,
		//! but should have an IP address. This will also allow us to build
		//! links to routers that are not directly connected to the current
		//! router, but are reachable through other routers. We can use the
		//! NextHop IP address to find the RouterID and Name of the next hop
		//! router in the database, and then use that information to build
		//! the link.
		/*--------------------------------------------------------------------------*/
		/* rtrNames := getRtrName(NextHop)                                          */
		/*                                                                          */
		/* if len(rtrNames) < 1 {                                                   */
		/*     log.Warn("No Router Name for Route Destination %s", NextHop)         */
		/*     link.ToRouterName = ""                                               */
		/* } else {                                                                 */
		/*     link.ToRouterName = rtrNames[0]                                      */
		/*       link.ToRouterID = getRtrIDByName(log, link.ToRouterName, database) */
		/*     link.ToRouterID = getRtrIDByIP(log, link.ToRouterIP, database)       */
		/* }                                                                        */
		/*                                                                          */
		/*--------------------------------------------------------------------------*/
		link.ToRouterID = getRtrIDByIP(log, NextHop, database)
		if link.ToRouterID == 0 {
			log.Warn("No Router ID for Route Destination %s. Link not added to database.", NextHop)
			continue
		}

		link.ToRouterName = getRtrNameByID(log, link.ToRouterID, database)
		link.ToRouterIP = NextHop

		// calculate LinkID
		link.LinkID = int(crc32.ChecksumIEEE([]byte(link.FromRouterIP + link.ToRouterIP)))

		if link.FromRouterName == link.ToRouterName {
			log.Info("From and To Routers are the same. Link not added to database.")
		} else {
			links = append(links, link)
		}

	}

	for i := 0; i < len(links); i++ {
		statement, err := database.Prepare("INSERT INTO Links (LinkID, FromRouterID, FromRouterName, FromRouterIP, FromRouterIfIndex, ToRouterID, ToRouterName, ToRouterIP) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			//log.Fatal("Links Insert Prepare err %v", err)
			log.Fatal(fmt.Sprintf("Links Insert Prepare err %v", err))
		}
		_, err = statement.Exec(links[i].LinkID, links[i].FromRouterID, links[i].FromRouterName, links[i].FromRouterIP, links[i].FromRouterIfIndex, links[i].ToRouterID, links[i].ToRouterName, links[i].ToRouterIP)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Info("Link already exists. Continue building links.")
			} else {
				//log.Fatal("Link INSERT error %v", err, "\n"+"Link Info: ", links[i])
				//log.Fatal(fmt.Sprintf("Link INSERT error %v \n- Link Info: %+v", err, links[i]))
				log.Warn(fmt.Sprintf("Link INSERT error %v \n- Link Info: %+v", err, links[i]))

			}
		}
		statement.Close()
	}

	log.Debug(fmt.Sprintf("func buildLinks version %s ending", BUILD_LINKS_VERSION))

	return database
}

/*-------------------------------------------------------------------------------------------------------------*/
/*                                                                                                             */
/*  * getRtrIDByName retrieves the RouteID from the database, given a Router Name                              */
/*                                                                                                             */
/* func getRtrIDByName(log *logger.Logger, Name string, database *sql.DB) int {                                */
/*        Retrive Router from the database                                                                     */
/*     var RouterID int                                                                                        */
/*       routerRows, queryErr := database.Query("SELECT RouterID, Name FROM Routers WHERE RouterID = ?", Name) */
/*     routerRows, queryErr := database.Query("SELECT RouterID, Name FROM Routers WHERE Name = ?", Name)       */
/*     if queryErr != nil {                                                                                    */
/*         log.Fatal("databaseForRead Query Router error %v", queryErr)                                        */
/*     }                                                                                                       */
/*     defer routerRows.Close()                                                                                */
/*     log.Debug("Successful Routers table Select")                                                            */
/*     for routerRows.Next() {                                                                                 */
/*         routerRows.Scan(&RouterID)                                                                          */
/*         return RouterID                                                                                     */
/*     }                                                                                                       */
/*     return RouterID                                                                                         */
/* }                                                                                                           */
/*-------------------------------------------------------------------------------------------------------------*/

/*
 * getRtrIDByIP retrieves the RouterID from the database, given a Router IP address
 */
func getRtrIDByIP(log *logger.Logger, IP string, database *sql.DB) int {
	// Retrive Router from the database
	var RouterID int
	routerRows, queryErr := database.Query("SELECT RouterID FROM RouterIp WHERE IpAddr = ?", IP)
	if queryErr != nil {
		log.Fatal("databaseForRead Query Router error %v", queryErr)
	}
	defer routerRows.Close()
	log.Debug("Successful RouterIp table Select")
	for routerRows.Next() {
		routerRows.Scan(&RouterID)
		return RouterID
	}

	return RouterID
}

/*
 * getRtrNameByID retrieves the Router Name from the database, given a RouterID
 */
func getRtrNameByID(log *logger.Logger, RouterID int, database *sql.DB) string {
	// Retrive Router from the database
	var Name string
	routerRows, queryErr := database.Query("SELECT Name FROM Routers WHERE RouterID = ?", RouterID)
	if queryErr != nil {
		log.Fatal("databaseForRead Query Router error %v", queryErr)
	}
	defer routerRows.Close()
	log.Debug("Successful Routers table Select")
	for routerRows.Next() {
		routerRows.Scan(&Name)
		return Name
	}

	return Name
}
