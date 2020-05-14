package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// instruction set of this virtual machine
// opcodes begins since 0
const (
	PSH = iota
	ADD
	POP
	JMP
	HLT
	MUL
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
		JMP, 0,
		PSH, 1,
		PSH, 2,
		MUL,
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
	case MUL:
		arg1 := stack[sp]
		sp--
		arg2 := stack[sp]
		stack[sp] = arg1 * arg2
	case JMP:
		pc++
		pc = pc+1
		pc++
	}
}

type Line struct {
	command int
	args []int
}

func recognizeInstr(instr string) int {
	switch instr {
	case "PSH":
		return PSH
	case "POP":
		return POP
	case "ADD":
		return ADD
	case "MUL":
		return MUL
	case "JMP":
		return JMP
	case "HLT":
		return HLT
	default:
		return HLT
	}
}

func tokenizer(input []string) []Line {
	result := make([]Line, len(input))

	for i := range input {
		temp := strings.Split(input[i], " ")

		line := Line{}

		line.command = recognizeInstr(temp[0])

		if len(temp) > 1 {
			temp2 := temp[1:]
			line.args = make([]int, len(temp2))

			var err error
			for i := range temp2 {
				line.args[i], err = strconv.Atoi(temp2[i])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}

		result[i] = line
	}


	return result
}

func readFile(fileName string) []string {
	result := make([]string, 0)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return result
}

func main() {
	// explicit assigning default values:
	// registers[SP] = -1 'cause execution stack is empty
	// registers[PC] = 0 'case execution of program begins from 0 address
	registers[SP] = -1
	registers[PC] = 0

	var programz []Line
	programz = tokenizer(readFile("test.mac"))
	fmt.Println(programz)


	for running {
		eval(fetch())
		pc++
	}
}
