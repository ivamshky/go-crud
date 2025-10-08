package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ivamshky/go-crud/handler"
	"github.com/ivamshky/go-crud/middlewares"
	"github.com/ivamshky/go-crud/repository"
	"github.com/ivamshky/go-crud/service"
	"google.golang.org/grpc"
)

func main() {

	RESTApiServer(8080)
}

func serveAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/api.yaml")
}

func RESTApiServer(port int) {
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

	log.Printf("Listening on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), middlewares.LogRequest(mux)))
}

func GrpcServer() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	grpcServer := grpc.NewServer()

	// register our server struct as a handle for the CoffeeShopService rpc calls that come in through grpcServer
	//pb.RegisterCoffeeShopServer(grpcServer, &server{})

	// Serve traffic
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
