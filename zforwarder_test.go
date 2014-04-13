package zforwarder_test

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/sublee/zforwarder"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	// set sockets
	pub, _ := zmq.NewSocket(zmq.PUB)
	sub, _ := zmq.NewSocket(zmq.SUB)
	xpub, _ := zmq.NewSocket(zmq.XPUB)
	xsub, _ := zmq.NewSocket(zmq.XSUB)
	xpub.Bind("inproc://xpub")
	xsub.Bind("inproc://xsub")
	pub.Connect("inproc://xsub")
	sub.Connect("inproc://xpub")
	// run forwarder
	go zforwarder.Forwarder(xpub, xsub)
	// communicate between PUB/SUB via XPUB/XSUB
	sub.SetSubscribe("")
	time.Sleep(time.Nanosecond)
	msgExpected := "Hello"
	pub.Send(msgExpected, 0)
	msg, err := sub.Recv(0)
	// check result
	if err != nil {
		t.Fatalf("An error occurs by SUB socket: %s", err)
	}
	if msg != msgExpected {
		t.Fatalf("Received message: %s", msg)
	}
}
