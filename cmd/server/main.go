package main

import (
	"context"
	"go-news-comments/pkg/api"
	"go-news-comments/pkg/config"
	"go-news-comments/pkg/storage"
	"go-news-comments/pkg/storage/pg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	AppName = "go-news-comments"
)

type server struct {
	db  storage.Storage
	api *api.API
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if cfg == nil {
		os.Exit(0)
	}

	log.Printf("starting %s service\n", AppName)
	log.Println(config.VersionString())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	database, err := pg.New(cfg.DBConnString)
	if err != nil {
		log.Fatal(err)
	}

	server := api.New(cfg, database)
	if err := server.Run(ctx); err != nil {
		log.Println("Got error:", err)
		os.Exit(0)
	}
}
