package main

import (
	"bufio"
	"flag"
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
	DIV
	SUB
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

type Token struct {
	command int
	args    []int
}

var (
	program []Token

	// running - flag that allows us to automate fetch cycle
	running bool = true

	// if current instruction is jmp then don't need to increment pc
	isJump bool

	stack     [256]int
	registers [NumOfRegisters]int

	// sp - stack pointer
	sp = registers[SP]
	// pc - program counter or instruction pointer
	pc = registers[PC]
)

// fetch - returns current instruction from program
func fetch() Token {
	return program[pc]
}

// eval - evaluates the given instruction
func eval(instr Token) {
	isJump = false

	switch instr.command {
	case HLT:
		running = false
	case PSH:
		sp++
		stack[sp] = instr.args[0]
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
		isJump = true
		pc = instr.args[0]
	case DIV:
		arg1 := stack[sp]
		sp--
		arg2 := stack[sp]
		stack[sp] = arg1 / arg2
	case SUB:
		arg1 := stack[sp]
		sp--
		arg2 := stack[sp]
		stack[sp] = arg1 - arg2
	}
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
	case "DIV":
		return DIV
	case "SUB":
		return SUB
	default:
		return HLT
	}
}

func tokenizer(input []string) []Token {
	result := make([]Token, len(input))

	for i := range input {
		temp := strings.Split(input[i], " ")

		token := Token{}

		token.command = recognizeInstr(temp[0])

		if len(temp) > 1 {
			temp2 := temp[1:]
			token.args = make([]int, len(temp2))

			var err error
			for i := range temp2 {
				token.args[i], err = strconv.Atoi(temp2[i])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}

		result[i] = token
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
	source := flag.String("file", "", "")

	flag.Parse()
	if *source == "" {
		fmt.Fprintln(os.Stderr, "-flag=[file] required")
//		return
	}

	// explicit assigning default values:
	// registers[SP] = -1 'cause execution stack is empty
	// registers[PC] = 0 'case execution of program begins from 0 address
	registers[SP] = -1
	registers[PC] = 0

//	program = tokenizer(readFile(*source))
	program = tokenizer(readFile("test.mac"))

	for running {
		eval(fetch())
		if !isJump {
			pc++
		}
	}

}
