package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Функция для инициализации подключения к базе данных
func InitDB(logError func(string)) (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logError(fmt.Sprintf("Ошибка при создании подключения к базе данных: %v", err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logError(fmt.Sprintf("Ошибка при пинге базы данных: %v", err))
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id BIGINT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            state VARCHAR(255)
        );
    `)
	if err != nil {
		logError(fmt.Sprintf("Ошибка при создании таблицы: %v", err))
		return nil, err
	}

	return db, nil
}

// Сохранение или обновление пользователя в базе данных
func SaveUserToDB(db *sql.DB, userID int64, userName string, userState string) error {
	query := `INSERT INTO users (id, name, state) VALUES (?, ?, ?)
              ON DUPLICATE KEY UPDATE name = VALUES(name), state = VALUES(state)`
	_, err := db.Exec(query, userID, userName, userState)
	if err != nil {
		log.Printf("Ошибка при сохранении пользователя: %v", err)
		return err
	}
	return nil
}

// Сохранение состояния пользователя в базе данных
func SetStateInDB(db *sql.DB, userID int64, state string) error {
	query := `UPDATE users SET state = ? WHERE id = ?`
	_, err := db.Exec(query, state, userID)
	if err != nil {
		log.Printf("Ошибка при сохранении состояния: %v", err)
		return err
	}
	return nil
}

// Получение состояния пользователя из базы данных
func GetStateFromDB(db *sql.DB, userID int64) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}
	var state string
	query := `SELECT state FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&state)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		log.Printf("Ошибка при получении состояния: %v", err)
		return "", err
	}
	return state, nil
}

// Создание нового пользователя в базе данных
func CreateUserInDB(db *sql.DB, userID int64, userName string, initialState string) error {
	query := `INSERT INTO users (id, name, state) VALUES (?, ?, ?)`
	_, err := db.Exec(query, userID, userName, initialState)
	if err != nil {
		log.Printf("Ошибка при добавлении нового пользователя: %v", err)
		return err
	}
	log.Printf("Новый пользователь %s (ID: %d) успешно добавлен с состоянием '%s'", userName, userID, initialState)
	return nil
}

// Получение имени пользователя из базы данных
func GetUserNameFromDB(db *sql.DB, userID int64) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var name string
	query := `SELECT name FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Пользователь с ID %d не найден", userID)
		}
		log.Printf("Ошибка при получении имени пользователя: %v", err)
		return "", err
	}
	return name, nil
}

// Удаление данных пользователя
func ResetUserData(db *sql.DB, userID int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, userID)
	if err != nil {
		log.Printf("Ошибка при удалении данных пользователя: %v", err)
		return err
	}
	log.Printf("Данные пользователя с ID %d успешно удалены", userID)
	return nil
}
