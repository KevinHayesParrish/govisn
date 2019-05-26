package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	_ "github.com/mattn/go-sqlite3"
)

//ViewnetVersion is the file version number
const ViewnetVersion = "0.1.6"

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("d", false, "Print Debug statements.")

//DbName is the name of the discovered network database file
var DbName = flag.String("f", "samplenetwork.db", "Name of the discovered network database")

// Router is the structure representing a network router
type Router struct {
	System struct {
		RouterID    int
		Name        string
		Description string
		UpTime      string
		Contact     string
		Location    string
		GPS         struct {
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

func main() {
	flag.Parse() // Scan the arguments list
	fmt.Println("viewnet version:", ViewnetVersion)
	if *versionFlag {
		return
	}
	if *debugFlag {
		fmt.Println("Debug option selected")
	}

	// Open the database containing the discovered network
	database, openErr := sql.Open("sqlite3", *DbName)
	if openErr != nil {
		fmt.Println("Error opening database", *DbName)
		log.Fatal(openErr)
	}
	// Retrieve the Routers table
	routers, queryErr := database.Query("SELECT RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		fmt.Println("Database Query error", queryErr)
		log.Fatal(openErr)
	}

	// Initialize the 3D space
	app, appErr := application.Create(application.Options{
		Title:  "GoVisn - 3D Network Visualization",
		Width:  1200,
		Height: 1000,
	})
	if appErr != nil {
		fmt.Println("Error Creating 3D g3n app", *DbName)
		log.Fatal(openErr)
	}

	var RouterID int
	var SystemName string
	var SystemDesc string
	var UpTime string
	var Contact string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	var router Router
	var x float32
	var y float32 = 1.0
	var z float32 = 1.0

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	app.Scene().Add(pointLight)
	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	app.CameraPersp().SetPosition(4, 0, 15)

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(1, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globeMat := material.NewPhong(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(-1, -1, -1)
	app.Scene().Add(globeMesh)

	for routers.Next() {
		routers.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt)

		router.System.RouterID = RouterID
		router.System.Name = SystemName
		router.System.UpTime = UpTime
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt
		if *debugFlag {
			//			fmt.Println(strconv.Itoa(RouterID) + ": " + SystemName + " " + SystemDesc + " " + UpTime)
			fmt.Println("router =", router)
		}
		// Create a blue cylinder to represent the router and adds it to the scene
		rtr3D := geometry.NewCylinder(1.0, 1.0, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		cylinderMesh.SetPosition(x, y, z)
		app.Scene().Add(cylinderMesh)
		x = x + 2.0
	}
	app.Run()
}
