// Copyright © 2017-2018. TIBCO Software Inc.
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package main

import (
	"log"
	"time"

	"github.com/TIBCOSoftware/eftl"
)

func main() {

	// eFTL library version.
	log.Printf("eFTL version %s\n", eftl.Version)

	// Channel on which to receive connection errors.
	errChan := make(chan error, 1)

	// Set connection options.
	opts := &eftl.Options{
		Username: "user",
		Password: "pass",
	}

	// Connect.
	conn, err := eftl.Connect("ws://localhost:9191/channel", opts, errChan)
	if err != nil {
		log.Println("connect failed:", err)
		return
	}

	// Close the connection when done.
	defer conn.Disconnect()

	log.Println("connected")

	// Listen for connection errors.
	go func() {
		for err := range errChan {
			log.Println("connection error:", err)
		}
	}()

	// Create a message and populate its contents.
	//
	// Message fields can be of type string, long,
	// double, date, and message, along with arrays
	// of each of the types.
	//
	// Subscribers will create matchers that match
	// on one or more of the message fields. Only
	// string and long fields are supported by
	// matchers.
	//
	msg := eftl.Message{
		"type": "example",
		"text": "This is an example message",
		"long": int64(101),
		"time": time.Now(),
	}

	// Publish the message.
	err = conn.Publish(msg)
	if err != nil {
		log.Println("publish failed:", err)
		return
	}

	log.Println("published message:", msg)
}
