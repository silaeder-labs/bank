package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

// Глобальная переменная для БД (в реальных проектах лучше передавать через контекст или структуру)
var db *sql.DB

func usersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, balance FROM users")
	if err != nil {
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User // Слайс для хранения всех найденных пользователей

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Balance); err != nil {
			log.Printf("Ошибка сканирования: %v", err)
			continue
		}
		users = append(users, u)
	}

	// Устанавливаем заголовок, чтобы браузер понял, что это JSON
	w.Header().Set("Content-Type", "application/json")

	// Кодируем слайс в JSON и отправляем в ResponseWriter
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Ошибка при формировании JSON", http.StatusInternalServerError)
	}
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение "id" из пути
	idStr := r.PathValue("id")

	// Превращаем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Запрос к БД
	var u User
	err = db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Balance)

	if err == sql.ErrNoRows {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func changeUserBalanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	// Структура для входящих данных
	var input struct {
		Value int `json:"value"`
	}

	// Читаем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	changeBalance(id, input.Value)
	w.WriteHeader(http.StatusNoContent)
}

func setUserBalanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значение "id" из пути
	idStr := r.PathValue("id")
	idValue := r.PathValue("value")

	// Превращаем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}
	value, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, "Неверный формат суммы", http.StatusBadRequest)
		return
	}

	setBalance(id, value)
}

func main() {
	dsn := "postgres://test_superuser:password@localhost:5432/testdb?sslmode=disable"

	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к драйверу: %v", err)
	}
	defer db.Close()

	// Проверяем соединение один раз при старте
	if err := db.Ping(); err != nil {
		log.Fatalf("База недоступна: %v", err)
	}

	http.HandleFunc("GET /users", usersHandler)
	http.HandleFunc("GET /users/{id}", getUserByIDHandler)
	http.HandleFunc("POST /change-balance/{id}", changeUserBalanceByIDHandler)
	http.HandleFunc("POST /set-balance/{id}", setUserBalanceByIDHandler)

	fmt.Println("Сервер запущен на http://localhost:8080/users")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setBalance(id int, summ int) {
	_, err := db.Exec(
		"UPDATE users SET balance = $1 WHERE id = $2",
		summ,
		id,
	)

	if err != nil {
		log.Printf("Internal error: %v", err)
	}
}

func changeBalance(id int, summ int) {
	_, err := db.Exec(
		"UPDATE users SET balance = balance + $1 WHERE id = $2",
		summ,
		id,
	)

	if err != nil {
		log.Printf("Internal error: %v", err)
	}
}
