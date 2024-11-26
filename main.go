package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Yamaguchi-Katsuya/golang-gacha-app-api/handler"
	"github.com/Yamaguchi-Katsuya/golang-gacha-app-api/handler/middleware"
	"github.com/Yamaguchi-Katsuya/golang-gacha-app-api/service"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}

	db, err := sql.Open("mysql", "gacha-app:password@tcp(localhost:3306)/gacha-app")
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}

	defer db.Close()

	userHandler := handler.NewUserHandler(service.NewUserService(db))
	gachaHandler := handler.NewGachaHandler(service.NewGachaService(db))
	characterHandler := handler.NewCharacterHandler(service.NewCharacterService(db))

	mux := http.NewServeMux()
	mux.Handle("/user/", middleware.ActionLog(userHandler))
	mux.Handle("/gacha/", middleware.ActionLog(gachaHandler))
	mux.Handle("/character/", middleware.ActionLog(characterHandler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("Server is starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
