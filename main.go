package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	NumOfRegisters
)

type Token struct {
	command int
	args    []int
}

var (
	version = 0.1
	overflow = 256

	program []Token

	// running - flag that allows us to automate fetch cycle
	running = true

	// if current instruction is jmp then don't need to increment pc
	isJump bool

	stack = NewStack()
	registers [NumOfRegisters]int

	// sp - stack pointer
//	sp = stack.Length()
	// pc - program counter or instruction pointer
	pc = 0
)

// fetch - returns current instruction from program
func fetch() Token {
	return program[pc]
}

// eval - evaluates the given instruction
func eval(instr Token) {
	isJump = false

	// if stack is out of established boundaries
	if stack.Length() > overflow {
		fmt.Printf("Stack Overflow!\n")
		running = false
		return
	}

	switch instr.command {
	case HLT:
		running = false
	case PSH:
		stack.Push(instr.args[0])
	case POP:
		// tempVal - value that popped from top of the stack
		tempVal := stack.Pop()
		fmt.Printf("popped %d\n", tempVal)
	case ADD:
		arg1 := stack.Pop()
		arg2 := stack.Pop()
		stack.Push(arg1 + arg2)
	case MUL:
		arg1 := stack.Pop()
		arg2 := stack.Pop()
		stack.Push(arg1 * arg2)
	case JMP:
		isJump = true
		pc = instr.args[0]
	case DIV:
		arg1 := stack.Pop()
		arg2 := stack.Pop()
		stack.Push(arg1 / arg2)
	case SUB:
		arg1 := stack.Pop()
		arg2 := stack.Pop()
		stack.Push(arg1 - arg2)
	case MOV:
		reg1 := instr.args[0]
		reg2 := instr.args[1]
		registers[reg1] = registers[reg2]
		registers[reg2] = 0
	case GLD:
		reg := instr.args[0]
		stack.Push(registers[reg])
	case GPT:
		reg := instr.args[0]
		registers[reg] = stack.Top()
	case -1:
		return
	}
}

// recognizeInstr - returns opcode for given string representation of instruction
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
	case "HLT":
		return HLT
	default:
		return -1
	}
}

func recognizeRegister(register string) int {
	switch register	{
	case "A":
		return A
	case "B":
		return B
	case "C":
		return C
	case "D":
		return D
	case "E":
		return E
	case "F":
		return F
	default:
		return -1
	}
}

// tokenizer - returns slice of Tokens for given input slice of strings
func tokenizer(input []string) []Token {
	result := make([]Token, len(input))

	for i := range input {
		temp := strings.Split(input[i], " ")

		token := Token{}

		token.command = recognizeInstr(temp[0])

		if len(temp) > 1 {
			arguments := temp[1:]
			token.args = make([]int, len(arguments))

			if token.command == GPT || token.command == GLD {
				token.args[0] = recognizeRegister(arguments[0])
			} else if token.command == MOV {
				token.args[0] = recognizeRegister(arguments[0])
				token.args[1] = recognizeRegister(arguments[1])
			} else {
				var err error
				for i := range arguments {
					token.args[i], err = strconv.Atoi(arguments[i])
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}
			}
		}

		result[i] = token
	}

	return result
}

// readFile - returns slice of strings populated by lines from file
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

// help - prints help info about each instruction for MacVM
func help() {
	fmt.Printf("The commands are: \n")
	fmt.Printf("PSH <v> - pushes <v> value to the stack\n")
	fmt.Printf("POP - pops value from the stack\n")
	fmt.Printf("ADD - adds top two values on stack\n")
	fmt.Printf("SUB - subtracts top two values on stack\n")
	fmt.Printf("DIV - divides top tow values on stack\n")
	fmt.Printf("MUL - multiplies top two values on stack\n")
	fmt.Printf("MOV <R1> <R2> - move the value from <R2> register to <R1> register\n")
	fmt.Printf("JMP <addr> - jump to the <addr> address\n")
	fmt.Printf("GLD <R> - loads <R> register to the top of stack\n")
	fmt.Printf("GPT <R> - loads top of the stack to <R> register\n")
	fmt.Printf("HLT - halts program\n")
}

// registerDump - prints current state of all registers for MacVM
func registerDump() {
	fmt.Printf("REGISTER DUMP:\n")
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("General purpose registers:\n")
	fmt.Printf("Register[A] = %d;\tRegister[B] = %d;\n", registers[A], registers[B])
	fmt.Printf("Register[C] = %d;\tRegister[D] = %d;\n", registers[C], registers[D])
	fmt.Printf("Register[E] = %d;\tRegister[F] = %d;\n", registers[E], registers[F])
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("Special purpose registers:\n")
	fmt.Printf("Program Counter = %d;\nStack Pointer = 0x%X;\n", pc, stack.Length())
	fmt.Printf("==========================================\n")
}

func printStack(slice []int) {
	fmt.Printf("Stack Pointer = 0x%X\n", stack.Length())
	for i, k := range slice {
		if i == 0 {
			fmt.Printf(" -> %d\n", k)
			continue
		}
		fmt.Printf("    %d\n", k)
	}
}

// stackDump - prints current state of execution stack with given depth
// if depth equals -1 then stackDump prints all elements from stack
func stackDump(depth int) {
	fmt.Printf("STACK DUMP:\n")
	fmt.Printf("------------------------------------------\n")
	tstack := stack.ToSlice()
	if stack.Length() == 0 {
		fmt.Printf("Stack is empty!\n")
	} else if depth == -1 || depth >= stack.Length() {
		printStack(tstack)
	} else {
		printStack(tstack[:depth+1])
	}
	fmt.Printf("==========================================\n")
}

// runFile - runs execution from given file
func runFile(file string) {
	program = tokenizer(readFile(file))
	for running {
		eval(fetch())
		if !isJump {
			pc++
		}
	}
}

// funPrompt - runs execution by interactive mode
func runPrompt() {
	program = make([]Token, 0)

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("MacVM (build %.2f)\n", version)
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
		} else if cleanString == "registers" {
			registerDump()
		} else {
			matched, err := regexp.MatchString("stack -*[0-9]+", cleanString)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if matched {
				num, err := strconv.Atoi(cleanString[6:])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				stackDump(num)
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
}

func main() {
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
