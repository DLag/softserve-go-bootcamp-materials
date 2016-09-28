package main

import (
	"net/http"
	"time"
	"context"
	"io/ioutil"
	"io"
	"strconv"
	"log"
)

type IdGenerator interface {
	Generate() int32
	Current() int32
	Alive() bool
}

type peer struct{
	Addr string
	Alive bool
}

type IdGeneratorHTTP struct {
	client http.Client
	peers []peer
	timeout time.Duration
}

func (idgen *IdGeneratorHTTP) AddPeer(addr string){
	for _, v := range idgen.peers {
		if v.Addr == addr {
			return
		}
	}
	idgen.peers = &peer{Addr: addr}
}

func (idgen *IdGeneratorHTTP) checkAlive() uint32 {
	for {
		for k, v:= range idgen.peers {
		ctx, _ := context.WithTimeout(context.Background(), idgen.timeout)
		req, err := http.NewRequest("GET", peer.Addr+"/generate", nil)
		if err != nil {
			log.Fatal("Wrong request")
		}
		req.WithContext(ctx)
		r, err := idgen.client.Do(req)
		if err != nil {
			log.Printf("Error on request: %v, marking peer %q as dead", err, peer.Addr)
			idgen[].Alive = false
			continue
		}

		time.Sleep(time.Second*10)
	}
}

func (idgen *IdGeneratorHTTP) tryGenerate(peer *peer) uint32 {
	ctx, _ := context.WithTimeout(context.Background(), idgen.timeout)
	req, err := http.NewRequest("GET", peer.Addr+"/generate", nil)
	if err != nil {
		log.Fatal("Wrong request")
	}
	req.WithContext(ctx)
	r, err := idgen.client.Do(req)
	if err != nil {
		log.Printf("Error on request: %v, marking peer %q as dead", err, peer.Addr)
		peer.Alive = false
		return 0
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	id, err := strconv.Atoi(string(buf))
	if err !=nil {
		log.Printf("Error on converting: %v, marking peer %q as dead", err, peer.Addr)
		peer.Alive = false
		return 0
	}
	return int32(id)
}

func (idgen *IdGeneratorHTTP) Generate() uint32 {
	for k, v:= range idgen.peers {
		if peer.Alive {
			log.Printf("Trying to generate id on peer %q", v.Alive)
			id:=idgen.tryGenerate(&idgen.peers[k])
			if id>0 {
				return id
			}
		}
	}
	return 0
}