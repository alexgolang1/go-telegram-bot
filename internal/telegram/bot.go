package telegram

import (
	"context"
	"log"
	"os"
	"os/signal"
	"tg/bot/internal/config"

	"github.com/go-telegram/bot"
	"gorm.io/gorm"
)

type BotManager struct {
	tgBot *bot.Bot
	Cfg   *config.Config
	DB    *gorm.DB
}

func NewBotManager(Cfg *config.Config, db *gorm.DB) *BotManager {
	return &BotManager{
		Cfg: Cfg,
		DB:  db,
	}
}

// START BOT
func (bm *BotManager) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(bm.HandleMessage),
		bot.WithServerURL(bm.Cfg.ApiURL),
	}

	var err error

	bm.tgBot, err = bot.New(bm.Cfg.BotToken, opts...)
	if err != nil {
		log.Fatalf("Критическая ошибка создания бота: %v", err)
	}

	bm.tgBot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		bm.HandleStart,
	)

	bm.tgBot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"show_prices",
		bot.MatchTypeExact,
		bm.HandleShowPrices,
	)

	bm.tgBot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"go_back",
		bot.MatchTypeExact,
		bm.HandleBackMenu,
	)

	log.Println("Бот на Go успешно запущен...")
	bm.tgBot.Start(ctx)
}
