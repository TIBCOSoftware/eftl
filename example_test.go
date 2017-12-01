// Copyright © 2017. TIBCO Software Inc.
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package eftl_test

import (
	"log"
	"time"

	"tibco.com/eftl"
)

// Connect to the server.
func ExampleConnect() {
	errChan := make(chan error)

	opts := &eftl.Options{
		Username: "username",
		Password: "password",
	}

	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", opts, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// listen for asynnchronous errors
	go func() {
		for err := range errChan {
			log.Printf("connection error: %s", err)
		}
	}()
}

// Reconnect to the server.
func ExampleConnection_Reconnect() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// disconnect from the server
	conn.Disconnect()

	// reconnect to the server
	err = conn.Reconnect()
	if err != nil {
		log.Printf("reconnect failed: %s", err)
		return
	}
}

// Publish messages.
func ExampleConnection_Publish() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// publish a message
	err = conn.Publish(eftl.Message{
		"_dest":        "sample",
		"field-int":    99,
		"field-float":  0.99,
		"field-string": "hellow, world!",
		"field-time":   time.Now(),
		"field-message": eftl.Message{
			"field-bytes": []byte("this is an embedded message"),
		},
	})
	if err != nil {
		log.Printf("publish failed: %s", err)
		return
	}
}

// Publish messages asynchronously.
func ExampleConnection_PublishAsync() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		log.printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	compChan := make(chan *eftl.Completion, 1)

	// publish a message
	err = conn.PublishAsync(eftl.Message{
		"_dest":        "sample",
		"field-int":    99,
		"field-float":  0.99,
		"field-string": "hellow, world!",
		"field-time":   time.Now(),
		"field-message": eftl.Message{
			"field-bytes": []byte("this is an embedded message"),
		},
	}, compChan)
	if err != nil {
		log.printf("publish failed: %s", err)
		return
	}

	// wait for publish operation to complete
	comp := <-compChan

	if comp.Error != nil {
		log.Printf("publish completion failed: %s", err)
	} else {
		log.Printf("published message: %s", comp.Message)
	}
}

// Subscribe to messages.
func ExampleConnection_Subscribe() {
	errChan := make(chan error)

	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	msgChan := make(chan eftl.Message)

	// subscribe to messages
	sub, err := conn.Subscribe("{\"_dest\": \"sample\"}", "", msgChan)
	if err != nil {
		log.Printf("subscribe failed: %s", err)
		return
	}

	// receive messages
	for {
		select {
		case msg := <-msgChan:
			log.Println(msg)
		case err := <-errChan:
			log.Printf("connection error: %s", err)
			return
		}
	}
}

// Subscribe to messages asynchronously.
func ExampleConnection_SubscribeAsync() {
	errChan := make(chan error)

	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, errChan)
	if err != nil {
		log.Printf("connect failed: %s", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	subChan := make(chan *eftl.Subscription)
	msgChan := make(chan eftl.Message)

	// subscribe to messages
	err := conn.SubscribeAsync("{\"_dest\": \"sample\"}", "", msgChan, subChan)
	if err != nil {
		log.Printf("subscribe failed: %s", err)
		return
	}

	// wait for subsribe operation to complete and receive messages
	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				log.Printf("subscription failed: %s", sub.Error)
				return
			}
			log.Println("subscription succeeded")
		case msg := <-msgChan:
			log.Println(msg)
		case err := <-errChan:
			log.Printf("connection error: %s", err)
			return
		}
	}
}
