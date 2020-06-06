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

	if err != nil {
		log.Fatalf("error while setting up repository %s", err)
	}

	log.Println("repository initialized")

	log.Println("sync the repo..")
	synced, err := repo.Sync()

	if err != nil {
		log.Fatalf("error while syncing the repo %s", err)
	}

	if synced {
		log.Println("repository is synced")
	} else {
		log.Println("repository was not synced")
	}

	// Testing

	// user := models.User{
	// 	Username: "sepp",
	// 	Password: "",
	// }
	// userCreated, err := repo.CreateUser(&user)

	// group := models.Group{
	// 	Name: "testgroup",
	// 	Type: "ldap",
	// }

	// groupCreated, err := repo.CreateGroup(&group)

	// if err != nil {
	// 	log.Fatalf("error creating group %s", err)
	// }

	// err = repo.AddUserToGroup(userCreated.ID, groupCreated.ID)

	return nil
}
