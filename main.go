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
	DIV
	SUB
	MOV
	GLD
	GPT
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
	version = 0.1

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
	case MOV:
		reg1 := instr.args[0]
		reg2 := instr.args[1]
		registers[reg1] = registers[reg2]
	case GLD:
		sp++
		reg := instr.args[0]
		stack[sp] = registers[reg]
	case GPT:
		reg := instr.args[0]
		registers[reg] = stack[sp]
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
	case "MOV":
		return MOV
	case "GLD":
		return GLD
	case "GPT":
		return GPT
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

func help() {
	fmt.Printf("The commands are: \n")
	fmt.Printf("PSH - pushes value to the stack\n")
	fmt.Printf("POP - pops value from the stack\n")
	fmt.Printf("ADD - adds top two values on stack\n")
	fmt.Printf("SUB - subtracts top two values on stack\n")
	fmt.Printf("DIV - divides top tow values on stack\n")
	fmt.Printf("MUL - multiplies top two values on stack\n")
	fmt.Printf("MOV - move the value from one register to another\n")
	fmt.Printf("JMP - jump to the passing address\n")
	fmt.Printf("GLD - loads given register to the top of stack\n")
	fmt.Printf("GPT - pushes top of the stack to given register\n")
	fmt.Printf("HLT - halts program\n")
}

func registerDump() {
	fmt.Printf("REGISTER DUMP:\n")
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("General purpose registers:\n")
	fmt.Printf("Register[A] = %d;\tRegister[B] = %d;\n", registers[A], registers[B])
	fmt.Printf("Register[C] = %d;\tRegister[D] = %d;\n", registers[C], registers[D])
	fmt.Printf("Register[E] = %d;\tRegister[F] = %d;\n", registers[E], registers[F])
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("Special purpose registers:\n")
	fmt.Printf("Program Counter = %d;\nStack Pointer = %d;\n", pc, sp)
	fmt.Printf("==========================================\n")
}

func stackDump(depth int) {
	fmt.Printf("STACK DUMP:\n")
	fmt.Printf("------------------------------------------\n")
	for i, k := range stack {
		if (i+1) > depth {
			break
		}
		fmt.Printf(" -> %d\n", k)
	}
	fmt.Printf("==========================================\n")
}

func runFile(file string) {
	program = tokenizer(readFile(file))
	for running {
		eval(fetch())
		if !isJump {
			pc++
		}
	}
}

func runPrompt() {
	program = make([]Token, 0)

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Mac VM (build %f)\n", version)
	fmt.Printf("Type \"help\" for more information\n")

	// interactive cycle
	for {
		fmt.Print(">>> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		cleanString := strings.Trim(input, "\r\n")
		if cleanString == "" {
			continue
		} else if cleanString == "exit" {
			break
		} else if cleanString == "help" {
			fmt.Println("Hello! It's help :D How are you?")
			help()
		} else {
			// run <input>
			program = append(program, tokenizer([]string{cleanString})[0])
			if running {
				eval(fetch())
				if !isJump {
					pc++
				}
			}
		}
	}
}

func main() {
	// explicit assigning default values:
	// registers[SP] = -1 'cause execution stack is empty
	// registers[PC] = 0 'case execution of program begins from 0 address
	registers[SP] = -1
	registers[PC] = 0

	var cmdArgs []string = os.Args[1:]
	var argsLen int = len(cmdArgs)

	if argsLen > 1 {
		fmt.Println("Usage: vm.exe [script]")
		os.Exit(64)
	} else if argsLen == 1 {
		runFile(cmdArgs[0])
	} else {
		runPrompt()
	}
}
