package main

import (
	"strconv"
)

// Get Router Coordinates from routerArray
func getRouterCoordinatesName(debug bool, routers []Router, routerName string) (float32, float32, float32) {
	var x float32
	var y float32
	var z float32

	//	for i := 0; i < len(routerArray); i++ {
	for i := 0; i < len(routers); i++ {
		//		if routerArray[i].System.Name == "" {
		//			break // end of routerArray entries
		//		}
		//		if routerArray[i].System.Name == routerName {
		if routers[i].System.Name == routerName {
			//			x = routerArray[i].System.Coordinates.X
			//			x1, err := strconv.ParseFloat(routerArray[i].System.GPS.Longitude, 32)
			x1, err := strconv.ParseFloat(routers[i].System.GPS.Longitude, 32)
			if err != nil {
				panic(err)
			}
			x = (float32)(x1)
			//			y = routerArray[i].System.Coordinates.Y
			//			y1, err := strconv.ParseFloat(routerArray[i].System.GPS.Latitude, 32)
			y1, err := strconv.ParseFloat(routers[i].System.GPS.Latitude, 32)
			if err != nil {
				panic(err)
			}
			y = (float32)(y1)
			//			z = routerArray[i].System.Coordinates.Z
			//			z1, err := strconv.ParseFloat(routerArray[i].System.GPS.Altitude, 32)
			z1, err := strconv.ParseFloat(routers[i].System.GPS.Altitude, 32)
			if err != nil {
				panic(err)
			}
			z = (float32)(z1)
			break
		}
	}
	return x, y, z
}
