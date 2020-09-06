package main

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/twpayne/go-kml"

	"github.com/g3n/engine/util/logger"
)

// EXPORTKMLVERSION is the version number of this function's source code file.
const EXPORTKMLVERSION = "0.1.0"

func exportKML(log *logger.Logger, kmlFilename string, DbName string) {
	log.Info("exportKML version %s ", EXPORTKMLVERSION+" started")
	log.Debug("kmlFilename = %s", kmlFilename)
	log.Debug("DbName = %s", DbName)

	databaseForRead, openErr := sql.Open("sqlite3", DbName)
	if openErr != nil {
		log.Fatal("Error opening databaseForRead %v", DbName)
	}
	defer databaseForRead.Close()

	// Retrieve the Routers table
	//	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt FROM Routers")
	routerRows, queryErr := databaseForRead.Query("SELECT Name, Description, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		log.Fatal("databaseForRead Query error %v", queryErr)
	}
	log.Debug("Successful Routers table Select")

	//var routers []Router
	//var router Router
	var Name string
	var Description string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	//	routerArrayIndex := 0

	//enc := xml.NewEncoder(os.Stdout)
	//enc.Indent("  ", "    ")

	var GpsLongFloat float64
	var GpsLatFloat float64
	var GpsAltFloat float64

	for routerRows.Next() {
		routerRows.Scan(&Name, &Description, &Location, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		//router.System.Name = Name
		//router.System.Location = Location
		//router.System.GPS.Latitude = GpsLat
		//router.System.GPS.Longitude = GpsLong
		//router.System.GPS.Altitude = GpsAlt
		s, _ := strconv.ParseFloat(GpsLong, 64)
		GpsLongFloat = s
		s, _ = strconv.ParseFloat(GpsLat, 64)
		GpsLatFloat = s
		s, _ = strconv.ParseFloat(GpsAlt, 64)
		GpsAltFloat = s
		k := kml.KML(
			kml.Placemark(
				kml.Name(Name),
				kml.Description(Description),
				kml.Coordinates(kml.Coordinate{Lon: GpsLongFloat, Lat: GpsLatFloat, Alt: GpsAltFloat}),
			),
		)

		if err := k.WriteIndent(os.Stdout, "", " "); err != nil {
			log.Fatal(err.Error())
		}
	}

	log.Info("exportKML version %s ", EXPORTKMLVERSION+" ending")
	return
}
