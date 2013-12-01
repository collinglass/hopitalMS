package main

import (
	"bytes"
	"encoding/json"
	"github.com/aybabtme/session"
	"github.com/collinglass/moustacheMS/server/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	sessionExpire = time.Second * 15
	tokenName     = "token"
)

func sessionStateHandler(mngr *session.SessionManager) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "DELETE":
			for _, token := range req.URL.Query()[tokenName] {
				mngr.DeleteSession(token)
			}
			rw.WriteHeader(http.StatusNoContent)
		case "POST":
			createUserSession(mngr, rw, req)
		default:
			http.Error(rw, "Not allowed:"+req.Method, http.StatusMethodNotAllowed)
		}
	}
}

func createUserSession(mngr *session.SessionManager, rw http.ResponseWriter, req *http.Request) {
	emplID, password, err := extractUserCredentials(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	empl, ok, err := models.FindEmployee(emplID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}
	if err := empl.ValidatePassword(password); err != nil {
		log.Printf("Error validating password of %v, %v", empl, err)
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}
	token, err := mngr.NewSession(strconv.Itoa(empl.EmployeeID), req.RemoteAddr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(struct {
		tokenName string `json:"paramName"`
		token     string `json:"token"`
	}{tokenName, token})

	if err != nil {
		mngr.DeleteSession(token)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	n, err := io.Copy(rw, buf)

	// If copy wasn't properly done, discard token
	if err != nil {
		mngr.DeleteSession(token)
		log.Printf("Error writing answer, %v", err)
	} else if n != int64(buf.Len()) {
		mngr.DeleteSession(token)
		log.Printf("Error writing answer, want %d bytes wrote %d",
			buf.Len(), n)
	}
}

func extractUserCredentials(req *http.Request) (int, []byte, error) {
	cred := struct {
		ID       int    `json:"employeeId"`
		Password string `json:"password"`
	}{}
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&cred)
	return cred.ID, []byte(cred.Password), err
}
