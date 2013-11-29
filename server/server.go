package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func init() {
	mime.AddExtensionType(".json", "application/json; charset=utf-8")
}

func main() {
	log.Println("Starting Server")

	stubFileserver := logHandler(restStubHandler("./stub/"))

	http.Handle("/api/", stubFileserver)
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

func restStubHandler(prefix string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		verb := req.Method

		dir, file := path.Split(req.URL.String())
		fmt.Printf("Method %v Path: %v File: %v\n", verb, dir, file)

		prefixDir := filepath.Join(filepath.Clean(prefix), filepath.Clean(dir))
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
}

func restMemberStub(verb, dir, file string, rw http.ResponseWriter, req *http.Request) {
	filename := filepath.Join(dir, file)
	log.Printf("Trying to open file %s", filename)
	f, err := os.Open(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error, %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil && os.IsNotExist(err) {
		log.Printf("Not found, %v", err)
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

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Writing response")
	setContentTypeFromReq(rw, req)
	n, err := rw.Write(bytes)
	if err != nil {
		log.Printf("Error during response write, %v", err)
	} else if n != len(bytes) {
		log.Printf("Short write, %d/%d", n, len(bytes))
	}
}

func tryFindFile(dirname, name string) (*os.File, error) {
	log.Printf("Trying to find file matching %s in %s", name, dirname)
	finfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, fi := range finfos {
		if fi.IsDir() {
			continue
		}
		if name == stripExtension(fi.Name()) {
			log.Printf("Found %s", fi.Name())
			return os.Open(filepath.Join(dirname, fi.Name()))
		}
	}
	return nil, fmt.Errorf("no match found in dir %s for name %s", dirname, name)
}

func restCollectionStub(verb, dir string, rw http.ResponseWriter, req *http.Request) {

	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	var names []string
	for _, fi := range finfos {
		if fi.IsDir() {
			continue
		}
		filename := fi.Name()
		ext := path.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		names = append(names, name)
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

func stripExtension(filename string) string {
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

func setContentTypeFromReq(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", mime.TypeByExtension("json"))
}
