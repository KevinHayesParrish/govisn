package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	//	"github.com/g3n/g3nd/material"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/application"
	_ "github.com/mattn/go-sqlite3"
)

/*
* TODO:
 */

//ViewnetVersion is the file version number
const ViewnetVersion = "0.8.3"
const maxRouters int = 1000

// The flag package provides a default help printer via -h switch
var versionFlag = flag.Bool("v", false, "Print the version number.")
var debugFlag = flag.Bool("de", false, "Print Debug statements.")
var sampleNetworkDB = flag.Bool("cr", false, "Create a sample database.")
var loadDBFlag = flag.Bool("l", false, "Load a database from an XML document.")

//var networkXML = "discoverednetwork.xml"
var dbName = "govisnDiscoveredNet.db"

// DbName is the name of the discovered network database file or name of XML input file
//var DbName = flag.String("f", "discoverednetwork.db", "Name of the discovered network database -or-\nName of the XML input file, if combined with -l option.")
var DbName = flag.String("f", "govisnDiscoveredNet.db", "Name of the discovered network database -or-\nName of the XML input file, if combined with -l option.")

//testArangodb is the startup option to test accessing an ArangoDB database
var testArangoDb = flag.Bool("a", false, "Test opening an ArangoDB database")

//discoverFlag is the option to discover a network
//var discoverFlag = flag.String("di", "127.0.0.1", "Discover a network using seed IP Address")
var discoverFlag = flag.String("di", "", "Discover a network using seed IP Address")
var seed = "127.0.0.1"
var community = flag.String("co", "public", "SNMP Community ReadOnly String")
var maxHops = flag.String("m", "0", "Scope of discovery. Maximum number of Hops from seed")

//routerRadius is the radius of the 3D object representing a network router
const routerRadius float64 = 0.5

//globeRadius is the radius of the 3D object representing the earth
//const globeRadius float64 = 1.5
const globeRadius float64 = 63.7

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
	if *testArangoDb {
		testarango()
	}
	if *DbName != "govisnDiscoveredNet.db" {
		dbName = *DbName
	}
	if *loadDBFlag {
		loaddb(*debugFlag, dbName)
		return
	}
	//	if discoverFlag != nil {
	if *discoverFlag != "" {
		seed = *discoverFlag
		if *debugFlag {
			fmt.Println("seed=", seed, "community=", *community)
		}
		discover(*debugFlag, dbName, seed, *community, *maxHops)
	}

	// Open the database containing the discovered network
	databaseForRead, openErr := sql.Open("sqlite3", *DbName)
	if openErr != nil {
		fmt.Println("Error opening databaseForRead", *DbName)
		log.Fatal(openErr)
	}
	defer databaseForRead.Close()

	databaseForUpdate, openErr := sql.Open("sqlite3", *DbName)
	if openErr != nil {
		fmt.Println("Error opening databaseForUpdate", *DbName)
		log.Fatal(openErr)
	}
	defer databaseForUpdate.Close()

	// Retrieve the Routers table
	//	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, SystemName, SystemDesc, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, Description, UpTime, Contact, Location, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(openErr)
	}
	if *debugFlag {
		fmt.Println("Successful Routers table Select")
	}

	// Retrieve the Links table
	//	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, FromRouter, ToRouter FROM Links")
	//	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, FromRouterName, FromRouterIP, ToRouterName, FromRouterIP FROM Links")
	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP FROM Links")
	if queryErr != nil {
		fmt.Println("databaseForRead Query error", queryErr)
		log.Fatal(openErr)
	}
	if *debugFlag {
		fmt.Println("Successful Links table Select")
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
	var Name string
	var Description string
	//	var UpTime string
	var UpTime uint32
	var Contact string
	var Location string
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	//	var X3D float32
	//	var Y3D float32
	//	var Z3D float32
	//	var router Router
	var routerArray [1000]Router
	var link Link
	var LinkID int
	//	var FromRouterName string
	//	var FromRouterIP string
	//	var ToRouterName string
	//	var ToRouterIP string
	var RouterName, DestinationName, DestinationIP, NextHopName, NextHopIP string
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
	app.CameraPersp().SetPosition(0.0, 0.0, (float32)(globeRadius*2.0))

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(globeRadius, 16, 16, 0, math.Pi*2, 0, math.Pi)
	globeMat := material.NewPhong(&math32.Color{R: 0.5, G: 0.5, B: 0.5})
	globeMat.SetTransparent(true)
	globeMat.SetOpacity(0.25)
	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetPosition(0, 0, 0)
	app.Scene().Add(globeMesh)

	if *debugFlag {
		fmt.Println("Beginning routerRows.Next loop; adding routers to 3D scene.")
	}
	/*
	* Add the routers to the 3D scene
	 */
	routerArrayIndex := 0
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &Name, &Description, &UpTime, &Contact, &Location, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		routerArray[routerArrayIndex].System.RouterID = RouterID
		routerArray[routerArrayIndex].System.UpTime = UpTime
		routerArray[routerArrayIndex].System.Name = Name
		routerArray[routerArrayIndex].System.Contact = Contact
		routerArray[routerArrayIndex].System.Location = Location
		routerArray[routerArrayIndex].System.GPS.Latitude = GpsLat
		routerArray[routerArrayIndex].System.GPS.Longitude = GpsLong
		routerArray[routerArrayIndex].System.GPS.Altitude = GpsAlt

		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		mat := material.NewPhong(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		/*
		 * Set coordinates and altitude
		 */
		var GpsLatFloat64 = 0.0
		//		GpsLatFloat64, parseErr := strconv.ParseFloat(GpsLat, 64)
		var parseErr error
		if GpsLat != "" {
			GpsLatFloat64, parseErr = strconv.ParseFloat(GpsLat, 64)
		}
		if parseErr != nil {
			fmt.Println("Error parsing GpsLat =", GpsLat)
			log.Fatal(parseErr)
		}
		xRadianLat := Rad(GpsLatFloat64)

		var GpsLongFloat64 = 0.0
		//		GpsLongFloat64, parseErr := strconv.ParseFloat(GpsLong, 64)
		if GpsLong != "" {
			GpsLongFloat64, parseErr = strconv.ParseFloat(GpsLong, 64)
			if parseErr != nil {
				fmt.Println("Error parsing GpsLong", GpsLong)
				log.Fatal(parseErr)
			}
		}

		xRadianLong := Rad(GpsLongFloat64)
		x = (float32)(globeRadius * math.Sin(xRadianLat) * math.Cos(xRadianLong))
		//		routerArray[routerArrayIndex].System.Coordinates.X = x // update router struc with x coordinate
		routerArray[routerArrayIndex].System.GPS.Longitude = fmt.Sprintf("%f", x) // update router struc with x coordinate
		yRadianLat := Rad(GpsLatFloat64)
		yRadianLong := Rad(GpsLongFloat64)
		y = (float32)(globeRadius * math.Sin(yRadianLat) * math.Sin(yRadianLong))
		//		routerArray[routerArrayIndex].System.Coordinates.Y = y // update router struc with y coordinate
		routerArray[routerArrayIndex].System.GPS.Latitude = fmt.Sprintf("%f", y) // update router struc with y coordinate

		GpsAltFloat64, parseErr := strconv.ParseFloat(GpsAlt, 64)
		GpsAltFloat64 = GpsAltFloat64 / 100000.0
		z = (float32)(globeRadius*(math.Cos(yRadianLat)) + GpsAltFloat64)
		//		routerArray[routerArrayIndex].System.Coordinates.Z = z // update router struc with z coordinate.
		routerArray[routerArrayIndex].System.GPS.Altitude = fmt.Sprintf("%f", z) // update router struc with z coordinate.

		if *debugFlag {
			fmt.Println("x =", x, "y =", y, "z", z)
			fmt.Println("router =", routerArray[routerArrayIndex])
			//			fmt.Println("router.System.Coordinates =", routerArray[routerArrayIndex].System.Coordinates)
			fmt.Println("router.System.GPS =", routerArray[routerArrayIndex].System.GPS)
			fmt.Println("RouterID=", RouterID, "Name=", Name)
		}

		// Add Router object to 3D scene.
		//		cylinderMesh.SetPosition(routerArray[routerArrayIndex].System.Coordinates.X, routerArray[routerArrayIndex].System.Coordinates.Y, routerArray[routerArrayIndex].System.Coordinates.Z)
		//		cylinderMesh.SetPosition(routerArray[routerArrayIndex].System.GPS.Longitude, routerArray[routerArrayIndex].System.GPS.Latitude, routerArray[routerArrayIndex].System.GPS.Altitude)
		cylinderMesh.SetPosition(x, y, z)
		app.Scene().Add(cylinderMesh)

		// Add router name to scene
		// Creates Font
		//		fontfile := app.DirData() + "/fonts/FreeSans.ttf"
		fontfile := os.Getenv("GOPATH") + "/src/govisn/data/fonts/FreeSans.ttf"
		font, err := text.NewFont(fontfile)
		if err != nil {
			app.Log().Fatal(err.Error())
		}

		font.SetLineSpacing(1.0)
		font.SetPointSize(28)
		font.SetDPI(72)
		font.SetFgColor(&math32.Color4{R: 0, G: 0, B: 1, A: 1})
		font.SetBgColor(&math32.Color4{R: 1, G: 1, B: 0, A: 0.8})
		canvas := text.NewCanvas(300, 200, &math32.Color4{R: 0, G: 1, B: 0, A: 0.8})
		rtext := "RouterID: " + strconv.Itoa(routerArray[routerArrayIndex].System.RouterID) + "\nHostname: " + routerArray[routerArrayIndex].System.Name
		swidth, sheight := font.MeasureText(rtext)
		canvas = text.NewCanvas(swidth, sheight, &math32.Color4{R: 0, G: 1, B: 1, A: 1})
		canvas.DrawText(0, 0, rtext, font)
		tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
		mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat3.AddTexture(tex3)
		aspect := float32(swidth) / float32(sheight)
		mesh3 := graphic.NewSprite(aspect, 1, mat3)
		//		mesh3.SetPosition(routerArray[routerArrayIndex].System.Coordinates.X, routerArray[routerArrayIndex].System.Coordinates.Y, routerArray[routerArrayIndex].System.Coordinates.Z+1.0)
		mesh3.SetPosition(x, y, z+1.0)
		app.Scene().Add(mesh3)

		queryErr = routerRows.Err()
		if queryErr != nil {
			log.Fatal(queryErr)
		}

		routerArrayIndex++
	}
	defer routerRows.Close()

	if *debugFlag {
		fmt.Println("Beginning linkRows.Next loop; adding links to the 3D scene")
	}
	/*
	* Add the links to the 3D scene
	 */
	var FromRouterX, FromRouterY, FromRouterZ, ToRouterX, ToRouterY, ToRouterZ float32
	for linkRows.Next() {
		//		err := linkRows.Scan(&LinkID, &FromRouter, &ToRouter)
		//		err := linkRows.Scan(&LinkID, &FromRouterName, &FromRouterIP, &ToRouterName, &ToRouterIP)
		err := linkRows.Scan(&LinkID, &RouterName, &DestinationName, &DestinationIP, &NextHopName, &NextHopIP)
		if err != nil {
			log.Fatal(err)
		}

		// Exclude false routes
		//		if FromRouter == "127.0.0.0" || FromRouter == "127.0.0.1" || FromRouter == "224.0.0.0" || FromRouter == "0.0.0.0" || ToRouter == "127.0.0.0" || ToRouter == "127.0.0.1" || ToRouter == "224.0.0.0" || ToRouter == "0.0.0.0" {
		//		if FromRouterIP == "127.0.0.0" || FromRouterIP == "127.0.0.1" || FromRouterIP == "224.0.0.0" || FromRouterIP == "0.0.0.0" || ToRouterIP == "127.0.0.0" || ToRouterIP == "127.0.0.1" || ToRouterIP == "224.0.0.0" || ToRouterIP == "0.0.0.0" {
		if DestinationIP == "127.0.0.0" || DestinationIP == "127.0.0.1" || DestinationIP == "224.0.0.0" || DestinationIP == "0.0.0.0" || NextHopIP == "127.0.0.0" || NextHopIP == "127.0.0.1" || NextHopIP == "224.0.0.0" || NextHopIP == "0.0.0.0" {
			continue
		}

		// Load link struct from DB fields
		link.LinkID = LinkID
		//		link.FromRouter = FromRouter
		//		link.ToRouter = ToRouter
		//		link.FromRouterName = FromRouterName
		//		link.FromRouterIP = FromRouterIP
		//		link.ToRouterName = ToRouterName
		//		link.ToRouterIP = ToRouterIP
		link.RouterName = RouterName
		link.DestinationName = DestinationName
		link.DestinationIP = DestinationIP
		link.NextHopName = NextHopName
		link.NextHopIP = NextHopIP

		// retrieve FromRouter coordinates from router struc
		if *debugFlag {
			fmt.Println("link =", link)
			//fmt.Println("From routername=", link.FromRouter)
			fmt.Println("From DestinationName=", link.DestinationName)
			fmt.Println("From routername=", link.RouterName)
		}
		//		FromRouterX, FromRouterY, FromRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.FromRouter)
		//		FromRouterX, FromRouterY, FromRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.FromRouterName)
		//		FromRouterX, FromRouterY, FromRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.Destination)
		//FromRouterX, FromRouterY, FromRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.DestinationName)
		FromRouterX, FromRouterY, FromRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.RouterName)
		if *debugFlag {
			//			fmt.Println("router coordinates =", routerArray[routerArrayIndex].System.Coordinates)
			fmt.Println("router coordinates =", routerArray[routerArrayIndex].System.GPS)
			fmt.Println("returned from getRouterCoordinates func: FromRouterX=", FromRouterX, "FromRouterY=", FromRouterY, "FromRouterZ=", FromRouterZ)
		}

		// retrieve ToRouter coordinates from router struc
		if *debugFlag {
			//			fmt.Println("To routername=", link.ToRouter)
			fmt.Println("To routername=", link.NextHopIP)
		}
		//		ToRouterX, ToRouterY, ToRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.ToRouter)
		//		ToRouterX, ToRouterY, ToRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.ToRouterName)
		ToRouterX, ToRouterY, ToRouterZ = getRouterCoordinates(*debugFlag, routerArray, link.NextHopName)
		if *debugFlag {
			//			fmt.Println("router coordinates =", routerArray[routerArrayIndex].System.Coordinates)
			fmt.Println("router coordinates =", routerArray[routerArrayIndex].System.GPS)
			fmt.Println("returned from getRouterCoordinates func: ToRouterX=", ToRouterX, "ToRouterY=", ToRouterY, "ToRouterZ=", ToRouterZ)
		}

		// Add link object to the 3D scene
		// <add gen code for line here>
		//		link3D := geometry.NewCylinder(routerRadius, routerRadius, 0.5, 16, 2, 0, 2*math.Pi, true, true)
		//		linkMat := material.NewPhong(math32.NewColor("Blue"))
		//		cylinderMesh := graphic.NewMesh(link3D, linkMat)
		//		cylinderMesh.SetPosition(FromRouterX, FromRouterY, FromRouterZ)
		linkGeom := geometry.NewGeometry()
		vertices := math32.NewArrayF32(0, 0)
		vertices.Append(
			FromRouterX, FromRouterY, FromRouterZ,
			ToRouterX, ToRouterY, ToRouterZ,
		)
		if *debugFlag {
			fmt.Println("link vertices=", vertices)
		}
		colors := math32.NewArrayF32(0, 0)
		colors.Append(
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
		)
		linkGeom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

		linkGeom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))

		/*
			linkFromVector := math32.NewVector3(FromRouterX, FromRouterY, FromRouterZ)
			linkToVector := math32.NewVector3(ToRouterX, ToRouterY, ToRouterZ)
			linkLine := math32.NewLine3(linkFromVector, linkToVector)
			linkGeom.AddVBO(gls.NewVBO(linkLine))
		*/

		// Creates basic material
		mat := material.NewBasic()

		// Creates lines with the specified geometry and material
		link3D := graphic.NewLines(linkGeom, mat)

		//		app.Scene().Add(cylinderMesh)
		app.Scene().Add(link3D)

		err = linkRows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Run()
}

const constX = math.Pi / 180

// Rad converts degrees to radians
func Rad(d float64) float64 { return d * constX }

// Deg converts radians to degrees
func Deg(r float64) float64 { return r / constX }

// Get Router Coordinates from routerArray
func getRouterCoordinates(debug bool, routerArray [1000]Router, routerName string) (float32, float32, float32) {
	var x float32
	var y float32
	var z float32

	for i := 0; i < len(routerArray); i++ {
		//		if routerArray[i].System.Name == "" {
		//			break // end of routerArray entries
		//		}
		if routerArray[i].System.Name == routerName {
			//			x = routerArray[i].System.Coordinates.X
			x1, err := strconv.ParseFloat(routerArray[i].System.GPS.Longitude, 32)
			if err != nil {
				panic(err)
			}
			x = (float32)(x1)
			//			y = routerArray[i].System.Coordinates.Y
			y1, err := strconv.ParseFloat(routerArray[i].System.GPS.Latitude, 32)
			if err != nil {
				panic(err)
			}
			y = (float32)(y1)
			//			z = routerArray[i].System.Coordinates.Z
			z1, err := strconv.ParseFloat(routerArray[i].System.GPS.Altitude, 32)
			if err != nil {
				panic(err)
			}
			z = (float32)(z1)
			break
		}
	}
	return x, y, z
}
