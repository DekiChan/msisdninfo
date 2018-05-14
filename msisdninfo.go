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
	startServer(router, ":8080")
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.HomeHandler).Methods("GET")
	router.HandleFunc("/transform", controllers.TransformHandler).Methods("GET")
}

func startServer(router *mux.Router, host string) {
	fmt.Println(fmt.Sprintf("starting to listen at localhost%s", host))
	http.ListenAndServe(host, router)
}
