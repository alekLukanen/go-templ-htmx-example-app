package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/alekLukanen/go-templ-htmx-example-app/database"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("expected a command")
	}

	m, err := migrate.New(
		"file://database/migrations",
		database.DatabaseConnectionURL("pgx"),
	)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "version":
		log.Println("- getting version")

		if version, dirty, err := m.Version(); err != nil {
			log.Println(err)
		} else {
			log.Println("- version:", version, "dirty:", dirty)
		}

	case "upOne":
		log.Println("- applying next migration")

		if err := m.Steps(1); err != nil {
			log.Fatal(err)
		}

	case "downOne":
		log.Println("- reverting migrations")
		log.Println("Are you sure you want to revert the last migration? (yes/no)")

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			log.Println("aborting revert process")
			return
		}

		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}

	case "force":
		log.Println("- forcing version number")
		log.Println("Enter version number to force to:")

		var version int
		fmt.Scanln(&version)

		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}

		if version, dirty, err := m.Version(); err != nil {
			log.Println(err)
		} else {
			log.Println("- version:", version, "dirty:", dirty)
		}

	default:
		log.Fatal("expected a command")
	}

}
