package models

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	patientAll    = "patients:all"
	patientPrefix = "patients:%d"
)

type NOK struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Address      string `json:"address"`
	PhoneNum     string `json:"phoneNum"`
}

type Patient struct {
	PatientID     int    `json:"patientId"`
	LastName      string `json:"lastName"`
	FirstName     string `json:"firstName"`
	HealthInsNum  string `json:"healthInsNum"`
	Address       string `json:"address"`
	PhoneNum      string `json:"phoneNum"`
	DateOfBirth   string `json:"dateOfBirth"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"maritalStatus"`
	NextOfKin     NOK    `json:"nextOfKin"`
}

func NewPatient(patientID int, lastName, firstName string) *Patient {
	return &Patient{
		PatientID: patientID,
		LastName:  lastName,
		FirstName: firstName,
	}
}

// Create an patient IFF it doesn't already exist
func (p *Patient) Create() error {
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf(patientPrefix, p.PatientID)
	// Create a record
	changes, err := redis.Int(conn.Do("HSETNX", patientAll, p.PatientID, key))
	if err != nil {
		return fmt.Errorf("error verifying for existence, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("patient with ID %d already exists, changes %d", p.PatientID, changes)
	}

	patientByte, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshalling patient, %v", err)
	}

	ok, err := redis.String(conn.Do("SET", key, patientByte))
	if err != nil {
		return fmt.Errorf("error setting patient, %v", err)
	}
	if ok != "OK" {
		return fmt.Errorf("expected 1 patient change, got %d", 1)
	}
	return nil
}

// Update
func (p *Patient) Update() error {
	conn := pool.Get()
	defer conn.Close()
	key := fmt.Sprintf(patientPrefix, p.PatientID)

	patient, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshalling patient, %v", err)
	}

	changes, err := redis.Int(conn.Do("SET", key, patient))
	if err != nil {
		return fmt.Errorf("error setting patient, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("expected 1 patient change, got %d", 1)
	}
	return nil
}

func (p *Patient) Delete() error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", patientAll, p.PatientID)
	if err != nil {
		return fmt.Errorf("deleting patient from record list, %v", err)
	}
	_, err = conn.Do("DEL", fmt.Sprintf(patientPrefix, p.PatientID))
	if err != nil {
		return fmt.Errorf("deleting actual patient, %v", err)
	}
	return nil
}

func FindPatient(patientID int) (*Patient, bool, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.String(conn.Do("HGET", patientAll, patientID))
	if err != nil {
		return nil, false, fmt.Errorf("getting patient ID %d, %v", patientID, err)
	}

	if id == "" {
		// doesn't exist
		return nil, false, nil
	}

	patientByte, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil, false, fmt.Errorf("getting bytes for patient %d, %v", patientID, err)
	}

	var patient Patient

	if err := json.Unmarshal(patientByte, &patient); err != nil {
		return nil, false, fmt.Errorf("unmarshalling bytes '%v', %v", string(patientByte), err)
	}
	return &patient, true, nil
}

func FindAllPatients() ([]*Patient, error) {
	conn := pool.Get()
	defer conn.Close()

	members, err := redis.Strings(conn.Do("HVALS", patientAll))
	if err != nil {
		return nil, fmt.Errorf("getting vals for %s, %v", patientAll, err)
	}

	if err := conn.Send("MULTI"); err != nil {
		return nil, fmt.Errorf("preparing pipelined GET over all patients, %v", err)
	}

	for _, id := range members {
		err := conn.Send("GET", id)
		if err != nil {
			return nil, fmt.Errorf("piping GET for %s, %v", id, err)
		}
	}

	replies, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		return nil, fmt.Errorf("executing pipeline, %v", err)
	}

	patients := make([]*Patient, 0, len(replies))
	for _, reply := range replies {
		var patient Patient
		err = json.Unmarshal([]byte(reply), &patient)
		if err != nil {
			return patients, err
		}
		patients = append(patients, &patient)
	}
	return patients, nil
}
