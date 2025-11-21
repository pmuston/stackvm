package stackvm

import (
	"errors"
	"fmt"
)

// Standard VM errors.
var (
	ErrStackOverflow        = errors.New("stack overflow")
	ErrStackUnderflow       = errors.New("stack underflow")
	ErrInvalidMemoryAddress = errors.New("invalid memory address")
	ErrReadOnlyMemory       = errors.New("memory is read-only")
	ErrInvalidInstruction   = errors.New("invalid instruction")
	ErrInvalidOpcode        = errors.New("invalid opcode")
	ErrInstructionLimit     = errors.New("instruction limit exceeded")
	ErrDivisionByZero       = errors.New("division by zero")
	ErrTypeMismatch         = errors.New("type mismatch")
	ErrTimeout              = errors.New("execution timeout")
	ErrInvalidOperand       = errors.New("invalid operand")
	ErrInvalidProgram       = errors.New("invalid program")
	ErrUnresolvedLabel      = errors.New("unresolved label")
)

// VMError wraps errors with execution context.
type VMError struct {
	// Err is the underlying error
	Err error

	// PC is the program counter at failure
	PC int

	// InstructionCount is the number of instructions executed before failure
	InstructionCount uint32

	// StackDepth is the stack depth at failure
	StackDepth int

	// Opcode is the instruction that failed
	Opcode Opcode

	// Message provides additional context
	Message string
}

// Error implements the error interface.
func (e *VMError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("VM error at PC=%d (opcode=%d, instructions=%d, stack=%d): %s: %v",
			e.PC, e.Opcode, e.InstructionCount, e.StackDepth, e.Message, e.Err)
	}
	return fmt.Sprintf("VM error at PC=%d (opcode=%d, instructions=%d, stack=%d): %v",
		e.PC, e.Opcode, e.InstructionCount, e.StackDepth, e.Err)
}

// Unwrap returns the underlying error.
func (e *VMError) Unwrap() error {
	return e.Err
}

// Is implements error matching for errors.Is.
func (e *VMError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// IsStackError returns true if the error is a stack overflow or underflow.
func IsStackError(err error) bool {
	return errors.Is(err, ErrStackOverflow) || errors.Is(err, ErrStackUnderflow)
}

// IsMemoryError returns true if the error is a memory-related error.
func IsMemoryError(err error) bool {
	return errors.Is(err, ErrInvalidMemoryAddress) || errors.Is(err, ErrReadOnlyMemory)
}

// IsLimitError returns true if the error is an instruction limit or timeout error.
func IsLimitError(err error) bool {
	return errors.Is(err, ErrInstructionLimit) || errors.Is(err, ErrTimeout)
}
