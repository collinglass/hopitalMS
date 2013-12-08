package main

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session"
	emplIDCookieKey   = "emplID"
)

func StartSessions(authKey, cryptKey []byte) (*sessions.CookieStore, http.Handler) {
	store := sessions.NewCookieStore(authKey, cryptKey)
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int(time.Hour * 24),
		// Secure: true,
	}

	sessionHandler := handlers.MethodHandler{
		"POST":   postSession(store),
		"DELETE": deleteSession(store),
	}
	return store, sessionHandler
}

type userCreds struct {
	EmployeeID int    `json:"employeeId"`
	Password   string `json:"password"`
}

func postSession(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var cred userCreds
		dec := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := dec.Decode(&cred)
		if err != nil {
			msg := fmt.Sprintf("Malformed credentials, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}

		empl, ok, err := models.FindEmployee(cred.EmployeeID)
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Finding employeeID %d, %v", cred.EmployeeID, err),
				"Oups, error",
				http.StatusInternalServerError)
			return
		}
		if !ok {
			errorResponse(rw,
				fmt.Sprintf("Unknown user, %#v", cred),
				"Invalid username/password",
				http.StatusForbidden)
			return
		}
		err = empl.ValidatePassword([]byte(cred.Password))
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Validate password failed for user %#v, %v", cred, err),
				"Invalid username/password",
				http.StatusForbidden)
			return
		}

		session, err := store.New(req, sessionCookieName)
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Creating cookie session for user %d, %v",
					empl.EmployeeID, err),
				"Oups, error",
				http.StatusInternalServerError)
			return
		}
		session.Values[emplIDCookieKey] = empl.EmployeeID

		err = session.Save(req, rw)
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Saving cookie session for user %d, %v",
					empl.EmployeeID, err),
				"Oups, error",
				http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	}
}

func deleteSession(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, err := req.Cookie(sessionCookieName)
		if err == http.ErrNoCookie {
			errorResponse(rw,
				fmt.Sprintf("Couldn't get session from cookie, %v", err),
				"Invalid session",
				http.StatusUnauthorized)
			return
		}

		session, err := store.Get(req, sessionCookieName)
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Error getting/saving session, %v", err),
				"Invalid session",
				http.StatusBadRequest)
			return
		}

		// old version of gorilla/sessions have a bug with Options, if
		// you have a problem, try updating gorilla/sessions
		session.Options.MaxAge = -1
		delete(session.Values, emplIDCookieKey)

		err = session.Save(req, rw)
		if err != nil {
			errorResponse(rw,
				fmt.Sprintf("Saving deleted session, %v", err),
				"Oups, error",
				http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
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
