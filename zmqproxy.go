package main

import "fmt"
import "github.com/droundy/goopt"
import zmq "github.com/pebbe/zmq4"
import "log"
import "time"

var queueDevice = goopt.Flag(
	[]string{"-Q", "--queue"}, nil, "choose 'queue' device", "")
var forwarderDevice = goopt.Flag(
	[]string{"-F", "--forwarder"}, nil, "choose 'forwarder' device", "")
var streamerDevice = goopt.Flag(
	[]string{"-S", "--streamer"}, nil, "choose 'streamer' device", "")
var frontendPort = goopt.Int(
	[]string{"-f", "--frontend"}, 5561,
	"listening port the frontend socket binds to.")
var backendPort = goopt.Int(
	[]string{"-b", "--backend"}, 5562,
	"listening port the backend socket binds to.")
var trafficDisabled = goopt.Flag(
	[]string{"--no-traffic"}, nil, "disable traffic reporting.", "")

type Traffic struct {
	Socket  *zmq.Socket
	bytes   int
	msgs    int
	ResetAt time.Time
}

func NewTraffic(socket *zmq.Socket) Traffic {
	traffic := Traffic{Socket: socket}
	traffic.Reset()
	return traffic
}

func (s *Traffic) Reset() {
	s.bytes = 0
	s.msgs = 0
	s.ResetAt = time.Now()
}

func (s *Traffic) Collect(bytes []byte) {
	s.bytes += len(bytes)
	s.msgs += 1
}

func (s *Traffic) Report(print func(float32, float32)) {
	now := time.Now()
	seconds := float32(now.Sub(s.ResetAt) / time.Second)
	bps := float32(s.bytes) / seconds
	mps := float32(s.msgs) / seconds
	print(bps, mps)
	s.Reset()
}

func (s *Traffic) CollectForever() {
	for {
		bytes, _ := s.Socket.RecvBytes(0)
		s.Collect(bytes)
	}
}

func (s *Traffic) ReportForever(
	duration time.Duration, print func(float32, float32)) {
	for {
		time.Sleep(duration)
		s.Report(print)
	}
}

func main() {
	// parse argv
	goopt.Summary = "Runs ZeroMQ proxy."
	goopt.Parse(nil)
	// init sockets by device
	var device string
	var types []zmq.Type
	switch {
	case *queueDevice:
		device = "queue"
		types = []zmq.Type{zmq.ROUTER, zmq.DEALER}
		break
	case *forwarderDevice:
		device = "forwarder"
		types = []zmq.Type{zmq.XSUB, zmq.XPUB}
		break
	case *streamerDevice:
		device = "streamer"
		types = []zmq.Type{zmq.PULL, zmq.PUSH}
		break
	default:
		log.Fatal("Choose device by (-Q|-F|-S) option")
	}
	frontend, _ := zmq.NewSocket(types[0])
	defer frontend.Close()
	backend, _ := zmq.NewSocket(types[1])
	defer backend.Close()
	log.Printf("ZeroMQ '%s' device chosen", device)
	// bind to the ports
	frontend.Bind(fmt.Sprintf("tcp://*:%d", *frontendPort))
	backend.Bind(fmt.Sprintf("tcp://*:%d", *backendPort))
	// collect and report traffic
	var capture *zmq.Socket
	if *trafficDisabled {
		capture = nil
		log.Println("Traffic reporting disabled")
	} else {
		capture, _ = zmq.NewSocket(zmq.PAIR)
		defer capture.Close()
		captured, _ := zmq.NewSocket(zmq.PAIR)
		defer captured.Close()
		capture.Bind("inproc://capture")
		captured.Connect("inproc://capture")
		traffic := NewTraffic(captured)
		print := func(mps float32, bps float32) {
			log.Printf("Traffic: %.2f msgs/sec (%.2f bytes/sec)", mps, bps)
		}
		go traffic.CollectForever()
		go traffic.ReportForever(time.Minute, print)
		log.Println("Traffic reporting enabled")
	}
	// run proxy
	frontendType, _ := frontend.GetType()
	backendType, _ := backend.GetType()
	log.Printf("Proxying between %d[%s] and %d[%s]...",
		*frontendPort, frontendType, *backendPort, backendType)
	err := zmq.Proxy(frontend, backend, capture)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Done")
	}
}
