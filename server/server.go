package main

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/ctrl"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
	apiRoutes := mainRouter.PathPrefix("/api/").Subrouter()

	// Sessions are handled in this package
	apiRoutes.Handle("/sessions", sessionHandler).Methods("POST", "DELETE")

	// Controllers
	apiRoutes.Handle("/employees/{id:[0-9]+}",
		ctrl.EmployeeCtrl).Methods("GET")
	apiRoutes.Handle("/employees/",
		ctrl.EmployeeCtrl).Methods("GET", "POST")

	apiRoutes.Handle("/patients/{id:[0-9]+}",
		ctrl.PatientCtrl).Methods("GET", "PUT", "DELETE")
	apiRoutes.Handle("/patients/",
		ctrl.PatientCtrl).Methods("GET", "POST")

	apiRoutes.Handle("/wards/{id:[0-9]+}",
		ctrl.wardCtrl).Methods("GET", "PUT")
	apiRoutes.Handle("/wards/",
		ctrl.wardCtrl).Methods("GET")

	// Static files
	mainRouter.Handle("/", http.FileServer(http.Dir("../app/")))

	// Pass back to stdlib http
	http.Handle("/", handlers.LoggingHandler(os.Stdout, mainRouter))

	log.Println("Starting Server")
	log.Printf("Listening on %s", spec.ListenOn)
	panic(http.ListenAndServe(spec.ListenOn, nil))
}

func apiRouter(store *sessions.CookieStore) http.HandlerFunc {

	return func(rw http.ResponseWriter, req *http.Request) {

		empl, ok := extractEmployee(rw, req, store)
		if !ok {
			// extractEmployee dealt with error messages
			return
		}
	}
}
