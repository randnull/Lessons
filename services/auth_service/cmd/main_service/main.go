package main

import (
	"github.com/randnull/Lessons/internal/app"
	"github.com/randnull/Lessons/internal/config"
	lg "github.com/randnull/Lessons/internal/logger"
)

func main() {
	NewConfig, err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	err = lg.InitLogger()

	if err != nil {
		panic(err)
	}

	a := app.NewApp(NewConfig)
	a.Run()
}
