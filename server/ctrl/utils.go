package ctrl

import (
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
)

const (
	// Must match ../sessions.go, but can't import (cycle)
	sessionCookieName = "session"
	emplIDCookieKey   = "emplID"
)

func extractID(req *http.Request) (int, bool) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		return -1, false
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("Couldn't parse int from id string, %v", err)
	}
	return idInt, true
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
