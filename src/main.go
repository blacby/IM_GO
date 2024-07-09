package main

func main() {

	s := NewSever("127.0.0.1", 8888, make(map[string]*User), make(chan string))
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
