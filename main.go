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
	// sp - stack pointer
	sp int = -1
	stack [256]int
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
	case PSH:
		sp++
		// incr pc for give args of PSH instruction
		pc++
		// push args to stack
		stack[sp] = program[pc]
	case POP:
		// tempVal - value that popped from top of the stack
		tempVal := stack[sp]
		// decrementing stack pointer
		sp--
		fmt.Printf("popped %d\n", tempVal)
	}
}

func main() {
	for running {
		eval(fetch())
		pc++
	}
}
