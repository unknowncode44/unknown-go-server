package main

import (
	"log"
	"os"

	"github.com/unknowncode44/unknown-go-server/config"
	database "github.com/unknowncode44/unknown-go-server/db"
	"github.com/unknowncode44/unknown-go-server/pkg/user"
	"github.com/unknowncode44/unknown-go-server/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)

	// bootstrap del primer admin: solo actúa si SEED_ADMIN_EMAIL y
	// SEED_ADMIN_PASSWORD están seteadas y todavía no existe ningún
	// ADMIN activo (idempotente, se puede dejar siempre)
	if err := user.SeedInitialAdmin(
		user.NewRepo(db.GetDb()),
		os.Getenv("SEED_ADMIN_EMAIL"),
		os.Getenv("SEED_ADMIN_PASSWORD"),
	); err != nil {
		log.Fatalf("Falla creando el admin inicial: %v", err)
	}

	server.NewFiberServer(conf, db).Start()
}
