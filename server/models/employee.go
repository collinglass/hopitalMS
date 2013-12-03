package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	employeeAll    = "employees:all"
	employeePrefix = "employees:%d"
)

// Employee holds the data of a PMS employee.  The passwordHash is not
// exposed to JSON decoder/encoders
type Employee struct {
	EmployeeID   int           `json:"employeeId"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Email        string        `json:"email"`
	Roles        map[Role]bool `json:"roles"`
	passwordHash []byte
}

// Compare this password to the actual employees' password
func (e *Employee) ValidatePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(e.passwordHash, password)
}

func NewEmployee(employeeID int, password []byte) (*Employee, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("couldn't hash password, %v", err)
	}
	return &Employee{
		EmployeeID:   employeeID,
		passwordHash: hash,
	}, nil
}

/*
	Helpers to serialize `passwordHash` for the DB without exposing it
	to clients
*/

func (e *Employee) marshal() ([]byte, error) {

	return json.Marshal(struct {
		EmployeeID   int           `json:"employeeId"`
		FirstName    string        `json:"firstName"`
		LastName     string        `json:"lastName"`
		Email        string        `json:"email"`
		Roles        map[Role]bool `json:"roles"`
		PasswordHash []byte        `json:"passwordHash"`
	}{
		EmployeeID:   e.EmployeeID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Email:        e.Email,
		Roles:        e.Roles,
		PasswordHash: e.passwordHash,
	})
}

func unmarshal(data []byte) (*Employee, error) {
	e := struct {
		EmployeeID   int           `json:"employeeId"`
		FirstName    string        `json:"firstName"`
		LastName     string        `json:"lastName"`
		Email        string        `json:"email"`
		Roles        map[Role]bool `json:"roles"`
		PasswordHash []byte        `json:"passwordHash"`
	}{}
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("unmarshalling, %v", err)
	}
	return &Employee{
		EmployeeID:   e.EmployeeID,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Email:        e.Email,
		Roles:        e.Roles,
		passwordHash: e.PasswordHash,
	}, nil
}

// Create an employee IFF it doesn't already exist
func (e *Employee) Create() error {
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf(employeePrefix, e.EmployeeID)
	// Create a record
	changes, err := redis.Int(conn.Do("HSETNX", employeeAll, e.EmployeeID, key))
	if err != nil {
		return fmt.Errorf("error verifying for existence, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("employee with ID %d already exists, changes %d", e.EmployeeID, changes)
	}

	emplByte, err := e.marshal()
	if err != nil {
		return fmt.Errorf("marshalling employee, %v", err)
	}

	ok, err := redis.String(conn.Do("SET", key, emplByte))
	if err != nil {
		return fmt.Errorf("error setting employee, %v", err)
	}
	if ok != "OK" {
		return fmt.Errorf("expected 1 employee change, got %d", 1)
	}
	return nil
}

// Update
func (e *Employee) Update() error {
	conn := pool.Get()
	defer conn.Close()
	key := fmt.Sprintf(employeePrefix, e.EmployeeID)

	emplByte, err := e.marshal()
	if err != nil {
		return fmt.Errorf("marshalling employee, %v", err)
	}

	changes, err := redis.Int(conn.Do("SET", key, emplByte))
	if err != nil {
		return fmt.Errorf("error setting employee, %v", err)
	}
	if changes != 1 {
		return fmt.Errorf("expected 1 employee change, got %d", 1)
	}
	return nil
}

func (e *Employee) Delete() error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", employeeAll, e.EmployeeID)
	if err != nil {
		return fmt.Errorf("deleting employee from record list, %v", err)
	}
	_, err = conn.Do("DEL", fmt.Sprintf(employeePrefix, e.EmployeeID))
	if err != nil {
		return fmt.Errorf("deleting actual employee, %v", err)
	}
	return nil
}

func FindEmployee(employeeID int) (*Employee, bool, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.String(conn.Do("HGET", employeeAll, employeeID))
	if err != nil {
		return nil, false, fmt.Errorf("getting employee ID %d, %v", employeeID, err)
	}

	if id == "" {
		// doesn't exist
		return nil, false, nil
	}

	emplByte, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil, false, fmt.Errorf("getting bytes for employee %d, %v", employeeID, err)
	}

	empl, err := unmarshal(emplByte)
	return empl, true, err
}

func FindAllEmployees() ([]*Employee, error) {
	conn := pool.Get()
	defer conn.Close()

	members, err := redis.Strings(conn.Do("HVALS", employeeAll))
	if err != nil {
		return nil, fmt.Errorf("getting vals for %s, %v", employeeAll, err)
	}

	if err := conn.Send("MULTI"); err != nil {
		return nil, fmt.Errorf("preparing pipelined GET over all employees, %v", err)
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

	var empl *Employee
	employees := make([]*Employee, 0, len(replies))
	for _, reply := range replies {
		empl, err = unmarshal([]byte(reply))
		if err != nil {
			return employees, err
		}
		employees = append(employees, empl)
	}
	return employees, nil
}
