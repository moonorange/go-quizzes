// embedding struct quiz

package main

import "fmt"

type Engine struct{}

func (Engine) Start() {
	fmt.Println("Engine starting...")
}

type Car struct {
	Engine
}

func (Car) Start() {
	fmt.Println("Car starting...")
}

func (Car) End() {
	fmt.Println("Car ending...")
}

func main() {
	c := Car{}
	c.Start() // Which one runs?
}
