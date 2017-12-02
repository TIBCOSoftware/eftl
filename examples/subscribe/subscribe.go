// Copyright © 2017. TIBCO Software Inc.
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package main

import (
	"flag"
	"log"

	"github.com/TIBCOSoftware/eftl"
)

func main() {

	// Durable subscriptions require a unique client identifier.
	clientIDPtr := flag.String("clientid", "client-go", "unique client identifier")

	flag.Parse()

	// eFTL library version.
	log.Printf("eFTL version %s\n", eftl.Version)

	// Channel on which to receive connection errors.
	errChan := make(chan error, 1)

	// Set connection options.
	opts := &eftl.Options{
		Username: "user",
		Password: "pass",
		ClientID: *clientIDPtr,
	}

	// Connect. 
	conn, err := eftl.Connect("wss://localhost:9191/channel", opts, errChan)
	if err != nil {
		log.Println("connect failed:", err)
		return
	}

	// Close the connection when done.
	defer conn.Disconnect()

	log.Println("connected")

	// Channel on which to listen for published messages.
	msgChan := make(chan eftl.Message, 100)

	// Subscribe to messages.
	//
	// A subscription receives only those messages
	// whose fields match the subscription's matcher
	// string. A subscription can only match on
	// string and long fields.
	//
	// To match all messages use the empty matcher "{}".
	//
	// This subscription matches messages containing
	// a string field with name "type" and value "example".
	//
	// The durable name "example" is being used for this
	// subscription.
	//
	sub, err := conn.Subscribe("{\"type\":\"example\"}", "example", msgChan)
	if err != nil {
		log.Println("subscription failed:", err)
	}

	// Unsubscribe when done.
	defer conn.Unsubscribe(sub)

	// Listen for messages and connection errors.
	for {
		select {
		case msg := <-msgChan:
			log.Println("received message:", msg)
		case err := <-errChan:
			log.Println("connection error:", err)
		}
	}
}
