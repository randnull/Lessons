package main

import (
	"context"
	"github.com/randnull/Lessons/internal/app"
	"github.com/randnull/Lessons/internal/config"
	lg "github.com/randnull/Lessons/internal/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Graceful shutdown started")
		stop()
	}()

	application := app.NewApp(NewConfig)

	err = application.Run(ctx)

	if err != nil {
		log.Fatal("app run error: " + err.Error())
	}
}
