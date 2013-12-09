package ctrl

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

func PatientCtrl(store *sessions.CookieStore) http.HandlerFunc {
	routes := map[string]http.Handler{
		"GET":    getPatient(store),
		"POST":   savePatient(store),
		"DELETE": deletePatient(store),
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		h, ok := routes[req.Method]
		if !ok {
			badRouteResponse(rw, req)
			return
		}
		empl, ok := extractEmployee(rw, req, store)
		if !ok {
			return // extractEmployee wrote the error message
		}
		if !empl.Roles[models.MedicalStaff] {
			unauthorizedResponse(rw, "access this location")
		}

		h.ServeHTTP(rw, req)
	}
}

func getPatient(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id, ok := extractID(req)
		if !ok {
			allPatient, err := models.FindAllPatients()
			if err != nil {
				dbErrorResponse(rw, err)
				return
			}
			jsonResponse(rw, allPatient)
			return
		}
		patient, ok, err := models.FindPatient(id)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		if !ok {
			msg := fmt.Sprintf("No patient with id %d", id)
			errorResponse(rw, msg, msg, http.StatusNotFound)
			return
		}
		jsonResponse(rw, patient)
	}
}

func savePatient(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id, ok := extractID(req)
		if !ok {
			msg := "can't update a collection"
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}

		patient, ok, err := models.FindPatient(id)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}

		err = json.NewDecoder(req.Body).Decode(patient)
		if err != nil {
			msg := fmt.Sprintf("error unmarshalling JSON, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		if !ok {
			id := createPatientID(patient)
			log.Printf("Creating new patient with ID %d", id)
			patient.PatientID = id
			err = patient.Create()
		} else {
			err = patient.Update()
		}
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func deletePatient(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id, ok := extractID(req)
		if !ok {
			msg := "can't delete a collection"
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}

		patient, ok, err := models.FindPatient(id)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		if !ok {
			msg := fmt.Sprintf("No patient with ID %d was found", id)
			errorResponse(rw, msg, msg, http.StatusNotFound)
			return
		}
		err = patient.Delete()
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func createPatientID(patient *models.Patient) int {
	h := sha1.New()
	fmt.Fprintf(h, "%s%s%s%s%v",
		patient.FirstName,
		patient.LastName,
		patient.HealthInsNum,
		patient.DateOfBirth,
		time.Now())
	idHash := h.Sum(make([]byte, 8))
	id, size := binary.Uvarint(idHash)
	if size <= 0 {
		log.Fatalf("Hash too small to produce uint64: %s", idHash)
	}
	return int(id)
}
