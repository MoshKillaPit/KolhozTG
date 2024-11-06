package handlers

import (
	"FarmTG/db"
	"FarmTG/utils"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

// Обработка пользовательского ответа
func HandleUserResponse(update tgbotapi.Update, bot *tgbotapi.BotAPI, dbConn *sql.DB) {
	if update.Message == nil {
		return
	}
	userID := update.Message.Chat.ID

	// Удаляем сообщение пользователя спустя некоторое время
	if !update.Message.Chat.IsPrivate() {
		utils.DeleteMessageAfter(bot, userID, update.Message.MessageID, 15*time.Second)
	}

	// Получаем текущее состояние пользователя
	currentState, err := db.GetStateFromDB(dbConn, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			reply := "Пожалуйста, начните с команды /start."
			utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second, nil)
			return
		} else {
			reply := "Произошла ошибка при получении состояния, попробуйте позже."
			utils.SendMessageWithDeletion(bot, userID, reply, 15*time.Second, nil)
			return
		}
	}

	switch currentState {
	// ... предыдущие case ...

	case "waiting_for_adventure_answer":
		userResponse := update.Message.Text
		// Убираем клавиатуру после ответа пользователя
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)

		if userResponse == "Да" {
			reply := "Отлично! Приключения ждут тебя!"
			utils.SendMessageWithDeletion(bot, userID, reply, 30*time.Second, removeKeyboard)
			// Обновляем состояние пользователя или продолжаем сценарий
			// Например, можно запустить следующую часть истории
			db.SetStateInDB(dbConn, userID, "adventure_started")
			// Запуск следующей части истории
			// storyflow.StartNextAdventure(bot, dbConn, userID)
		} else if userResponse == "Нет" {
			reply := "Очень жаль. Если передумаешь, напиши мне снова."
			utils.SendMessageWithDeletion(bot, userID, reply, 30*time.Second, removeKeyboard)
			// Сбрасываем состояние пользователя
			db.SetStateInDB(dbConn, userID, "")
		} else {
			reply := "Пожалуйста, выбери 'Да' или 'Нет'."
			utils.SendMessageWithDeletion(bot, userID, reply, 15*time.Second, nil)
		}

		// ... остальные case ...
	}
}
