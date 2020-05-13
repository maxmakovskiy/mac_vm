package main

import "fmt"

const (
	PSH = iota + 1
	ADD
	POP
	SET
	HLT
)

func main() {
	fmt.Println("Hello vm!")
}
