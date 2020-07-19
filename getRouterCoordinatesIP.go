package main

import (
	"database/sql"
	"strconv"

	"github.com/g3n/engine/util/logger"
)

// Get Router Coordinates from routerArray
func getRouterCoordinatesIP(debugFlag bool, database *sql.DB, ToRouterIPIn string) (float32, float32, float32) {
	var log *logger.Logger
	log.Debug("getRouterCoordiantesIP starting")

	var x float32
	var y float32
	var z float32

	LinksTableRows, err := database.Query("SELECT RouterID  FROM RouterIP WHERE IpAddr = $1", ToRouterIPIn)
	if err != nil {
		log.Fatal("LinksTableRows error")
	}
	log.Debug("Successful LinksTableRows Query")

	defer LinksTableRows.Close()

	var RouterIDLinks string
	var RouterIDRouter, GpsLat, GpsLong, GpsAlt string
	for LinksTableRows.Next() {
		LinksTableRows.Scan(&RouterIDLinks)

		RouterTableRows, err := database.Query("SELECT RouterID, GpsLat, GpsLong, GpsAlt FROM Routers WHERE RouterID = $1", RouterIDLinks)
		if err != nil {
			log.Fatal("RouterTableRows error: %v", err)
		}
		defer LinksTableRows.Close()

		log.Debug("Successful RouterTable Query")

		for RouterTableRows.Next() {
			RouterTableRows.Scan(&RouterIDRouter, &GpsLat, &GpsLong, &GpsAlt)
		}

		x1, parseErr := strconv.ParseFloat(GpsLong, 32)
		if parseErr != nil {
			log.Fatal("x1 ParseFloat error: %v", parseErr)
		}
		x = (float32)(x1)

		y1, parseErr := strconv.ParseFloat(GpsLat, 32)
		if parseErr != nil {
			log.Fatal("y1 ParseFloat error: %v", parseErr)
		}
		y = (float32)(y1)

		z1, parseErr := strconv.ParseFloat(GpsAlt, 32)
		if parseErr != nil {
			log.Fatal("z1 ParseFloat error %v", parseErr)
		}
		z = (float32)(z1)
	}
	log.Debug("getRouterCoordiantesIP ending")

	return x, y, z
}
