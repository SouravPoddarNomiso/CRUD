package router

import (
	"demo/controller"
	"fmt"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controller.Home).Methods("GET")
	router.HandleFunc("/user", controller.GetUsers).Methods("GET")
	router.HandleFunc("/user", controller.InsertUser).Methods("POST")
	router.HandleFunc("/user/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user", controller.UpdateUser).Methods("PUT")
	fmt.Println("Routes are Loaded.")
	return router
}
