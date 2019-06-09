package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type management interface {
	Start(chan<- bool, <-chan bool)
	stop()
	handleConnection()
}

type Server struct {
	Port     int
	Host     string
	listener net.Listener
	Error    error
}

type Connection struct {
	conn net.Conn
}

func (s *Server) Start(result chan<- bool, isStopped <-chan bool) {
	s.listener, s.Error = net.Listen("tcp", strings.Join([]string{":", strconv.Itoa(s.Port)}, ""))
	if s.Error == nil {
		result <- true
		for {
			select {
			case cancel := <-isStopped:
				if cancel {
					s.stop()
					break
				}
			default:
				c, err := s.listener.Accept()
				if err != nil {
					continue
				}
				Connection.handleConnection(Connection{conn: c})
			}
		}
	} else {
		result <- false
		s.stop()
	}
}

func (s *Server) stop() {
	s.Error = s.listener.Close()
}

func (c Connection) handleConnection() {
	var msg string
	err := json.NewDecoder(c.conn).Decode(&msg)
	if err != nil {
		fmt.Println("ERROR: ", err)
	} else {
		fmt.Println("Message: ", msg)
	}
	c.conn.Close()
}
