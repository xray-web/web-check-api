package main

import (
	"log"

	"github.com/xray-web/web-check-api/config"
	"github.com/xray-web/web-check-api/server"
)

func main() {
	s := server.New(config.New())
	log.Println(s.Run())
}
