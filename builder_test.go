package stackvm

import (
	"testing"
)

func TestNewProgramBuilder(t *testing.T) {
	builder := NewProgramBuilder()
	if builder == nil {
		t.Fatal("NewProgramBuilder() returned nil")
	}

	program, err := builder.Build()
	if err != nil {
		t.Errorf("Build() on empty builder failed: %v", err)
	}
	if len(program.Instructions()) != 0 {
		t.Errorf("Empty builder should produce empty program, got %d instructions", len(program.Instructions()))
	}
}

func TestBuilderBasicOperations(t *testing.T) {
	// Test fluent API with basic arithmetic
	builder := NewProgramBuilder()
	program, err := builder.
		Push(10).
		Push(5).
		Add().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
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

func TestBuilderStackOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(1).
		Push(2).
		Dup().
		Swap().
		Over().
		Rot().
		Pop().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if len(program.Instructions()) != 8 {
		t.Errorf("Expected 8 instructions, got %d", len(program.Instructions()))
	}
}

func TestBuilderArithmeticOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(10).
		Push(3).
		Add().
		Push(2).
		Mul().
		Push(1).
		Sub().
		Neg().
		Abs().
		Inc().
		Dec().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	opcodes := []Opcode{
		OpPUSH, OpPUSH, OpADD,
		OpPUSH, OpMUL,
		OpPUSH, OpSUB,
		OpNEG, OpABS, OpINC, OpDEC,
		OpHALT,
	}

	instructions := program.Instructions()
	if len(instructions) != len(opcodes) {
		t.Fatalf("Expected %d instructions, got %d", len(opcodes), len(instructions))
	}

	for i, inst := range instructions {
		if inst.Opcode != opcodes[i] {
			t.Errorf("Instruction %d: opcode = %d, want %d", i, inst.Opcode, opcodes[i])
		}
	}
}

func TestBuilderLogicOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(1).
		Push(1).
		And().
		Push(0).
		Or().
		Not().
		Push(1).
		Xor().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if len(program.Instructions()) != 9 {
		t.Errorf("Expected 9 instructions, got %d", len(program.Instructions()))
	}
}

func TestBuilderComparisonOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(5).
		Push(3).
		Gt().
		Pop().
		Push(5).
		Push(3).
		Lt().
		Pop().
		Push(5).
		Push(5).
		Eq().
		Pop().
		Push(5).
		Push(3).
		Ne().
		Pop().
		Push(5).
		Push(5).
		Ge().
		Pop().
		Push(3).
		Push(5).
		Le().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) < 10 {
		t.Errorf("Expected at least 10 instructions, got %d", len(instructions))
	}
}

func TestBuilderMathOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(16).
		Sqrt().
		Push(0).
		Sin().
		Push(0).
		Cos().
		Push(0).
		Tan().
		Push(5).
		Push(10).
		Min().
		Push(5).
		Push(10).
		Max().
		Push(3.7).
		Floor().
		Push(3.2).
		Ceil().
		Push(3.5).
		Round().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	if len(program.Instructions()) < 15 {
		t.Errorf("Expected at least 15 instructions, got %d", len(program.Instructions()))
	}
}

func TestBuilderLabels(t *testing.T) {
	t.Run("Forward jump", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Push(1).
			Jmp("skip").
			Push(999).       // This should be skipped
			Label("skip").
			Push(2).
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		// Verify JMP instruction has correct operand (target address)
		instructions := program.Instructions()
		jmpInst := instructions[1]
		if jmpInst.Opcode != OpJMP {
			t.Errorf("Instruction 1 should be JMP, got %d", jmpInst.Opcode)
		}
		if jmpInst.Operand != 3 { // Should point to instruction 3 (label "skip")
			t.Errorf("JMP operand = %d, want 3", jmpInst.Operand)
		}
	})

	t.Run("Backward jump", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Label("loop").
			Push(1).
			Push(0).
			JmpZ("end").
			Jmp("loop").
			Label("end").
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		instructions := program.Instructions()
		backJmpInst := instructions[3] // Jmp("loop")
		if backJmpInst.Opcode != OpJMP {
			t.Errorf("Instruction 3 should be JMP, got %d", backJmpInst.Opcode)
		}
		if backJmpInst.Operand != 0 { // Should point back to instruction 0
			t.Errorf("JMP operand = %d, want 0", backJmpInst.Operand)
		}
	})

	t.Run("Conditional jumps", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Push(0).
			JmpZ("zero").
			Push(999).
			Label("zero").
			Push(1).
			JmpNZ("nonzero").
			Push(888).
			Label("nonzero").
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		instructions := program.Instructions()
		jmpzInst := instructions[1]
		jmpnzInst := instructions[4]

		if jmpzInst.Opcode != OpJMPZ {
			t.Errorf("Instruction 1 should be JMPZ, got %d", jmpzInst.Opcode)
		}
		if jmpnzInst.Opcode != OpJMPNZ {
			t.Errorf("Instruction 4 should be JMPNZ, got %d", jmpnzInst.Opcode)
		}
	})
}

func TestBuilderUnresolvedLabel(t *testing.T) {
	builder := NewProgramBuilder()
	_, err := builder.
		Push(1).
		Jmp("nonexistent").
		Halt().
		Build()

	if err == nil {
		t.Error("Build() should fail with unresolved label")
	}

	if err != nil && err.Error() != "unresolved label: nonexistent" {
		t.Errorf("Expected unresolved label error, got: %v", err)
	}
}

func TestBuilderCallAndReturn(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Call("function").
		Halt().
		Label("function").
		Push(42).
		Ret().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	callInst := instructions[0]
	retInst := instructions[3]

	if callInst.Opcode != OpCALL {
		t.Errorf("Instruction 0 should be CALL, got %d", callInst.Opcode)
	}
	if callInst.Operand != 2 { // Should point to "function" label
		t.Errorf("CALL operand = %d, want 2", callInst.Operand)
	}
	if retInst.Opcode != OpRET {
		t.Errorf("Instruction 3 should be RET, got %d", retInst.Opcode)
	}
}

func TestBuilderMemoryOperations(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(42).
		Store(0).
		Load(0).
		Push(1).
		Push(99).
		StoreD().
		Push(1).
		LoadD().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 9 {
		t.Errorf("Expected 9 instructions, got %d", len(instructions))
	}

	// Verify STORE has correct operand
	storeInst := instructions[1]
	if storeInst.Opcode != OpSTORE {
		t.Errorf("Instruction 1 should be STORE, got %d", storeInst.Opcode)
	}
	if storeInst.Operand != 0 {
		t.Errorf("STORE operand = %d, want 0", storeInst.Operand)
	}

	// Verify LOAD has correct operand
	loadInst := instructions[2]
	if loadInst.Opcode != OpLOAD {
		t.Errorf("Instruction 2 should be LOAD, got %d", loadInst.Opcode)
	}
	if loadInst.Operand != 0 {
		t.Errorf("LOAD operand = %d, want 0", loadInst.Operand)
	}
}

func TestBuilderCustomInstruction(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Push(5).
		Custom(128, 42). // Custom opcode with operand
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 3 {
		t.Errorf("Expected 3 instructions, got %d", len(instructions))
	}

	customInst := instructions[1]
	if customInst.Opcode != 128 {
		t.Errorf("Custom instruction opcode = %d, want 128", customInst.Opcode)
	}
	if customInst.Operand != 42 {
		t.Errorf("Custom instruction operand = %d, want 42", customInst.Operand)
	}
}

func TestBuilderIntegrationWithVM(t *testing.T) {
	t.Run("Simple arithmetic", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Push(10).
			Push(5).
			Add().
			Push(2).
			Mul().
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		vm := New()
		memory := NewSimpleMemory(0)
		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
		if !result.Halted {
			t.Error("Program should have halted")
		}
	})

	t.Run("Conditional jump", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Push(1).
			JmpZ("else").
			Push(100).
			Jmp("end").
			Label("else").
			Push(200).
			Label("end").
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		vm := New()
		memory := NewSimpleMemory(0)
		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}
		// Push(1), JmpZ pops it and doesn't jump (1 is truthy), Push(100)
		// Final stack: [100]
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}
	})

	t.Run("Memory operations", func(t *testing.T) {
		builder := NewProgramBuilder()
		program, err := builder.
			Push(42).
			Store(5).
			Load(5).
			Halt().
			Build()

		if err != nil {
			t.Fatalf("Build() failed: %v", err)
		}

		vm := New()
		memory := NewSimpleMemory(10)
		result, err := vm.Execute(program, memory, ExecuteOptions{})

		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}
		if result.StackDepth != 1 {
			t.Errorf("StackDepth = %d, want 1", result.StackDepth)
		}

		// Verify value was stored and loaded
		val, err := memory.Load(5)
		if err != nil {
			t.Fatalf("Memory Load failed: %v", err)
		}
		f, err := val.AsFloat()
		if err != nil {
			t.Fatalf("AsFloat() failed: %v", err)
		}
		if f != 42.0 {
			t.Errorf("Memory[5] = %f, want 42.0", f)
		}
	})
}

func TestBuilderComplexProgram(t *testing.T) {
	// Build a program that calculates factorial using a loop
	// This tests labels, jumps, and complex control flow
	builder := NewProgramBuilder()
	program, err := builder.
		PushInt(5).          // n = 5
		PushInt(1).          // result = 1
		Label("loop").
		Over().              // Copy n to top
		PushInt(1).
		Le().                // n <= 1?
		JmpNZ("done").       // If yes, done
		Over().              // Copy n
		Mul().               // result *= n
		Swap().              // Swap to get n on top
		Dec().               // n--
		Swap().              // Swap back
		Jmp("loop").
		Label("done").
		Swap().              // Get result on top
		Pop().               // Remove n
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	vm := New()
	memory := NewSimpleMemory(0)
	result, err := vm.Execute(program, memory, ExecuteOptions{
		MaxInstructions: 1000,
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}
	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
	if !result.Halted {
		t.Error("Program should have halted")
	}
}

func TestBuilderPushInt(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		PushInt(42).
		PushInt(-10).
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 3 {
		t.Errorf("Expected 3 instructions, got %d", len(instructions))
	}

	if instructions[0].Opcode != OpPUSHI {
		t.Errorf("First instruction should be PUSHI, got %d", instructions[0].Opcode)
	}
	if instructions[0].Operand != 42 {
		t.Errorf("First PUSHI operand = %d, want 42", instructions[0].Operand)
	}

	if instructions[1].Opcode != OpPUSHI {
		t.Errorf("Second instruction should be PUSHI, got %d", instructions[1].Opcode)
	}
	if instructions[1].Operand != -10 {
		t.Errorf("Second PUSHI operand = %d, want -10", instructions[1].Operand)
	}
}

func TestBuilderNop(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Nop().
		Push(1).
		Nop().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}

	if instructions[0].Opcode != OpNOP {
		t.Errorf("First instruction should be NOP, got %d", instructions[0].Opcode)
	}
	if instructions[2].Opcode != OpNOP {
		t.Errorf("Third instruction should be NOP, got %d", instructions[2].Opcode)
	}
}

func TestBuilderSymbolTable(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Label("start").
		Push(1).
		Label("loop").
		Dup().
		Push(10).
		Lt().
		JmpZ("end").
		Inc().
		Jmp("loop").
		Label("end").
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Verify symbol table was populated
	symbols := program.SymbolTable()
	if len(symbols) != 3 {
		t.Errorf("SymbolTable should have 3 entries, got %d", len(symbols))
	}

	// Check specific labels
	if symbols[0] != "start" {
		t.Errorf("symbols[0] = %s, want 'start'", symbols[0])
	}
	if symbols[1] != "loop" {
		t.Errorf("symbols[1] = %s, want 'loop'", symbols[1])
	}
}

func TestBuilderMetadata(t *testing.T) {
	metadata := ProgramMetadata{
		Name:        "test-program",
		Version:     "1.0",
		Author:      "tester",
		Description: "A test program",
	}

	builder := NewProgramBuilder()
	program, err := builder.
		SetMetadata(metadata).
		Push(1).
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	meta := program.Metadata()
	if meta.Name != "test-program" {
		t.Errorf("Metadata.Name = %s, want 'test-program'", meta.Name)
	}
	if meta.Version != "1.0" {
		t.Errorf("Metadata.Version = %s, want '1.0'", meta.Version)
	}
	if meta.Author != "tester" {
		t.Errorf("Metadata.Author = %s, want 'tester'", meta.Author)
	}
	if meta.Description != "A test program" {
		t.Errorf("Metadata.Description = %s, want 'A test program'", meta.Description)
	}
}
