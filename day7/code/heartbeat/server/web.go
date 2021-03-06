package main

import (
	"fmt"
	"net/http"
)

type heartbeatHandler struct {
	model idGenModel
}

func (h *heartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.model.Alive() {
		w.WriteHeader(http.StatusOK) // 200
		w.Write([]byte("alive"))
	} else {
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write([]byte("dead"))
	}
}

func newHeartbeatHandler(model idGenModel) *heartbeatHandler {
	return &heartbeatHandler{model: model}
}

type idGenHandler struct {
	f func() uint32
}

func (h *idGenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := h.f()
	if id > 0 {
		w.WriteHeader(http.StatusOK) // 200
	} else {
		w.WriteHeader(http.StatusInternalServerError) //500
	}
	fmt.Fprintf(w, "%d", id)
}

func newIdGenHandler(f func() uint32) *idGenHandler {
	return &idGenHandler{f: f}
}
