package serve

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Serve() {
	router := mux.NewRouter()

	setupRoutes(router)
	startServer(router, 8080)
}

// Initializes routes
func setupRoutes(router *mux.Router) {
	router.HandleFunc("/transform", TransformHandler).Methods("GET")
}

func startServer(router *mux.Router, port int) {
	portStr := fmt.Sprintf(":%d", port)

	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", port))
	http.ListenAndServe(portStr, router)
}
