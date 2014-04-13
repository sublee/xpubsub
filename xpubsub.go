package main

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
)

var xpubPort = flag.Int("xpub", 5561, "listening port for XPUB socket")
var xsubPort = flag.Int("xsub", 5562, "listening port for XSUB socket")

func Proxy(frontend *zmq.Socket, backend *zmq.Socket, errChan chan error) {
	errChan <- zmq.Proxy(frontend, backend, nil)
}

func main() {
	flag.Parse()

	xpub, _ := zmq.NewSocket(zmq.XPUB)
	defer xpub.Close()
	xsub, _ := zmq.NewSocket(zmq.XSUB)
	defer xsub.Close()
	xpub.Bind(fmt.Sprintf("tcp://*:%d", *xpubPort))
	xsub.Bind(fmt.Sprintf("tcp://*:%d", *xsubPort))
	log.Printf("Forwarding messages from %d to %d...", *xpubPort, *xsubPort)

	errChan := make(chan error)
	go Proxy(xpub, xsub, errChan)
	go Proxy(xsub, xpub, errChan)
	err := <-errChan
	log.Fatal(err)
}
