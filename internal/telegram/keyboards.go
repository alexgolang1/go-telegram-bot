package telegram

import (
	"github.com/go-telegram/bot/models"
)

func MainKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Наш сайт", URL: "https://google.com"},
				{Text: " Услуги и цены", CallbackData: "show_prices"},
				{Text: "📝 Оставить заявку", CallbackData: "create_order"},
			},
		},
	}
}

func BackKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "⬅️ Назад в меню", CallbackData: "go_back"},
			},
		},
	}
}
