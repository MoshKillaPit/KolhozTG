package handlers

import (
	"FarmTG/db"
	"FarmTG/utils"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

// Обработка команды /start
func HandleStartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, dbConn *sql.DB) {
	userID := update.Message.Chat.ID
	userName := update.Message.From.FirstName

	// Удаляем сообщение пользователя спустя некоторое время
	utils.DeleteMessageAfter(bot, userID, update.Message.MessageID, 10*time.Second)

	// Проверяем, есть ли пользователь в базе данных
	_, err := db.GetStateFromDB(dbConn, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаем нового пользователя
			err := db.CreateUserInDB(dbConn, userID, userName, "waiting_for_answer")
			if err != nil {
				reply := "Произошла ошибка при создании нового пользователя, попробуйте позже."
				utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
				return
			}
		} else {
			log.Printf("Ошибка при получении состояния: %v", err)
			reply := "Произошла ошибка, попробуйте позже."
			utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
			return
		}
	} else {
		// Обновляем состояние до начального
		err = db.SetStateInDB(dbConn, userID, "waiting_for_answer")
		if err != nil {
			reply := "Произошла ошибка при обновлении состояния, попробуйте позже."
			utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
			return
		}
	}

	// Приветственное сообщение
	reply := fmt.Sprintf("Привет, %s! Готов ли ты познать мир колхоза?", userName)

	// Создаем клавиатуру с кнопками "Да" и "Нет"
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Да"),
			tgbotapi.NewKeyboardButton("Нет"),
		),
	)

	// Отправляем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(userID, reply)
	msg.ReplyMarkup = keyboard
	sentMessage, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
		return
	}

	// Удаляем сообщение бота спустя некоторое время
	utils.DeleteMessageAfter(bot, userID, sentMessage.MessageID, 5*time.Second)
}
