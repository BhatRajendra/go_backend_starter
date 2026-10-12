package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"sandhi/internal/app"
	"sandhi/internal/db"
	"sandhi/internal/logging"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DB"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	defer dbPool.Close()
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping to db")
	}

	queries := db.New(dbPool)
	logger := logging.New()

	app := app.NewApp(dbPool, queries, logger)

	logger.Info("Server has started listening", "port", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), app.Router())
}

// init() here does only loading of env variables and check if env exists
func init() {
	// Load .env when available.
	// in prod, there won't be any .env file
	// Production can supply environment variables directly.
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Failed to load .env: %v", err)
		}
	}
	// check each of the env var exists or not.
	envVariables := []string{"DB", "PORT"} // add all the future env variables here for validation
	for _, variable := range envVariables {
		if os.Getenv(variable) == "" {
			log.Fatalf("%s is not configured", variable)
		}
	}
}
