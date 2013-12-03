package models_test

import (
	"github.com/collinglass/moustacheMS/server/models"
	"reflect"
	"testing"
)

func Test_Ward_CanCreateAndDelete(t *testing.T) {

	ward := models.NewWard(1, "Ward East", 1, 2)

	err := ward.Create()
	if err != nil {
		t.Errorf("ward shouldn't exist, %v", err)
	}

	err = ward.Delete()
	if err != nil {
		t.Errorf("wardoye should be deletanle, %v", err)
	}
}

func Test_Ward_CanCreateAndRetrieve(t *testing.T) {
	wardID := 1
	want := models.NewWard(wardID, "Ward West", 4, 8)

	err := want.Create()
	if err != nil {
		t.Errorf("ward shouldn't exist, %v", err)
	}
	defer want.Delete()

	got, ok, err := models.FindWard(wardID)
	if err != nil {
		t.Fatalf("Finding back ward %d, %v", wardID, err)
	}
	if !ok {
		t.Error("should have found ward that just created")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

var manyWard = []struct {
	id     int
	name   string
	nurse  int
	doctor int
}{
	{1, "Ward #1", 34, 544},
	{2, "Ward #2", 34, 544},
	{4, "Ward #4", 34, 544},
	{9999, "Ward #9999", 34, 544},
	{-1, "Ward #-1", 34, 544},
	{-9292, "Ward #-9292", 34, 544},
}

func Test_Ward_CanCreateManyAndRetrieveMany(t *testing.T) {
	wantWard := make(map[int]*models.Ward)
	for _, val := range manyWard {
		ward := models.NewWard(val.id, val.name, val.nurse, val.doctor)

		if err := ward.Create(); err != nil {
			t.Fatalf("cant' create ward, %v", err)
		}
		defer ward.Delete()
		wantWard[ward.WardID] = ward
	}

	gotAll, err := models.FindAllWards()
	if err != nil {
		t.Fatalf("should have gotAll wards, %v", err)
	}

	if len(gotAll) != len(wantWard) {
		t.Fatalf("want %d ward, got %d", len(wantWard), len(gotAll))
	}

	gotWard := make(map[int]*models.Ward)
	for _, got := range gotAll {

		t.Logf("Got WardID %d", got.WardID)
		gotWard[got.WardID] = got

		want, ok := wantWard[got.WardID]

		if !ok {
			t.Errorf("Got back ward we didn't put, %#v", got)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}

	for key, want := range wantWard {
		got, ok := gotWard[key]
		if !ok {
			t.Errorf("Missing a ward, want, %#v", want)
		} else if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v, got %#v", want, got)
		}
	}
}

func Test_Ward_CantCreateWardAlreadyExists(t *testing.T) {

	exists := models.NewWard(1, "Ward East", 1, 2)

	err := exists.Create()
	if err != nil {
		t.Errorf("ward shouldn't exist, %v", err)
	}
	defer exists.Delete()

	extraWard := models.NewWard(1, "Ward East", 1, 2)

	err = extraWard.Create()
	if err == nil {
		t.Errorf("ward shouldn't be created (already exists)")
		defer extraWard.Delete()
	}
}
