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
		lg.Info("Graceful shutdown signal received")
		stop()
	}()

	a := app.NewApp(NewConfig)

	err = a.Run(ctx)

	if err != nil {
		log.Fatalf("app run error: %v", err.Error())
	}
}
