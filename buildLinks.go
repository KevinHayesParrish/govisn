package main

import (
	"database/sql"
	"hash/crc32"

	//"log"
	"strings"

	//"log"
	//"math"
	"github.com/g3n/engine/util/logger"
	_ "github.com/mattn/go-sqlite3"
)

//buildLinksVersion is the file version sequence number
const buildLinksVersion = "0.0.4"

func buildLinks(debugFlag bool, log *logger.Logger, database *sql.DB) *sql.DB {
	//	fmt.Println("func buildLinks version", buildLinksVersion, "started")
	log.Debug("func buildLinks version %s", buildLinksVersion+" started")

	/* TODO
	 *
	 * Populate RouterName, DestinationName, DestinationIP, NextHopName and NextHopIP from RouteTable elements.
	 * calculate LinkID using CRC of RouterName/DestinationIP/NextHopIP
	 * Write Links Row to database
	 *
	 */

	routeTableRows, err := database.Query("SELECT RouterID, Name, DestAddr, IPRouteIfIndex, NextHop FROM Routers INNER JOIN RouteTable USING (RouterID)")
	if err != nil {
		//		log.Fatalln("databaseForRead JOIN error", err.Error())
		log.Fatal("databaseForRead JOIN error")
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
		link.FromRouterIP = fromIPs[0]
		if len(fromIPs) < 1 {
			//			fmt.Println("No Router IP Address the link from Router", Name)
			log.Warn("No Router IP Address the link from Router %s", Name)
			link.FromRouterIP = ""
		} else {
			link.FromRouterIP = fromIPs[0]
		}

		rtrNames := getRtrName(NextHop)
		if len(rtrNames) < 1 {
			//			fmt.Println("No Router Name for Route Destination", NextHop)
			log.Warn("No Router Name for Route Destination %s", NextHop)
			link.ToRouterName = ""
		} else {
			link.ToRouterName = rtrNames[0]
			// TODO
			//	1) query Router table where RouterName = link.ToRouterName
			//	2) link.ToRouterID = result.rtrIDs[0]
			rtrID := getRtrID(log, link.ToRouterName, database)
			link.ToRouterID = rtrID
		}

		link.ToRouterIP = NextHop

		// calculate LinkID
		link.LinkID = int(crc32.ChecksumIEEE([]byte(link.FromRouterIP + link.ToRouterIP)))

		if link.FromRouterName == link.ToRouterName {
			//			fmt.Println("From and To Routers are the same. Link not added to database.")
			log.Info("From and To Routers are the same. Link not added to database.")
		} else {
			links = append(links, link)
		}

	}
	routeTableRows.Close()

	for i := 0; i < len(links); i++ {
		//		statement, err := database.Prepare("INSERT INTO Links (LinkID, FromRouterName, FromRouterIP, ToRouterName, ToRouterIP) VALUES (?, ?, ?, ?, ?)")
		//		statement, err := database.Prepare("INSERT INTO Links (LinkID, FromRouterID, FromRouterName, FromRouterIP, ToRouterID, ToRouterName, ToRouterIP) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		statement, err := database.Prepare("INSERT INTO Links (LinkID, FromRouterID, FromRouterName, FromRouterIP, FromRouterIfIndex, ToRouterID, ToRouterName, ToRouterIP) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			//			log.Fatalln("Links Insert Prepare err:", err.Error())
			log.Fatal("Links Insert Prepare err %v", err)
		}
		//		_, err = statement.Exec(links[i].LinkID, links[i].FromRouterName, links[i].FromRouterIP, links[i].ToRouterName, links[i].ToRouterIP)
		//		_, err = statement.Exec(links[i].LinkID, links[i].FromRouterID, links[i].FromRouterName, links[i].FromRouterIP, links[i].ToRouterID, links[i].ToRouterName, links[i].ToRouterIP)
		_, err = statement.Exec(links[i].LinkID, links[i].FromRouterID, links[i].FromRouterName, links[i].FromRouterIP, links[i].FromRouterIfIndex, links[i].ToRouterID, links[i].ToRouterName, links[i].ToRouterIP)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				//				fmt.Println("Link already exists. Continue building links.")
				log.Info("Link already exists. Continue building links.")
			} else {
				//				log.Fatalln("Link INSERT error:", err.Error())
				log.Fatal("Link INSERT error %v", err)
			}
		}
		defer statement.Close()
	}

	//	fmt.Println("func buildLinks version", buildLinksVersion, "ending")
	log.Debug("func buildLinks version %s", buildLinksVersion+" ending")

	return database
}

// getRtrID retrieves the RouteID from the database, given a Router Name
func getRtrID(log *logger.Logger, Name string, database *sql.DB) int {
	// Retrive Router from the database
	var RouterID int
	routerRows, queryErr := database.Query("SELECT RouterID, Name FROM Routers WHERE RouterID = ?", Name)
	if queryErr != nil {
		log.Fatal("databaseForRead Query Router error %v", queryErr)
	}
	log.Debug("Successful Routers table Select")
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &Name)
		return RouterID
	}
	return RouterID
}
