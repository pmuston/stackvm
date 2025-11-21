package stackvm

import (
	"os"
	"testing"
)

func TestNewAssembler(t *testing.T) {
	asm := NewAssembler()
	if asm == nil {
		t.Fatal("NewAssembler() returned nil")
	}
}

func TestAssembleSimple(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 10
		PUSH 5
		ADD
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}

	// Verify instruction sequence
	expected := []Opcode{OpPUSH, OpPUSH, OpADD, OpHALT}
	for i, inst := range instructions {
		if inst.Opcode != expected[i] {
			t.Errorf("Instruction %d: opcode = %d, want %d", i, inst.Opcode, expected[i])
		}
	}
}

func TestAssembleWithLabels(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 1
		JMPZ ELSE
		PUSH 100
		JMP END
	ELSE:
		PUSH 200
	END:
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 6 {
		t.Errorf("Expected 6 instructions, got %d", len(instructions))
	}

	// Verify JMPZ points to correct label
	jmpzInst := instructions[1]
	if jmpzInst.Opcode != OpJMPZ {
		t.Errorf("Instruction 1 should be JMPZ, got %d", jmpzInst.Opcode)
	}
	if jmpzInst.Operand != 4 { // Should point to ELSE label
		t.Errorf("JMPZ operand = %d, want 4", jmpzInst.Operand)
	}
}

func TestAssembleWithComments(t *testing.T) {
	asm := NewAssembler()

	source := `
		; This is a comment
		PUSH 10     ; Push 10
		# Another comment style
		PUSH 5      # Push 5
		ADD         ; Add them
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}
}

func TestAssembleMemoryOperations(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 42
		STORE 0
		LOAD 0
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}

	// Verify STORE operand
	storeInst := instructions[1]
	if storeInst.Opcode != OpSTORE {
		t.Errorf("Instruction 1 should be STORE, got %d", storeInst.Opcode)
	}
	if storeInst.Operand != 0 {
		t.Errorf("STORE operand = %d, want 0", storeInst.Operand)
	}
}

func TestAssembleCaseInsensitive(t *testing.T) {
	asm := NewAssembler()

	source := `
		push 10
		Push 5
		ADD
		add
		halt
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 5 {
		t.Errorf("Expected 5 instructions, got %d", len(instructions))
	}
}

func TestAssembleNegativeNumbers(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH -10
		PUSHI -5
		ADD
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}

	// Verify negative operands
	pushInst := instructions[0]
	if pushInst.Operand != -10 {
		t.Errorf("PUSH operand = %d, want -10", pushInst.Operand)
	}

	pushiInst := instructions[1]
	if pushiInst.Operand != -5 {
		t.Errorf("PUSHI operand = %d, want -5", pushiInst.Operand)
	}
}

func TestAssembleFloats(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 3.14
		PUSH 2.5
		MUL
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}
}

func TestAssembleUnknownOpcode(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 10
		BADOPCODE
		HALT
	`

	_, err := asm.Assemble(source)
	if err == nil {
		t.Error("Assemble() should fail with unknown opcode")
	}
}

func TestAssembleUnresolvedLabel(t *testing.T) {
	asm := NewAssembler()

	source := `
		PUSH 1
		JMP NONEXISTENT
		HALT
	`

	_, err := asm.Assemble(source)
	if err == nil {
		t.Error("Assemble() should fail with unresolved label")
	}
}

func TestAssembleFile(t *testing.T) {
	asm := NewAssembler()

	// Test with simple_add.asm
	program, err := asm.AssembleFile("testdata/programs/simple_add.asm")
	if err != nil {
		t.Fatalf("AssembleFile() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}

	// Verify it's the correct program
	expected := []Opcode{OpPUSH, OpPUSH, OpADD, OpHALT}
	for i, inst := range instructions {
		if inst.Opcode != expected[i] {
			t.Errorf("Instruction %d: opcode = %d, want %d", i, inst.Opcode, expected[i])
		}
	}
}

func TestAssembleFileNotFound(t *testing.T) {
	asm := NewAssembler()

	_, err := asm.AssembleFile("nonexistent.asm")
	if err == nil {
		t.Error("AssembleFile() should fail for non-existent file")
	}
}

func TestAssembleAndExecute(t *testing.T) {
	testCases := []struct {
		name          string
		source        string
		expectedStack int
		expectedValue float64
	}{
		{
			name: "Simple addition",
			source: `
				PUSH 10
				PUSH 5
				ADD
				HALT
			`,
			expectedStack: 1,
			expectedValue: 15,
		},
		{
			name: "Subtraction",
			source: `
				PUSH 10
				PUSH 3
				SUB
				HALT
			`,
			expectedStack: 1,
			expectedValue: 7,
		},
		{
			name: "Multiplication",
			source: `
				PUSH 6
				PUSH 7
				MUL
				HALT
			`,
			expectedStack: 1,
			expectedValue: 42,
		},
		{
			name: "Conditional true",
			source: `
				PUSH 15
				PUSH 10
				GT
				JMPZ ELSE
				PUSH 1
				JMP END
			ELSE:
				PUSH 0
			END:
				HALT
			`,
			expectedStack: 1,
			expectedValue: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			asm := NewAssembler()
			program, err := asm.Assemble(tc.source)
			if err != nil {
				t.Fatalf("Assemble() failed: %v", err)
			}

			vm := New()
			memory := NewSimpleMemory(0)
			result, err := vm.Execute(program, memory, ExecuteOptions{})
			if err != nil {
				t.Fatalf("Execute() failed: %v", err)
			}

			if result.StackDepth != tc.expectedStack {
				t.Errorf("StackDepth = %d, want %d", result.StackDepth, tc.expectedStack)
			}
		})
	}
}

func TestAssembleTestdataPrograms(t *testing.T) {
	programs := []struct {
		file          string
		expectedStack int
	}{
		{"simple_add.asm", 1},
		{"conditional.asm", 1},
		{"memory.asm", 1},
		{"loop.asm", 1},
		{"math.asm", 1},
	}

	asm := NewAssembler()
	vm := New()

	for _, prog := range programs {
		t.Run(prog.file, func(t *testing.T) {
			path := "testdata/programs/" + prog.file

			// Check if file exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Skipf("Test file %s does not exist", path)
			}

			program, err := asm.AssembleFile(path)
			if err != nil {
				t.Fatalf("AssembleFile(%s) failed: %v", path, err)
			}

			memory := NewSimpleMemory(10)
			result, err := vm.Execute(program, memory, ExecuteOptions{
				MaxInstructions: 1000,
			})
			if err != nil {
				t.Fatalf("Execute(%s) failed: %v", path, err)
			}

			if result.StackDepth != prog.expectedStack {
				t.Errorf("StackDepth = %d, want %d", result.StackDepth, prog.expectedStack)
			}

			if !result.Halted {
				t.Error("Program should have halted")
			}

			vm.Reset()
		})
	}
}

func TestAssembleAllOpcodes(t *testing.T) {
	source := `
		; Stack operations
		PUSH 1
		PUSHI 2
		DUP
		POP
		SWAP
		OVER
		ROT
		POP
		POP
		POP

		; Arithmetic
		PUSH 10
		PUSH 5
		ADD
		PUSH 3
		SUB
		PUSH 2
		MUL
		PUSH 2
		DIV
		PUSH 3
		MOD
		NEG
		ABS
		INC
		DEC

		; Logic
		PUSH 1
		PUSH 1
		AND
		PUSH 0
		OR
		NOT
		PUSH 1
		XOR

		; Comparison
		PUSH 5
		PUSH 3
		GT
		POP
		PUSH 3
		PUSH 5
		LT
		POP
		PUSH 5
		PUSH 5
		EQ
		POP
		PUSH 5
		PUSH 3
		NE
		POP
		PUSH 5
		PUSH 5
		GE
		POP
		PUSH 3
		PUSH 5
		LE
		POP

		; Math
		PUSH 16
		SQRT
		PUSH 0
		SIN
		POP
		PUSH 0
		COS
		POP
		PUSH 5
		PUSH 10
		MIN
		PUSH 5
		PUSH 10
		MAX
		PUSH 3.7
		FLOOR
		PUSH 3.2
		CEIL
		PUSH 3.5
		ROUND

		HALT
	`

	asm := NewAssembler()
	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	vm := New()
	memory := NewSimpleMemory(0)
	result, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if !result.Halted {
		t.Error("Program should have halted")
	}
}

func TestAssembleWithRegistry(t *testing.T) {
	// Create a custom instruction
	registry := NewInstructionRegistry()

	doubleHandler := &testInstructionHandler{
		name: "DOUBLE",
		fn: func(ctx ExecutionContext, operand int32) error {
			val, err := ctx.Pop()
			if err != nil {
				return err
			}
			f, _ := val.AsFloat()
			return ctx.Push(FloatValue(f * 2))
		},
	}

	err := registry.Register(128, doubleHandler)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	asm := NewAssembler()
	asm.SetRegistry(registry)

	source := `
		PUSH 5
		DOUBLE
		HALT
	`

	program, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	vm := NewWithConfig(Config{
		StackSize:           256,
		InstructionRegistry: registry,
	})

	memory := NewSimpleMemory(0)
	result, err := vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
}

// testInstructionHandler is a test implementation of InstructionHandler.
type testInstructionHandler struct {
	name string
	fn   func(ExecutionContext, int32) error
}

func (h *testInstructionHandler) Execute(ctx ExecutionContext, operand int32) error {
	if h.fn != nil {
		return h.fn(ctx, operand)
	}
	return nil
}

func (h *testInstructionHandler) Name() string {
	return h.name
}
