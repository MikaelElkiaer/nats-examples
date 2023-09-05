package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

func main() {
	natsHost, exists := os.LookupEnv("NATS_HOST")
	if !exists {
		logger.Fatal("NATS_HOST is not set.")
	}

	natsSubject, exists := os.LookupEnv("NATS_SUBJECT")
	if !exists {
		logger.Fatal("NATS_SUBJECT is not set.")
	}

	clientType, exists := os.LookupEnv("NATS_CLIENT_TYPE")
	if !exists {
		logger.Fatal("NATS_CLIENT_TYPE is not set, must be one of: responder, requester")
	}

	connectionClosed := make(chan bool, 1)
	nc, err := nats.Connect(natsHost,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ClosedHandler(func(c *nats.Conn) {
			connectionClosed <- true
		}),
		nats.ReconnectWait(time.Second))

	if err != nil {
		logger.Fatalf("Failed to connect: %v", err)
	}

	if clientType == "responder" {
		doSubscribeAndRespond(nc, natsSubject)
	} else if clientType == "requester" {
		doRequest(nc, natsSubject)
	} else {
		logger.Fatalf("Unsupported type: %s", clientType)
	}

	nc.Drain()

	// Wait for the connection to finish draining
	<-connectionClosed
}

func doSubscribeAndRespond(nc *nats.Conn, natsSubject string) {
	logger.Printf("Subscribing to: %s", natsSubject)
	sub, err := nc.Subscribe(natsSubject, subHandler)

	if err != nil {
		logger.Fatalf("Failed to queue subscribe: %v", err)
	}

	// Block until an interrupt is received
	c := make(chan os.Signal, 1)

	// SIGTERM is what is actually sent in kubernetes pods, but let's also do a graceful shutdown for other signals
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for shutdown signal
	<-c

	// Drain and handle remaining messages
	sub.Drain()
}

func doRequest(nc *nats.Conn, natsSubject string) {
	request := "doUpdate"
	logger.Printf("Sending request: %s", request)
	msg, err := nc.Request(natsSubject, []byte(request), time.Second)
	if err != nil {
		log.Fatal(err)
	}
	logger.Printf("Reply: %s", msg.Data)
}

func subHandler(msg *nats.Msg) {
	logger.Printf("Message recieved: %s", msg.Data)
	logger.Printf("Replying to: %s", msg.Reply)
	msg.Respond([]byte("updateDone"))
}
