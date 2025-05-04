package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var users []UserData

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal(errEnv)
	}
	path := os.Getenv("PATH_TO_JSON")

	//GetJson(w, &users)
	buf, errRead := os.ReadFile(path)
	if errRead != nil {
		log.Printf("failed to read the file: %v\n", errRead)
		http.Error(w, "failed to read the file", http.StatusInternalServerError)
		return
	}
	errUnmarsh := json.Unmarshal(buf, &users)
	if errUnmarsh != nil {
		log.Printf("failed to Unmarshal the file: %v\n", errUnmarsh)
		http.Error(w, "failed to Unmarshal the file", http.StatusInternalServerError)
		return
	}

	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode the json data: %v", err), http.StatusBadRequest)
		return
	}

	for _, curUser := range users {
		if curUser.Email == user.Email {
			if curUser.User.Password == user.User.Password {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Добро пожаловать! :)"))
				return
			}
			http.Error(w, "Неверный пароль", http.StatusBadRequest)
			return
		}
	}
	http.Error(w, "Пользователь не зарегестрирован", http.StatusBadRequest)
}
