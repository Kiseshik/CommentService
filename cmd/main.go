package main

import (
	"log"

	"github.com/Kiseshik/CommentService.git/internal/config"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Config loaded: %+v", cfg)
}
