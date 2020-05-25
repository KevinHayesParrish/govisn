package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	//"log"
	//"math"
	_ "github.com/mattn/go-sqlite3"
)

//buildLinksVersion is the file version sequence number
const buildLinksVersion = "0.0.1"

func buildLinks(debugFlag bool, database *sql.DB) *sql.DB {
	fmt.Println("func buildLinks version", buildLinksVersion, "started")

	/* TODO
	 *
	 * Populate RouterName, DestinationName, DestinationIP, NextHopName and NextHopIP from RouteTable elements.
	 * calculate LinkID using CRC of RouterName/DestinationIP/NextHopIP
	 * Write Links Row to database
	 *
	 */

	//	routeTableRows, err := database.Query("SELECT RouterID, Name, DestAddr, NextHop FROM Routers INNER JOIN RouteTable USING (RouterID)")
	routeTableRows, err := database.Query("SELECT RouterID, Name, DestAddr, IPRouteIfIndex, NextHop FROM Routers INNER JOIN RouteTable USING (RouterID)")
	if err != nil {
		log.Fatalln("databaseForRead JOIN error", err.Error())
	}
	if debugFlag {
		fmt.Println("Successful Routers/RouteTable JOIN")
	}
	defer routeTableRows.Close()

	//var routers []Router
	var router Router
	//routerArrayIndex := 0
	var links []Link
	var link Link
	var RouterID int
	var Name string
	//var IpIfIndex string
	var DestAddr string
	var IPRouteIfIndex string
	var NextHop string
	for routeTableRows.Next() {
		//		routeTableRows.Scan(&RouterID, &Name, &DestAddr, &NextHop)
		routeTableRows.Scan(&RouterID, &Name, &DestAddr, &IPRouteIfIndex, &NextHop)
		router.System.RouterID = RouterID
		router.System.Name = Name

		//   Determine router interface using ipRouteIfIndex. This is the index of the interface. We can use this to get the interface IP address.

		link.FromRouterName = Name // Current router

		// Find FromRouterIP by DNS lookup by name
		fromIPs := getHostIP(Name)
		link.FromRouterIP = fromIPs[0]
		if len(fromIPs) < 1 {
			fmt.Println("No Router IP Address the link from Router", Name)
			link.FromRouterIP = ""
		} else {
			link.FromRouterIP = fromIPs[0]
		}

		// Find ToRouterName
		var routerID string
		var ipAddr string
		var ifIndex string
		queryRouterRows, queryRtrErr := database.Query("SELECT RouterID, IpAddr, IfIndex FROM RouterIp WHERE ipAddr = $1 AND IfIndex = $2", DestAddr, IPRouteIfIndex)
		if queryRtrErr != nil {
			fmt.Println("Query where RouterIp.IPRouteIfIndex = NextHop", queryRtrErr)
			log.Fatal(queryRtrErr)
		}
		defer queryRouterRows.Close()

		for queryRouterRows.Next() {
			queryRouterRows.Scan(&routerID, &ipAddr, ifIndex)
			rtrNames := getRtrName(ipAddr)
			if len(rtrNames) < 1 {
				fmt.Println("No Router Name for Route Destination", ipAddr)
				link.ToRouterName = ""
			} else {
				link.ToRouterName = rtrNames[0]
			}
		}

		link.ToRouterIP = NextHop

		// calculate LinkID
		//		link.LinkID = int(crc32.ChecksumIEEE([]byte(Name)))
		link.LinkID = int(crc32.ChecksumIEEE([]byte(link.FromRouterIP + link.ToRouterIP)))

		links = append(links, link)

		//		statement, err := database.Prepare("INSERT INTO Links (LinkID, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP) VALUES (?, ?, ?, ?, ?, ?)")
		//		if err != nil {
		//			log.Fatalln("Links Insert Prepare err:", err.Error())
		//		}
		//		_, err = statement.Exec(link.LinkID, link.RouterName, link.DestinationName, link.DestinationIP, link.NextHopName, link.NextHopIP)
		//		if err != nil {
		//			log.Fatalln("Link INSERT error:", err.Error())
		//		}
		//		defer statement.Close()

	}
	routeTableRows.Close()

	for i := 0; i < len(links); i++ {
		//SELECT LinkID, FromRouterName, FromRouterIP, ToRouterName, FromRouterIP FROM Links
		//	statement, err := database.Prepare("INSERT INTO Links (LinkID, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP) VALUES (?, ?, ?, ?, ?, ?)")
		statement, err := database.Prepare("INSERT INTO Links (LinkID, FromRouterName, FromRouterIP, ToRouterName, ToRouterIP) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalln("Links Insert Prepare err:", err.Error())
		}
		//		_, err = statement.Exec(links[i].LinkID, links[i].RouterName, links[i].DestinationName, links[i].DestinationIP, links[i].NextHopName, links[i].NextHopIP)
		_, err = statement.Exec(links[i].LinkID, links[i].FromRouterName, links[i].FromRouterIP, links[i].ToRouterName, links[i].ToRouterIP)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				fmt.Println("Link alrady exists. Continue building links.")
			} else {
				log.Fatalln("Link INSERT error:", err.Error())
			}
		}
		defer statement.Close()
	}

	fmt.Println("func buildLinks version", buildLinksVersion, "stopped")
	return database
}
