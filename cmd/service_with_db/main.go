package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"homework3/internal/pkg/db"
	"homework3/internal/pkg/repository/postgresql"
	"homework3/internal/pkg/server"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDBWithDSN(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	mainServer := server.Server{ArticleRepo: postgresql.NewArticles(database), CommentRepo: postgresql.NewComments(database)}
	http.Handle("/", server.CreateRouter(mainServer))
	if err := http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
