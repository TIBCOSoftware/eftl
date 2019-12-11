// Copyright © 2017-2019. TIBCO Software Inc.
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package eftl_test

import (
	"fmt"

	"github.com/TIBCOSoftware/eftl"
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
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// listen for asynnchronous connection errors
	go func() {
		for err := range errChan {
			fmt.Println("connection error:", err)
		}
	}()
}

// Reconnect to the server.
func ExampleConnection_Reconnect() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// disconnect from the server
	conn.Disconnect()

	// reconnect to the server
	err = conn.Reconnect()
	if err != nil {
		fmt.Println("reconnect failed:", err)
	}
}

// Publish messages.
func ExampleConnection_Publish() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	// publish a message
	err = conn.Publish(eftl.Message{
		"_dest": "sample",
		"text":  "Hello, World!",
	})
	if err != nil {
		fmt.Println("publish failed:", err)
	}
}

// Publish messages asynchronously.
func ExampleConnection_PublishAsync() {
	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	compChan := make(chan *eftl.Completion, 1)

	// publish a message
	err = conn.PublishAsync(eftl.Message{
		"_dest": "sample",
		"text":  "Hello, World!",
	}, compChan)
	if err != nil {
		fmt.Println("publish failed:", err)
		return
	}

	// wait for publish operation to complete
	comp := <-compChan

	if comp.Error != nil {
		fmt.Println("publish completion failed:", err)
	} else {
		fmt.Println("published message:", comp.Message)
	}
}

// Subscribe to messages.
func ExampleConnection_Subscribe() {
	errChan := make(chan error)

	// connect to the server
	conn, err := eftl.Connect("ws://localhost:9191/channel", nil, errChan)
	if err != nil {
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	msgChan := make(chan eftl.Message)

	// subscribe to messages
	sub, err := conn.Subscribe("{\"_dest\": \"sample\"}", "", msgChan)
	if err != nil {
		fmt.Println("subscribe failed:", err)
		return
	}

	// unsubscribe when done
	conn.Unsubscribe(sub)

	// receive messages
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg)
		case err := <-errChan:
			fmt.Println("connection error:", err)
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
		fmt.Println("connect failed:", err)
		return
	}

	// disconnect from the server when done with the connection
	defer conn.Disconnect()

	subChan := make(chan *eftl.Subscription)
	msgChan := make(chan eftl.Message)

	// subscribe to messages
	err = conn.SubscribeAsync("{\"_dest\": \"sample\"}", "", msgChan, subChan)
	if err != nil {
		fmt.Println("subscribe failed:", err)
		return
	}

	// wait for subsribe operation to complete and receive messages
	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				fmt.Println("subscription failed:", sub.Error)
				return
			}
			// unsubscribe when done
			defer conn.Unsubscribe(sub)
		case msg := <-msgChan:
			fmt.Println(msg)
		case err := <-errChan:
			fmt.Println("connection error:", err)
			return
		}
	}
}
