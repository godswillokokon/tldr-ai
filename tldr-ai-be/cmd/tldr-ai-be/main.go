package main

import (
	"log"
	"os"

	"tldr-ai-be/internal/app"
	"tldr-ai-be/internal/config"
)

func main() {
	if err := config.LoadDotEnvOverride(".env"); err != nil {
		log.Fatal(err)
	}
	if err := config.LoadDotEnvOverride("tldr-ai-be/.env"); err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
