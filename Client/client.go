package main

// sample places
/*

	state:"AL" name:"Abbeville city" lat:31.566367 lon:-85.2513
	state:"AL" name:"Adamsville city" lat:33.590411 lon:-86.949166
	state:"AL" name:"Addison town" lat:34.200042 lon:-87.177851
	state:"AL" name:"Akron town" lat:32.876425 lon:-87.740978

*/

import (
	"../AirportServer/airportdata"
	"../PlaceServer/placedata"
	"code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	"log"
	"math"
	"net/rpc"
)

const (
	RadToDeg = 180 / math.Pi
	DegToRad = math.Pi / 180
	R        = 6371 // earth's radius in km
)

type PlaceQuery struct {
	Name, State string
}

type AirportQuery struct {
	Lat, Lon float64
}

func printUsage() {
	fmt.Println("Usage:\n\tclient [--place-host host] [--place-port port] [--airport-host host] [--airport-port port] [--help] [--kilometers] city state\nDefaults:")
	flag.PrintDefaults()
}

func main() {
	var placeServerHost string
	var placeServerPort int
	var airportServerHost string
	var airportServerPort int
	var help bool
	var kilometers bool
	flag.StringVar(&placeServerHost, "place-host", "localhost", "The PlaceServer host")
	flag.IntVar(&placeServerPort, "place-port", 1080, "The PlaceServer listening port")
	flag.StringVar(&airportServerHost, "airport-host", "localhost", "The AirportServer host")
	flag.IntVar(&airportServerPort, "airport-port", 1082, "The AirportServer listening port")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.BoolVar(&kilometers, "kilometers", false, "Display distances in kilometers")
	flag.Parse()

	if help || len(flag.Args()) != 2 {
		printUsage()
		return
	}

	// try to connect to the PlaceServer
	fmt.Printf("Connecting to PlaceServer at %s:%d\n", placeServerHost, placeServerPort)
	placeServer, err := rpc.DialHTTP("tcp", placeServerHost+":"+fmt.Sprint(placeServerPort))
	if err != nil {
		log.Fatalf("place server: could not connect to the place server service. message: ", err)
	}

	// try to connect to the AirportServer
	fmt.Printf("Connecting to AirportServer at %s:%d\n\n", airportServerHost, airportServerPort)
	airportServer, err := rpc.DialHTTP("tcp", airportServerHost+":"+fmt.Sprint(airportServerPort))
	if err != nil {
		log.Fatalf("airport server: could not connect to the airport server service. message: ", err)
	}

	// create a query to send to the PlaceServer
	place, err := getPlaceDetails(placeServer, &PlaceQuery{flag.Arg(0), flag.Arg(1)})
	if err != nil {
		log.Fatalln("could not get place details: " + err.Error())
		return
	}

	// lets get the airports
	airportList, err := getNearestAirports(airportServer, &AirportQuery{place.GetLat(), place.GetLon()})
	if err != nil {
		log.Fatalln("could not get nearest airports: " + err.Error())
		return
	}

	// print out the returned airports
	for _, airport := range airportList.GetAirport() {
		fmt.Printf("%s: %s, %s ",
			airport.GetCode(),
			airport.GetName(),
			airport.GetState())

		var distance = calculateDistance(place.GetLat(), place.GetLon(), airport.GetLat(), airport.GetLon())
		if kilometers {
			fmt.Printf("distance: %.2f km\n", distance*1.85200)
		} else {
			fmt.Printf("distance: %.2f miles\n", distance*1.1507794)
		}
	}
}

func getPlaceDetails(placeServer *rpc.Client, query *PlaceQuery) (*placedata.Place, error) {
	var placeData []byte

	// use the Places.Find service to get details about a place
	err := placeServer.Call("Places.Find", query, &placeData)
	if err != nil {
		placeServer.Close()
		return nil, err
	}

	// unmarshal the returned place and close the socket
	var place = &placedata.Place{}
	proto.Unmarshal(placeData, place)
	placeServer.Close()

	return place, nil
}

func getNearestAirports(airportServer *rpc.Client, query *AirportQuery) (*airportdata.AirportList, error) {
	var airportListData []byte

	// use the Airports.Find service to get nearest airports
	err := airportServer.Call("Airports.Find", query, &airportListData)
	if err != nil {
		airportServer.Close()
		return nil, err
	}

	// unmarshal the airport list and close the socket
	var airportList = &airportdata.AirportList{}
	proto.Unmarshal(airportListData, airportList)
	airportServer.Close()

	return airportList, nil
}

// haversine formula for getting greater-circle distance between two points over the earth's surface
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var dLat = (lat2 - lat1) * DegToRad
	var dLon = (lon2 - lon1) * DegToRad
	lat1 = lat1 * DegToRad
	lat2 = lat2 * DegToRad

	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c

	return d
}
