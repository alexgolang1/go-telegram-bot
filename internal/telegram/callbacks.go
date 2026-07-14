package telegram

import (
	"context"
	"log"
	"tg/bot/internal/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bm *BotManager) HandleShowPrices(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if err != nil {
		log.Printf("Не совпадение ID действия: %v", err)
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: BackKeyboard(),
		Text: "💰 Наши услуги и цены:\n\n" +
			"1. Разработка ТГ-бота на Go — от 5 000 руб.\n" +
			"2. Проектирование архитектуры API — от 10 000 руб.\n" +
			"3. Настройка серверов и конфигурации — от 3 000 руб.",
	})

	if err != nil {
		log.Printf("[ERROR] Ошибка отправки прайс-листа: %v", err)
	}

}

// EMPTY FIELD BEETWEN 2 FUNCTIONS

func (bm *BotManager) HandleBackMenu(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if err != nil {
		log.Printf("Не совпадение ID действия: %v", err)
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: MainKeyboard(),
		Text:        "Выберите интересующий раздел меню ниже:",
	})

	if err != nil {
		log.Printf("[ERROR] Ошибка отправки прайс-листа: %v", err)
	}
}

// ORDER CALLBACK CHECKER

func (bm *BotManager) HandleCreateOrder(ctx context.Context, b *bot.Bot, update *models.Update) {

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if err != nil {
		log.Printf("Не совпадение ID действия: %v", err)
	}

	ChatID := update.CallbackQuery.Message.Message.Chat.ID
	MessageID := update.CallbackQuery.Message.Message.ID

	err = bm.DB.Model(&models.User{}).Where("id = ?", ChatID).Update("state", domain.StateWaitName).Error
	if err != nil {
		log.Printf("[ERROR] Не удалось обновить состояние юзера %d: %v", ChatID, err)
		return
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    ChatID,
		MessageID: MessageID,
		Text:      "📝 **Запуск создания заявки**\n\nШаг 1: Пожалуйста, введите ваше **Имя** или никнейм для обратной связи:",
	})

	if err != nil {
		log.Printf("[ERROR] Ошибка EditMessageText в заявке: %v", err)
	}
}
