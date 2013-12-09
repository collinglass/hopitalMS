package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/moustacheMS/server/models"
	"github.com/gorilla/sessions"
	"net/http"
	"reflect"
)

func WardCtrl(store *sessions.CookieStore) http.HandlerFunc {
	routes := map[string]http.Handler{
		"GET":  getWard(store),
		"POST": updateWard(store),
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

func getWard(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		empl, ok := extractEmployee(rw, req, store)
		if !ok {
			return // extractEmployee wrote the error message
		}

		isNurse := false
		if empl.Roles[models.ChargeNurseRole] {
			// is it true?
			_, ok, err := models.FindChargeNurse(empl.EmployeeID)
			if err != nil {
				dbErrorResponse(rw, err)
				return
			}
			if !ok {
				// lied about roles
				unauthorizedResponse(rw, "access this location")
				return
			}
			isNurse = true
		}

		id, ok := extractID(req)
		if !ok {
			allWards, err := models.FindAllWards()
			if err != nil {
				dbErrorResponse(rw, err)
				return
			}

			if !isNurse {
				// Strip out details
				for i := 0; i < len(allWards); i++ {
					allWards[i].AdmissionRequests = nil
					allWards[i].AdmissionResponses = nil
				}
			}

			jsonResponse(rw, allWards)
			return
		}
		ward, ok, err := models.FindWard(id)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		if !ok {
			http.NotFoundHandler().ServeHTTP(rw, req)
			return
		}
		if !isNurse {
			ward.AdmissionRequests = nil
			ward.AdmissionResponses = nil
		}

		jsonResponse(rw, ward)
	}
}

func updateWard(store *sessions.CookieStore) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		empl, ok := extractEmployee(rw, req, store)
		if !ok {
			return // extractEmployee wrote the error message
		}

		if !empl.Roles[models.MedicalStaff] {
			unauthorizedResponse(rw, "access this location as a non-medical staff")
			return
		}

		var wantWard models.Ward
		err := json.NewDecoder(req.Body).Decode(&wantWard)
		if err != nil {
			msg := fmt.Sprintf("Bad JSON format, %v", err)
			errorResponse(rw, msg, msg, http.StatusBadRequest)
			return
		}

		// Can only possibly change a ward you belong to
		if wantWard.WardID != empl.WardID {
			format := "update wards, employeeID %d doesn't belong to wardID %d"
			msg := fmt.Sprintf(format, empl.EmployeeID, wantWard.WardID)
			unauthorizedResponse(rw, msg)
			return
		}

		currentWard, ok, err := models.FindWard(wantWard.WardID)
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}
		if !ok {
			msg := fmt.Sprintf("WardID %d is not a Ward", wantWard.WardID)
			errorResponse(rw, msg, msg, http.StatusNotFound)
			return
		}

		if isChangeAdmReqList(&wantWard, currentWard) ||
			isChangeAdmRespList(&wantWard, currentWard) {
			if !empl.Roles[models.ChargeNurseRole] {
				unauthorizedResponse(rw, "modify response lists")
				return
			}
		}

		err = wantWard.Update() // Overwrite currentWard
		if err != nil {
			dbErrorResponse(rw, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func isChangePatientList(want, got *models.Ward) bool {
	return !reflect.DeepEqual(want.Patients, got.Patients)
}

func isChangeBedList(want, got *models.Ward) bool {
	return !reflect.DeepEqual(want.Beds, got.Beds)
}

func isChangeAdmReqList(want, got *models.Ward) bool {
	return !reflect.DeepEqual(want.AdmissionRequests, got.AdmissionRequests)
}
func isChangeAdmRespList(want, got *models.Ward) bool {
	return !reflect.DeepEqual(want.AdmissionResponses, got.AdmissionResponses)
}
