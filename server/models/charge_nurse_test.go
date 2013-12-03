package models_test

import (
	"github.com/collinglass/moustacheMS/server/models"
	"reflect"
	"testing"
)

func Test_ChargeNurse_CanCreateAndDelete(t *testing.T) {

	nurse := models.NewChargeNurse(1, "x3456", "x62276")

	err := nurse.Create()
	if err != nil {
		t.Errorf("nurse shouldn't exist, %v", err)
	}

	err = nurse.Delete()
	if err != nil {
		t.Errorf("nurseoye should be deletanle, %v", err)
	}
}

func Test_ChargeNurse_CanCreateAndRetrieve(t *testing.T) {
	nurseID := 1
	want := models.NewChargeNurse(nurseID, "x3456", "x62276")

	err := want.Create()
	if err != nil {
		t.Errorf("nurse shouldn't exist, %v", err)
	}
	defer want.Delete()

	got, ok, err := models.FindChargeNurse(nurseID)
	if err != nil {
		t.Fatalf("Finding back nurse %d, %v", nurseID, err)
	}
	if !ok {
		t.Error("should have found nurse that just created")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

var manyChargeNurse = []struct {
	id     int
	telExt string
	bipExt string
}{
	{1, "x3456", "x62276"},
	{2, "x2", "x112"},
	{4, "x4", "x114"},
	{9999, "x9999", "x119999"},
	{-1, "x1", "x111"},
	{-9292, "x9292", "x119292"},
}

func Test_ChargeNurse_CanCreateManyAndRetrieveMany(t *testing.T) {
	wantChargeNurse := make(map[int]*models.ChargeNurse)
	for _, val := range manyChargeNurse {
		nurse := models.NewChargeNurse(val.id, val.telExt, val.bipExt)

		if err := nurse.Create(); err != nil {
			t.Fatalf("cant' create nurse, %v", err)
		}
		defer nurse.Delete()
		wantChargeNurse[nurse.EmployeeID] = nurse
	}

	gotAll, err := models.FindAllChargeNurses()
	if err != nil {
		t.Fatalf("should have gotAll nurses, %v", err)
	}

	if len(gotAll) != len(wantChargeNurse) {
		t.Fatalf("want %d nurse, got %d", len(wantChargeNurse), len(gotAll))
	}

	gotChargeNurse := make(map[int]*models.ChargeNurse)
	for _, got := range gotAll {

		t.Logf("Got EmployeeID %d", got.EmployeeID)
		gotChargeNurse[got.EmployeeID] = got

		want, ok := wantChargeNurse[got.EmployeeID]

		if !ok {
			t.Errorf("Got back nurse we didn't put, %#v", got)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}

	for key, want := range wantChargeNurse {
		got, ok := gotChargeNurse[key]
		if !ok {
			t.Errorf("Missing a nurse, want, %#v", want)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}
}

func Test_ChargeNurse_CantCreatenurseAlreadyExists(t *testing.T) {

	exists := models.NewChargeNurse(1, "x3456", "x62276")

	err := exists.Create()
	if err != nil {
		t.Errorf("nurse shouldn't exist, %v", err)
	}
	defer exists.Delete()

	extraChargeNurse := models.NewChargeNurse(1, "x3456", "x62276")

	err = extraChargeNurse.Create()
	if err == nil {
		t.Errorf("nurse shouldn't be created (already exists)")
		defer extraChargeNurse.Delete()
	}
}
