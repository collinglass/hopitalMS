package models_test

import (
	"github.com/collinglass/moustacheMS/server/models"
	"reflect"
	"testing"
)

func Test_CanValidatePassword(t *testing.T) {
	emplID := 1
	pass := []byte("hello!")

	empl, err := models.NewEmployee(emplID, pass)
	if err != nil {
		t.Fatal(err)
	}

	if err := empl.ValidatePassword(pass); err != nil {
		t.Errorf("same password but: %v", err)
	}

	notPass := []byte("not the pass")
	err = empl.ValidatePassword(notPass)
	if err == nil {
		t.Errorf("different password but no error")
	}
}

func Test_CanCreateAndDelete(t *testing.T) {
	emplID := 1
	pass := []byte("hello!")
	empl, err := models.NewEmployee(emplID, pass)
	if err != nil {
		t.Fatal(err)
	}

	err = empl.Create()
	if err != nil {
		t.Errorf("Employee shouldn't exist, %v", err)
	}

	err = empl.Delete()
	if err != nil {
		t.Errorf("Employe should be deletanle, %v", err)
	}
}

func Test_CanCreateAndRetrieve(t *testing.T) {
	emplID := 1
	pass := []byte("hello!")
	want, err := models.NewEmployee(emplID, pass)
	if err != nil {
		t.Fatal(err)
	}

	err = want.Create()
	if err != nil {
		t.Errorf("Employee shouldn't exist, %v", err)
	}
	defer want.Delete()

	got, ok, err := models.FindEmployee(emplID)
	if err != nil {
		t.Errorf("Finding back employee %d, %v", emplID, err)
	}
	if !ok {
		t.Error("should have found employee that just created")
	}

	if err := got.ValidatePassword(pass); err != nil {
		t.Errorf("same password but: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

var manyEmpl = []struct {
	id   int
	pass []byte
}{
	{1, []byte("hello")},
	{2, []byte("hello!")},
	{4, []byte("hello!!!!")},
	{9999, []byte("h@@@@ello")},
	{-1, []byte("hell###$o")},
	{-9292, []byte("hel#!#@!lo")},
}

func Test_CanCreateManyAndRetrieveMany(t *testing.T) {
	wantEmpl := make(map[int]*models.Employee)
	for _, val := range manyEmpl {
		empl, err := models.NewEmployee(val.id, val.pass)
		if err != nil {
			t.Fatalf("can't make new employee, %v", err)
		}

		if err := empl.Create(); err != nil {
			t.Fatalf("cant' create employee, %v", err)
		}
		defer empl.Delete()
		wantEmpl[empl.EmployeeID] = empl
	}

	gotAll, err := models.FindAllEmployees()
	if err != nil {
		t.Errorf("should have gotAll employees, %v", err)
	}

	gotEmpl := make(map[int]*models.Employee)
	for _, got := range gotAll {
		gotEmpl[got.EmployeeID] = got

		want, ok := wantEmpl[got.EmployeeID]

		if !ok {
			t.Errorf("Got back employee we didn't put, %#v", got)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}

	for key, want := range wantEmpl {
		got, ok := gotEmpl[key]
		if !ok {
			t.Errorf("Missing an employee, %#v", want)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}
}

func Test_CantCreateEmployeeAlreadyExists(t *testing.T) {
	emplID := 1
	pass := []byte("hello!")
	exists, err := models.NewEmployee(emplID, pass)
	if err != nil {
		t.Fatal(err)
	}

	err = exists.Create()
	if err != nil {
		t.Errorf("Employee shouldn't exist, %v", err)
	}
	defer exists.Delete()

	extraEmpl, err := models.NewEmployee(emplID, pass)
	if err != nil {
		t.Fatal(err)
	}

	err = extraEmpl.Create()
	if err == nil {
		t.Errorf("Employee shouldn't be created (already exists)")
		defer extraEmpl.Delete()
	}
}
