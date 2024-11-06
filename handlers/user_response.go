package handlers

import (
	"FarmTG/db"
	"FarmTG/storyflow"
	"FarmTG/utils"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

// Обработка пользовательского ответа
func HandleUserResponse(update tgbotapi.Update, bot *tgbotapi.BotAPI, dbConn *sql.DB) {
	userID := update.Message.Chat.ID

	// Удаляем сообщение пользователя спустя некоторое время
	utils.DeleteMessageAfter(bot, userID, update.Message.MessageID, 5*time.Second)

	// Получаем текущее состояние пользователя
	currentState, err := db.GetStateFromDB(dbConn, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			reply := "Пожалуйста, начните с команды /start."
			utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second)
			return
		} else {
			reply := "Произошла ошибка при получении состояния, попробуйте позже."
			utils.SendMessageWithDeletion(bot, userID, reply, 15*time.Second)
			return
		}
	}

	switch currentState {
	case "waiting_for_answer":
		userResponse := update.Message.Text
		if userResponse == "Да" {
			// Обновляем состояние и просим ввести имя
			err := db.SetStateInDB(dbConn, userID, "registrationName")
			if err != nil {
				reply := "Произошла ошибка, попробуйте позже."
				utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second)
				return
			}
			reply := "Отлично! Пожалуйста, введи своё имя:"
			utils.SendMessageWithDeletion(bot, userID, reply, 20*time.Second)
		} else if userResponse == "Нет" {
			reply := "Очень жаль. Если передумаешь, напиши /start."
			utils.SendMessageWithDeletion(bot, userID, reply, 15*time.Second)
			// Сбрасываем состояние
			db.SetStateInDB(dbConn, userID, "")
		} else {
			reply := "Пожалуйста, выбери 'Да' или 'Нет'."
			utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second)
		}

	case "registrationName":
		userName := update.Message.Text
		if err := db.SaveUserToDB(dbConn, userID, userName, "storyTelling1"); err != nil {
			reply := "Произошла ошибка при сохранении данных, попробуйте позже."
			utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second)
		} else {
			reply := fmt.Sprintf("Отлично, %s! Ты успешно зарегистрирован на ферме!", userName)
			utils.SendMessageWithDeletion(bot, userID, reply, 20*time.Second)

			// Запускаем сценарий истории после регистрации
			storyflow.StartStorySequence(bot, dbConn, userID)
		}

	default:
		reply := "Пожалуйста, начни с команды /start."
		utils.SendMessageWithDeletion(bot, userID, reply, 10*time.Second)
	}
}
