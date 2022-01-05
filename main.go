package main

import "fmt"

func main() {
	s := ServerStatusManager{}
	s.Init()
	defer s.Start(":9123")
	fmt.Println("Server started")
}
