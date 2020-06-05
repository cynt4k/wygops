package main

import (
	"log"

	evbus "github.com/asaskevich/EventBus"
	"github.com/cynt4k/wygops/internal/mysql"
	"github.com/cynt4k/wygops/internal/repository"
)

func serve() error {

	// New Message bus
	bus := evbus.New()

	log.Println("connecting to database...")
	db, err := mysql.Init()

	if err != nil {
		log.Fatalf("error while connecting to the database %s", err)
	}
	log.Println("connection to database established")

	log.Println("setting up repository...")
	repo, err := repository.NewGormRepository(db, &bus)

	log.Println("repository initialized")

	// Testing

	// user := models.User{
	// 	Username: "sepp",
	// 	Password: "asdf",
	// }
	// repo.CreateUser(&user)

	return nil
}
