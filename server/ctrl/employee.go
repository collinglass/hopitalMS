package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/sessions"
	"net/http"
)

func EmployeeCtrl(store *sessions.CookieStore) http.HandlerFunc {
	routes := map[string]http.Handler{
		"GET":  getEmployee(store),
		"POST": createEmployee(store),
		"PUT":  updateEmployee(store),
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		h, ok := routes[req.Method]
		if !ok {
			badRouteResponse(rw, req)
			return
		}
		h.ServeHTTP(rw, req)
	}
}

func createEmployee(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		emplDetails := struct {
			EmployeeID int                  `json:"employeeId"`
			WardID     int                  `json:"wardId"`
			FirstName  string               `json:"firstName"`
			LastName   string               `json:"lastName"`
			Email      string               `json:"email"`
			Roles      map[models.Role]bool `json:"roles"`
			Password   string               `json:"password"`
		}{}

		err := json.NewDecoder(req.Body).Decode(&emplDetails)
		if err != nil {
			msg := fmt.Sprintf("Couldn't decode JSON employee to create, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		empl, err := models.NewEmployee(
			emplDetails.EmployeeID,
			[]byte(emplDetails.Password),
		)
		if err != nil {
			msg := fmt.Sprintf("Error creating employee, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}
		empl.WardID = emplDetails.WardID
		empl.FirstName = emplDetails.FirstName
		empl.LastName = emplDetails.LastName
		empl.Email = emplDetails.Email
		empl.Roles = emplDetails.Roles

		err = empl.Create()
		if err != nil {
			msg := fmt.Sprintf("Error creating employee, %v", err)
			public := "Couldn't save employee to DB"
			errorResponse(rw, msg, public, http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusCreated)
	}
}

func getEmployee(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		_, ok := extractEmployee(rw, req, store)
		if !ok {
			return // extractEmployee wrote the error message
		}

		// No roles for getting an employee

		id, ok := extractID(req)
		if !ok {
			allEmpl, err := models.FindAllEmployees()
			if err != nil {
				dbErrorResponse(rw, err)
				return
			}
			jsonResponse(rw, allEmpl)
			return
		}
		empl, ok, err := models.FindEmployee(id)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		if !ok {
			http.NotFoundHandler().ServeHTTP(rw, req)
			return
		}
		jsonResponse(rw, empl)
	}
}

func updateEmployee(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		empl, ok := extractEmployee(rw, req, store)
		if !ok {
			return // extractEmployee wrote the error message
		}

		err := json.NewDecoder(req.Body).Decode(empl)
		if err != nil {
			msg := fmt.Sprintf("Bad JSON format, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}

		err = empl.Update()
		if err != nil {
			msg := fmt.Sprintf("Error updating DB, %v", err)
			errorResponse(rw, msg, msg, http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}
