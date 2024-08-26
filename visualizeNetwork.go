package main

import (
	"database/sql"
	"fmt"
	"strings"

	//"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	//logger "github.com/alouca/gologger"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/util/stats"

	//g "github.com/soniah/gosnmp"
	g "github.com/gosnmp/gosnmp"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/window"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"

	//	"github.com/g3n/engine/util/application"

	"github.com/g3n/engine/experimental/collision"

	_ "github.com/mattn/go-sqlite3"
)

// VISUALIZENETWORKVERSION is the version number of the visualizeNetwork func
const VISUALIZENETWORKVERSION = "0.3.3"

// App contains the application state
type App struct {
	*app.Application                // Embedded standard application object
	log              *logger.Logger // Application logger
	//currentDemo      IDemo            // Current demo
	dirData    string           // Full path of the data directory
	scene      *core.Node       // Scene rendered
	demoScene  *core.Node       // Scene populated by individual demos
	ambLight   *light.Ambient   // Scene ambient light
	frameRater *util.FrameRater // Render loop frame rater

	// GUI
	mainPanel  *gui.Panel
	demoPanel  *gui.Panel
	labelFPS   *gui.Label         // header FPS label
	treeTests  *gui.Tree          // tree with test names
	stats      *stats.Stats       // statistics object
	statsTable *stats.StatsTable  // statistics table panel
	control    *gui.ControlFolder // Pointer to gui control panel

	// Camera and orbit control
	camera *camera.Camera       // Camera
	orbit  *camera.OrbitControl // Orbit control
}

// GuiMenu is the structure containing the Menus for the Gui
type GuiMenu struct {
}

// Raycast is the structure containing the Raycaster
type Raycast struct {
	rayCast *collision.Raycaster
}

// func visualizeNetwork(debugFlag bool, log *logger.Logger, databaseForRead *sql.DB, snmpTarget string, community string, params *g.GoSNMP) *sql.DB {
// func visualizeNetwork(log *logger.Logger, databaseForRead *sql.DB, snmpTarget string, community string, params *g.GoSNMP) *sql.DB {
func visualizeNetwork(log *logger.Logger, databaseForRead *sql.DB, snmpTarget string, params *g.GoSNMP) *sql.DB {
	//	const VISUALIZENETWORKVERSION = "0.3.1"
	log.Debug("visualizeNetwork %s", VISUALIZENETWORKVERSION+" started")

	// Retrieve the Routers table
	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, Description, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt FROM Routers")
	if queryErr != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query error %v", queryErr)
	}
	log.Debug("Successful Routers table Select")

	// Retrieve the Links table
	linkRows, queryErr := databaseForRead.Query("SELECT LinkID, FromRouterName, FromRouterIP, ToRouterName, ToRouterIP FROM Links")
	if queryErr != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query error %v", queryErr)
	}
	log.Debug("Successful Links table Select")

	// Initialize the 3D space

	// Create application and scene
	gv := new(gvapp)
	a := app.App()
	gv.Application = a
	gv.scene = core.NewNode()
	gv.scene.SetName("GoVisnScene")

	// Create perspective camera
	gv.camPos = math32.Vector3{X: 0, Y: 0, Z: (float32)(globeRadius * 2.0)}
	gv.cam = camera.New(1) // perspective camera with defaults
	gv.cam.SetName("camera")
	gv.cam.SetPosition(0, 0, (float32)(globeRadius*2.0))

	// Setup orbit control for the camera
	gv.orbit = camera.NewOrbitControl(gv.cam)

	gv.scene.Add(gv.cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		gv.cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Create and Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1.0)
	ambientLight.SetName("ambient")
	gv.scene.Add(ambientLight)

	//	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	//  pointLight.SetName("pointLight")
	//	pointLight.SetPosition((float32)(globeRadius+10), (float32)(globeRadius+10), (float32)(globeRadius+20))
	//	gv.scene.Add(pointLight)

	//	dirLight := light.NewDirectional(math32.NewColor("white"), 0.8)
	//	dirLight.SetPosition((float32)(globeRadius+100), 0, 0)
	//	gv.scene.Add(dirLight)

	// Add an axis helper to the scene
	//axes := helper.NewAxes(1)
	//axes.SetName("helperAxes")
	//gv.scene.Add(axes)

	// Set background color to black
	a.Gls().ClearColor(0.0, 0.0, 0.0, 0.0)

	gv = addTitle(log, gv)

	// Build Menus
	//	buildMenus(debugFlag, gv, a, databaseForRead)
	//buildMenus(gv, a, databaseForRead)
	buildMenus(gv, a)

	var RouterID int
	var Name string
	var Description string
	var UpTime uint32
	var Contact string
	var Location string
	var Services int
	var GpsLat string
	var GpsLong string
	var GpsAlt string
	var link Link
	var LinkID int
	var FromRouterName, FromRouterIP, ToRouterName, ToRouterIP string
	var x float32
	var y float32 = 1.0
	var z float32 = 1.0

	// Setup Mouse clicking of objects within the 3D scene
	var t Raycast
	//	t.Initialize(debugFlag, gv.scene, gv.cam, gv, a, databaseForRead)
	t.Initialize(gv.scene, gv.cam, gv, a, databaseForRead)

	// Create Globe texture
	gobinDir := os.Getenv("GOBIN")
	texfile := gobinDir + "/data/images/earth_clouds_big.jpg"
	globeTex, err := texture.NewTexture2DFromImage(texfile)
	if err != nil {
		databaseForRead.Close()
		log.Fatal("Error loading texture.\n Insure govisn /data/images is copied to GOBIN \n GOBIN env variable must be set.")
	}
	globeTex.SetFlipY(false)

	// Create a sphere representing the globe
	globe3D := geometry.NewSphere(globeRadius, 16, 16)
	//globeMat := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0}) // White 255, 255, 255
	globeMat := material.NewStandard(math32.NewColor("grey"))
	//globeMat.AddTexture(globeTex)
	globeMat.SetTransparent(true)
	globeMat.SetOpacity(.30)

	globeMesh := graphic.NewMesh(globe3D, globeMat)
	globeMesh.SetName("globe")
	globeMesh.SetPosition(0, 0, 0)
	gv.scene.Add(globeMesh)

	log.Debug("Beginning routerRows.Next loop; adding routers to 3D scene.")

	//
	// Add the routers to the 3D scene
	//
	var routers []Router
	var router Router
	routerArrayIndex := 0
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &Name, &Description, &UpTime, &Contact, &Location, &Services, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		router.System.RouterID = RouterID
		router.System.UpTime = UpTime
		router.System.Name = Name
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.Services = Services
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt

		routers = append(routers, router)

		rtr3D := geometry.NewCylinder(routerRadius, routerRadius, 16, 2, true, true)
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		cylinderMesh := graphic.NewMesh(rtr3D, mat)
		cylinderMesh.SetName(router.System.Name)
		cylinderMesh.SetUserData(strconv.Itoa(router.System.RouterID))

		// Set coordinates and altitude
		x, y, z = calcCoordinates(GpsLat, GpsLong, GpsAlt)

		//		log.Debug("x = %s", strconv.FormatFloat(float64(x), 'f', 5, 32)+"y = %s"+strconv.FormatFloat(float64(y), 'f', 5, 32)+"z = %s"+strconv.FormatFloat(float64(z), 'f', 5, 32))
		log.Debug("x = %s", strconv.FormatFloat(float64(x), 'f', 5, 32)+"y = "+strconv.FormatFloat(float64(y), 'f', 5, 32)+"z = "+strconv.FormatFloat(float64(z), 'f', 5, 32))
		log.Debug("router = %v", routers[routerArrayIndex])
		log.Debug("router.System.GPS = %s", routers[routerArrayIndex].System.GPS)
		//		log.Debug("RouterID= %s", strconv.Itoa(RouterID)+"Name= %s"+Name)
		log.Debug("RouterID= %s", strconv.Itoa(RouterID)+"Name= "+Name)

		// Add Router object to 3D scene.
		cylinderMesh.SetPosition(x, y, z)
		log.Debug("cylinderMesh Name= %s", cylinderMesh.Name())
		log.Debug("cylinderMesh UserData= %s", cylinderMesh.UserData())

		gv.scene.Add(cylinderMesh)

		//
		// Add router name to scene
		//
		fontfile := os.Getenv("GOBIN") + "/data/fonts/FreeSans.ttf"
		font, err := text.NewFont(fontfile)
		if err != nil {
			databaseForRead.Close()
			log.Fatal("Error loading font %s" + err.Error() + "\n Insure govisn /data/fonts is copied to GOBIN \n GOBIN env variable must be set.")
		}

		font.SetLineSpacing(1.0)
		font.SetPointSize(28)
		font.SetDPI(72)
		font.SetFgColor(&math32.Color4{R: 0, G: 0, B: 1, A: 1})
		font.SetBgColor(&math32.Color4{R: 1, G: 1, B: 0, A: 0.8})
		canvas := text.NewCanvas(300, 200, &math32.Color4{R: 0, G: 1, B: 0, A: 0.8})
		rtext := "RouterID: " + strconv.Itoa(routers[routerArrayIndex].System.RouterID) + "\nHostname: " + routers[routerArrayIndex].System.Name
		swidth, sheight := font.MeasureText(rtext)
		canvas = text.NewCanvas(swidth, sheight, &math32.Color4{R: 0, G: 1, B: 1, A: 1})
		canvas.DrawText(0, 0, rtext, font)
		tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
		mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat3.AddTexture(tex3)
		aspect := float32(swidth) / float32(sheight)
		mesh3 := graphic.NewSprite(aspect, 1, mat3)
		mesh3.SetPosition(x, y+1.0, z)
		gv.scene.Add(mesh3)

		queryErr = routerRows.Err()
		if queryErr != nil {
			databaseForRead.Close()
			log.Fatal(queryErr.Error())
		}

		routerArrayIndex++
	}
	defer routerRows.Close()

	log.Debug("Beginning linkRows.Next loop; adding links to the 3D scene")
	//
	// Add the links to the 3D scene
	//
	var FromRouterX, FromRouterY, FromRouterZ, ToRouterX, ToRouterY, ToRouterZ string
	for linkRows.Next() {
		err := linkRows.Scan(&LinkID, &FromRouterName, &FromRouterIP, &ToRouterName, &ToRouterIP)
		if err != nil {
			databaseForRead.Close()
			log.Fatal(err.Error())
		}

		// Exclude false routes
		if FromRouterIP == "127.0.0.0" || FromRouterIP == "127.0.0.1" || FromRouterIP == "224.0.0.0" || FromRouterIP == "0.0.0.0" || ToRouterIP == "127.0.0.0" || ToRouterIP == "127.0.0.1" || ToRouterIP == "224.0.0.0" || ToRouterIP == "0.0.0.0" {
			continue
		}

		// Load link struct from DB fields
		link.LinkID = LinkID
		link.FromRouterName = FromRouterName
		link.FromRouterIP = FromRouterIP
		link.ToRouterName = ToRouterName
		link.ToRouterIP = ToRouterIP

		// retrieve FromRouter coordinates
		log.Debug("link = %v", link)
		log.Debug("FromRouterName= %s", link.FromRouterName)
		log.Debug("FromRouterIP= %s", link.FromRouterIP)

		//  Query database for FromRouter GPS coordinates

		routerGpsRows, err := databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name = $1", FromRouterName)
		if err != nil {
			databaseForRead.Close()
			log.Fatal("databaseForRead Query error %s", err.Error())
		}
		log.Debug("Successful Query for FromRouter GPS Coordinates")

		defer routerGpsRows.Close()
		var linksFromRouterName string
		var linksFromRouterGpsLat, linksFromRouterGpsLong, linksFromRouterGpsAlt string
		for routerGpsRows.Next() {
			routerGpsRows.Scan(&linksFromRouterName, &linksFromRouterGpsLat, &linksFromRouterGpsLong, &linksFromRouterGpsAlt)
		}

		FromRouterX = linksFromRouterGpsLat
		FromRouterY = linksFromRouterGpsLong
		FromRouterZ = linksFromRouterGpsAlt
		//		log.Debug("returned from getRouterCoordinatesName func: FromRouterX= %s", FromRouterX+" FromRouterY= %s"+FromRouterY+" FromRouterZ= %s"+FromRouterZ)
		log.Debug("returned from getRouterCoordinatesName func: FromRouterX= %s", FromRouterX+" FromRouterY= "+FromRouterY+" FromRouterZ= "+FromRouterZ)

		//  Query database for FromRouter GPS coordinates
		log.Debug("ToRouterName= %s", link.ToRouterName)
		log.Debug("ToRouterIP= %s", link.ToRouterIP)

		routerGpsRows, err = databaseForRead.Query("SELECT Name, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name = $1", ToRouterName)
		if err != nil {
			databaseForRead.Close()
			log.Fatal("databaseForRead Query error %v", err)
		}
		log.Debug("Successful Query for ToRouter GPS Coordinates")

		var linksToRouterName string
		var linksToRouterGpsLat, linksToRouterGpsLong, linksToRouterGpsAlt string
		for routerGpsRows.Next() {
			routerGpsRows.Scan(&linksToRouterName, &linksToRouterGpsLat, &linksToRouterGpsLong, &linksToRouterGpsAlt)
		}
		ToRouterX = linksToRouterGpsLat
		ToRouterY = linksToRouterGpsLong
		ToRouterZ = linksToRouterGpsAlt

		log.Debug("router %s", Name+" GPS coordinates = %s, %s, %s"+GpsLat+GpsLong+GpsAlt)
		log.Debug("returned from getRouterCoordinatesIP func: ToRouterX= %s", ToRouterX+"ToRouterY= %s"+ToRouterY+"ToRouterZ= %s"+ToRouterZ)

		// Add link object to the 3D scene
		fromX, fromY, fromZ := calcCoordinates(FromRouterX, FromRouterY, FromRouterZ)
		toX, toY, toZ := calcCoordinates(ToRouterX, ToRouterY, ToRouterZ)
		log.Debug("fromX= %s", strconv.FormatFloat(float64(fromX), 'f', 5, 64)+
			" fromY= %s"+strconv.FormatFloat(float64(fromY), 'f', 5, 64)+
			" fromZ= %s"+strconv.FormatFloat(float64(fromZ), 'f', 5, 64))
		log.Debug("toX= %s", strconv.FormatFloat(float64(toX), 'f', 5, 64)+
			" toY= %s"+strconv.FormatFloat(float64(toY), 'f', 5, 64)+
			" toZ="+strconv.FormatFloat(float64(toZ), 'f', 5, 64))

		// Build Link using glLine - BEGIN
		linkGeom := geometry.NewGeometry()
		vertices := math32.NewArrayF32(0, 0)
		vertices.Append(
			fromX, fromY, fromZ,
			toX, toY, toZ,
		)
		log.Debug("link vertices= %v", vertices)
		log.Debug("")

		linkGeom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

		//mat := material.NewStandard(math32.NewColor("White"))
		mat := material.NewStandard(math32.NewColor("grey"))

		// Check Runtime environment.
		// OpenGL Implementation on MacOS will only accept Line width of 1.0
		if runtime.GOOS == "darwin" {
			mat.SetLineWidth(1.0)
			log.Info("*** Link SetLineWidth() request ignored. OpenGL Implementation on MacOS will only accept lineWidth of 1.0 ***")
		}
		link3D := graphic.NewLines(linkGeom, mat)
		link3D.SetName(strconv.Itoa(link.LinkID))
		// Build Link using glLine - END

		// Build Link using Polygon
		//posA := math32.NewVector3(fromX, fromY, fromZ)
		//		posB := math32.NewVector3(toX, toY, toZ)
		//		vertices, normals, indices := calcLinkVBOs(debugFlag, gv.camPos, *posA, *posB, float32(1.00))
		//		if debugFlag {
		//			fmt.Println("calcLinkVBOs returned: ", vertices, normals, indices)
		//		}
		//		linkGeom.SetIndices(indices)
		//		linkGeom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
		//		linkGeom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))

		// Creates basic material
		//		mat := material.NewStandard(math32.NewColor("White"))
		//		mat.SetSide(material.SideDouble)
		//		link3D := graphic.NewMesh(linkGeom, mat)
		//		link3D.SetVisible(true)
		// Build Link using Polygon - END

		// Build Link using Cylinder - BEGIN
		//vertices := math32.NewArrayF32(0, 0)
		//vertices.Append(
		//			fromX, fromY, fromZ,
		//			toX, toY, toZ,
		//		)
		//		if debugFlag {
		//			fmt.Println("link vertices=", vertices)
		//			fmt.Println()
		//		}

		//		linkGeom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

		//posA := math32.NewVector3(fromX, fromY, fromZ)
		//posB := math32.NewVector3(toX, toY, toZ)
		//cylHeight := calcDistance(debugFlag, posA, posB)
		//linkGeom := geometry.NewCylinder(1.0, cylHeight, 16, 2, true, true)
		//mat := material.NewStandard(math32.NewColor("White"))
		// Check Runtime environment.
		// OpenGL Implementation on MacOS will only accept Line width of 1.0
		//		if runtime.GOOS == "darwin" {
		//			mat.SetLineWidth(1.0)
		//			fmt.Println("*** Link SetLineWidth() request ignored. OpenGL Implementation on MacOS will only accept 1.0 ***")
		//		} else {
		//			mat.SetLineWidth(3.0)
		//		}
		//link3D := graphic.NewMesh(linkGeom, mat)
		//link3D.SetName(string(link.LinkID))
		// Build Link using Cylinder - END

		// Creates lines with the specified geometry and material

		gv.scene.Add(link3D)

		err = linkRows.Err()
		if err != nil {
			databaseForRead.Close()
			log.Fatal(err.Error())
		}
	}

	// Creating a 60 second timer for auto link update feature
	linkUpdateTimer := time.NewTimer(60 * time.Second)
	updateLinksOK := false

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(gv.scene, gv.cam)

		// Notifying channel under go function
		go func() {
			<-linkUpdateTimer.C

			// set NetPollingEnabled switch when timer is fired
			//NetPollingEnabled = true
			if NetPollingEnabled {
				updateLinksOK = true
			}

			// Reset the linkUpdateTimer to 60 seconds
			linkUpdateTimer.Reset(60 * time.Second)
			log.Info("linkUpdateTimer Reset")
		}()

		//		if NetPollingEnabled {
		if updateLinksOK {
			//gv = updateLinks(log, gv, databaseForRead, snmpTarget, community, params)
			gv = updateLinks(log, gv, databaseForRead, snmpTarget, params)
			//			NetPollingEnabled = false
			updateLinksOK = false
		}
	})

	log.Debug("visualizeNetwork %s", VISUALIZENETWORKVERSION+" func ending.")

	return databaseForRead

}

func calcCoordinates(GpsLat string, GpsLong string, GpsAlt string) (float32, float32, float32) {
	var x, y, z float32

	var GpsLatFloat64 = 0.0
	var parseErr error
	if GpsLat != "" {
		GpsLatFloat64, parseErr = strconv.ParseFloat(GpsLat, 64)
	}
	if parseErr != nil {
		fmt.Println("Error parsing GpsLat =", GpsLat)
		log.Fatal(parseErr.Error())
	}
	xRadianLat := Rad(GpsLatFloat64)

	var GpsLongFloat64 = 0.0
	if GpsLong != "" {
		GpsLongFloat64, parseErr = strconv.ParseFloat(GpsLong, 64)
		if parseErr != nil {
			log.Fatal("Error parsing GpsLong %s", GpsLong)
		}
	}

	xRadianLong := Rad(GpsLongFloat64)
	x = (float32)(globeRadius * math.Sin(xRadianLat) * math.Cos(xRadianLong))

	yRadianLat := Rad(GpsLatFloat64)
	yRadianLong := Rad(GpsLongFloat64)
	y = (float32)(globeRadius * math.Sin(yRadianLat) * math.Sin(yRadianLong))

	GpsAltFloat64, parseErr := strconv.ParseFloat(GpsAlt, 64)
	GpsAltFloat64 = GpsAltFloat64 / 100000.0
	z = (float32)(globeRadius*(math.Cos(yRadianLat)) + GpsAltFloat64)

	return x, y, z
}

// FileExitSelected indicates the File Exit menuitem was selected
var FileExitSelected bool = false

// NetPollingEnabled indicates the Network Traffic Polling state
var NetPollingEnabled bool = false

// buildmenus creates the Gui menus and menuitems for the application
// func buildMenus(debugFlag bool, gv *gvapp, a *app.Application, databaseForRead *sql.DB) *app.Application {
// func buildMenus(gv *gvapp, a *app.Application, databaseForRead *sql.DB) *app.Application {
func buildMenus(gv *gvapp, a *app.Application) *app.Application {
	log.Debug("Starting func buildMenus")

	// Event handler for menu clicks
	onClick := func(evname string, ev interface{}) {
		switch ev.(*gui.MenuItem).Id() {
		case "Reset":
			{
				log.Debug("Resetting Camera to initial view.")
				gv.cam.SetPositionVec(&gv.camPos)
				gv.cam.LookAt(&math32.Vector3{X: 0, Y: 0, Z: 0}, &math32.Vector3{X: 0, Y: 1, Z: 0})
				gv.orbit.Reset()
			}
		case "Print":
			{
				Dump3dScene(gv)
			}
		case "Exit":
			{
				FileExitSelected = true
				log.Info("GoVisn terminating. File/Exit selected.")
				gv.Exit()
			}
		case "Enable Polling":
			{
				NetPollingEnabled = true
				log.Debug("Network Traffic Polling %t", NetPollingEnabled)
			}
		case "Disable Polling":
			{
				NetPollingEnabled = false
				log.Debug("Network Traffic Polling %t", NetPollingEnabled)
			}
		}
	}

	gui.Manager().Set(gv.scene)

	// Create menu bar
	mb := gui.NewMenuBar()
	mb.Subscribe(gui.OnClick, onClick)
	mb.SetPosition(10, 10)

	// Create fileMenu and add it to the menu bar
	m1 := gui.NewMenu()
	m1.AddOption("Reset Camera to Initial View").
		SetId("Reset")
	m1.AddOption("Print 3D Scene graph").
		SetId("Print")
	m1.AddOption("Exit").
		SetId("Exit")
	mb.AddMenu("File", m1).
		SetId("File").
		SetShortcut(window.ModAlt, window.Key1)

	// Create linksMenu and add it to the menu bar
	m2 := gui.NewMenu()
	m2.AddOption("Enable Auto-Update Links").
		SetId("Enable Polling")
	m2.AddOption("Disable Auto-Update Links").
		SetId("Disable Polling")
	mb.AddMenu("Links", m2).
		SetId("Enable").
		SetShortcut(window.ModAlt, window.Key1)

	mb.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		material.NewStandard(math32.NewColor("DarkRed"))
	})
	gv.scene.Add(mb)

	gui.Manager().SetKeyFocus(mb)

	log.Debug("func buildMenus ended")

	return (a)
}

// Initialize the raycaster
// func (t *Raycast) Initialize(debugFlag bool, scene *core.Node, cam *camera.Camera, gv *gvapp, app *app.Application, databaseForRead *sql.DB) {
func (t *Raycast) Initialize(scene *core.Node, cam *camera.Camera, gv *gvapp, app *app.Application, databaseForRead *sql.DB) {
	log.Debug("Initializing the raycaster")

	// Creates the raycaster
	t.rayCast = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	t.rayCast.LinePrecision = 0.05
	t.rayCast.PointPrecision = 0.05

	// Subscribe to mouse button down events
	app.SubscribeID(window.OnMouseDown, app, func(evname string, ev interface{}) {
		//		t.onMouse(debugFlag, scene, cam, gv, app, databaseForRead, ev)
		t.onMouse(scene, cam, gv, app, databaseForRead, ev)
	})
}

// onMouse is executed when an object in the 3D scene is selected with a mouse click
// func (t *Raycast) onMouse(debugFlag bool, scene *core.Node, cam *camera.Camera, gv *gvapp, app *app.Application, databaseForRead *sql.DB, ev interface{}) {
func (t *Raycast) onMouse(scene *core.Node, cam *camera.Camera, gv *gvapp, app *app.Application, databaseForRead *sql.DB, ev interface{}) {
	// Convert mouse coordinates to normalized device coordinates
	mev := ev.(*window.MouseEvent)
	width, height := app.GetSize()
	x := 2*(mev.Xpos/float32(width)) - 1
	y := -2*(mev.Ypos/float32(height)) + 1
	log.Debug("onMouse x= %f", x)
	log.Debug("onMouse y= %f", y)

	// Set the raycaster from the current camera and mouse coordinates
	t.rayCast.SetFromCamera(cam, x, y)
	log.Debug("rayCast:%+v\n", t.rayCast.Ray)

	// Checks intersection with all objects in the scene
	intersects := t.rayCast.IntersectObjects(scene.Children(), true)
	log.Debug("intersects:%+v\n", intersects)

	if len(intersects) == 0 {
		return
	}

	// Get first intersection
	obj := intersects[0].Object
	router3D := obj.GetNode()
	router3DName := router3D.Name()
	if router3DName == "" {
		log.Debug("No Router selected. Try again.")
	} else {
		log.Debug("Picked object Name= %s", router3DName)
		log.Debug("Picked object UserData= %s", router3D.UserData())
	}

	// Retrieve Router info from database
	//	router := RetrieveRouter(debugFlag, router3DName, databaseForRead, app)
	router := RetrieveRouter(router3DName, databaseForRead, app)
	log.Debug("router= %v", router)

	// Add Router info to 3D scene
	fontfile := os.Getenv("GOBIN") + "/data/fonts/FreeSans.ttf"
	font, err := text.NewFont(fontfile)
	if err != nil {
		log.Fatal("Error loading font.\n Insure govisn /data/fonts is copied to GOBIN. \n GOBIN env variable must be set.")
	}

	font.SetLineSpacing(1.0)
	font.SetPointSize(28)
	font.SetDPI(72)
	font.SetFgColor(&math32.Color4{R: 0, G: 0, B: 1, A: 1})
	font.SetBgColor(&math32.Color4{R: 1, G: 1, B: 0, A: 0.8})
	canvas := text.NewCanvas(300, 200, &math32.Color4{R: 0, G: 1, B: 0, A: 0.8})

	x, y, z := calcCoordinates(router.System.GPS.Latitude, router.System.GPS.Longitude, router.System.GPS.Altitude)

	for i := 0; i < len(router.Addresses.NetworkAddresses.IPAddress); i++ {

		rtext := "\nIP Address: " + router.Addresses.NetworkAddresses.IPAddress[i]
		swidth, sheight := font.MeasureText(rtext)
		canvas = text.NewCanvas(swidth, sheight, &math32.Color4{R: 0, G: 1, B: 1, A: 1})
		canvas.DrawText(0, 0, rtext, font)
		tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
		mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat3.AddTexture(tex3)
		aspect := float32(swidth) / float32(sheight)
		mesh3 := graphic.NewSprite(aspect, 1, mat3)

		mesh3.SetPosition(x-3.0, y-1.0-float32(i), z)
		gv.scene.Add(mesh3)
	}
	for j := 0; j < len(router.Addresses.NetworkAddresses.IPAddress); j++ {
		rtext := "\nMAC Address: " + router.Addresses.MediaAddresses.MediaAddress[j]
		swidth, sheight := font.MeasureText(rtext)
		canvas = text.NewCanvas(swidth, sheight, &math32.Color4{R: 0, G: 1, B: 1, A: 1})
		canvas.DrawText(0, 0, rtext, font)
		tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
		mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
		mat3.AddTexture(tex3)
		aspect := float32(swidth) / float32(sheight)
		mesh3 := graphic.NewSprite(aspect, 1, mat3)

		mesh3.SetPosition(x+3.0, y-1.0-float32(j), z)
		gv.scene.Add(mesh3)
	}

	// TESTING ONLY - BEGIN
	// Convert INode to IGraphic
	ig, ok := obj.(graphic.IGraphic)
	if !ok {
		return
	}
	// Get graphic object
	gr := ig.GetGraphic()
	imat := gr.GetMaterial(0)

	type matI interface {
		EmissiveColor() math32.Color
		SetEmissiveColor(*math32.Color)
	}

	if v, ok := imat.(matI); ok {
		if em := v.EmissiveColor(); em.R == 1 && em.G == 1 && em.B == 1 {
			v.SetEmissiveColor(&math32.Color{R: 0, G: 0, B: 0})
		} else {
			v.SetEmissiveColor(&math32.Color{R: 1, G: 1, B: 1})
		}
	}
	// TESTING ONLY - END

}

// Dump3dScene writes the Collada file representing the 3D Scene
func Dump3dScene(gv *gvapp) {
	//	fmt.Println("Dumping 3D Scene")
	log.Debug("Dumping 3D Scene")
	//var decoder collada.Decoder
	//var out io.Writer
	//decoder.Dump(out, 4)
	//var scene = gv.scene
	//scene.Dump(out, 4)
}

// RetrieveRouter is called when an object in the 3D scene is mouse clicked. It retrieve's the
//
//	routers information from the database and opens a new window to display it.
//
// func RetrieveRouter(debugFlag bool, router3DName string, databaseForRead *sql.DB, app *app.Application) Router {
func RetrieveRouter(router3DName string, databaseForRead *sql.DB, app *app.Application) Router {
	var router Router
	var RouterID, Services int
	var Name, Contact, Location, GpsLat, GpsLong, GpsAlt string
	var MacAddr, IPAddr string

	var UpTime uint32
	// Retrive Router from the database
	//	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt FROM Routers WHERE RouterID = ?", router3DName)
	routerRows, queryErr := databaseForRead.Query("SELECT RouterID, Name, UpTime, Contact, Location, Services, GpsLat, GpsLong, GpsAlt FROM Routers WHERE Name = ?", router3DName)
	if queryErr != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query Router error %v", queryErr)
	}
	log.Debug("Successful Routers table Select")
	for routerRows.Next() {
		routerRows.Scan(&RouterID, &Name, &UpTime, &Contact, &Location, &Services, &GpsLat, &GpsLong, &GpsAlt)

		// Load router struct from DB fields
		router.System.RouterID = RouterID
		router.System.UpTime = UpTime
		router.System.Name = Name
		router.System.Contact = Contact
		router.System.Location = Location
		router.System.Services = Services
		router.System.GPS.Latitude = GpsLat
		router.System.GPS.Longitude = GpsLong
		router.System.GPS.Altitude = GpsAlt
	}

	// Retrieve MAC Addresses from the database
	macRows, queryErr := databaseForRead.Query("SELECT RouterID, MacAddr FROM RouterMac WHERE RouterID = ?", router.System.RouterID)
	if queryErr != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query MAC error %v", queryErr)
	}
	i := 0
	for macRows.Next() {
		macRows.Scan(&RouterID, &MacAddr)

		// Load router struct from DB fields
		router.Addresses.MediaAddresses.MediaAddress = append(router.Addresses.MediaAddresses.MediaAddress, MacAddr)

		i++
	}

	// Retrieve IP Addresses from the database
	ipRows, queryErr := databaseForRead.Query("SELECT RouterID, IpAddr FROM RouterIp WHERE RouterID = ?", router.System.RouterID)
	if queryErr != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query IP error %v", queryErr)
	}
	j := 0
	for ipRows.Next() {
		ipRows.Scan(&RouterID, &IPAddr)

		// Load router struct from DB fields
		router.Addresses.NetworkAddresses.IPAddress = append(router.Addresses.NetworkAddresses.IPAddress, IPAddr)

		log.Debug("OnMouse Router= %v", router)
		j++
	}
	return router
}

// calcLinkVBOs calculates the vertices of the polygon representing the network link.
// func calcLinkVBOs(debugFlag bool, camPos math32.Vector3, posA math32.Vector3, posB math32.Vector3, scalar float32) (
/* func calcLinkVBOs(camPos math32.Vector3, posA math32.Vector3, posB math32.Vector3, scalar float32) (
	vertices math32.ArrayF32,
	normals math32.ArrayF32,
	indices math32.ArrayU32) {

	var linkVertex1,
		linkVertex2,
		linkVertex3,
		linkVertex4 math32.Vector3

	// PosA vertices
	linkVertex1.SetX(posA.Component(0))
	linkVertex1.SetY(posA.Component(1))
	linkVertex1.SetZ(posA.Component(2))

	linkVertex2.SetX(posA.Component(0))
	linkVertex2.SetY(posA.Component(1) + scalar)
	linkVertex2.SetZ(posA.Component(2))

	//PosB vertices
	linkVertex3.SetX(posA.Component(0))
	linkVertex3.SetY(posA.Component(1) + scalar)
	linkVertex3.SetZ(posA.Component(2))

	linkVertex4.SetX(posB.Component(0))
	linkVertex4.SetY(posB.Component(1))
	linkVertex4.SetZ(posB.Component(2))

	vertices = math32.NewArrayF32(0, 0)
	vertices.Append(
		linkVertex1.Component(0), linkVertex1.Component(1), linkVertex1.Component(2),
		linkVertex2.Component(0), linkVertex2.Component(1), linkVertex2.Component(2),
		linkVertex3.Component(0), linkVertex3.Component(1), linkVertex3.Component(2),
		linkVertex4.Component(0), linkVertex4.Component(1), linkVertex4.Component(2),
	)
	normals = math32.NewArrayF32(0, 0)
	normals.Append(
		camPos.Component(0), camPos.Component(1), camPos.Component(2),
		camPos.Component(0), camPos.Component(1), camPos.Component(2),
		camPos.Component(0), camPos.Component(1), camPos.Component(2),
		camPos.Component(0), camPos.Component(1), camPos.Component(2),
	)

	indices = math32.NewArrayU32(0, 0)
	indices.Append(
		0, 1, 2,
		0, 2, 3,
	)

	return vertices, normals, indices
}
*/

/* func calcDistance(debugFlag bool, posA *math32.Vector3, posB *math32.Vector3) (distance float64) {
	x2 := posB.Component(0)
	x1 := posA.Component(0)
	y2 := posB.Component(1)
	y1 := posA.Component(1)
	z2 := posB.Component(2)
	z1 := posA.Component(2)

	distance = math.Sqrt(math.Pow(float64(x2-x1), 2.0) + math.Pow(float64(y2-y1), 2.0) + math.Pow(float64(z2-z1), 2.0))
	log.Debug("distance= %f", distance)
	return distance
}
*/

// updateLinks queries the router objects' interfaces and calculates the bitsPerSec. It then updates the links'
//
//	lineWidth and color to reflect the amount of traffic flowing over each link.
//
// func updateLinks(log *logger.Logger, gv *gvapp, databaseForRead *sql.DB, snmpTarget string, community string, params *g.GoSNMP) *gvapp {
func updateLinks(log *logger.Logger, gv *gvapp, databaseForRead *sql.DB, snmpTarget string, params *g.GoSNMP) *gvapp {
	log.Info("Updating Links")
	// TODO
	//	1) Add RouterID to Links DB table - DONE
	//	2) For each Link in Links DB table - DONE
	//		2.1) query the From Routers applicable interfaces outbound traffic - DONE
	//		2.2) calculate bitsPerSec on the interface - DONE
	//	3) Update the link lineWidth and color to indicate the calculated link Utilization
	var LinkID int
	var FromRouterID, FromRouterName, FromRouterIfIndex, ToRouterName string

	//
	// Retrieve the links from the database
	//

	linkRows, err := databaseForRead.Query("SELECT LinkID, FromRouterID, FromRouterName, FromRouterIfIndex, ToRouterName FROM Links")
	if err != nil {
		databaseForRead.Close()
		log.Fatal("databaseForRead Query Router error %v", err)
	}
	log.Debug("Successful Links table Query")
	defer linkRows.Close()

	//
	// Loop through all the links
	//
	for linkRows.Next() {

		//
		// Calculate link bits per second
		//
		linkRows.Scan(&LinkID, &FromRouterID, &FromRouterName, &FromRouterIfIndex, &ToRouterName)

		// SNMP Get Router Interface Outbound traffic (ifOutOctets 1)
		oids := []string{
			ifOutOctets + "." + FromRouterIfIndex, // ifOutOctets
			ifSpeed + "." + FromRouterIfIndex,     // ifSpeed
		}

		snmpTarget = FromRouterName
		params.Target = snmpTarget
		err := params.Connect()
		if err != nil {
			databaseForRead.Close()
			log.Fatal("Connect() err: %v", err)
		}
		defer params.Conn.Close()

		result, err := params.Get(oids) // Get() accepts up to g.MAX_OIDS
		if err != nil {
			if strings.Contains(err.Error(), "Request timeout") {
				log.Warn("Router %s", snmpTarget+" SNMP Timeout. Continuing.")
			} else {
				log.Warn("SNMP Get() error: %v", err)
			}
		}
		if result == nil {
			log.Warn("No Router Interface found for Link definition. Error in LinkID: %d", LinkID)
		} else {
			ifOutOctets1 := result.Variables[0].Value.(uint)
			ifSpeed1 := result.Variables[1].Value.(uint)

			// Sleep for 1 second
			time.Sleep(1 * time.Second)

			// SNMP Get Router Interface Outbound (ifOutOctets 2)
			result, err = params.Get(oids) // Get() accepts up to g.MAX_OIDS
			if err != nil {
				databaseForRead.Close()
				log.Fatal("Get() err %v", err)
			}
			ifOutOctets2 := result.Variables[0].Value.(uint)

			// Calculate differnce of ifOutOctets 2 and ifOutOctets 1) and multiply by 8.
			// This is the approx. number of bits per second.
			bitsPerSec := (ifOutOctets2 - ifOutOctets1) * 8
			log.Debug("LinkID %s", strconv.Itoa(LinkID)+" (From "+FromRouterName+"To "+ToRouterName+"): bps= "+strconv.FormatInt(int64(bitsPerSec), 10))

			// set linkColor and linkWidth, depending on linkUtilization
			var linkUtil = float32(bitsPerSec / ifSpeed1)
			var linkColor = math32.ColorName("white")
			var linkWidth float32 = 1.0

			//linkUtil = 0.76 // **** TESTING ONLY. REMOVE WHEN TESTING COMPLETED ****

			if linkUtil < 0.75 {
				linkColor = math32.ColorName("lime") // Lime Green
				linkWidth = 1.0
			} else if linkUtil >= 0.75 && linkUtil < 0.90 {
				linkColor = math32.ColorName("yellow") // Yellow
				linkWidth = 2.5
			} else if linkUtil >= 0.90 && linkUtil <= 1.0 {
				linkColor = math32.ColorName("red") // Red
				linkWidth = 5.0
			}
			log.Debug("linkColor = %v", linkColor)
			log.Debug("linkWidth = %f", linkWidth)

			//
			// Update Link lineWidth and color, depending on link utilization
			//

			// Find 3D line object
			sceneChildren := gv.scene.Children()
			log.Debug("scene.Name %s", gv.scene.Name())
			log.Debug("sceneChildren: %v", sceneChildren)

			// loop: parse scene Children getting each node
			link := getRouterFromScene(log, sceneChildren, LinkID)
			log.Debug("link found: %v", link)
			if link == nil {
				log.Warn("updateLinks: link not found in 3D scene.")
				continue
			}

			//
			// Set line object color
			//
			// Convert INode to IGraphic
			ig, ok := link.(graphic.IGraphic)
			if !ok {
				databaseForRead.Close()
				log.Fatal("Error when converting link INode to IGraphic")
			}
			// Get graphic object
			gr := ig.GetGraphic()
			imat := gr.GetMaterial(0)

			type matI interface {
				EmissiveColor() math32.Color
				SetEmissiveColor(*math32.Color)
				//AmientColor() math32.Color
				//SetAmbientColor(*math32.Color)
				SetLineWidth(float32)
			}
			v := imat.(matI)
			//v.SetEmissiveColor(&math32.Color{R: 0, G: 1, B: 0})
			v.SetEmissiveColor(&linkColor)
			//v.SetAmbientColor(&linkColor)

			// Set line object width
			// Check Runtime environment.
			// OpenGL Implementation on MacOS will only accept Line width of 1.0
			if runtime.GOOS == "darwin" {
				v.SetLineWidth(1.0)
				log.Info("*** Link SetLineWidth() request ignored. OpenGL Implementation on MacOS will only accept lineWidth of 1.0 ***")
			} else {
				v.SetLineWidth(linkWidth)
			}

			gr.SetChanged(true)
			gr.Render(gv.Application.Gls())
		}
	}

	log.Info("Links Updated.")

	return (gv)
}

func getRouterFromScene(log *logger.Logger, sceneChildren []core.INode, LinkID int) core.INode {
	log.Debug("getRouterFromScene Started.")
	var link core.INode
	for i := 0; i < len(sceneChildren); i++ {
		if sceneChildren[i].Name() == strconv.Itoa(LinkID) {
			link = sceneChildren[i]
			break
		}
	}
	return link
}

// Render renders the mouse pick action
func (t *Raycast) Render(a *app.Application) {
}

// Update is called every frame.
func (t *Raycast) Update(a *app.Application, deltaTime time.Duration) {}

// Cleanup is called once at the end of the demo.
func (t *Raycast) Cleanup(a *app.Application) {}

// Add a title to the scene
func addTitle(log *logger.Logger, gv *gvapp) *gvapp {

	//	fontfile := os.Getenv("GOBIN") + "/data/fonts/FreeSans.ttf"
	//	font, err := text.NewFont(fontfile)
	//	if err != nil {
	//		databaseForRead.Close()
	//		log.Fatal("Error loading font %s" + err.Error() + "\n Insure govisn /data/fonts is copied to GOBIN \n GOBIN env variable must be set.")
	//	}

	//	font.SetLineSpacing(1.0)
	//	font.SetPointSize(100)
	//	font.SetDPI(72)
	//	font.SetFgColor(&math32.Color4{R: 0, G: 0, B: 0, A: 1})
	//	font.SetBgColor(&math32.Color4{R: 1, G: 1, B: 1, A: 0.8})
	//canvas := text.NewCanvas(300, 200, &math32.Color4{R: 1, G: 1, B: 1, A: 0.8})
	//	rtext := "GoVisn version " + GOVISNVERSION
	//	swidth, sheight := font.MeasureText(rtext)
	//	canvas := text.NewCanvas(swidth, sheight, &math32.Color4{R: 1, G: 1, B: 1, A: 1})
	//	canvas.DrawText(0, 0, rtext, font)
	//	tex3 := texture.NewTexture2DFromRGBA(canvas.RGBA)
	//	mat3 := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
	//	mat3.AddTexture(tex3)
	//	aspect := float32(swidth) / float32(sheight)
	//	mesh3 := graphic.NewSprite(aspect, 1, mat3)
	//	mesh3.SetPosition(10.0, 10.0, 0.0)
	//	gv.scene.Add(title)

	log.Debug("addTitle Started.")
	titleLines := []string{
		"                  GoVisn version " + GOVISNVERSION,
		"\nNetwork Visualization in 3D from " + *DbName,
		//		"Copyright 2020 Kevin Hayes Parrish",
		//		"All rights reserved.",
	}
	title := gui.NewLabel(strings.Join(titleLines, "  "))
	//title := gui.NewLabel("GoVisn version " + GOVISNVERSION)
	title.SetPosition(100, 0)
	title.SetBordersColor(math32.NewColor("grey"))
	title.SetBgColor(math32.NewColor("blue"))
	title.SetColor(math32.NewColor("white"))
	title.SetBorders(1, 1, 1, 1)
	title.SetPaddings(5, 5, 5, 5)
	title.SetFontSize(15)
	gv.scene.Add(title)

	return gv
}
