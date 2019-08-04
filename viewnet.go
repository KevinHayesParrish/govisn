package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	_ "github.com/mattn/go-sqlite3"
)

/*
* TODO:
*	1. Add the router's 3D coordinates to the Router table.
*		This is for use when adding the links to the 3D scene.
 */

//ViewnetVersion is the file version number
const ViewnetVersion = "0.3.4"

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("d", false, "Print Debug statements.")
var sampleNetworkDB = flag.Bool("c", false, "Create a sample database.")

//DbName is the name of the discovered network database file
var DbName = flag.String("f", "discoverednetwork.db", "Name of the discovered network database")

//routerRadius is the radius of the 3D object representing a network router
const routerRadius float64 = 0.5

//globeRadius is the radius of the 3D object representing the earth
//const globeRadius float64 = 1.5
const globeRadius float64 = 65.0

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
		Coordinates struct {
			X float32
			Y float32
			Z float32
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

// Link is the structure representing a network link between two routers
type Link struct {
	LinkID     int
	FromRouter string
	ToRouter   string
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
	if *sampleNetworkDB {
		createsampledb()
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

	// Retrieve the Links table
	links, queryErr := database.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")
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
	var link Link
	var LinkID int
	var FromRouter string
	var ToRouter string
	var x float32
	var y float32 = 1.0
	var z float32 = 1.0

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	//	pointLight.SetPosition(1, 0, 2)
	//	pointLight.SetPosition(10, 10, 20)
	pointLight.SetPosition((float32)(globeRadius+10), (float32)(globeRadius+10), (float32)(globeRadius+20))
	app.Scene().Add(pointLight)
	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	// Set initial camera position, i.e. viewing point
	//	app.CameraPersp().SetPosition(4, 0, 15)
	//	app.CameraPersp().SetPosition(0, 0, (float32)(globeRadius+40))
	app.CameraPersp().SetPosition(0, 0, (float32)(globeRadius*3.0))

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(globeRadius, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globeMat := material.NewPhong(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(-1, -1, -1)
	app.Scene().Add(globeMesh)

	/*
	* Add the routers to the 3D scene
	 */
	for routers.Next() {
		routers.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		router.System.RouterID = RouterID
		router.System.Name = SystemName
		router.System.UpTime = UpTime
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt
		//		if *debugFlag {
		//			//			fmt.Println(strconv.Itoa(RouterID) + ": " + SystemName + " " + SystemDesc + " " + UpTime)
		//			fmt.Println("router =", router)
		//		}
		// Create a blue cylinder to represent the router and adds it to the scene
		//		rtr3D := geometry.NewCylinder(1.0, 1.0, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		/*
		 * Set coordinates and altitude
		 */
		GpsLatFloat64, parseErr := strconv.ParseFloat(GpsLat, 64)
		if parseErr != nil {
			fmt.Println("Error parsing GpsLat", GpsLat)
			log.Fatal(openErr)
		}
		xRadianLat := Rad(GpsLatFloat64)
		GpsLongFloat64, parseErr := strconv.ParseFloat(GpsLong, 64)
		if parseErr != nil {
			fmt.Println("Error parsing GpsLong", GpsLong)
			log.Fatal(openErr)
		}
		xRadianLong := Rad(GpsLongFloat64)
		x = (float32)(globeRadius * math.Sin(xRadianLat) * math.Cos(xRadianLong))
		router.System.Coordinates.X = x // update router struc with x coordinate
		yRadianLat := Rad(GpsLatFloat64)
		yRadianLong := Rad(GpsLongFloat64)
		y = (float32)(globeRadius * math.Sin(yRadianLat) * math.Sin(yRadianLong))
		router.System.Coordinates.Y = y // update router struc with y coordinate

		GpsAltFloat64, parseErr := strconv.ParseFloat(GpsAlt, 64)
		//		zRadianAlt := Rad(GpsAltFloat64)
		//		z = (float)(radius * (java.lang.Math.cos(java.lang.Math.toRadians(Float.valueOf(routerLatitude))))) + (Float.valueOf(routerAltitude));
		z = (float32)(globeRadius*(math.Cos(yRadianLat)) + GpsAltFloat64)
		router.System.Coordinates.Z = z // update router struc with z coordinate

		if *debugFlag {
			//			fmt.Println(strconv.Itoa(RouterID) + ": " + SystemName + " " + SystemDesc + " " + UpTime)
			fmt.Println("router =", router)

		}

		// TODO: write 3D coordinates to router DB row for later retrieval
		//	var stringArray []string
		//		var concatString []string
		var xFloat64 = float64(x)
		var yFloat64 = float64(y)
		var zFloat64 = float64(z)
		//stringArray = append(stringArray, "UPDATE Routers SET X3D =", strconv.FormatFloat(xFloat64, 'f', -1, 32), "Y3D = ", strconv.FormatFloat(yFloat64, 'f', -1, 32), "Z3D =", strconv.FormatFloat(zFloat64, 'f', -1, 32))
		//		concatString = string(concatStringArray)
		updateStatement, _ := database.Prepare("UPDATE Routers SET (X3D, Y3D, Z3D) WHERE SystemName= VALUES (?, ?, ?, ?)")
		updateStatement.Exec(strconv.FormatFloat(xFloat64, 'f', -1, 64), strconv.FormatFloat(yFloat64, 'f', -1, 64), strconv.FormatFloat(zFloat64, 'f', -1, 64), SystemName)

		cylinderMesh.SetPosition(x, y, z)
		app.Scene().Add(cylinderMesh)

		//x = x + 2.0
	}

	/*
	* Add the links to the 3D scene
	 */
	for links.Next() {
		links.Scan(&LinkID, &FromRouter, &ToRouter)
		// Load link struct from DB fields
		link.LinkID = LinkID
		link.FromRouter = FromRouter
		link.ToRouter = ToRouter
		if *debugFlag {
			fmt.Println("link =", link)
		}

		// retrieve FromRouter coordinates from router struc
		//		routers, queryErr = database.Query("SELECT RouterID, SystemName FROM Routers WHERE SystemName =", FromRouter)
		routers.Scan(&RouterID, &SystemName, &SystemDesc, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt)
		x = router.System.Coordinates.X
		y = router.System.Coordinates.Y
		z = router.System.Coordinates.Z
		if *debugFlag {
			fmt.Println("router coordinates =", router.System.Coordinates)
		}

	}
	app.Run()
}

const constX = math.Pi / 180

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * constX }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / constX }
