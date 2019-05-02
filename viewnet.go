package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type DiscoveredNetwork struct {
	XMLName xml.Name `xml:"Discovered_Network"`
	Routers Router[]	`xml:"Router"`
}
type Router  struct {
		Text   string `xml:",chardata"`
		System struct {
			Text        string `xml:",chardata"`
			Name        string `xml:"Name"`
			Description string `xml:"Description"`
			UpTime      string `xml:"Up_Time"`
			Contact     string `xml:"Contact"`
			Location    string `xml:"Location"`
			GPS         struct {
				Text      string `xml:",chardata"`
				Latitude  string `xml:"Latitude"`
				Longitude string `xml:"Longitude"`
				Altitude  string `xml:"Altitude"`
			} `xml:"GPS"`
		} `xml:"System"`
		Addresses struct {
			Text             string `xml:",chardata"`
			NetworkAddresses struct {
				Text      string   `xml:",chardata"`
				IPAddress []string `xml:"IP_Address"`
			} `xml:"Network_Addresses"`
			MediaAddresses struct {
				Text         string `xml:",chardata"`
				MediaAddress string `xml:"Media_Address"`
			} `xml:"Media_Addresses"`
		} `xml:"Addresses"`
		Neighbors struct {
			Text     string `xml:",chardata"`
			Neighbor []struct {
				Text               string `xml:",chardata"`
				DestinationAddress string `xml:"Destination_Address"`
				NextHop            string `xml:"Next_Hop"`
			} `xml:"Neighbor"`
		} `xml:"Neighbors"`
	} `xml:"Router"`
}

func main() {

	// Open the network xmlFile
	xmlFile, err := os.Open("V15N-GPS-5.xml")
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened V15N-GPS-5.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var network DiscoveredNetwork
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'netework' which we defined above
	xml.Unmarshal(byteValue, &network)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(network.DiscoveredNetwork.Router); i++ {
		fmt.Println("Router Name: " + network.Router[i].System.Name)
	}
}
