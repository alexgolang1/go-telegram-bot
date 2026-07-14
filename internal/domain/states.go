package domain

import (
	"errors"
	"strings"
)

type State int

const (
	StateNone State = iota
	StateWaitName
	StateWaitDesc
)

func (s State) String() string {

	switch s {

	case StateNone:
		return "StateNone (Обычный режим)"

	case StateWaitName:
		return "StateWaitName (Ожидание имени)"

	case StateWaitDesc:
		return "StateWaitDesc (Ожидание ТЗ)"

	default:
		return "UnknownState (Неизвестный шаг)"

	}
}

func ValidateMessages(text string, minLenght int, maxLenght int) error {
	cleanText := strings.TrimSpace(text)

	if len(cleanText) < minLenght || len(cleanText) > maxLenght {

		return errors.New("Превышение/Недостаток символов в ответе для успешной интеграции в базу")
	}

	return nil
}
