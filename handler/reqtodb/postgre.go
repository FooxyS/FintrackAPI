package reqtodb

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() {
	//подключение к базе данных
	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatalf("error with loading env file: %v\n", errLoadEnv)
	}

	pgpool, errConnDB := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if errConnDB != nil {
		log.Fatalf("error with connecting to postgresql: %v\n", errConnDB)
	}
	defer pgpool.Close()
}
