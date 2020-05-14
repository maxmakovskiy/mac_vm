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
		// incr sp 'cause we add new value to the top of the stack
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
	case ADD:
		// pop first value from the stack
		// and decrementing stack pointer
		arg1 := stack[sp]
		sp--
		// pop next value from the stack
		arg2 := stack[sp]
		// theoretically, we need to do
		// sp-- decrementing sp after popped the second value from the stack
		// sp++ incrementing sp for push result of arg1+arg2
		// but (sp--) + (sp++) = 0 and we just stay it out
		stack[sp] = arg1 + arg2
	}
}

func main() {
	for running {
		eval(fetch())
		pc++
	}
}
