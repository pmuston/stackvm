package stackvm

import (
	"context"
	"testing"
	"time"
)

func TestVMBasicExecution(t *testing.T) {
	t.Run("PUSH and HALT", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 42),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if !result.Halted {
			t.Error("Expected program to be halted")
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
		if result.InstructionCount != 2 {
			t.Errorf("InstructionCount = %d, want 2", result.InstructionCount)
		}
	})

	t.Run("PUSHI and HALT", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSHI, 100),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if !result.Halted {
			t.Error("Expected program to be halted")
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
	})

	t.Run("PUSH POP HALT", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 10),
			NewInstruction(OpPUSH, 20),
			NewInstruction(OpPOP, 0),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if !result.Halted {
			t.Error("Expected program to be halted")
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1 (one value remaining after POP)", result.StackDepth)
		}
	})

	t.Run("DUP duplicates top value", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 42),
			NewInstruction(OpDUP, 0),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 2 {
			t.Errorf("StackDepth = %d, want 2 (original + duplicate)", result.StackDepth)
		}
	})

	t.Run("SWAP exchanges top two values", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 10),
			NewInstruction(OpPUSH, 20),
			NewInstruction(OpSWAP, 0),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 2 {
			t.Errorf("StackDepth = %d, want 2", result.StackDepth)
		}
	})

	t.Run("Multiple operations", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 1),
			NewInstruction(OpPUSH, 2),
			NewInstruction(OpPUSH, 3),
			NewInstruction(OpDUP, 0),   // Stack: 1 2 3 3
			NewInstruction(OpPOP, 0),   // Stack: 1 2 3
			NewInstruction(OpSWAP, 0),  // Stack: 1 3 2
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 3 {
			t.Errorf("StackDepth = %d, want 3", result.StackDepth)
		}
		if result.InstructionCount != 7 {
			t.Errorf("InstructionCount = %d, want 7", result.InstructionCount)
		}
	})
}

func TestVMErrors(t *testing.T) {
	t.Run("Stack underflow on POP", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPOP, 0), // Nothing to pop
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != ErrStackUnderflow {
			t.Errorf("Expected ErrStackUnderflow, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected non-nil result")
		}
		if result.Error != ErrStackUnderflow {
			t.Errorf("Result.Error = %v, want ErrStackUnderflow", result.Error)
		}
	})

	t.Run("Stack underflow on DUP", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpDUP, 0), // Nothing to duplicate
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		_, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != ErrStackUnderflow {
			t.Errorf("Expected ErrStackUnderflow, got %v", err)
		}
	})

	t.Run("Stack underflow on SWAP", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 1),
			NewInstruction(OpSWAP, 0), // Only one value on stack
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		_, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != ErrStackUnderflow {
			t.Errorf("Expected ErrStackUnderflow, got %v", err)
		}
	})

	t.Run("Stack overflow", func(t *testing.T) {
		vm := New()
		// Create a program that tries to overflow the stack
		instructions := make([]Instruction, 0, 300)
		for i := 0; i < 300; i++ {
			instructions = append(instructions, NewInstruction(OpPUSH, int32(i)))
		}
		instructions = append(instructions, NewInstruction(OpHALT, 0))
		program := NewProgram(instructions)
		memory := NewSimpleMemory(0)

		_, err := vm.Execute(program, memory, ExecuteOptions{
			MaxStackDepth: 256,
		})

		if err != ErrStackOverflow {
			t.Errorf("Expected ErrStackOverflow, got %v", err)
		}
	})

	t.Run("Invalid opcode", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(Opcode(99), 0), // Invalid opcode
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		_, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != ErrInvalidOpcode {
			t.Errorf("Expected ErrInvalidOpcode, got %v", err)
		}
	})
}

func TestVMExecuteOptions(t *testing.T) {
	t.Run("MaxInstructions limit", func(t *testing.T) {
		vm := New()
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 1),
			NewInstruction(OpPUSH, 2),
			NewInstruction(OpPUSH, 3),
			NewInstruction(OpPUSH, 4),
			NewInstruction(OpPUSH, 5),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{
			MaxInstructions: 3,
		})

		if err != ErrInstructionLimit {
			t.Errorf("Expected ErrInstructionLimit, got %v", err)
		}
		if result.InstructionCount != 3 {
			t.Errorf("InstructionCount = %d, want 3", result.InstructionCount)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		vm := New()
		// Create a long program
		instructions := make([]Instruction, 0, 10000)
		for i := 0; i < 10000; i++ {
			instructions = append(instructions, NewInstruction(OpPUSH, int32(i)))
			instructions = append(instructions, NewInstruction(OpPOP, 0))
		}
		instructions = append(instructions, NewInstruction(OpHALT, 0))
		program := NewProgram(instructions)
		memory := NewSimpleMemory(0)

		result, err := vm.Execute(program, memory, ExecuteOptions{
			Timeout: 1 * time.Nanosecond, // Very short timeout
		})

		if err != ErrTimeout {
			t.Errorf("Expected ErrTimeout, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected non-nil result")
		}
	})

	t.Run("Context cancellation", func(t *testing.T) {
		vm := New()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 1),
			NewInstruction(OpHALT, 0),
		})
		memory := NewSimpleMemory(0)

		_, err := vm.Execute(program, memory, ExecuteOptions{
			Context: ctx,
		})

		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestVMReset(t *testing.T) {
	vm := New()
	program := NewProgram([]Instruction{
		NewInstruction(OpPUSH, 1),
		NewInstruction(OpPUSH, 2),
		NewInstruction(OpHALT, 0),
	})
	memory := NewSimpleMemory(0)

	// Execute once
	result1, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("First Execute() error = %v", err)
	}

	// Reset
	vm.Reset()

	// Execute again - should work the same
	result2, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Second Execute() error = %v", err)
	}

	if result1.StackDepth != result2.StackDepth {
		t.Errorf("After reset, stack depths differ: %d vs %d", result1.StackDepth, result2.StackDepth)
	}
	if result1.InstructionCount != result2.InstructionCount {
		t.Errorf("After reset, instruction counts differ: %d vs %d", result1.InstructionCount, result2.InstructionCount)
	}
}

func TestProgramInterface(t *testing.T) {
	t.Run("NewProgram", func(t *testing.T) {
		instructions := []Instruction{
			NewInstruction(OpPUSH, 1),
			NewInstruction(OpHALT, 0),
		}
		program := NewProgram(instructions)

		if program == nil {
			t.Fatal("NewProgram returned nil")
		}
		if len(program.Instructions()) != 2 {
			t.Errorf("Instructions() length = %d, want 2", len(program.Instructions()))
		}
	})

	t.Run("NewProgramWithMetadata", func(t *testing.T) {
		instructions := []Instruction{
			NewInstruction(OpPUSH, 1),
		}
		metadata := ProgramMetadata{
			Name:        "test",
			Version:     "1.0",
			Author:      "tester",
			Description: "test program",
			Created:     time.Now(),
		}
		program := NewProgramWithMetadata(instructions, metadata)

		if program == nil {
			t.Fatal("NewProgramWithMetadata returned nil")
		}
		meta := program.Metadata()
		if meta.Name != "test" {
			t.Errorf("Metadata.Name = %s, want test", meta.Name)
		}
	})

	t.Run("Symbol table", func(t *testing.T) {
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 1),
		})

		if program.SymbolTable() != nil {
			t.Error("Initial symbol table should be nil")
		}

		program.AddSymbol(0, "START")
		if program.SymbolTable() == nil {
			t.Error("Symbol table should not be nil after AddSymbol")
		}
		if program.SymbolTable()[0] != "START" {
			t.Errorf("Symbol at 0 = %s, want START", program.SymbolTable()[0])
		}

		symbols := map[int]string{
			0: "BEGIN",
			5: "END",
		}
		program.SetSymbolTable(symbols)
		if program.SymbolTable()[0] != "BEGIN" {
			t.Errorf("After SetSymbolTable, symbol at 0 = %s, want BEGIN", program.SymbolTable()[0])
		}
	})
}

func TestEmptyProgram(t *testing.T) {
	vm := New()
	program := NewProgram([]Instruction{})
	memory := NewSimpleMemory(0)

	result, err := vm.Execute(program, memory, ExecuteOptions{})

	if err != nil {
		t.Errorf("Empty program should not error, got %v", err)
	}
	if !result.Halted {
		t.Error("Empty program should be marked as halted")
	}
	if result.InstructionCount != 0 {
		t.Errorf("InstructionCount = %d, want 0", result.InstructionCount)
	}
}

func TestProgramWithoutHalt(t *testing.T) {
	vm := New()
	program := NewProgram([]Instruction{
		NewInstruction(OpPUSH, 1),
		NewInstruction(OpPUSH, 2),
		// No HALT instruction
	})
	memory := NewSimpleMemory(0)

	result, err := vm.Execute(program, memory, ExecuteOptions{})

	if err != nil {
		t.Errorf("Program without HALT should not error, got %v", err)
	}
	if !result.Halted {
		t.Error("Program should be marked as halted when reaching end")
	}
	if result.StackDepth != 2 {
		t.Errorf("StackDepth = %d, want 2", result.StackDepth)
	}
}
