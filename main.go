package main

import (
	"web-check-go/config"
	"web-check-go/routers"
)

func main() {
	router := routers.Routes()
	config.ServerConfig(router)
}
