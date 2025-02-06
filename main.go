package main

import (
	"log"
	"net/http"
	"recipt-processor/configuration"
	"recipt-processor/controller"
	"recipt-processor/repository"

	"github.com/gorilla/mux"
)

func main() {
	db, err := configuration.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	repository.SetDB(db)

	r := mux.NewRouter()
	controller.RegisterStudentRoutes(r)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
