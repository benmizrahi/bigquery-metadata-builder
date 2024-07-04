package main

import (
	"log"

	"github.com/benmizrahi/bigquery-metadata-builder/internal/datasources"
	"github.com/benmizrahi/bigquery-metadata-builder/internal/intelligence"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	ai := intelligence.Resolver(nil)

	ds := datasources.Resolver(nil).
		Explore().
		BuildVector()

	ai.SuggestMetadata(ds)

}
