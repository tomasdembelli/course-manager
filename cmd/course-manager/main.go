package main

import (
	"github.com/tomasdembelli/course-manager/server"
	"log"
)

func main() {
	server.StartServer(&server.Config{
		Port:   8000,
		Logger: log.Default(),
	})
}
