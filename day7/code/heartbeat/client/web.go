package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type GenerateHandler struct {
	generate func() uint32
}

func (h *GenerateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, h.generate())
}

func newGenerateHandler(generate func() uint32) *GenerateHandler {
	return &GenerateHandler{generate: generate}
}

type HealthHandler struct {
	health aliveFunc
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	online, all := h.health()
	if online != 0 {
		w.WriteHeader(http.StatusOK) // 200
		fmt.Fprintf(w, "%d/%d", online, all)
	} else {
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write([]byte("dead"))
	}
}

func newHealthHandler(health aliveFunc) *HealthHandler {
	return &HealthHandler{health: health}
}

type AddPeerHandler struct {
	idgen *IdGeneratorHTTP
}

func getFromMap(key string, m map[string]interface{}) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getPeerFromMap(m map[string]interface{}) peer {
	b, _ := strconv.ParseBool(getFromMap("alive", m))
	return peer{Addr: getFromMap("addr", m), Alive: b}
}

func (h *AddPeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var j map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&j) // map[string]interface{} {
	// "addr": "adfasdf",
	// "alive: false,
	// }

	if h.idgen.AddPeer(getPeerFromMap(j)) {
		w.WriteHeader(http.StatusOK) // 200
		fmt.Fprint(w, "Added")
	} else {
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write([]byte("Can't add peer"))
	}
}

func newAddPeerHandler(idgen *IdGeneratorHTTP) *AddPeerHandler {
	return &AddPeerHandler{idgen: idgen}
}
