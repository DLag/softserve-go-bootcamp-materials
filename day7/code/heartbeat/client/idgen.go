package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type IdGenerator interface {
	Generate() int32
	Current() int32
	Alive() bool
}

type peer struct {
	Addr  string
	Alive bool
}

type IdGeneratorHTTP struct {
	client  http.Client
	peers   []peer
	timeout time.Duration
	sync.RWMutex
}

func (idgen *IdGeneratorHTTP) AddPeer(addr string) {
	for _, v := range idgen.peers {
		if v.Addr == addr {
			return
		}
	}
	idgen.peers = append(idgen.peers, peer{Addr: addr})
	log.Printf("Added peer %q", addr)
}

func (idgen *IdGeneratorHTTP) checkAlive() uint32 {
	for {
		log.Println("Checking alive nodes...")
		for k := range idgen.peers {
			ctx, _ := context.WithTimeout(context.Background(), idgen.timeout)
			req, err := http.NewRequest("GET", idgen.peers[k].Addr+"/generate", nil)
			if err != nil {
				log.Printf("Wrong request %q", req.URL)
			}
			req.WithContext(ctx)
			_, err = idgen.client.Do(req)
			if err != nil {
				log.Printf("Error on request: %v, marking peer %q as dead", err, idgen.peers[k].Addr)
				idgen.peers[k].Alive = false
				continue
			}
			idgen.peers[k].Alive = true
		}
		time.Sleep(time.Second * 10)
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
	defer func() {
		buf := make([]byte, 1024)
		r.Body.Read(buf)
		r.Body.Close()
	}()
	buf, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	id, err := strconv.Atoi(string(buf))
	if err != nil {
		log.Printf("Error on converting: %v, marking peer %q as dead", err, peer.Addr)
		peer.Alive = false
		return 0
	}
	return uint32(id)
}

func (idgen *IdGeneratorHTTP) Generate() uint32 {
	for k := range idgen.peers {
		idgen.RLock()
		alive := idgen.peers[k].Alive
		idgen.RUnlock()
		if alive {
			//log.Printf("Trying to generate id on peer %q", idgen.peers[k].Addr)
			id := idgen.tryGenerate(&idgen.peers[k])
			if id > 0 {
				log.Printf("ID: %d, peer %q", id, idgen.peers[k].Addr)
				return id
			}
		}
	}
	log.Println("No nodes online")
	return 0
}

func newIdGeneratorHTTP(timeout time.Duration) *IdGeneratorHTTP {
	g := &IdGeneratorHTTP{
		timeout: timeout,
		client:  http.Client{Timeout: timeout},
		peers:   make([]peer, 0),
	}
	go g.checkAlive() //Starting checker
	return g
}
