package main

import (
	"encoding/json"
	"fmt"
	"io"

	"strings"
	"time"

	clnt "../internal/client"
	srv "../internal/server"
)

// Datetime format
const DatetimeFormat string = "2006.01.02 15:04:05"

// Send message to server
func Send(c *clnt.Client) bool {
	message := strings.Join([]string{"{", "\"time\": ", "'", time.Now().Format(DatetimeFormat), "' }"}, "")
	err := json.NewEncoder(c.Conn).Encode(message)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return false
	}
	return true
}

// Create clients
func Clients() {
	for {
		c := clnt.Client{Host: "127.0.0.1", Port: 9999}
		open, err := c.IsOpen()
		if !open {
			fmt.Println("Host does not open! Error: ", err)
			time.Sleep(1 * time.Second)
			continue
		}
		// Connect
		if c.Connect() {
			// Send message
			Send(&c)
			// Disconnect
			c.Disconnect()
		} else {
			fmt.Println("Client connection error: ", c.Error)
		}
		time.Sleep(1 * time.Second)
	}
}

// Alias for type Server
type Serv srv.Server

// Processing message from clients
func CustomProcessingFunc(c *srv.Connection) {
	var msg string
	err := json.NewDecoder(&io.LimitedReader{R: c.Conn, N: c.SizeBuffer}).Decode(&msg)
	if err != nil {
		if err.Error() != "EOF" {
			fmt.Println("Error: ", err)
		}
	} else {
		fmt.Println("Msg: ", msg)
	}
}

func main() {
	s := srv.Server{Host: "127.0.0.1", Port: 9999, ProcessingFunc: CustomProcessingFunc}
	isRunning := make(chan bool)

	go s.Start(isRunning)
	time.Sleep(10 * time.Millisecond)
	// check if server was running
	select {
	case val := <-isRunning:
		if !val {
			fmt.Println(s.Error)
			return
		}
	default:
		fmt.Println("No value inside channel!")
	}

	go Clients()
	time.Sleep(10 * time.Second)
	s.Stop()
	time.Sleep(50 * time.Millisecond)
}
