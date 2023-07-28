package helpers

import "strings"

func DistinctBy[T any](slice []T, fn func(x T) []string) []T {
	distinctMap := make(map[string]T)

	// Iterate over the slice and add elements to the map
	for _, elem := range slice {
		key := fn(elem)
		_, ok := distinctMap[strings.Join(key, ",")]

		if !ok {
			distinctMap[strings.Join(key, ",")] = elem
		}
	}

	// Create a new slice with distinct elements
	distinctSlice := make([]T, 0, len(distinctMap))
	for _, val := range distinctMap {
		distinctSlice = append(distinctSlice, val)
	}

	return distinctSlice
}
