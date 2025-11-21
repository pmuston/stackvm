package stackvm

// ExecutionContext provides access to VM state during instruction execution.
// This interface is used by custom instruction handlers to interact with the VM.
type ExecutionContext interface {
	// Stack Operations

	// Push adds a value to the top of the stack.
	// Returns ErrStackOverflow if the stack is full.
	Push(value Value) error

	// Pop removes and returns the value from the top of the stack.
	// Returns ErrStackUnderflow if the stack is empty.
	Pop() (Value, error)

	// Peek returns the value at the top of the stack without removing it.
	// Returns ErrStackUnderflow if the stack is empty.
	Peek() (Value, error)

	// PeekN returns the nth value from the top of the stack without removing it.
	// n=0 returns the top value, n=1 returns the second value, etc.
	// Returns ErrStackUnderflow if n is beyond the stack depth.
	PeekN(n int) (Value, error)

	// StackDepth returns the current number of values on the stack.
	StackDepth() int

	// Program Counter

	// PC returns the current program counter value.
	PC() int

	// SetPC sets the program counter to the specified value.
	SetPC(pc int)

	// Jump sets the program counter to the specified offset.
	// This is equivalent to SetPC(offset).
	Jump(offset int)

	// Memory

	// Memory returns the memory provider associated with this execution.
	Memory() Memory

	// Execution Control

	// InstructionCount returns the number of instructions executed so far.
	InstructionCount() uint32

	// IncrementInstructionCount increments the instruction counter by one.
	IncrementInstructionCount()

	// Halt stops execution. The VM will terminate after the current instruction.
	Halt()

	// IsHalted returns true if execution has been halted.
	IsHalted() bool
}
