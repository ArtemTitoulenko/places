Artem Titoulenko  
CS 417: Distributed Systems  
Homework 5  

# Distributed Airport Locator Service

This collection of servers and a client allows a user to find the five nearest airports to a place. There is a list of places that is not updated and will not be updated.

## Installation

In either `AirportServer` or `PlaceServer`, run:

	$ go get

This should install the `code.google.com/p/goprotobuf/proto` package which is obviously a protobuffer package for Go.

## Running the PlaceServer

In the `PlaceServer` directory, you may either run:

	$ go run ./place-server.go

or you may compile the program and run it like so:

	$ go build
	$ ./PlaceServer

You can check what options are available by passing the `--help` flag. You should see something like:

	Usage:
		place-server [-p port]
	Defaults:
	  -help=false: Print usage information
	  -p=1080: The port on which to listen for connections
	  -port=1080: The port on which to listen for connections

## Running the AirportServer

In the `AirportServer` directory, you may either run:

	$ go run ./airport-server.go

or you may compile the program and run it like so:

	$ go build
	$ ./AirportServer

You can check what options are available by passing the `--help` flag. You should see something like:

	Usage:
		airport-server [-p port]
	Defaults:
	  -help=false: Print usage information
	  -p=1082: The port on which to listen for connections
	  -port=1082: The port on which to listen for connections

## Running the Client

In the `Client` directory, you may either run:

	$ go run ./client.go

or you may compile the program and run it like so:

	$ go build
	$ ./Client

You can check what options are available by passing the `--help` flag. You should see something like:
	
	Usage:
		client [--place-host host] [--place-port port] [--airport-host host] [--airport-port port] [--help] [--kilometers] city state
	Defaults:
	  -airport-host="localhost": The AirportServer host
	  -airport-port=1082: The AirportServer listening port
	  -help=false: Print help
	  -kilometers=false: Display distances in kilometers
	  -place-host="localhost": The PlaceServer host
	  -place-port=1080: The PlaceServer listening port

### Usage

The Client application takes a city and a state as parameters and returns the five nearest airports, ordered by distancce. The units of measure default to miles but can be displayed in kilometers by passing the `--kilometers` flag. City names must be exact, and state's are denoted by their two-letter abbreviations.

An example run may look like:

	$ go run client.go "Abbeville city" AL
	Connecting to PlaceServer at localhost:1080
	Connecting to AirportServer at localhost:1082
	
	DHN: Dothan, AL distance: 38.27 miles
	OZR: Fort Rucker, AL distance: 62.94 miles
	TOI: Troy, AL distance: 92.25 miles
	MAI: Marianna, FL distance: 93.27 miles
	LSF: Fort Benning, GA distance: 101.45 miles

## Bugs and Discrepancies

Distances calculated by the airport server and conversely the client application are off by