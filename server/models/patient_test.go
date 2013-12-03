package models_test

import (
	"github.com/collinglass/moustacheMS/server/models"
	"reflect"
	"testing"
)

func Test_Patient_CanCreateAndDelete(t *testing.T) {

	patient := models.NewPatient(1, "John", "Smith")

	err := patient.Create()
	if err != nil {
		t.Errorf("patient shouldn't exist, %v", err)
	}

	err = patient.Delete()
	if err != nil {
		t.Errorf("patientoye should be deletanle, %v", err)
	}
}

func Test_Patient_CanCreateAndRetrieve(t *testing.T) {
	patientID := 1
	want := models.NewPatient(patientID, "John", "Smith")

	err := want.Create()
	if err != nil {
		t.Errorf("patient shouldn't exist, %v", err)
	}
	defer want.Delete()

	got, ok, err := models.FindPatient(patientID)
	if err != nil {
		t.Fatalf("Finding back patient %d, %v", patientID, err)
	}
	if !ok {
		t.Error("should have found patient that just created")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

var manyPatient = []struct {
	id    int
	first string
	last  string
}{
	{1, "John", "Smith"},
	{2, "2John", "Smith2"},
	{4, "4John", "Smith4"},
	{9999, "9999John", "Smith9999"},
	{-1, "1John", "Smith1"},
	{-9292, "9292John", "Smith9292"},
}

func Test_Patient_CanCreateManyAndRetrieveMany(t *testing.T) {
	wantPatient := make(map[int]*models.Patient)
	for _, val := range manyPatient {
		patient := models.NewPatient(val.id, val.first, val.last)

		if err := patient.Create(); err != nil {
			t.Fatalf("cant' create patient, %v", err)
		}
		defer patient.Delete()
		wantPatient[patient.PatientID] = patient
	}

	gotAll, err := models.FindAllPatients()
	if err != nil {
		t.Fatalf("should have gotAll patients, %v", err)
	}

	if len(gotAll) != len(wantPatient) {
		t.Fatalf("want %d patient, got %d", len(wantPatient), len(gotAll))
	}

	gotPatient := make(map[int]*models.Patient)
	for _, got := range gotAll {

		t.Logf("Got PatientID %d", got.PatientID)
		gotPatient[got.PatientID] = got

		want, ok := wantPatient[got.PatientID]

		if !ok {
			t.Errorf("Got back patient we didn't put, %#v", got)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}

	for key, want := range wantPatient {
		got, ok := gotPatient[key]
		if !ok {
			t.Errorf("Missing a patient, want, %#v", want)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}
}

func Test_Patient_CantCreatePatientAlreadyExists(t *testing.T) {

	exists := models.NewPatient(1, "John", "Smith")

	err := exists.Create()
	if err != nil {
		t.Errorf("patient shouldn't exist, %v", err)
	}
	defer exists.Delete()

	extraPatient := models.NewPatient(1, "John", "Smith")

	err = extraPatient.Create()
	if err == nil {
		t.Errorf("patient shouldn't be created (already exists)")
		defer extraPatient.Delete()
	}
}
