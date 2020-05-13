package main

import "fmt"

const (
	PSH = iota + 1
	ADD
	POP
	SET
	HLT
)

var (
	program []int = []int{
		PSH, 5,
		PSH, 6,
		ADD,
		POP,
		HLT,
	}
)

func main() {
	fmt.Println("Hello vm!")
}
