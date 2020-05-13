package main

import "fmt"

// instruction set of this virtual machine
// opcodes begins since 0
const (
	PSH = iota
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
	// pc - program counter or instruction pointer
	pc int = 0
	running bool = true
)

// fetch - returns current instruction from program
func fetch() int {
	return program[pc]
}

func main() {
	for running {
		x := fetch()
		if x == HLT {
			running = false
		} else {
			fmt.Println(x)
			pc++
		}
	}
}
