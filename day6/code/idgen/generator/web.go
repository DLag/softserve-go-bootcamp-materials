package generator

import (
	"fmt"
	"net/http"
)

type idGenHandler struct {
	idgen idGenerator
}

func (h *idGenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.idgen.Generate())
}

func NewIdGenHandler(idgen idGenerator) *idGenHandler {
	return &idGenHandler{idgen: idgen}
}
