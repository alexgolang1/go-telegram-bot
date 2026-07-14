package main

import (
	"log"
	"tg/bot/internal/config"
	"tg/bot/internal/repository"
	"tg/bot/internal/telegram"
)

func main() {
	log.Println("Инициализация приложения...")

	cfg := config.LoadConfig()

	db, err := repository.InitDB(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	botManager := telegram.NewBotManager(cfg, db)

	botManager.Start()
}
