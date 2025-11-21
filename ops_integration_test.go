package stackvm

import (
	"testing"
)

func TestArithmeticIntegration(t *testing.T) {
	vm := New()
	memory := NewSimpleMemory(0)

	t.Run("Basic arithmetic", func(t *testing.T) {
		// Program: (10 + 5) * 2 - 3 = 27
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 10),
			NewInstruction(OpPUSH, 5),
			NewInstruction(OpADD, 0),   // 15
			NewInstruction(OpPUSH, 2),
			NewInstruction(OpMUL, 0),   // 30
			NewInstruction(OpPUSH, 3),
			NewInstruction(OpSUB, 0),   // 27
			NewInstruction(OpHALT, 0),
		})

		result, err := vm.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
	})

	t.Run("Division and modulo", func(t *testing.T) {
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 17),
			NewInstruction(OpPUSH, 5),
			NewInstruction(OpDIV, 0),   // 3.4
			NewInstruction(OpHALT, 0),
		})

		vm.Reset()
		_, err := vm.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
}

func TestLogicAndComparisonIntegration(t *testing.T) {
	vm := New()
	memory := NewSimpleMemory(0)

	t.Run("Comparison with logic", func(t *testing.T) {
		// Program: (10 > 5) && (3 < 8)
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 10),
			NewInstruction(OpPUSH, 5),
			NewInstruction(OpGT, 0),    // true
			NewInstruction(OpPUSH, 3),
			NewInstruction(OpPUSH, 8),
			NewInstruction(OpLT, 0),    // true
			NewInstruction(OpAND, 0),   // true
			NewInstruction(OpHALT, 0),
		})

		result, err := vm.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
	})
}

func TestMathFunctionsIntegration(t *testing.T) {
	vm := New()
	memory := NewSimpleMemory(0)

	t.Run("Pythagorean theorem", func(t *testing.T) {
		// Program: sqrt(3^2 + 4^2) = 5
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 3),
			NewInstruction(OpDUP, 0),
			NewInstruction(OpMUL, 0),   // 9
			NewInstruction(OpPUSH, 4),
			NewInstruction(OpDUP, 0),
			NewInstruction(OpMUL, 0),   // 16
			NewInstruction(OpADD, 0),   // 25
			NewInstruction(OpSQRT, 0),  // 5
			NewInstruction(OpHALT, 0),
		})

		vm.Reset()
		result, err := vm.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
	})

	t.Run("Trigonometry", func(t *testing.T) {
		// Program: sin(0) = 0
		program := NewProgram([]Instruction{
			NewInstruction(OpPUSH, 0),
			NewInstruction(OpSIN, 0),   // 0
			NewInstruction(OpHALT, 0),
		})

		vm.Reset()
		_, err := vm.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})
}

func TestComplexProgram(t *testing.T) {
	vm := New()
	memory := NewSimpleMemory(0)

	// Complex program combining multiple operations
	// Calculate: max(abs(-10), sqrt(16)) + floor(3.7) = max(10, 4) + 3 = 13
	program := NewProgram([]Instruction{
		NewInstruction(OpPUSH, -10),
		NewInstruction(OpABS, 0),       // 10
		NewInstruction(OpPUSH, 16),
		NewInstruction(OpSQRT, 0),      // 4
		NewInstruction(OpMAX, 0),       // 10
		NewInstruction(OpPUSH, 3),
		NewInstruction(OpPUSH, 7),
		NewInstruction(OpDIV, 0),       // 0.428...
		NewInstruction(OpPUSH, 10),
		NewInstruction(OpMUL, 0),       // 4.28...
		NewInstruction(OpFLOOR, 0),     // 4
		NewInstruction(OpADD, 0),       // 14
		NewInstruction(OpHALT, 0),
	})

	result, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
	if !result.Halted {
		t.Error("Program should have halted")
	}
}

func TestStackOperations(t *testing.T) {
	vm := New()
	memory := NewSimpleMemory(0)

	// Test OVER and ROT operations
	program := NewProgram([]Instruction{
		NewInstruction(OpPUSH, 1),
		NewInstruction(OpPUSH, 2),
		NewInstruction(OpPUSH, 3),
		NewInstruction(OpOVER, 0),  // Stack: 1 2 3 2
		NewInstruction(OpROT, 0),   // Stack: 1 3 2 2
		NewInstruction(OpHALT, 0),
	})

	result, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.StackDepth != 4 {
		t.Errorf("StackDepth = %d, want 4", result.StackDepth)
	}
}
