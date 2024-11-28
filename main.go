package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kitrop/songGO-lib/config"
	"github.com/Kitrop/songGO-lib/handlers"
	"github.com/Kitrop/songGO-lib/repository"
	"github.com/go-chi/chi"
)

func main() {
	// Подгрузка env и базы данных
	config.LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	db, err := config.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	
	handlers.SongHandler(r, &repository.SongRepository{})

	http.ListenAndServe(":3000", r)
}