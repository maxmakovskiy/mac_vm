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
	// running - flag that allows us to automate fetch cycle
	running bool = true
)

// fetch - returns current instruction from program
func fetch() int {
	return program[pc]
}

// eval - evaluates the given instruction
func eval(instr int) {
	switch instr {
	case HLT:
		running = false
	default:
		fmt.Println(instr)
	}
}

func main() {
	for running {
		eval(fetch())
		pc++
	}
}
