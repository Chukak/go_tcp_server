package server

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type management interface {
	Start(chan<- bool)
	Stop()
}

type Server struct {
	Port           int
	Host           string
	listener       net.Listener
	exited         bool
	Error          error
	MaxSizeBuffer  int64
	ProcessingFunc func(*Connection)
}

type Connection struct {
	Conn       net.Conn
	SizeBuffer int64
	AtHost     string
	AtPort     int
}

const defaultBuffer int64 = 1800
const NetworkType string = "tcp"

func (s *Server) Start(result chan<- bool) {
	// size buffer
	if s.MaxSizeBuffer == 0 {
		fmt.Println(strings.Join([]string{"MaxSizeBuffer = 0. Using ", strconv.FormatInt(defaultBuffer, 10), " by default."}, ""))
		s.MaxSizeBuffer = defaultBuffer
	}
	// processing function
	if s.ProcessingFunc == nil {
		fmt.Println("ProcessingFunc is null. Using encoding/god by default.")
		s.ProcessingFunc = defaultProcessingFunc
	}
	fmt.Println(strings.Join([]string{"Server was running at '", s.Host, ":", strconv.Itoa(s.Port), "'."}, ""))
	s.listener, s.Error = net.Listen(NetworkType, strings.Join([]string{":", strconv.Itoa(s.Port)}, ""))
	if s.Error == nil {
		result <- true
		s.exited = false
		for {
			c, err := s.listener.Accept()
			if err != nil {
				if s.exited {
					break
				}
				fmt.Println("Accept connection failed! What: ", err)
				continue
			}
			con := Connection{Conn: c, SizeBuffer: s.MaxSizeBuffer, AtHost: s.Host, AtPort: s.Port}
			s.ProcessingFunc(&con)
			con.Conn.Close()
		}
	} else {
		result <- false
		s.exited = true
	}
}

func (s *Server) Stop() {
	fmt.Print("Stopping server...")
	s.exited = true
	s.Error = s.listener.Close()
	fmt.Print("Done!\n")
}

func defaultProcessingFunc(c *Connection) {
	var msg string
	err := gob.NewDecoder(&io.LimitedReader{R: c.Conn, N: c.SizeBuffer}).Decode(&msg)
	fmt.Println("(Server) -> Message: '", msg, "' Error: ", err)
}
