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

const (
	A = iota
	B
	C
	D
	E
	F
	PC
	SP
	NumOfRegisters
)

var (
	program []int = []int{
		PSH, 5,
		PSH, 6,
		ADD,
		POP,
		HLT,
	}

	// running - flag that allows us to automate fetch cycle
	running bool = true

	stack [256]int
	registers [NumOfRegisters]int

	// sp - stack pointer
	sp = registers[SP]
	// pc - program counter or instruction pointer
	pc = registers[PC]
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
	// explicit assigning default values:
	// registers[SP] = -1 'cause execution stack is empty
	// registers[PC] = 0 'case execution of program begins from 0 address
	registers[SP] = -1
	registers[PC] = 0

	for running {
		eval(fetch())
		pc++
	}
}
