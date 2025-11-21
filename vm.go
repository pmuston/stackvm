package stackvm

import (
	"context"
	"time"
)

// VM is the virtual machine interface for executing programs.
type VM interface {
	// Execute runs a program with the given memory and options.
	// Returns execution results and statistics, or an error.
	Execute(program Program, memory Memory, opts ExecuteOptions) (*Result, error)

	// Reset clears the VM state for reuse.
	Reset()
}

// ExecuteOptions configures VM execution behavior.
type ExecuteOptions struct {
	// MaxInstructions limits the number of instructions executed (0 = unlimited).
	// Returns ErrInstructionLimit if exceeded.
	MaxInstructions uint32

	// MaxStackDepth sets the stack size limit (0 = default 256).
	// Returns ErrStackOverflow if exceeded.
	MaxStackDepth int

	// Timeout sets a wall-clock timeout for execution (0 = no timeout).
	// Returns ErrTimeout if exceeded.
	Timeout time.Duration

	// Context provides cancellation support (nil = no cancellation).
	// Returns the context error if cancelled.
	Context context.Context
}

// Result contains execution statistics and results.
type Result struct {
	// InstructionCount is the number of instructions executed.
	InstructionCount uint32

	// StackDepth is the final stack depth.
	StackDepth int

	// ExecutionTime is the total execution time.
	ExecutionTime time.Duration

	// Halted is true if a HALT instruction was reached.
	Halted bool

	// Error is the execution error, if any (nil if successful).
	Error error
}

// Config configures a VM instance.
type Config struct {
	// StackSize is the initial stack capacity (default 256).
	StackSize int

	// DefaultInstrLimit is the default instruction limit (0 = unlimited).
	DefaultInstrLimit uint32

	// InstructionRegistry provides custom instruction handlers (nil = standard only).
	InstructionRegistry InstructionRegistry

	// ValueConverter provides custom type conversions (nil = defaults).
	ValueConverter ValueConverter
}

// InstructionRegistry allows registration of custom instruction handlers.
// This will be implemented in a future phase.
type InstructionRegistry interface {
	// Register adds a handler for a custom opcode (128-255).
	Register(opcode Opcode, handler InstructionHandler) error

	// Unregister removes a handler for an opcode.
	Unregister(opcode Opcode) error

	// Get retrieves a handler for an opcode.
	Get(opcode Opcode) (InstructionHandler, bool)

	// List returns all registered custom opcodes.
	List() []Opcode

	// Names returns a mapping of opcodes to their names.
	Names() map[Opcode]string
}

// InstructionHandler executes a custom instruction.
// This will be implemented in a future phase.
type InstructionHandler interface {
	// Execute performs the instruction operation.
	Execute(ctx ExecutionContext, operand int32) error

	// Name returns the mnemonic for the instruction.
	Name() string
}

// ValueConverter provides custom type conversion logic.
// This will be implemented in a future phase.
type ValueConverter interface {
	// Convert converts a value to the target type.
	Convert(value Value, targetType ValueType) (Value, error)
}

// New creates a new VM with default configuration.
func New() VM {
	return NewWithConfig(Config{
		StackSize: 256,
	})
}

// NewWithConfig creates a new VM with custom configuration.
func NewWithConfig(config Config) VM {
	return newExecutor(config)
}
