package main

import "fmt"

// You can write any code in here to test anything

func main() {
	var i interface{} = 42

	v, ok := i.(int) // v=42, ok=true
	s, ok2 := i.(string)

	fmt.Println(v, ok)
	fmt.Println(s, ok2)

	s2 := i.(string) // this panics without comma, ok form
}
