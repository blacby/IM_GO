package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	IP      string
	Port    int
	UserMap map[string]*User
	Message chan string
	RWMux   sync.RWMutex
}

func NewSever(ip string, port int, message chan string, userMap map[string]*User) *Server {
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
		s,
	)
	user.UserOnline()
	//写这个发不出第二次消息了，没有初始化导致发不出来写不进去了
	//var isAlive chan bool
	isAlive := make(chan bool)
	//question1 ：var isAlive chan bool  当为初始化为nil时，
	//发送一次消息会导致 go协程卡死，
	//无论发什么都不能延长时间，user map 未能offline，仍旧存在，
	//当这个user 超时后，有客户端重新链接会导致 广播时用到了关闭了的这个用户的channel
	//question2：为什么select case 读取nil channel时不会直接卡死呢？
	//anwer2:在你的 select 语句中，有一个 <-isAlive 的 case，但由于 isAlive 是 nil，这个 case 实际上是一个无效的 case。在 Go 语言中，当 select 语句中的所有 case 都无法执行时，会阻塞在 select 语句处，直到有一个 case 可以执行为止。因此，尽管 <-isAlive 是一个无效的操作（因为 isAlive 是 nil），但它并不会导致整个 handleConnection 协程卡住，因为还有其他的 case（比如 time.After）可以执行。
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := user.con.Read(buf)
			if n == 0 {
				user.UserOffline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read client message err:" + err.Error())
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isAlive <- true
		}
	}()

	//test:
	//<-isAlive 为什么在这里写这个 handle协程永远卡住，但是瞎买可能的 select不会导致整个handle协程卡住？

	for {
		select {
		case <-isAlive:

		case <-time.After(600 * time.Second):
			user.SendMsg("your client overtime,unactivity so long")
			//question2： user.UserOffline()
			close(user.uc)
			conn.Close()
			return
		}
	}
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
