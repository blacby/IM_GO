package main

import (
	"fmt"
	"net"
)

type Client struct {
	Name string
	IP   string
	Port string
	Conn net.Conn
}

func NewClient(ip string, port string) (*Client, error) {
	conn, err := net.Dial("tcp", ip+port)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &Client{IP: ip, Port: port, Conn: conn}, nil
}

func main() {
	client, err := NewClient("localhost", ":8888")
	if err != nil || client == nil {
		return
	} else {
		fmt.Println(">>>>>>>>>>>>>>>>>connect to server successfully")
	}

	select {}
}
