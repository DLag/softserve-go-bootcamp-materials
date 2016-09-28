package main

import (
	"net/http"
)

type idGenServerApp interface {
	Run() error
}

type idGenServerAppMysql struct {
	cfg                 cfgModel
	idgen               idGenModel
	heartbeat, generate http.Handler
}

func (app *idGenServerAppMysql) Run() error {
	app.idgen.Alive()
	http.Handle("/generate", app.generate)
	http.Handle("/health", app.heartbeat)
	return http.ListenAndServe(":8080", nil)
}

func newIdGenServerAppMysql() *idGenServerAppMysql {
	app := new(idGenServerAppMysql)
	app.cfg = newCfgModelEnv()
	app.idgen = newIdGenModelMysql(app.cfg.Get("DSN"))
	app.heartbeat = newHeartbeatHandler(app.idgen)
	app.generate = newIdGenHandler(app.idgen)
	return app
}