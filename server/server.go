package main

import (
	"github.com/collinglass/moustacheMS/server/ctrl"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Parsing specification")
	spec := ParseSpec()

	log.Println("Preparing Redis pool")
	models.Start(spec.Redis)

	log.Println("Starting secret sessions store")
	store, sessionHandler := StartSessions([]byte(spec.AuthKey), []byte(spec.CryptKey))

	log.Println("Registering handlers")

	mainRouter := mux.NewRouter()

	apiRoutes := mainRouter.PathPrefix("/api/v0.1").Subrouter()

	// Sessions are handled in this package
	apiRoutes.Handle("/sessions", sessionHandler)

	// Controllers
	apiRoutes.Handle("/employees/{id:[0-9]+}", ctrl.EmployeeCtrl(store))
	apiRoutes.Handle("/employees", ctrl.EmployeeCtrl(store))

	apiRoutes.Handle("/patients/{id:[0-9]+}", ctrl.PatientCtrl(store))
	apiRoutes.Handle("/patients", ctrl.PatientCtrl(store))

	apiRoutes.Handle("/wards/{id:[0-9]+}", ctrl.WardCtrl(store))
	apiRoutes.Handle("/wards", ctrl.WardCtrl(store))

	// Heartbeat
	mainRouter.PathPrefix("/ping").Handler(pingHandler())
	// Static files
	mainRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("../app/")))

	// Pass back to stdlib http
	http.Handle("/", handlers.LoggingHandler(os.Stdout, mainRouter))

	log.Println("Starting Server")
	log.Printf("Listening on %s", spec.ListenOn)
	panic(http.ListenAndServe(spec.ListenOn, nil))
}

func pingHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, err := io.WriteString(rw, "pong")
		if err != nil {
			log.Fatalf("Pong failed, %v", err)
		}
	}
}
