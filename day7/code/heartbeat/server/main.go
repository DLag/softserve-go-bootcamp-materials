package main

import "log"

func main() {
	log.Fatal(newIdGenServerAppMysql().Run())
}
