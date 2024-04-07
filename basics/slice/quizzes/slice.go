package main

import "fmt"

func main() {
	// Q1
	slices := create2DSlice(8, 8)
	fmt.Printf("A1 2d array %d\n", slices)
	fmt.Printf("\n")

	// Q2
	copySlice()
}

// Implement the function to create 2d slice of length dy, each element of which is a slice of dx
func create2DSlice(dy int, dx int) [][]uint8 {
	return [][]uint8{}
}

// Change the code so that new slice doesn't affect "slice"'s underlying array
func copySlice() {
	slice := []string{"a", "b", "c"}
	new := slice
	new[0] = "CHANGED"
	fmt.Printf("A2 new: %s\n", new)
	fmt.Printf("A2 slice has been modified: %s\n", slice)
}
