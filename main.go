package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/phillipmugisa/go_resume_generator/app"
	"github.com/phillipmugisa/go_resume_generator/storage"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	PORT := os.Getenv("PORT")
	// storage service
	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.SetUpDB(); err != nil {
		log.Fatal(err)
	}

	a := app.NewAppServer(PORT, store)
	log.Fatal(a.Run())
}
