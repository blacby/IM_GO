package main

import (
	"fmt"
	"io"
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
	s.RWMux.Lock()
	s.UserMap[user.name] = user
	s.RWMux.Unlock()
	//广播上线消息
	s.BoradCast(user, "user online")
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := user.con.Read(buf)
			if n == 0 {
				s.BoradCast(user, "user offline")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read client message err:" + err.Error())
				return
			}
			s.BoradCast(user, string(buf[:n-1]))
		}

	}()

}

// BoradCast user online message to every user 触发用户上线消息
func (s *Server) BoradCast(user *User, msg string) {
	s.Message <- "name[" + user.name + "] addr[" + user.addr + "] : " + msg
}

// 所有listen方法listen的都是channel，广播的定义是由服务端广播到客户端
func (s *Server) ListenMessage() string {
	for {
		msg := <-s.Message
		s.RWMux.RLock()
		for _, user := range s.UserMap {
			user.uc <- msg
		}
		s.RWMux.RUnlock()
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.IP+":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	defer listen.Close()
	//bordcast online message
	go s.ListenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)

	}
}
