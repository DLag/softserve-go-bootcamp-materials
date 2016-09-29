package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	generator := newIdGeneratorHTTP(time.Second * 2)
	//generator.AddPeer("http://163.172.185.194:8080")
	//generator.AddPeer("http://163.172.184.133:8080")
	http.Handle("/generate", newGenerateHandler(generator.Generate))
	http.Handle("/current", newGenerateHandler(generator.Current))
	http.Handle("/health", newHealthHandler(generator.Alive))
	http.Handle("/add", newAddPeerHandler(generator))
	log.Fatal(http.ListenAndServe(":8090", nil))
}
