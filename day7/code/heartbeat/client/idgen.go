package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type IdGenerator interface {
	Generate() uint32
	Current() uint32
	Alive() (uint32, uint32)
}

type aliveFunc func() (uint32, uint32)

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

func (idgen *IdGeneratorHTTP) AddPeer(p peer) bool {
	if len(p.Addr) == 0 {
		log.Print("Peer: empty")
		return false
	}
	if _, err := url.Parse(p.Addr); err != nil {
		log.Print("Peer: Can't parse, error: ", err)
		return false
	}
	for _, v := range idgen.peers {
		if v.Addr == p.Addr {
			log.Printf("Peer: Already exist: %q", p.Addr)
			return false
		}
	}
	idgen.peers = append(idgen.peers, p)
	log.Printf("Added peer %q", p.Addr)
	return true
}

func (idgen *IdGeneratorHTTP) Alive() (uint32, uint32) {
	var online uint32
	for k := range idgen.peers {
		if idgen.peers[k].Alive {
			online++
		}
	}
	return online, uint32(len(idgen.peers))
}

func (idgen *IdGeneratorHTTP) checkAlive() {
	for {
		log.Println("Checking alive nodes...")
		for k := range idgen.peers {
			ctx, _ := context.WithTimeout(context.Background(), idgen.timeout)
			req, err := http.NewRequest("GET", idgen.peers[k].Addr+"/health", nil)
			if err != nil {
				log.Printf("Wrong request %q", req.URL)
			}
			req.WithContext(ctx)
			resp, err := idgen.client.Do(req)
			if err != nil {
				log.Printf("Error on request: %v, marking peer %q as dead", err, idgen.peers[k].Addr)
				idgen.peers[k].Alive = false
				continue
			}
			buf, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Error on reading request: %v, marking peer %q as dead", err, idgen.peers[k].Addr)
				idgen.peers[k].Alive = false
				continue
			}
			if bytes.Compare(buf, []byte("alive")) != 0 {
				log.Printf("Wrong peer response: %v, marking peer %q as dead", string(buf), idgen.peers[k].Addr)
				idgen.peers[k].Alive = false
				continue
			}
			idgen.peers[k].Alive = true
		}
		time.Sleep(time.Second * 10)
	}
}

func (idgen *IdGeneratorHTTP) tryHTTP(peer *peer, uri string) uint32 {
	ctx, _ := context.WithTimeout(context.Background(), idgen.timeout)
	req, err := http.NewRequest("GET", peer.Addr+uri, nil)
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

func (idgen *IdGeneratorHTTP) tryAllNodesHTTP(uri string) uint32 {
	for k := range idgen.peers {
		idgen.RLock()
		alive := idgen.peers[k].Alive
		idgen.RUnlock()
		if alive {
			//log.Printf("Trying to generate id on peer %q", idgen.peers[k].Addr)
			id := idgen.tryHTTP(&idgen.peers[k], uri)
			if id > 0 {
				log.Printf("ID: %d, peer %q", id, idgen.peers[k].Addr)
				return id
			}
		}
	}
	log.Println("No nodes online")
	return 0
}

func (idgen *IdGeneratorHTTP) Generate() uint32 {
	return idgen.tryAllNodesHTTP("/generate")
}

func (idgen *IdGeneratorHTTP) Current() uint32 {
	return idgen.tryAllNodesHTTP("/current")
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
