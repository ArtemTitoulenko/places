package main

import (
	"./placedata"
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

var placeDataFilename = "./placedata/places-proto.bin"
var placeList []*placedata.Place

type PlaceQuery struct {
	Name, State string
}

type Places int

func printUsage() {
	fmt.Println("Usage:\n\tplace-server [-p port]\nDefaults:")
	flag.PrintDefaults()
}

// load the list of places into placeList
func loadPlaces(file string) {
	fmt.Println("Loading proto data file")
	var ret, err = getPlaceList(file)
	if err != nil {
		log.Fatalf("Unable to load the proto data file %s", file)
	}
	placeList = ret.GetPlace()
}

func main() {
	var port int
	var help bool
	var placeDataFileLocation string
	flag.IntVar(&port, "p", 1080, "The port on which to listen for connections")
	flag.IntVar(&port, "port", 1080, "The port on which to listen for connections")
	flag.BoolVar(&help, "help", false, "Print usage information")
	flag.StringVar(&placeDataFileLocation, "place-data", placeDataFilename, "Which places file to use")
	flag.Parse()

	if help {
		printUsage()
		return
	}

	loadPlaces(placeDataFileLocation)

	fmt.Printf("Listening on port: %d\n", port)

	// register the Places service
	rpc.Register(new(Places))
	rpc.HandleHTTP()

	socket, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(socket, nil)
}

// Place Service
func getPlaceList(filename string) (*placedata.PlaceList, error) {
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
	var placeList = &placedata.PlaceList{}
	err = proto.Unmarshal(fileContents, placeList)
	if err != nil {
		return nil, err
	}

	return placeList, nil
}

// return a place to the caller if one is found
func (t *Places) Find(query *PlaceQuery, result *[]byte) error {
	for _, place := range placeList {
		if strings.EqualFold(place.GetName(), query.Name) &&
			strings.EqualFold(place.GetState(), query.State) {
			*result, _ = proto.Marshal(place)
			return nil
		}
	}

	return errors.New("could not find place")
}
