package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/ivamshky/go-crud/gen/grpc/user"
	"github.com/ivamshky/go-crud/handler"
	"github.com/ivamshky/go-crud/middlewares"
	"github.com/ivamshky/go-crud/repository"
	"github.com/ivamshky/go-crud/service"
	"google.golang.org/grpc"
)

func main() {
	if db, err := connectDB(); err != nil {
		slog.Error("cannot connect to db. exiting...", err)
	} else {
		defer db.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var wg sync.WaitGroup
		go RESTApiServer(ctx, &wg, db)
		go GrpcServer(ctx, &wg, db)
		slog.Info("Interrupt to close servers")
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		slog.Info(fmt.Sprintf("Interrupt signal %v received. Initiating shutdown", sig))
		cancel()
		wg.Wait()
		slog.Info("All servers closed.")
	}
}

func connectDB() (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"dbadmin", // MySQL username
		"1234",    // MySQL password
		"db",      // Host IP or service name
		"3306",    // MySQL port
		"userdb")  //db name
	return sql.Open("mysql", connString)
}

func serveAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/api.yaml")
}

func RESTApiServer(ctx context.Context, wg *sync.WaitGroup, db *sql.DB) {
	wg.Add(1)
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.html", serveAPISpec)

	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)))

	mux.Handle("GET /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleGetUserById)))
	mux.Handle("POST /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleCreateUser)))
	mux.Handle("GET /user", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleListUsers)))
	mux.Handle("DELETE /user/{userId}", middlewares.RequireAuth(http.HandlerFunc(userHandler.HandleDeleteUser)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: middlewares.LogRequest(mux),
	}

	serveErr := make(chan error, 1)
	go func() {
		slog.Info("REST Server Starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	select {
	case err := <-serveErr:
		slog.Error("REST server experienced an error: %v", err)
	case <-ctx.Done():
		slog.Info("Shutting down REST Server")

		shutDownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := server.Shutdown(shutDownCtx); err != nil {
			slog.Error("Error while graceful shutdown", err)
			if closeErr := server.Close(); closeErr != nil {
				slog.Error("Error while closing server", closeErr)
			}
		}
	}
}

func GrpcServer(ctx context.Context, wg *sync.WaitGroup, db *sql.DB) {
	wg.Add(1)
	defer wg.Done()
	serveErr := make(chan error, 1)

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		slog.Error("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	userService := &handler.UserServer{
		UserService: service.NewUserService(repository.NewUserRepository(db)),
	}
	pb.RegisterUserServiceServer(grpcServer, userService)

	go func() {
		log.Print("GRPC Server listening on port 9001")
		if err := grpcServer.Serve(lis); err != nil {
			serveErr <- fmt.Errorf("failed to serve: %s", err)
		}
	}()

	select {
	case err := <-serveErr:
		slog.Error("Error starting server", err)
	case <-ctx.Done():
		slog.Info("Received Shutdown signal. Shutting down..")
		grpcServer.GracefulStop()
		slog.Info("GRPC server stopped.")
	}
}
