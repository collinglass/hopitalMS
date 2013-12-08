package main

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"log"
	"mime"
	"net/http"
	"os"
)

func init() {
	mime.AddExtensionType(".json", "application/json; charset=utf-8")
}

func main() {

	log.Println("Parsing specification")
	spec := ParseSpec()
	log.Println("Preparing Redis pool")
	models.Start(spec.Redis)
	log.Println("Starting secret sessions store")
	store, sessionHandler := StartSessions([]byte(spec.AuthKey), []byte(spec.CryptKey))

	log.Println("Registering handlers")
	http.Handle("/api/sessions", logHandler(sessionHandler))
	http.Handle("/api/", logHandler(apiRouter(store)))
	http.Handle("/", logHandler(http.FileServer(http.Dir("../app/"))))

	log.Println("Starting Server")
	log.Printf("Listening on %s", spec.ListenOn)
	panic(http.ListenAndServe(spec.ListenOn, nil))
}

func logHandler(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
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

func extractEmployee(
	rw http.ResponseWriter,
	req *http.Request,
	store *sessions.CookieStore,
) (*models.Employee, bool) {
	curSession, err := store.Get(req, sessionCookieName)
	if err != nil {
		errorResponse(rw,
			"No session in cookie",
			"Session invalid",
			http.StatusBadRequest)
		return nil, false
	}
	emplRaw, ok := curSession.Values[emplIDCookieKey]
	if !ok {
		errorResponse(rw,
			"No employee in cookie",
			"Forbidden",
			http.StatusForbidden)
		return nil, false
	}
	emplID, ok := emplRaw.(int)
	if !ok {
		errorResponse(rw,
			fmt.Sprintf("Couldn't parse int from emplRaw '%#v', %v",
				emplRaw, err),
			"Bad request, invalid session",
			http.StatusBadRequest)
		return nil, false
	}

	empl, ok, err := models.FindEmployee(emplID)

	if err != nil {
		errorResponse(rw,
			fmt.Sprintf("Database find of emplID %d, %v", emplID, err),
			"Database error",
			http.StatusInternalServerError)
		return nil, false
	}

	if !ok {
		errorResponse(rw,
			fmt.Sprintf("Couldn't find employeeID %d", emplID),
			"Invalid employeeID",
			http.StatusBadRequest)
		return nil, false
	}
	return empl, true
}

func errorResponse(rw http.ResponseWriter, cause, public string, code int) {
	log.Printf("%d: %s, answered '%s'", code, cause, public)
	jsonErr := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		code,
		public,
	}
	jsonBytes, err := json.Marshal(&jsonErr)
	if err != nil {
		log.Printf("error serving error message, %v", err)
		http.Error(rw, "Error serving error message", http.StatusInternalServerError)
		return
	}
	http.Error(rw, string(jsonBytes), code)
}
