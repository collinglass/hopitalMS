package models

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
)

type Employee struct {
	EmployeeID   int    `json:"employeeId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Roles        []int  `json:"roles"`
	passwordHash []byte `json:"-"`
}

func NewEmployee(employeeID int, password []byte) (*Employee, error) {
	defer clearBytes(password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("couldn't hash password, %v", err)
	}
	return &Employee{passwordHash: hash}, nil
}

func (e *Employee) validatePassword(password []byte) error {
	defer clearBytes(password)
	return bcrypt.CompareHashAndPassword(e.passwordHash, password)
}

func clearBytes(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}
