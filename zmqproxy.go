package main

import (
	"fmt"
	"github.com/docopt/docopt.go"
	zmq "github.com/pebbe/zmq4"
	"log"
	"strconv"
)

var (
	frontend         *zmq.Socket
	backend          *zmq.Socket
	frontendTypeName string
	backendTypeName  string
)

func main() {
	usage := `Runs ZeroMQ proxy.

  Usage: zmqproxy (queue|forwarder|streamer) [options]

  Options:
    -f, --frontend PORT    Port the frontend socket binds to. [default: 5561]
    -b, --backend PORT     Port the backend socket binds to. [default: 5562]`
	args, _ := docopt.Parse(usage, nil, true, "", false)
	// set device
	if args["queue"].(bool) {
		frontend, _ = zmq.NewSocket(zmq.ROUTER)
		backend, _ = zmq.NewSocket(zmq.DEALER)
		frontendTypeName = "ROUTER"
		backendTypeName = "DEALER"
	} else if args["forwarder"].(bool) {
		frontend, _ = zmq.NewSocket(zmq.XSUB)
		backend, _ = zmq.NewSocket(zmq.XPUB)
		frontendTypeName = "XSUB"
		backendTypeName = "XPUB"
	} else if args["streamer"].(bool) {
		frontend, _ = zmq.NewSocket(zmq.PULL)
		backend, _ = zmq.NewSocket(zmq.PUSH)
		frontendTypeName = "PULL"
		backendTypeName = "PUSH"
	}
	// bind to the ports
	frontendPort, _ := strconv.ParseInt(args["--frontend"].(string), 10, 16)
	backendPort, _ := strconv.ParseInt(args["--backend"].(string), 10, 16)
	frontend.Bind(fmt.Sprintf("tcp://*:%d", frontendPort))
	backend.Bind(fmt.Sprintf("tcp://*:%d", backendPort))
	// run proxy
	log.Printf(
		"Running ZeroMQ proxy from %d[%s] to %d[%s]...",
		frontendPort, frontendTypeName, backendPort, backendTypeName)
	err := zmq.Proxy(frontend, backend, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Done")
	}
}
