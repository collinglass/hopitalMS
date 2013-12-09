package main

import (
	"encoding/json"
	"github.com/collinglass/moustacheMS/server/config"
	"github.com/collinglass/moustacheMS/server/models"
	"log"
	"os"
	"path/filepath"
)

func main() {

	failFast := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	log.Println("Parsing spec")
	spec := config.ParseSpec()

	log.Printf("Opening connection with Redis at %s", spec.Redis)
	models.Start(spec.Redis)

	log.Printf("Scanning for stub JSON wards")
	filepath.Walk("wards", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != "json" {
			return nil
		}

		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var ward models.Ward
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&ward))
		failFast(ward.Create())

		return nil
	})

	log.Printf("Scanning for stub JSON patients")
	filepath.Walk("patients", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != "json" {
			return nil
		}

		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var patient models.Patient
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&patient))
		failFast(patient.Create())

		return nil
	})

	log.Printf("Scanning for stub JSON employees")
	filepath.Walk("employees", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != "json" {
			return nil
		}

		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var employee models.Employee
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&employee))
		failFast(employee.Create())

		return nil
	})
}
