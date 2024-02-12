package main

import (
	"log"

	"github.com/phillipmugisa/go_resume_generator/app"
)

func main() {
	PORT := 8080

	a := app.NewAppServer(PORT)
	log.Fatal(a.Run())
}
