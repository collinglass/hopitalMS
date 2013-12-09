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

	models.FlushAll(spec.Redis, "Yes I am sure I want to flush all my Redis")

	log.Printf("Opening connection with Redis at %s", spec.Redis)
	models.Start(spec.Redis)

	log.Printf("Scanning for stub JSON wards")
	filepath.Walk("wards", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != ".json" {
			return nil
		}

		log.Printf("\tfound '%s'", name)
		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var ward models.Ward
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&ward))
		err = ward.Create()
		if err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}
		log.Printf("\t ✓ created '%s'", ward.Name)
		return nil
	})

	log.Printf("Scanning for stub JSON patients")
	filepath.Walk("patients", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != ".json" {
			return nil
		}

		log.Printf("\tfound '%s'", name)
		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var patient models.Patient
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&patient))
		err = patient.Create()
		if err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}
		log.Printf("\t ✓ created '%s'", patient.FirstName)
		return nil
	})

	log.Printf("Scanning for stub JSON employees")
	filepath.Walk("employees", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != ".json" {
			return nil
		}

		log.Printf("\tfound '%s'", name)
		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		dec := json.NewDecoder(fd)
		emplDetails := struct {
			EmployeeID int                  `json:"employeeId"`
			WardID     int                  `json:"wardId"`
			FirstName  string               `json:"firstName"`
			LastName   string               `json:"lastName"`
			Email      string               `json:"email"`
			Roles      map[models.Role]bool `json:"roles"`
			Password   string               `json:"password"`
		}{}
		if err := dec.Decode(&emplDetails); err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}

		employee, err := models.NewEmployee(
			emplDetails.EmployeeID,
			[]byte(emplDetails.Password),
		)
		if err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}

		employee.WardID = emplDetails.WardID
		employee.FirstName = emplDetails.FirstName
		employee.LastName = emplDetails.LastName
		employee.Email = emplDetails.Email
		employee.Roles = emplDetails.Roles

		if err := employee.Create(); err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}
		log.Printf("\t ✓ created '%s'", employee.FirstName)
		return nil
	})

	log.Printf("Scanning for stub JSON charge nurses")
	filepath.Walk("chargeNurses", func(name string, fi os.FileInfo, err error) error {
		if filepath.Ext(name) != ".json" {
			return nil
		}

		log.Printf("\tfound '%s'", name)
		fd, err := os.Open(name)
		failFast(err)
		defer fd.Close()

		var nurse models.ChargeNurse
		dec := json.NewDecoder(fd)
		failFast(dec.Decode(&nurse))

		empl, ok, _ := models.FindEmployee(nurse.EmployeeID)
		if !ok {
			log.Printf("\t ✗ error, charge nurse with ID %d is not a known employee", nurse.EmployeeID)
			return nil
		}

		err = nurse.Create()
		if err != nil {
			log.Printf("\t ✗ error, %v", err)
			return nil
		}
		log.Printf("\t ✓ linked '%s' to charge nurse", empl.FirstName)
		return nil
	})
}
