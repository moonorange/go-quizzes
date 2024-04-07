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
	res := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		xs := make([]uint8, dx)
		for x := 0; x < dx; x++ {
			xs[x] = uint8((x * y))
		}
		res[y] = xs
	}
	return res
}

func copySlice() {
	slice := []string{"a", "b", "c"}
	new := make([]string, len(slice))
	// copy slice to new
	copy(new, slice)
	// change new
	new[0] = "CHANGED"
	fmt.Printf("A new: %s\n", new)
	fmt.Printf("A slice has not been modified %s\n", slice)
}
