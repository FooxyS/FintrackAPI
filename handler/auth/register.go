package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type AccessToken struct {
	Access string `json:"access"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MyCustomClaims struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//parsing json user's data (email and password) and validate it
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		return
	}
	if user.Email == "" || user.Password == "" {
		http.Error(w, "has empty fields", http.StatusBadRequest)
		return
	}

	//check him in the data base
	//TRUE -> message(user already exists) (409 - Conflict)
	dataBase := make(map[string]string)
	if _, ok := dataBase[user.Email]; ok {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("user already exists"))
		return
	}

	//FALSE -> to hash the password, add it in the BD
	hash, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Printf("error with hashing: %v", errHash)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dataBase[user.Email] = string(hash)

	//generate the new UUID and return it as JWT tokens (access and refresh)
	id, errUUID := uuid.NewUUID()
	if errUUID != nil {
		log.Printf("error with generating UUID: %v", errUUID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//get jwtkey from .env
	errGotEnv := godotenv.Load()
	if errGotEnv != nil {
		log.Printf("error with loading env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwtkey := os.Getenv("JWT_KEY")

	//access token
	accessClaims := MyCustomClaims{
		UserID: id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, errWithAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString([]byte(jwtkey))
	if errWithAccess != nil {
		log.Printf("error with generating access token: %v", errWithAccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//sending tokens
}

/*
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserData struct {
	Email string `json:"email"`
	User  User   `json:"user"`
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

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal(errEnv)
	}

	//декодирование json запроса
	var user UserData
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode json request", http.StatusBadRequest)
		return
	}

	//валидация данных пользователя
	if user.User.Username == "" || user.User.Password == "" || user.Email == "" {
		http.Error(w, "empty fields in json", http.StatusBadRequest)
		return
	}

	//запись юзера в мапу
	var users []UserData

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

	for _, curUser := range users {
		if curUser.Email == user.Email {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Пользователь уже существует"))
			return
		}
	}

	users = append(users, user)

	//кодирование мапы
	compbuf, errMarsh := json.Marshal(users)
	if errMarsh != nil {
		http.Error(w, "failed to marshal data", http.StatusInternalServerError)
		return
	}
	os.WriteFile("C:/Users/FooxyS/Desktop/FintrackAPI/data/data.json", compbuf, 0644)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно создан"))
}
*/
