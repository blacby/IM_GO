package main

import "net"

type User struct {
	name string
	addr string
	uc   chan string
	con  net.Conn
}

func NewUser(uc chan string, con net.Conn) *User {
	user := &User{name: con.RemoteAddr().String(), addr: con.RemoteAddr().String(), uc: uc, con: con}
	return user
}

// 发送上线消息
func (user *User) SendOnlineMessage() {
	user.uc <- "user online \n"
}

// 接收管道消息发送给 服务端channel
func (user *User) ListenMessage(s *Server) {
	for {
		str := <-user.uc
		s.Message <- str
	}
}
