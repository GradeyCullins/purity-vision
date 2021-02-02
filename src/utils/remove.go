package utils

import "fmt"

// StringSliceRemove removes element at index i from slice arr.
func StringSliceRemove(arr []string, i int) ([]string, error) {
	if i > len(arr)-1 {
		return nil, fmt.Errorf("Index %d must be between 0 and the length of argument arr (%d)", i, len(arr)-1)
	}
	arr[i] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return arr, nil
}
