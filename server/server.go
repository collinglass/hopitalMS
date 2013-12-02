package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	//"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	//"time"
)

func init() {
	mime.AddExtensionType(".json", "application/json; charset=utf-8")
}

/*
type CookieValue struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
}

func AuthenticationHandler(w http.ResponseWriter, req *http.Request) {
	var cookieValue CookieValue

	cookieValue.Username = ""
	cookieValue.Role = 1

	// var authCookie Cookie

	/*
		// If request already contains a user
		// deserialize json then check if that mother fucker contains a user
		err := json.NewDecoder(req.Body).Decode(&authCookie)
		if err != nil {
			panic(err)
		}
		// if it does keep passing that shit around
		if authCookie.name == "user" {
			cookieValue.Username = authCookie.
		}
*/

// Serialize the response to JSON
/*
		j, err := json.Marshal(cookieValue)
		if err != nil {
			panic(err)
		}
	log.Println("Hello")
	var res []byte
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{"user", `{"username": "", "role": 1}`, "/", "www.dummy.com", expire, expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}}

	http.SetCookie(w, &cookie)
	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
*/
func main() {
	log.Println("Starting Server")

	// sessionMngr := session.NewSessionManager(sessionExpire)
	// http.Handle("/api/v0.1/sessions", sessionStateHandler(sessionMngr))

	// stubFileserver := session.EnsureHasSession(sessionMngr,
	// 	tokenName,
	// 	)
	//r := mux.NewRouter()

	//r.HandleFunc("/#/login", AuthenticationHandler).Methods("GET")
	http.Handle("/api/", newRESTStubHandler("./stub/"))
	http.Handle("/", logHandler(http.FileServer(http.Dir("../app/"))))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}

func jsonHandler(h http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(rw, req)
	}
}

func logHandler(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

type restStubHandler struct {
	prefix    string
	currentID string
}

func newRESTStubHandler(prefix string) restStubHandler {
	return restStubHandler{prefix: prefix}
}

func (r restStubHandler) SetID(id string) {
	r.currentID = id
}

func (r restStubHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	verb := req.Method

	dir, file := path.Split(req.URL.String())
	log.Printf("Method %v Collection: %v Member: %v\n", verb, dir, file)

	prefixDir := filepath.Join(filepath.Clean(r.prefix), filepath.Clean(dir))
	if file != "" {
		// it's a member request
		log.Printf("Serving member")
		restMemberStub(verb, prefixDir, file, rw, req)
	} else {
		// it's a collection request
		log.Printf("Serving collection")
		restCollectionStub(verb, prefixDir, rw, req)
	}
}

func restMemberStub(verb, dir, file string, rw http.ResponseWriter, req *http.Request) {
	filename := filepath.Join(dir, file)
	f, err := os.Open(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error, %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil && os.IsNotExist(err) {
		f, err = tryFindFile(dir, file)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("No match, %v", err)
			http.NotFound(rw, req)
			return
		} else if err != nil {
			log.Printf("Finding match went wrong, %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Will serve content of %s", f.Name())
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if isDir(f) {
		restCollectionStub(verb, filepath.Join(dir, file), rw, req)
		return
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	setContentTypeFromReq(rw, req)
	n, err := rw.Write(bytes)
	if err != nil {
		log.Printf("Error during response write, %v", err)
	} else if n != len(bytes) {
		log.Printf("Short write, %d/%d", n, len(bytes))
	}
}

func tryFindFile(dirname, name string) (*os.File, error) {
	finfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("reading directory at %s, %v", dirname, err)
	}
	for _, fi := range finfos {
		if name == stripExtension(fi.Name()) {
			return os.Open(filepath.Join(dirname, fi.Name()))
		}
	}
	return nil, fmt.Errorf("no match found in dir %s for name %s", dirname, name)
}

type memberId struct {
	ID string `json:"id"`
}

func restCollectionStub(verb, dir string, rw http.ResponseWriter, req *http.Request) {

	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	var names []memberId
	for _, fi := range finfos {
		if fi.IsDir() {
			continue
		}
		filename := fi.Name()
		ext := path.Ext(filename)
		name := filename[:len(filename)-len(ext)]

		names = append(names, memberId{ID: name})
	}
	respondWithJSON(rw, req, names)
}

func respondWithJSON(rw http.ResponseWriter, req *http.Request, data interface{}) {
	resp := bytes.NewBuffer(nil)
	err := json.NewEncoder(resp).Encode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	setContentTypeFromReq(rw, req)
	l := resp.Len()
	n, err := io.Copy(rw, resp)
	if err != nil {
		log.Printf("Error during response write, %v", err)
	} else if n != int64(l) {
		log.Printf("Short write, %d/%d", n, l)
	}
}

func isDir(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return stat.Mode().IsDir()
}

func stripExtension(filename string) string {
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

func setContentTypeFromReq(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", mime.TypeByExtension("json"))
}
