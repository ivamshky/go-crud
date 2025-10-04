package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ivamshky/go-crud/handler"
	"github.com/ivamshky/go-crud/middlewares"
	"github.com/ivamshky/go-crud/repository"
	"github.com/ivamshky/go-crud/service"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	db, err := sql.Open("mysql", "dbadmin:1234@/userdb")
	if err != nil {
		slog.Error("Error connecting to DB", "err", err)
		return
	}
	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)))

	mux.Handle("GET /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleGetUserById)))
	mux.Handle("POST /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleCreateUser)))
	mux.Handle("GET /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleListUsers)))
	mux.Handle("DELETE /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleDeleteUser)))

	log.Print("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LogRequest(mux)))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}
