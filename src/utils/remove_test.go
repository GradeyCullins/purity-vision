package utils

import (
	"testing"
)

func TestStringSliceRemove(t *testing.T) {
	var err error
	arr := []string{"hello", "world", "fubar"}
	size := len(arr)
	removed := 1

	for i := 0; i < size; i++ {
		arr, err = StringSliceRemove(arr, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(arr) != size-removed {
			t.Fatalf("result array should have size: %d\n", size-removed)
		}
		removed += 1
	}

	if len(arr) != 0 {
		t.Fatalf("result array should be empty")
	}

	if _, err := StringSliceRemove(arr, 100); err == nil {
		t.Fatal("an error should have been thrown")
	}

	arr = []string{"hello", "world", "fubar"}

	// Do a few arbitrary removes.
	if arr, err = StringSliceRemove(arr, 1); err != nil {
		t.Fatal(err)
	}
	if len(arr) != 2 {
		t.Fatalf("result array should have two elements")
	}

	if arr, err = StringSliceRemove(arr, 1); err != nil {
		t.Fatal(err)
	}
	if len(arr) != 1 {
		t.Fatalf("result array should have one element")
	}

	if arr, err = StringSliceRemove(arr, 0); err != nil {
		t.Fatal(err)
	}
	if len(arr) != 0 {
		t.Fatalf("result array should be empty")
	}

}
