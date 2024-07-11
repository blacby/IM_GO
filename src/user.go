package main

import "net"

type User struct {
	name   string
	addr   string
	uc     chan string
	con    net.Conn
	server *Server
}

func NewUser(uc chan string, con net.Conn, server *Server) *User {
	user := &User{name: con.RemoteAddr().String(), addr: con.RemoteAddr().String(), uc: uc, con: con, server: server}
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

func (user *User) UserOnline() {
	user.server.RWMux.Lock()
	user.server.UserMap[user.name] = user
	user.server.RWMux.Unlock()
	//广播上线消息
	user.server.BoradCast(user, "user online")
}

func (user *User) UserOffline() {
	user.server.RWMux.Lock()
	delete(user.server.UserMap, user.name)
	user.server.RWMux.Unlock()
	//广播上线消息
	user.server.BoradCast(user, "user offline")
}

func (user *User) SendMsg(msg string) {
	user.con.Write([]byte(msg))
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		//查询在线用户
		user.server.RWMux.Lock()
		for _, userEnty := range user.server.UserMap {
			whoOnlineMsg := "[" + userEnty.addr + "]" + userEnty.name + ":在线..\n"
			user.SendMsg(whoOnlineMsg)
		}
		user.server.RWMux.Unlock()

	} else {
		user.server.BoradCast(user, msg)
	}

}
