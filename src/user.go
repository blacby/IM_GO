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
	go user.ListenMessage()
	return user
}

// 接收管道消息发送给客户端界面
func (user *User) ListenMessage() {
	for {
		msg := <-user.uc
		user.con.Write([]byte(msg + "\n"))
	}
}
