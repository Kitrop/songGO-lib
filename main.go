package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kitrop/songGO-lib/config"
	"github.com/Kitrop/songGO-lib/handlers"
	"github.com/Kitrop/songGO-lib/repository"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi"
	_ "github.com/Kitrop/songGO-lib/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Song API
// @version 1.0
// @description API for managing songs.
// @host localhost:8080
// @BasePath /
func main() {
	// Настройка логгера
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[INFO] Запускаем приложение...")

	// Подключаем env
	config.LoadEnv()


	// Подключаем БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	db, err := config.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("[ERROR] Ошибка подключения к базе данных: %v", err)
	}

	// Закрытие подключения к БД
	defer func() {
		log.Println("[INFO] Закрытие соединения с базой данных....")
		db.Close()
	}()

	// Создаем репозиторий
	repo := repository.NewSongRepository()


	// Создаем обработчики
	handler := handlers.NewSongHandler(repo)


	// Настраиваем маршруты
	r := chi.NewRouter()

	// Встраиваем middleware
	r.Use(LoggerMiddleware)        // Логирование запросов
	r.Use(middleware.Recoverer)    // Восстановление после паники

	// CRUD операции
	r.Get("/songs", handler.GetAllSongs)        // Получить список всех песен
	r.Get("/songs/{id}", handler.GetSongByID)   // Получить песню по ID
	r.Post("/songs", handler.CreateSong)        // Добавить новую песню
	r.Put("/songs/{id}", handler.UpdateSong)    // Обновить существующую песню
	r.Delete("/songs/{id}", handler.DeleteSong) // Удалить песню
	
	// Подключение Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	
	// Запускаем сервер
	port := os.Getenv("PORT")
	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// LoggerMiddleware логирует информацию о каждом запросе.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[INFO] Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Передаём управление следующему обработчику
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[INFO] Completed request: %s %s in %v", r.Method, r.URL.Path, duration)
	})
}
