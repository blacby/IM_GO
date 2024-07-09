package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	IP   string
	Port int
}

func NewSever(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

func (s *Server) HandleConnection(conn net.Conn) {
	fmt.Println("...doing handleConnection")
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.IP+":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)

	}
}
