package main

import (
	"flag"
	"garagesale/internal/platform/conf"
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"github.com/pkg/errors"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := conf.Parse()

	db, err := database.Open(config.Database)
	if err != nil {
		return err
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "Migrating tables")
		}
		log.Println("Migration Completed")

	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "Applying seed")
		}
		log.Println("Seeding completed")
	}

	return nil
}
