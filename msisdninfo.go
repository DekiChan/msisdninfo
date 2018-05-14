package main

import (
	"fmt"
	"net/http"

	"github.com/dekichan/msisdninfo/controllers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("msisdninfo started")

	router := mux.NewRouter()

	setupRoutes(router)
	startServer(router, 8080)
}

// Initializes routes
func setupRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.HomeHandler).Methods("GET")
	router.HandleFunc("/transform", controllers.TransformHandler).Methods("GET")
}

func startServer(router *mux.Router, port int) {
	portStr := fmt.Sprintf(":%d", port)

	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", port))
	http.ListenAndServe(portStr, router)
}
