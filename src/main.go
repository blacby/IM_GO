package main

func main() {

	s := NewSever("127.0.0.1", 8888, make(chan string), make(map[string]*User))
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
