package telegram

import (
	"context"
	"fmt"
	"log"
	"tg/bot/internal/domain"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bm *BotManager) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message == nil || update.Message.Text == "" {
		return
	}

	chatID := update.Message.Chat.ID
	userText := update.Message.Text

	var user domain.User
	err := bm.DB.Where("id = ?", chatID).First(&user).Error
	if err != nil {
		log.Printf("[ERROR] Не удалось получить стейт юзера %d: %v", chatID, err)
		return
	}

	switch user.State {

	case domain.StateWaitName:
		err = domain.ValidateMessages(userText, 3, 30)
		if err != nil {
			log.Printf("[FSM_WARNING] Имя не прошло валидацию: %v", err)
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "⚠️ **Уведомление безопасности**\n\nВведённое имя нарушает правила системы.\n\n📐 **Лимиты:** от 3 до 30 символов.\nПожалуйста, введите имя заново:",
			})
			return
		}

		log.Printf("[B2B_FSM] Юзер %d прислал валидное имя: %s", chatID, userText)

		err = bm.DB.Model(&models.User{}).Where("id = ?", chatID).Updates(map[string]interface{}{
			"first_name": userText,
			"state":      domain.StateWaitDesc,
		}).Error
		if err != nil {
			log.Printf("[ERROR] Не удалось сохранить имя: %v", err)
			return
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "✅ **Имя зафиксировано в системе.**\n\nШаг 2: Теперь подробно опишите суть вашего проекта или ТЗ в одном сообщении:",
		})

	case domain.StateWaitDesc:
		err = domain.ValidateMessages(userText, 100, 1000)
		if err != nil {
			log.Printf("[FSM_WARNING] ТЗ не прошло валидацию: %v", err)
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "⚠️ **Уведомление безопасности**\n\nВаше ТЗ слишком короткое или превышает лимит.\n\n📐 **Лимиты:** от 100 до 1000 символов.\nПожалуйста, опишите проект подробнее в одном сообщении:",
			})
			return
		}

		log.Printf("[B2B_FSM] Юзер %d прислал валидное ТЗ", chatID)

		var updatedUser models.User
		if err := bm.DB.Where("id = ?", chatID).First(&updatedUser).Error; err != nil {
			log.Printf("[ERROR] Ошибка получения данных юзера: %v", err)
			return
		}

		err = bm.DB.Model(&models.User{}).Where("id = ?", chatID).Updates(map[string]interface{}{
			"description": userText,
			"state":       domain.StateNone,
		}).Error
		if err != nil {
			log.Printf("[ERROR] Не удалось сохранить ТЗ в базу: %v", err)
			return
		}

		var adminChatID int64 = 0000000001

		adminMessage := fmt.Sprintf(
			"🚨 **НОВАЯ ЗАЯВКА В ПОРТФОЛИО!**\n\n"+
				"👤 **Клиент:** %s\n"+
				"🆔 **TG ID:** `%d`\n"+
				"🌐 **Username:** @%s\n\n"+
				"📝 **Техническое Задание:**\n_%s_",
			updatedUser.FirstName, chatID, updatedUser.Username, userText,
		)

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    adminChatID,
			Text:      adminMessage,
			ParseMode: "Markdown",
		})

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        "🎉 **Отлично! Ваша коммерческая заявка успешно записана в базу данных и передана администратору.**\n\nСпасибо за обращение!",
			ReplyMarkup: MainKeyboard(),
		})

	default:
		if userText == "/start" {
			bm.HandleStart(ctx, b, update)
		}
	}
}
