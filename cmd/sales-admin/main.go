package main

import (
	"flag"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"log"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
