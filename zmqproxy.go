package main

import (
	"fmt"
	"github.com/droundy/goopt"
	zmq "github.com/pebbe/zmq4"
	"log"
)

var device = goopt.Alternatives(
	[]string{"-d", "--device"}, []string{"queue", "forwarder", "streamer"},
	"device type to run.")
var frontendPort = goopt.Int(
	[]string{"-f", "--frontend"}, 5561,
	"listening port the frontend socket binds to.")
var backendPort = goopt.Int(
	[]string{"-b", "--backend"}, 5562,
	"listening port the backend socket binds to.")

func main() {
	// parse argv
	goopt.Summary = "Runs ZeroMQ proxy."
	goopt.Parse(nil)
	// init sockets by device
	devices := map[string][]zmq.Type{
		"queue": []zmq.Type{zmq.ROUTER, zmq.DEALER},
		"forwarder": []zmq.Type{zmq.XSUB, zmq.XPUB},
		"streamer": []zmq.Type{zmq.PULL, zmq.PUSH},
	}
	types := devices[*device]
	frontend, _ := zmq.NewSocket(types[0])
	defer frontend.Close()
	backend, _ := zmq.NewSocket(types[1])
	defer backend.Close()
	log.Printf("ZeroMQ device '%s' selected", *device)
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
