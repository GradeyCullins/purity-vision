package utils

import "testing"

func TestStringSliceRemove(t *testing.T) {
	arr := []string{"hello", "world", "fubar"}
	expected := []string{"hello", "world"}

	res, err := StringSliceRemove(arr, 2)
	if err != nil {
		t.Fatal(err)
	}

	for i := range res {
		if res[i] != expected[i] {
			t.Fatalf("found %s but expected %s", res[i], expected[i])
		}
	}

	if _, err = StringSliceRemove(arr, 100); err == nil {
		t.Fatal("an error should have been thrown")
	}
}
