package models

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	wardAll    = "wards:all"
	wardPrefix = "wards:%d"
)

type Bed struct {
	BedID  int `json:"bedId"`
	RoomID int `json:"roomId"`
}

type AdmissionRequest struct {
	AdmRequestID int    `json:"admRequestId"`
	PatientID    int    `json:"wardId"`
	FromWardID   int    `json:"fromWardId"`
	Priority     string `json:"priority"`
	Rationale    string `json:"rationale"`
}

type AdmissionResponse struct {
	AdmResponseID int    `json:"admResponseID"`
	PatientID     int    `json:"wardId"`
	ToWardID      int    `json:"toWardId"`
	InProgress    bool   `json:"inProgress"`
	Refusal       string `json:"refusal"`
}

type InPatient struct {
	PatientID int    `json:"wardId"`
	BedID     int    `json:"bedId"`
	Status    string `json:"status"`
}

type Ward struct {
	WardID             int                 `json:"wardId"`
	Name               string              `json:"name"`
	ChargeNurseID      int                 `json:"chargeNurseId"`
	DoctorID           int                 `json:"doctorId"`
	Beds               []Bed               `json:"beds"`
	Patients           []InPatient         `json:"wards"`
	AdmissionResponses []AdmissionResponse `json:"admissionResponses"`
	AdmissionRequests  []AdmissionRequest  `json:"admissionRequests"`
}

func NewWard(wardID int, name string, chargeNurseID int, doctorID int) *Ward {
	return &Ward{
		WardID:        wardID,
		Name:          name,
		ChargeNurseID: chargeNurseID,
		DoctorID:      doctorID,
	}
}

// Create an ward IFF it doesn't already exist
func (w *Ward) Create() error {
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf(wardPrefix, w.WardID)
	// Create a record
	changes, err := redis.Int(conn.Do("HSETNX", wardAll, w.WardID, key))
	if err != nil {
		return fmt.Errorf("error verifying for existence, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("ward with ID %d already exists, changes %d", w.WardID, changes)
	}

	wardByte, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("marshalling ward, %v", err)
	}

	ok, err := redis.String(conn.Do("SET", key, wardByte))
	if err != nil {
		return fmt.Errorf("error setting ward, %v", err)
	}
	if ok != "OK" {
		return fmt.Errorf("expected 1 ward change, got %d", 1)
	}
	return nil
}

// Update
func (w *Ward) Update() error {
	conn := pool.Get()
	defer conn.Close()
	key := fmt.Sprintf(wardPrefix, w.WardID)

	ward, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("marshalling ward, %v", err)
	}

	changes, err := redis.Int(conn.Do("SET", key, ward))
	if err != nil {
		return fmt.Errorf("error setting ward, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("expected 1 ward change, got %d", 1)
	}
	return nil
}

func (w *Ward) Delete() error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", wardAll, w.WardID)
	if err != nil {
		return fmt.Errorf("deleting ward from record list, %v", err)
	}
	_, err = conn.Do("DEL", fmt.Sprintf(wardPrefix, w.WardID))
	if err != nil {
		return fmt.Errorf("deleting actual ward, %v", err)
	}
	return nil
}

func FindWard(wardID int) (*Ward, bool, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.String(conn.Do("HGET", wardAll, wardID))
	if err != nil {
		if err == redis.ErrNil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("getting ward ID %d, %v", wardID, err)
	}

	if id == "" {
		// doesn't exist
		return nil, false, nil
	}

	wardByte, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil, false, fmt.Errorf("getting bytes for ward %d, %v", wardID, err)
	}

	var ward Ward

	if err := json.Unmarshal(wardByte, &ward); err != nil {
		return nil, false, fmt.Errorf("unmarshalling bytes '%v', %v", string(wardByte), err)
	}
	return &ward, true, nil
}

func FindAllWards() ([]*Ward, error) {
	conn := pool.Get()
	defer conn.Close()

	members, err := redis.Strings(conn.Do("HVALS", wardAll))
	if err != nil {
		return nil, fmt.Errorf("getting vals for %s, %v", wardAll, err)
	}

	if err := conn.Send("MULTI"); err != nil {
		return nil, fmt.Errorf("preparing pipelined GET over all wards, %v", err)
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

	wards := make([]*Ward, 0, len(replies))
	for _, reply := range replies {
		var ward Ward
		err = json.Unmarshal([]byte(reply), &ward)
		if err != nil {
			return wards, err
		}
		wards = append(wards, &ward)
	}
	return wards, nil
}
