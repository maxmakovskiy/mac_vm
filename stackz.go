package main

type item struct {
	value int
	next  *item
}

type Stack struct {
	top  *item
	size int
}

// NewStack returns new instance of struct Stack
func NewStack() Stack {
	return Stack{}
}

// Length returns current size of stack
func (s *Stack) Length() int {
	return (*s).size
}

// Push adds new value to the top of stack
func (s *Stack) Push(v int) {
	// initialize new item
	// and add it to top of stack
	(*s).top = &item{
		value: v,
		next:  (*s).top,
	}

	// increase size
	(*s).size++
}

// Pop pops the top value from the stack
func (s *Stack) Pop() int {
	// get value from top element
	val := (*s).top.value
	// make top element next element
	(*s).top = (*s).top.next
	// decrease size
	(*s).size--

	return val
}

// ToSlice - converts current contents of stack to the slice
func (s *Stack) ToSlice() []int {
	result := make([]int, 0)
	iter := s.top

	for iter != nil {
		result = append(result, iter.value)
		iter = iter.next
	}

	return result
}

/*
func main() {
	stack := new(Stack)

	fmt.Printf("Length of stack before populating -> %d\n", stack.Length())

	for i := 1; i <= 10; i++ {
		stack.Push(i)
	}

	fmt.Printf("Length of stack after populating -> %d\n", stack.Length())
	fmt.Println("Elements of stack:")

	for i := 0; i < 10; i++ {
		fmt.Println(stack.Pop())
	}

}
*/
