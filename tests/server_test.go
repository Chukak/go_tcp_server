package tests

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"

	"net"
	"testing"
	"time"

	clnt "../internal/client"
	srv "../internal/server"
)

// Testing messages
const MsgTest1 string = "Message TEXT"
const MsgTest2 string = "{ \"text\": 'Message JSON' }"
const MsgTest3 string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<text>Message XML</text>"

/**
 * Client test functions
 */
func SendTest1(c *clnt.Client) error {
	message := "Message TEXT"
	return gob.NewEncoder(c.Conn).Encode(message)
}

func SendTest2(c *clnt.Client) error {
	message := "{ \"text\": 'Message JSON' }"
	return json.NewEncoder(c.Conn).Encode(message)
}

func SendTest3(c *clnt.Client) error {
	message := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<text>Message XML</text>"
	return xml.NewEncoder(c.Conn).Encode(message)
}

type Result struct {
	string
	error
}

// Global results
var CurrentResult Result

/**
 * Server test methods
 */
func ProcessingMsgTest1(c *srv.Connection) {
	var msg string
	err := gob.NewDecoder(&io.LimitedReader{R: c.Conn, N: c.SizeBuffer}).Decode(&msg)
	CurrentResult = Result{msg, err}
}

func ProcessingMsgTest2(c *srv.Connection) {
	var msg string
	err := json.NewDecoder(&io.LimitedReader{R: c.Conn, N: c.SizeBuffer}).Decode(&msg)
	CurrentResult = Result{msg, err}
}

func ProcessingMsgTest3(c *srv.Connection) {
	var msg string
	err := xml.NewDecoder(&io.LimitedReader{R: c.Conn, N: c.SizeBuffer}).Decode(&msg)
	CurrentResult = Result{msg, err}
}

// Constants
const Address string = "127.0.0.1:9999"

// Tests
func TestRunServer(t *testing.T) {
	s := srv.Server{Host: "127.0.0.1", Port: 9999, MaxSizeBuffer: 1800}
	isRunning := make(chan bool)

	go s.Start(isRunning)
	time.Sleep(10 * time.Millisecond)
	select {
	case val := <-isRunning:
		if !val {
			t.Errorf("isRunning = %t; want: true", val)
		}
	default:
		t.Errorf("isRunning = None; want: any value")
	}
	conn, err := net.Dial("tcp", Address)
	if err != nil {
		t.Errorf("err = %v; want nil", err)
	}
	conn.Close()
	s.Stop()
	time.Sleep(10 * time.Millisecond)
}

func TestServer1(t *testing.T) {
	s := srv.Server{Host: "127.0.0.1", Port: 9999, MaxSizeBuffer: 2800, ProcessingFunc: ProcessingMsgTest1}
	isRunning := make(chan bool, 1)

	go s.Start(isRunning)
	time.Sleep(10 * time.Millisecond)
	select {
	case val := <-isRunning:
		if !val {
			t.Errorf("isRunning = %t; want: true; err %v", val, s.Error)
		}
	default:
		t.Errorf("isRunning = None; want: any value")
	}

	c := clnt.Client{Host: "127.0.0.1", Port: 9999}
	if c.Connect() {
		SendTest1(&c)
		c.Disconnect()
	} else {
		t.Errorf("c.Connect() = false; want true")
	}

	time.Sleep(10 * time.Millisecond)
	if CurrentResult.string != MsgTest1 {
		t.Errorf("pair.string = %s; want: %s", CurrentResult.string, MsgTest1)
	}
	if CurrentResult.error != nil {
		t.Errorf("pair.error = %v; want: nil", CurrentResult.error)
	}

	s.Stop()
	time.Sleep(10 * time.Millisecond)
}

func TestServer2(t *testing.T) {
	s := srv.Server{Host: "127.0.0.1", Port: 9999, ProcessingFunc: ProcessingMsgTest2}
	isRunning := make(chan bool)

	go s.Start(isRunning)
	time.Sleep(10 * time.Millisecond)
	select {
	case val := <-isRunning:
		if !val {
			t.Errorf("isRunning = %t; want: true", val)
		}
	default:
		t.Errorf("isRunning = None; want: any value")
	}

	c := clnt.Client{Host: "127.0.0.1", Port: 9999}
	if c.Connect() {
		SendTest2(&c)
		c.Disconnect()
	} else {
		t.Errorf("c.Connect() = false; want true")
	}

	time.Sleep(10 * time.Millisecond)
	if CurrentResult.string != MsgTest2 {
		t.Errorf("pair.string = %s; want: %s", CurrentResult.string, MsgTest2)
	}
	if CurrentResult.error != nil {
		t.Errorf("pair.error = %v; want: nil", CurrentResult.error)
	}

	s.Stop()
	time.Sleep(10 * time.Millisecond)
}

func TestServer3(t *testing.T) {
	s := srv.Server{Host: "127.0.0.1", Port: 9999, ProcessingFunc: ProcessingMsgTest3}
	isRunning := make(chan bool)

	go s.Start(isRunning)
	time.Sleep(10 * time.Millisecond)
	select {
	case val := <-isRunning:
		if !val {
			t.Errorf("isRunning = %t; want: true", val)
		}
	default:
		t.Errorf("isRunning = None; want: any value")
	}

	c := clnt.Client{Host: "127.0.0.1", Port: 9999}
	if c.Connect() {
		SendTest3(&c)
		c.Disconnect()
	} else {
		t.Errorf("c.Connect() = false; want true")
	}

	time.Sleep(10 * time.Millisecond)
	if CurrentResult.string != MsgTest3 {
		t.Errorf("pair.string = %s; want: %s", CurrentResult.string, MsgTest3)
	}
	if CurrentResult.error != nil {
		t.Errorf("pair.error = %v; want: nil", CurrentResult.error)
	}

	s.Stop()
	time.Sleep(10 * time.Millisecond)
}
