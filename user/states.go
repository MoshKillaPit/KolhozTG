package user

import (
	"FarmTG/db"
	"database/sql"
)

// Установка состояния пользователя в базе данных
func SetState(dbConn *sql.DB, userID int64, state string) {
	db.SetStateInDB(dbConn, userID, state)
}

// Получение состояния пользователя из базы данных
func GetState(dbConn *sql.DB, userID int64) string {
	state, err := db.GetStateFromDB(dbConn, userID)
	if err != nil {
		return "waiting_for_answer" // Возвращаем начальное состояние в случае ошибки
	}
	return state
}
