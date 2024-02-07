package main

import "fmt"

func main() {
	slices := create2DSlice(8, 8)
	fmt.Printf("%d\n", slices)
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
