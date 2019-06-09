package main

import (
	"fmt"
	"strings"
	"time"

	cln "./internal/client"
	serv "./internal/server"
)

const DatetimeFormat string = "2006.01.02 15:04:05"

func clients() {
	for {
		c := cln.Client{Host: "127.0.0.1", Port: 9999}
		if c.Connect() {
			c.Receive(strings.Join([]string{"{", "\"time\": ", "'", time.Now().Format(DatetimeFormat), "' }"}, ""))
			c.Disconnect()
		} else {
			fmt.Println(c.Error)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	s := serv.Server{Host: "127.0.0.1", Port: 9999}
	isOk := make(chan bool)
	isStopped := make(chan bool, 2)
	isStopped <- false

	go s.Start(isOk, isStopped)
	time.Sleep(80 * time.Millisecond)
	select {
	case val := <-isOk:
		if !val {
			fmt.Println(s.Error)
			return
		} else {
			fmt.Println("Server was running at '127.0.0.1:9999'")
		}
	default:
		fmt.Println("No value in channel!")
	}

	go clients()
	time.Sleep(10 * time.Second)

	isStopped <- true
	time.Sleep(80 * time.Millisecond)
	select {
	case val := <-isOk:
		if !val {
			fmt.Println(s.Error)
		}
	default:
	}
}
