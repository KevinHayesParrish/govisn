package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"log"

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
		//		link.RouterName = Name
		//		link.DestinationName = ""
		//		link.DestinationIP = DestAddr
		//		link.NextHopName = ""
		//		link.NextHopIP = NextHop

		// link.FromRouterName = Name from scan routerrows
		link.FromRouterName = Name // Current router

		//   Determine router interface using ipRouteIfIndex. This is the index of the interface. We can use this to get the interface IP address.

		// Query Select IpAddr from RouterIp where IfIndex == IPRouteIfIndex
		//args := [1]string{IPRouteIfIndex}
		//args := IPRouteIfIndex
		queryRouterRows, queryRtrErr := database.Query("SELECT IpAddr, IfIndex FROM RouterIp WHERE IfIndex = $1", IPRouteIfIndex)
		if queryRtrErr != nil {
			fmt.Println("Query where RouterIp.IfIndex = IPRouteIfIndex", queryRtrErr)
			log.Fatal(queryRtrErr)
		}

		// link.FromRouterIP = IpAddr returned from Select statement
		var ipAddr string
		for queryRouterRows.Next() {
			queryRouterRows.Scan(&ipAddr)
			link.FromRouterIP = ipAddr
		}
		// Query Select router from RouterIP where IpAddr == NextHop
		queryRouterRows, queryRtrErr = database.Query("SELECT IpAddr, IfIndex FROM RouterIp WHERE IfIndex = $1", NextHop)
		if queryRtrErr != nil {
			fmt.Println("Query where RouterIp.IPRouteIfIndex = NextHop", queryRtrErr)
			log.Fatal(queryRtrErr)
		}
		// link.ToRouterName = Name from Select statement
		for queryRouterRows.Next() {
			queryRouterRows.Scan(&ipAddr)
			rtrNames := getRtrName(ipAddr)
			link.ToRouterName = rtrNames[0]
		}

		// link.ToRouterIP = NextHop from scan routerrows
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
			log.Fatalln("Link INSERT error:", err.Error())
		}
		defer statement.Close()
	}

	fmt.Println("func buildLinks version", buildLinksVersion, "stopped")
	return database
}
