package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type management interface {
	Connect() bool
	Receive(string) bool
	Disconnect() bool
}

type Client struct {
	Port       int
	Host       string
	connection net.Conn
	Error      error
}

const NetworkType string = "tcp"

func (c *Client) Connect() bool {
	c.connection, c.Error = net.Dial(NetworkType, strings.Join([]string{c.Host, ":", strconv.Itoa(c.Port)}, ""))
	return c.Error == nil
}

func (c *Client) Receive(message string) bool {
	err := json.NewEncoder(c.connection).Encode(message)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return false
	}
	return true
}

func (c *Client) Disconnect() bool {
	c.Error = c.connection.Close()
	return c.Error == nil
}
