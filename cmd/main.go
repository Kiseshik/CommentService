package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kiseshik/CommentService.git/internal/app"
	"github.com/Kiseshik/CommentService.git/internal/config"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Config loaded: %+v", cfg)

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	go func() {
		log.Printf("server starting on %s", cfg.ListenAddr)
		if err := application.Run(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	if err := application.Shutdown(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("server stopped")
}
