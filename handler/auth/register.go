package handler

import (
	"crypto/rand"
	"encoding/base64"
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

var dataBase map[string]string

type AccessToken struct {
	UserID string `json:"userid"`
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

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, errRand := rand.Read(b)
	if errRand != nil {
		return "", errRand
	}
	res := base64.URLEncoding.EncodeToString(b)
	return res, nil
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
	if _, ok := dataBase[user.Email]; ok {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("user already exists"))
		return
	}

	//FALSE -> to hashedPass the password, add it in the BD
	hashedPass, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Printf("error with hashing: %v", errHash)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dataBase[user.Email] = string(hashedPass)

	//get jwtkey from .env
	errGotEnv := godotenv.Load()
	if errGotEnv != nil {
		log.Printf("error with loading env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwtkey := os.Getenv("JWT_KEY")

	//generate the new UUID and return it as JWT tokens (access and refresh)
	id := uuid.New()

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

	//generating refresh token (base64)
	refreshToken, errRefresh := GenerateRefreshToken()
	if errRefresh != nil {
		log.Printf("error with gen refresh: %v", errRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//hashing refresh token and push it to data base
	HashedRefresh, errHashRefresh := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if errHashRefresh != nil {
		log.Printf("error with hashing refresh: %v", errHashRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//имитация добавления в базу данных
	log.Printf("добавление хэшированный рефреш токен в базу данных: %v", HashedRefresh)

	//sending tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	resp := AccessToken{
		UserID: string(id.String()),
		Access: accessToken,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(resp)
	if errEncode != nil {
		log.Printf("error with encoding json: %v", errEncode)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
