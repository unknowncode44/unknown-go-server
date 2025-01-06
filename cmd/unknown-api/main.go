package main

import (
	"github.com/unknowncode44/unknown-go-server/config"
	database "github.com/unknowncode44/unknown-go-server/db"
	"github.com/unknowncode44/unknown-go-server/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	server.NewFiberServer(conf, db).Start()
}
