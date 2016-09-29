package main

import (
	"time"
)

func main() {
	generator := newIdGeneratorHTTP(time.Second * 2)
	generator.AddPeer("http://163.172.185.194:8080")
	generator.AddPeer("http://163.172.184.133:8080")
	for {
		generator.Generate()
		time.Sleep(time.Second)
	}
}
