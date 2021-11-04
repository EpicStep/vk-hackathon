package main

import (
	"context"
	"errors"
	"github.com/EpicStep/vk-hackathon/internal/config"
	"github.com/EpicStep/vk-hackathon/internal/image"
	"github.com/EpicStep/vk-hackathon/internal/router"
	"github.com/EpicStep/vk-hackathon/pkg/database"
	"github.com/EpicStep/vk-hackathon/pkg/server"
	"github.com/go-chi/chi/v5"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	cfg, err := config.New()
	if err != nil {
		return err
	}

	db, err := database.New(cfg.MySQLURL)
	if err != nil {
		return err
	}

	r := router.New()

	service := image.New(db)

	r.Route("/", func(r chi.Router) {
		service.Routes(r)
	})

	addr := ":" + cfg.Port

	srv := server.New(addr, r)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		return errors.New("server shutdown failed")
	}

	if err := db.Close(); err != nil {
		return errors.New("database close failed")
	}

	return nil
}
