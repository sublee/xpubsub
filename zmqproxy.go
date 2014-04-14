package main

import (
	"fmt"
	"github.com/droundy/goopt"
	zmq "github.com/pebbe/zmq4"
	"log"
)

var device = goopt.Alternatives(
	[]string{"--device"}, []string{"queue", "forwarder", "streamer"},
	"device type to run.")
var frontendPort = goopt.Int(
	[]string{"-f", "--frontend"}, 5561,
	"listening port the frontend socket binds to.")
var backendPort = goopt.Int(
	[]string{"-b", "--backend"}, 5562,
	"listening port the backend socket binds to.")

func queue() (*zmq.Socket, *zmq.Socket) {
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	backend, _ := zmq.NewSocket(zmq.DEALER)
	return frontend, backend
}

func forwarder() (*zmq.Socket, *zmq.Socket) {
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	backend, _ := zmq.NewSocket(zmq.XPUB)
	return frontend, backend
}

func streamer() (*zmq.Socket, *zmq.Socket) {
	frontend, _ := zmq.NewSocket(zmq.PULL)
	backend, _ := zmq.NewSocket(zmq.PUSH)
	return frontend, backend
}

func main() {
	goopt.Summary = "Runs ZeroMQ proxy."
	goopt.Parse(nil)
	// init sockets by device
	devices := map[string]func() (*zmq.Socket, *zmq.Socket){
		"queue":     queue,
		"forwarder": forwarder,
		"streamer":  streamer,
	}
	frontend, backend := devices[*device]()
	defer frontend.Close()
	defer backend.Close()
	log.Printf("Device '%s' selected", *device)
	// bind to the ports
	frontend.Bind(fmt.Sprintf("tcp://*:%d", *frontendPort))
	backend.Bind(fmt.Sprintf("tcp://*:%d", *backendPort))
	// run proxy
	frontendType, _ := frontend.GetType()
	backendType, _ := backend.GetType()
	log.Printf("Proxying between %d[%s] and %d[%s]...",
		*frontendPort, frontendType, *backendPort, backendType)
	err := zmq.Proxy(frontend, backend, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Done")
	}
}
