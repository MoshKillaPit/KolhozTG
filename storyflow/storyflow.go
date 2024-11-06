package storyflow

import (
	"FarmTG/db"
	"FarmTG/utils"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

// StartStorySequence автоматически проходит через этапы истории для пользователя
func StartStorySequence(bot *tgbotapi.BotAPI, dbConn *sql.DB, userID int64) {
	go func() {
		var reply string

		// Этап storyTelling1
		reply = "Поздравляю с регистрацией, позволь я поделюсь с тобой кое-чем!"
		utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
		db.SetStateInDB(dbConn, userID, "storyTelling2")

		// Этап storyTelling2
		time.Sleep(5 * time.Second) // Пауза перед следующим этапом
		reply = "На этой ферме всё совсем не так просто, как кажется! Здесь даже земля иногда жужжит, когда слишком долго без работы."
		utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
		db.SetStateInDB(dbConn, userID, "storyTelling3")

		// Этап storyTelling3
		time.Sleep(5 * time.Second)
		reply = "Прежде чем ты вступишь на эту землю, тебе предстоит выбрать свою роль. От выбора зависит, кто из фермерской братии будет помогать тебе…"
		utils.SendMessageWithDeletion(bot, userID, reply, 5*time.Second)
		db.SetStateInDB(dbConn, userID, "difficult")

		time.Sleep(5 * time.Second)

		// Получение имени пользователя из базы данных
		userName, err := db.GetUserNameFromDB(dbConn, userID)
		if err != nil {
			log.Printf("Ошибка при получении имени пользователя: %v", err)
			return
		}

		// Этап difficult
		reply1 := "Старый фермер пожимает тебе руку и с улыбкой говорит:"
		reply2 := fmt.Sprintf("Эй %s, ну что, готов к фермерским приключениям?", userName)

		bot.Send(tgbotapi.NewMessage(userID, reply1))
		time.Sleep(3 * time.Second)
		bot.Send(tgbotapi.NewMessage(userID, reply2))

		// Сбрасываем состояние пользователя после завершения сценария
		db.SetStateInDB(dbConn, userID, "")
	}()
}
