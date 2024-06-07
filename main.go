package main

import (
	"github.com/xray-web/web-check-api/config"
	"github.com/xray-web/web-check-api/routers"
)

func main() {
	router := routers.Routes()
	config.ServerConfig(router)
}
