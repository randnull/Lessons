package main

import (
	"github.com/randnull/Lessons/internal/app"
	"github.com/randnull/Lessons/internal/config"
	"log"
)

func main() {
	NewConfig, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(NewConfig)
	a.Run()
}
