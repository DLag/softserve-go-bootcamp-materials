package main

import (
	"log"
	_ "net/http/pprof"
)

func main() {
	log.Fatal(newIdGenServerAppMysql().Run())
}
