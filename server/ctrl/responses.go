package ctrl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func jsonResponse(rw http.ResponseWriter, data interface{}) {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		errorResponse(rw,
			fmt.Sprintf("Encoding answer, %v", err),
			"Error preparing JSON =(",
			http.StatusInternalServerError)
		return
	}

	toWrite := body.Len()
	rw.WriteHeader(http.StatusOK)
	n, err := io.Copy(rw, body)
	if err != nil {
		log.Printf("Writing response, %v", err)
	} else if n != int64(toWrite) {
		log.Printf("Write wanted %d bytes, wrote %d", toWrite, n)
	}
}

func badRouteResponse(rw http.ResponseWriter, req *http.Request) {
	errorResponse(rw,
		fmt.Sprintf("Requested bad method: %s", req.RequestURI),
		"Method not allowed",
		http.StatusMethodNotAllowed)
}

func dbErrorResponse(rw http.ResponseWriter, err error) {
	msg := fmt.Sprintf("Error querying DB, %v", err)
	errorResponse(rw, msg, msg, http.StatusInternalServerError)
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
