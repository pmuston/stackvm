package stackvm

import (
	"errors"
	"testing"
)

func TestStandardErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrStackOverflow", ErrStackOverflow},
		{"ErrStackUnderflow", ErrStackUnderflow},
		{"ErrInvalidMemoryAddress", ErrInvalidMemoryAddress},
		{"ErrReadOnlyMemory", ErrReadOnlyMemory},
		{"ErrInvalidInstruction", ErrInvalidInstruction},
		{"ErrInvalidOpcode", ErrInvalidOpcode},
		{"ErrInstructionLimit", ErrInstructionLimit},
		{"ErrDivisionByZero", ErrDivisionByZero},
		{"ErrTypeMismatch", ErrTypeMismatch},
		{"ErrTimeout", ErrTimeout},
		{"ErrInvalidOperand", ErrInvalidOperand},
		{"ErrInvalidProgram", ErrInvalidProgram},
		{"ErrUnresolvedLabel", ErrUnresolvedLabel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("Error %s is nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("Error %s has empty message", tt.name)
			}
		})
	}
}

func TestVMError(t *testing.T) {
	t.Run("Error message with context", func(t *testing.T) {
		vmErr := &VMError{
			Err:              ErrStackOverflow,
			PC:               42,
			InstructionCount: 100,
			StackDepth:       256,
			Opcode:           OpPUSH,
			Message:          "stack limit reached",
		}

		msg := vmErr.Error()
		if msg == "" {
			t.Error("Error message should not be empty")
		}
		// Check that the message contains key information
		if !containsString(msg, "PC=42") {
			t.Errorf("Error message should contain PC: %s", msg)
		}
		if !containsString(msg, "stack limit reached") {
			t.Errorf("Error message should contain context: %s", msg)
		}
	})

	t.Run("Error message without context", func(t *testing.T) {
		vmErr := &VMError{
			Err:              ErrDivisionByZero,
			PC:               10,
			InstructionCount: 20,
			StackDepth:       5,
			Opcode:           OpDIV,
		}

		msg := vmErr.Error()
		if msg == "" {
			t.Error("Error message should not be empty")
		}
		if !containsString(msg, "PC=10") {
			t.Errorf("Error message should contain PC: %s", msg)
		}
	})
}

func TestVMErrorUnwrap(t *testing.T) {
	underlyingErr := ErrStackUnderflow
	vmErr := &VMError{
		Err:              underlyingErr,
		PC:               5,
		InstructionCount: 10,
		StackDepth:       0,
		Opcode:           OpPOP,
	}

	unwrapped := vmErr.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlyingErr)
	}
}

func TestVMErrorIs(t *testing.T) {
	tests := []struct {
		name   string
		vmErr  *VMError
		target error
		want   bool
	}{
		{
			"Matches underlying error",
			&VMError{Err: ErrStackOverflow},
			ErrStackOverflow,
			true,
		},
		{
			"Does not match different error",
			&VMError{Err: ErrStackOverflow},
			ErrStackUnderflow,
			false,
		},
		{
			"Matches with errors.Is",
			&VMError{Err: ErrDivisionByZero},
			ErrDivisionByZero,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.vmErr, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStackError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Stack overflow is stack error", ErrStackOverflow, true},
		{"Stack underflow is stack error", ErrStackUnderflow, true},
		{"Memory error is not stack error", ErrInvalidMemoryAddress, false},
		{"Division by zero is not stack error", ErrDivisionByZero, false},
		{"Wrapped stack overflow", &VMError{Err: ErrStackOverflow}, true},
		{"Wrapped stack underflow", &VMError{Err: ErrStackUnderflow}, true},
		{"Wrapped other error", &VMError{Err: ErrTimeout}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStackError(tt.err); got != tt.want {
				t.Errorf("IsStackError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMemoryError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Invalid memory address is memory error", ErrInvalidMemoryAddress, true},
		{"Read-only memory is memory error", ErrReadOnlyMemory, true},
		{"Stack error is not memory error", ErrStackOverflow, false},
		{"Division by zero is not memory error", ErrDivisionByZero, false},
		{"Wrapped invalid address", &VMError{Err: ErrInvalidMemoryAddress}, true},
		{"Wrapped read-only", &VMError{Err: ErrReadOnlyMemory}, true},
		{"Wrapped other error", &VMError{Err: ErrTimeout}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMemoryError(tt.err); got != tt.want {
				t.Errorf("IsMemoryError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Instruction limit is limit error", ErrInstructionLimit, true},
		{"Timeout is limit error", ErrTimeout, true},
		{"Stack error is not limit error", ErrStackOverflow, false},
		{"Division by zero is not limit error", ErrDivisionByZero, false},
		{"Wrapped instruction limit", &VMError{Err: ErrInstructionLimit}, true},
		{"Wrapped timeout", &VMError{Err: ErrTimeout}, true},
		{"Wrapped other error", &VMError{Err: ErrStackOverflow}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLimitError(tt.err); got != tt.want {
				t.Errorf("IsLimitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
