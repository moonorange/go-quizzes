package main

import "fmt"

// You can write any code in here to test anything

func main() {
	q1()
	q2()
}

func q1() {
	var i interface{} = 42

	v, ok := i.(int) // v=42, ok=true
	s, ok2 := i.(string)

	fmt.Println(v, ok)
	fmt.Println(s, ok2)

	// s2 := i.(string) // this panics without comma, ok form
}

// typed interface
type Stringer interface {
	String() string
}

type Person struct {
	Name string
}

func (p *Person) String() string {
	return p.Name
}

func q2() {
	// empty interface can hold any values
	var emptyInterface interface{}

	emptyInterface = 42
	fmt.Println(emptyInterface)

	emptyInterface = "hello"
	fmt.Println(emptyInterface)

	emptyInterface = struct {
		Field1 int
		Field2 string
	}{
		Field1: 42,
		Field2: "hello",
	}
	fmt.Println(emptyInterface)

	var s Stringer
	s = &Person{Name: "John"}
	fmt.Println(s)

	// s = 42 // this is not valid
}
