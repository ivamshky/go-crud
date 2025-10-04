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
	mux.HandleFunc("/swagger.html", serveAPISpec)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"dbadmin", // MySQL username
		"1234",    // MySQL password
		"db",      // Host IP or service name
		"3306",    // MySQL port
		"userdb")  //db name
	db, err := sql.Open("mysql", connString)
	if err != nil {
		slog.Error("Error connecting to DB", "err", err)
		return
	}
	defer db.Close()
	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)))

	mux.Handle("GET /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleGetUserById)))
	mux.Handle("POST /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleCreateUser)))
	mux.Handle("GET /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleListUsers)))
	mux.Handle("DELETE /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleDeleteUser)))

	log.Print("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LogRequest(mux)))
}

func serveAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/api.yaml")
}
