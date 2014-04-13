package main

import (
	zmq "github.com/pebbe/zmq4"
)

func Proxy(frontend *zmq.Socket, backend *zmq.Socket, errChan chan error) {
	errChan <- zmq.Proxy(frontend, backend, nil)
}

func Forwarder(xpub *zmq.Socket, xsub *zmq.Socket) error {
	errChan := make(chan error)
	go Proxy(xpub, xsub, errChan)
	go Proxy(xsub, xpub, errChan)
	return <-errChan
}
