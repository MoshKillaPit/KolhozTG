package utils

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendMessageWithDeletion отправляет сообщение, ждет указанное время, а затем удаляет его
func SendMessageWithDeletion(bot *tgbotapi.BotAPI, chatID int64, text string, delay time.Duration, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}
	sentMessage, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
		return
	}

	// Удаляем отправленное сообщение спустя задержку
	DeleteMessageAfter(bot, chatID, sentMessage.MessageID, delay)
}

// DeleteMessageAfter удаляет сообщение после заданной задержки
func DeleteMessageAfter(bot *tgbotapi.BotAPI, chatID int64, messageID int, delay time.Duration) {
	time.AfterFunc(delay, func() {
		DeleteMessage(bot, chatID, messageID)
	})
}

// DeleteMessage удаляет сообщение по chatID и messageID
func DeleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)

	_, err := bot.Request(delMsg)
	if err != nil {
		log.Println("Ошибка удаления сообщения:", err)
	} else {
		log.Println("Сообщение успешно удалено")
	}
}
