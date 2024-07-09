package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	IP      string
	Port    int
	UserMap map[string]*User
	Message chan string
	RWMux   sync.RWMutex
}

func NewSever(ip string, port int, userMap map[string]*User, message chan string) *Server {
	server := &Server{
		IP:      ip,
		Port:    port,
		UserMap: userMap,
		Message: message,
	}
	return server
}

func (s *Server) HandleConnection(conn net.Conn) {
	fmt.Println("...doing handleConnection")
	//create user
	user := NewUser(
		make(chan string),
		conn,
	)
	go user.ListenMessage(s)
	s.RWMux.Lock()
	s.UserMap[user.name] = user
	s.RWMux.Unlock()
	user.SendOnlineMessage()
}

// BoradCast user online message to every user
func (s *Server) BoradCast() {
	for {
		msg := <-s.Message
		for name, user := range s.UserMap {
			user.con.Write([]byte("name[" + name + "] addr[" + user.addr + "] :" + msg))
		}
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.IP+":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	defer listen.Close()
	//bordcast online message
	go s.BoradCast()
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)

	}
}
