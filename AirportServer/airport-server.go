package main

import (
	"./airportdata"
	"code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sort"
)

const (
	RadToDeg = 180.0 / math.Pi
	DegToRad = math.Pi / 180.0
	R        = 6371.0 // earth's radius in km
)

var airportDataFilename = "./airportdata/airports-proto.bin"
var airportList []*airportdata.Airport

type AirportQuery struct {
	Lat, Lon float64
}

type Airports int

func printUsage() {
	fmt.Println("Usage:\n\tairport-server [-p port]\nDefaults:")
	flag.PrintDefaults()
}

// load the list of airports before main() runs
func loadAirports(file string) {
	fmt.Println("Loading proto data file")
	var ret, err = getAirportList(file)
	if err != nil {
		log.Fatalf("Unable to load the proto data file %s", file)
	}
	airportList = ret.GetAirport()
}

func main() {
	var port int
	var help bool
	var airportDataFileLocation string
	flag.IntVar(&port, "p", 1082, "The port on which to listen for connections")
	flag.IntVar(&port, "port", 1082, "The port on which to listen for connections")
	flag.BoolVar(&help, "help", false, "Print usage information")
	flag.StringVar(&airportDataFileLocation, "airport-data", airportDataFilename, "Which airport location file to use")
	flag.Parse()

	if help {
		printUsage()
		return
	}

	loadAirports(airportDataFileLocation)

	fmt.Printf("Listening on port: %d\n", port)

	// register the Airports service
	rpc.Register(new(Airports))
	rpc.HandleHTTP()

	socket, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(socket, nil)
}

// Extract the list of airports from the binary file
func getAirportList(filename string) (*airportdata.AirportList, error) {
	// check if the file exists
	if _, err := os.Stat(string(filename)); os.IsNotExist(err) {
		return nil, err
	}

	// read the contents
	var fileContents, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// extract the PlaceList
	var airportList = &airportdata.AirportList{}
	err = proto.Unmarshal(fileContents, airportList)
	if err != nil {
		return nil, err
	}

	return airportList, nil
}

func (t *Airports) Find(query *AirportQuery, result *[]byte) error {
	var (
		currentLat         = query.Lat
		currentLon         = query.Lon
		airportDistances   = make([]float64, len(airportList))
		airportDistanceMap = make(map[float64]*airportdata.Airport, len(airportList))
	)

	// calculate distances to all airports, map distance to airport, collisions are unlikely
	for i, airport := range airportList {
		var d = calculateGreaterCircleDistance(currentLat, currentLon, airport.GetLat(), airport.GetLon())
		airportDistances[i] = d
		airportDistanceMap[d] = airport
	}

	// hold a list of the 5 closest airports
	airportList := make([]*airportdata.Airport, 5)

	// sort the list of distances
	sort.Float64s(airportDistances)
	for i := 0; i < 5; i++ {
		airportList[i] = airportDistanceMap[airportDistances[i]]
	}

	// make a new airport list and send it back to the client
	*result, _ = proto.Marshal(&airportdata.AirportList{airportList, nil})
	return nil
}

// calculate the greater circle distance between two coordinates
func calculateGreaterCircleDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var a = math.Sin(lat1*DegToRad) * math.Sin(lat2*DegToRad)
	var b = math.Cos(lat1*DegToRad) * math.Cos(lat2*DegToRad) * math.Cos(lon2*DegToRad-lon1*DegToRad)
	return 60.0 * math.Acos(a+b) * RadToDeg
}
