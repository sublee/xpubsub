package main

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
)

var xpubPort = flag.Int("xpub", 5561, "listening port for XPUB socket")
var xsubPort = flag.Int("xsub", 5562, "listening port for XSUB socket")

func main() {
	// parse command-line arguments
	flag.Parse()
	// init sockets
	xpub, _ := zmq.NewSocket(zmq.XPUB)
	defer xpub.Close()
	xsub, _ := zmq.NewSocket(zmq.XSUB)
	defer xsub.Close()
	// bind to ports
	xpub.Bind(fmt.Sprintf("tcp://*:%d", *xpubPort))
	xsub.Bind(fmt.Sprintf("tcp://*:%d", *xsubPort))
	log.Printf("Forwarding messages from %d to %d...", *xpubPort, *xsubPort)
	// run forwarder
	err := Forwarder(xpub, xsub)
	if err != nil {
		log.Fatal(err)
	}
}
