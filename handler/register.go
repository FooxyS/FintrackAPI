package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserData struct {
	Email string `json:"email"`
	User  `json:"user"`
}

type UserStore interface {
	Register(user User) error
	Exists(username string) bool
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//проверка на правильный метод
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//декодирование json запроса
	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode json request", http.StatusBadRequest)
		return
	}

	//валидация данных пользователя
	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "empty fields in json", http.StatusBadRequest)
		return
	}

	//запись юзера в мапу
	var users map[string]User

	buf, errRead := os.ReadFile("C:/Users/FooxyS/Desktop/FintrackAPI/data/data.json")
	if errRead != nil {
		log.Printf("failed to read the file: %v\n", errRead)
		http.Error(w, "failed to read the file", http.StatusInternalServerError)
		return
	}
	errUnmarsh := json.Unmarshal(buf, &users)
	if errUnmarsh != nil {
		http.Error(w, "failed to Unmarshal the file", http.StatusInternalServerError)
		return
	}

	if _, ok := users[user.Email]; ok {
		http.Error(w, "User is already created", http.StatusOK)
		return
	}
	users[user.Email] = user.User

	//кодирование мапы
	compbuf, errMarsh := json.Marshal(users)
	if errMarsh != nil {
		http.Error(w, "failed to marshal data", http.StatusInternalServerError)
		return
	}
	os.WriteFile("data.json", compbuf, 0644)

	http.Error(w, "user is successfully created", http.StatusCreated)
}
