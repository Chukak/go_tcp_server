package client

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type management interface {
	Connect() bool
	DefaultSend(string) bool
	Disconnect() bool
	IsOpen() (bool, error)
}

type Client struct {
	Port  int
	Host  string
	Conn  net.Conn
	Error error
}

const NetworkType string = "tcp"

func (c *Client) Connect() bool {
	c.Conn, c.Error = net.Dial(NetworkType, strings.Join([]string{c.Host, ":", strconv.Itoa(c.Port)}, ""))
	return c.Error == nil
}

func (c *Client) DefaultSend(message string) bool {
	err := gob.NewEncoder(c.Conn).Encode(message)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	return true
}

func (c *Client) Disconnect() bool {
	c.Error = c.Conn.Close()
	return c.Error == nil
}

func (c *Client) IsOpen() (bool, error) {

	conn, err := net.Dial(NetworkType, strings.Join([]string{c.Host, ":", strconv.Itoa(c.Port)}, ""))
	defer func() {
		if err == nil {
			conn.Close()
		}
	}()
	return err == nil, err
}
