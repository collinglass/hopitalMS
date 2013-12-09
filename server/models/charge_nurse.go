package models

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	nurseAll    = "chargeNurses:all"
	nursePrefix = "chargeNurses:%d"
)

type ChargeNurse struct {
	EmployeeID int    `json:"employeeId"`
	PhoneExt   string `json:"phoneExt"`
	BipperExt  string `json:"bipperExt"`
}

func NewChargeNurse(nurseID int, phoneExt, bipperExt string) *ChargeNurse {
	return &ChargeNurse{
		EmployeeID: nurseID,
		PhoneExt:   phoneExt,
		BipperExt:  bipperExt,
	}
}

// Create an nurse IFF it doesn't already exist
func (cn *ChargeNurse) Create() error {
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf(nursePrefix, cn.EmployeeID)
	// Create a record
	changes, err := redis.Int(conn.Do("HSETNX", nurseAll, cn.EmployeeID, key))
	if err != nil {
		return fmt.Errorf("error verifying for existence, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("nurse with ID %d already exists, changes %d", cn.EmployeeID, changes)
	}

	nurseByte, err := json.Marshal(cn)
	if err != nil {
		return fmt.Errorf("marshalling nurse, %v", err)
	}

	ok, err := redis.String(conn.Do("SET", key, nurseByte))
	if err != nil {
		return fmt.Errorf("error setting nurse, %v", err)
	}
	if ok != "OK" {
		return fmt.Errorf("expected 1 nurse change, got %d", 1)
	}
	return nil
}

// Update
func (cn *ChargeNurse) Update() error {
	conn := pool.Get()
	defer conn.Close()
	key := fmt.Sprintf(nursePrefix, cn.EmployeeID)

	nurse, err := json.Marshal(cn)
	if err != nil {
		return fmt.Errorf("marshalling nurse, %v", err)
	}

	changes, err := redis.String(conn.Do("SET", key, nurse))
	if err != nil {
		return fmt.Errorf("error setting nurse, %v", err)
	}
	if changes != "OK" {
		return fmt.Errorf("expected 1 nurse change, got %d", 1)
	}
	return nil
}

func (cn *ChargeNurse) Delete() error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", nurseAll, cn.EmployeeID)
	if err != nil {
		return fmt.Errorf("deleting nurse from record list, %v", err)
	}
	_, err = conn.Do("DEL", fmt.Sprintf(nursePrefix, cn.EmployeeID))
	if err != nil {
		return fmt.Errorf("deleting actual nurse, %v", err)
	}
	return nil
}

func FindChargeNurse(nurseID int) (*ChargeNurse, bool, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.String(conn.Do("HGET", nurseAll, nurseID))
	if err != nil {
		if err == redis.ErrNil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("getting nurse ID %d, %v", nurseID, err)
	}

	if id == "" {
		// doesn't exist
		return nil, false, nil
	}

	nurseByte, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil, false, fmt.Errorf("getting bytes for nurse %d, %v", nurseID, err)
	}

	var nurse ChargeNurse

	if err := json.Unmarshal(nurseByte, &nurse); err != nil {
		return nil, false, fmt.Errorf("unmarshalling bytes '%v', %v", string(nurseByte), err)
	}
	return &nurse, true, nil
}

func FindAllChargeNurses() ([]*ChargeNurse, error) {
	conn := pool.Get()
	defer conn.Close()

	members, err := redis.Strings(conn.Do("HVALS", nurseAll))
	if err != nil {
		return nil, fmt.Errorf("getting vals for %s, %v", nurseAll, err)
	}

	if err := conn.Send("MULTI"); err != nil {
		return nil, fmt.Errorf("preparing pipelined GET over all nurses, %v", err)
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

	nurses := make([]*ChargeNurse, 0, len(replies))
	for _, reply := range replies {
		var nurse ChargeNurse
		err = json.Unmarshal([]byte(reply), &nurse)
		if err != nil {
			return nurses, err
		}
		nurses = append(nurses, &nurse)
	}
	return nurses, nil
}
