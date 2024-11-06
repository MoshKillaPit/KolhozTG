package handlers

import (
	"FarmTG/db"
	"FarmTG/utils"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

// Обработка команды /reset
func HandleResetCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI, dbConn *sql.DB) {
	userID := update.Message.Chat.ID

	// Удаляем сообщение пользователя спустя некоторое время
	utils.DeleteMessageAfter(bot, userID, update.Message.MessageID, 5*time.Second)

	err := db.ResetUserData(dbConn, userID)
	if err != nil {
		reply := "Произошла ошибка при сбросе данных, попробуйте позже."
		utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
		return
	}

	reply := "Ваши данные были успешно сброшены. Вы можете начать сначала, отправив команду /start."
	utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
	log.Printf("Данные пользователя с ID %d были сброшены", userID)
}
