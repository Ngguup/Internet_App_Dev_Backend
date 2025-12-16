package main

import (
	"context"
	"log"
	"os"
	"lab1/internal/pkg/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Println("cant create application", err)
		os.Exit(2)
	}

	log.Println("Application start!")

	// HTTPS запуск
	application.RunApp()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Application terminated!")
}
