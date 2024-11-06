package main

import (
	"FarmTG/bot"
	"FarmTG/db"
	"FarmTG/handlers"
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Настройка формата логов
func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Функции логирования с уровнями
func logInfo(message string) {
	log.Printf("[INFO] %s", message)
}

func logWarning(message string) {
	log.Printf("[WARNING] %s", message)
}

func logError(message string) {
	log.Printf("[ERROR] %s", message)
}

// Обработка обновлений с логированием
func processUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, dbConn *sql.DB) {
	for update := range updates {
		if update.Message != nil {
			logInfo(fmt.Sprintf("Получено сообщение от пользователя %s: %s", update.Message.From.UserName, update.Message.Text))

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					logInfo("Обработка команды /start")
					handlers.HandleStartCommand(update, bot, dbConn)
				case "reset":
					logInfo("Обработка команды /reset")
					handlers.HandleResetCommand(update, bot, dbConn)
				default:
					logWarning(fmt.Sprintf("Неизвестная команда: %s", update.Message.Command()))
				}
			} else {
				logInfo("Обработка пользовательского сообщения")
				handlers.HandleUserResponse(update, bot, dbConn)
			}
		}
	}
}

func main() {
	// Инициализация базы данных
	dbConn, err := db.InitDB(logError)
	if err != nil {
		logError(fmt.Sprintf("Ошибка подключения к базе данных: %v", err))
		return
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logWarning(fmt.Sprintf("Ошибка при закрытии подключения к базе данных: %v", err))
		} else {
			logInfo("Подключение к базе данных успешно закрыто")
		}
	}()

	// Инициализация бота
	botInstance := bot.InitBot()
	logInfo("Бот успешно инициализирован")

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botInstance.GetUpdatesChan(u)

	// Логирование успешного начала процесса обработки обновлений
	logInfo("Начало обработки обновлений")

	// Обработка обновлений
	processUpdates(updates, botInstance, dbConn)
}
