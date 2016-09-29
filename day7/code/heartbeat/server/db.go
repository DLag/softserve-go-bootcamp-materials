package main

import (
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"log"
	_ "net/http/pprof"
	"time"
)

type idGenModel interface {
	Generate() uint32
	Current() uint32
	Alive() bool
}

type idGenModelMysql struct {
	dsn string
	db  *sql.DB
}

func (v *idGenModelMysql) checkConnection() bool {
	var err error
	if v.db == nil {
		log.Print("DB: Not connected, connecting: ", v.dsn)
		v.db, err = sql.Open("mymysql", v.dsn)
		if err != nil {
			log.Fatal("DB: Error connecting to DB: ", err)
		}
	}
	if err = v.db.Ping(); err != nil {
		log.Print("DB: Ping error: ", err)
		return false
	}
	return true
}

func (v *idGenModelMysql) waitConnection() bool {
	reconnectTimeout := time.After(5 * time.Second)
	for {
		select {
		case <-reconnectTimeout:
			log.Print("DB reconnect timed out")
			return false
		default:
			if v.checkConnection() {
				return true
			} else {
				log.Print("No DB connection, retrying...")
				time.Sleep(time.Second)
			}
		}
	}
}

func (v *idGenModelMysql) Alive() bool {
	return v.checkConnection()
}

func (v *idGenModelMysql) Generate() uint32 {
	if v.waitConnection() {
		res, err := v.db.Exec("REPLACE INTO idgen32 (stub) VALUES ('a');")
		if err != nil {
			log.Println("Error on REPLACE: ", err)
			return 0
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Println("Error on LastInsertId: ", err)
			return 0
		}
		return uint32(id)
	}
	return 0
}

func (v *idGenModelMysql) Current() uint32 {
	if v.waitConnection() {
		var id int
		row := v.db.QueryRow("SELECT id FROM idgen32 WHERE stub=? LIMIT 1", "a")
		err := row.Scan(&id)
		if err != nil {
			log.Println("Error on SELECT: ", err)
			return 0
		}
		return uint32(id)
	}
	return 0
}

func newIdGenModelMysql(dsn string) *idGenModelMysql {
	return &idGenModelMysql{dsn: dsn}
}
