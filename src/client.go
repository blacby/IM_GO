package main

import (
	"flag"
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
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &Client{IP: ip, Port: port, Conn: conn}, nil
}

var ServerIp string
var ServerPort string

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "input server ip")
	flag.StringVar(&ServerPort, "port", "8888", "input server port")
}

func main() {
	flag.Parse() //这个经常容易忘记
	client, err := NewClient(ServerIp, ServerPort)
	if err != nil || client == nil {
		return
	} else {
		fmt.Println(">>>>>>>>>>>>>>>>>connect to server successfully")
	}

	select {}
}
