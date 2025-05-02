package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var users []UserData

	GetJson(w, &users)

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
