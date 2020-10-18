package main

import (
	"flag"
	"garagesale/internal/platform/conf"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"log"
)

func main() {
	config := conf.Parse()

	db, err := database.Open(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Applying migrations: ", err)
		}
		log.Println("Migration Completed")

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Applying seed: ", err)
		}
		log.Println("Seeding completed")
	}
}
