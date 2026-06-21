package main

import (
	"context"
	"log"

	"canarias.run/internal/app"
	"canarias.run/internal/web"
)

func main() {
	log.Println("=== Inicializace canarias.run ===")

	application, err := app.New()
	if err != nil {
		log.Fatalf("Chyba konfigurace/aplikace: %v", err)
	}

	application.StartScraperDaemon(context.Background())

	server, err := web.New(application)
	if err != nil {
		log.Fatalf("Chyba inicializace web vrstvy: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Chyba pri spousteni serveru: %v", err)
	}
}
