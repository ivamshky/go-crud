package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ivamshky/go-crud/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.Handle("/user", middlewares.RequireAuth(http.HandlerFunc(handleUser)))
	log.Print("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", middlewares.LogRequest(mux)))
}

func handleUser(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		fmt.Fprint(writer, "GET /user called")
		break
	case "POST":
		fmt.Fprint(writer, "POST /user called")
		break
	default:
		panic("Unsupported method")
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}
