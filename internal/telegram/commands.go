package telegram

import (
	"context"
	"log"
	"tg/bot/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bm *BotManager) HandleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	firstName := update.Message.From.FirstName

	exists, err := repository.UserExists(bm.DB, chatID)
	if err != nil {
		log.Printf("Ошибка проверки юзера в БД: %v", err)
	}

	var welcomeText string

	if exists {

		log.Printf("[DB] Пользователь %d уже существует. Выдаем кастомный текст.", chatID)
		welcomeText = "👋 **С возвращением! Рады видеть вас снова.**\n\nВыберите интересующий раздел меню ниже:"

	} else {

		log.Printf("[DB] Новый пользователь %d. Записываем в базу...", chatID)

		newUser := models.User{
			ID:        chatID,
			Username:  username,
			FirstName: firstName,
		}

		err = bm.DB.Create(&newUser).Error
		if err != nil {
			log.Printf("[ERROR] Не удалось сохранить нового пользователя %d: %v", chatID, err)
		}

		welcomeText = "Приветствуем! Мы рады видеть вас. Выберите интересующий раздел меню ниже:"

	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        welcomeText,
		ReplyMarkup: MainKeyboard(),
	})

	if err != nil {
		log.Printf("Не удалось отправить приветственное сообщение: %v", err)
		return
	}
}
