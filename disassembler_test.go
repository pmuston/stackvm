package stackvm

import (
	"strings"
	"testing"
)

func TestNewDisassembler(t *testing.T) {
	disasm := NewDisassembler()
	if disasm == nil {
		t.Fatal("NewDisassembler() returned nil")
	}
}

func TestDisassembleSimple(t *testing.T) {
	// Create a simple program
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

	disasm := NewDisassembler()
	output, err := disasm.Disassemble(program)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Verify output contains expected instructions
	expectedInstructions := []string{"PUSH 10", "PUSH 5", "ADD", "HALT"}
	for _, instr := range expectedInstructions {
		if !strings.Contains(output, instr) {
			t.Errorf("Output missing instruction: %s\nOutput:\n%s", instr, output)
		}
	}
}

func TestDisassembleWithLabels(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Label("START").
		Push(1).
		JmpZ("END").
		Push(100).
		Label("END").
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	disasm := NewDisassembler()
	output, err := disasm.Disassemble(program)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Verify labels are present
	if !strings.Contains(output, "START:") {
		t.Error("Output missing START label")
	}
	if !strings.Contains(output, "END:") {
		t.Error("Output missing END label")
	}
}

func TestDisassembleWithMetadata(t *testing.T) {
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

	disasm := NewDisassembler()
	output, err := disasm.Disassemble(program)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Verify metadata is in output
	if !strings.Contains(output, "test-program") {
		t.Error("Output missing program name")
	}
	if !strings.Contains(output, "1.0") {
		t.Error("Output missing version")
	}
}

func TestDisassembleAllOpcodes(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		// Stack
		Push(1).
		PushInt(2).
		Dup().
		Pop().
		Swap().
		Over().
		Rot().
		Pop().
		Pop().
		Pop().
		// Arithmetic
		Push(10).
		Push(5).
		Add().
		Sub().
		Mul().
		Div().
		Mod().
		Neg().
		Abs().
		Inc().
		Dec().
		// Logic
		Push(1).
		Push(1).
		And().
		Or().
		Not().
		Xor().
		// Comparison
		Push(5).
		Push(3).
		Gt().
		Pop().
		Push(3).
		Push(5).
		Lt().
		Pop().
		// Math
		Push(16).
		Sqrt().
		Push(0).
		Sin().
		Pop().
		// Control
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	disasm := NewDisassembler()
	output, err := disasm.Disassemble(program)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Just verify it produces output
	if len(output) == 0 {
		t.Error("Disassemble() produced empty output")
	}
}

func TestDisassembleAndReassemble(t *testing.T) {
	// Create a program
	source := `
		PUSH 10
		PUSH 5
		ADD
		HALT
	`

	asm := NewAssembler()
	program1, err := asm.Assemble(source)
	if err != nil {
		t.Fatalf("Assemble() failed: %v", err)
	}

	// Disassemble it
	disasm := NewDisassembler()
	disassembled, err := disasm.Disassemble(program1)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Reassemble the disassembled code
	program2, err := asm.Assemble(disassembled)
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	// Verify both programs have the same instructions
	instr1 := program1.Instructions()
	instr2 := program2.Instructions()

	if len(instr1) != len(instr2) {
		t.Fatalf("Instruction count mismatch: %d vs %d", len(instr1), len(instr2))
	}

	for i := range instr1 {
		if instr1[i].Opcode != instr2[i].Opcode {
			t.Errorf("Instruction %d opcode mismatch: %d vs %d",
				i, instr1[i].Opcode, instr2[i].Opcode)
		}
	}
}

func TestDisassembleCustomInstructions(t *testing.T) {
	// Create a custom instruction
	registry := NewInstructionRegistry()

	testHandler := &testInstructionHandler{
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

	err := registry.Register(128, testHandler)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	// Build a program using the custom instruction
	builder := NewProgramBuilder()
	program, err := builder.
		Push(5).
		Custom(128, 0).
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Disassemble with registry
	disasm := NewDisassembler()
	disasm.SetRegistry(registry)

	output, err := disasm.Disassemble(program)
	if err != nil {
		t.Fatalf("Disassemble() failed: %v", err)
	}

	// Verify custom instruction name appears
	if !strings.Contains(output, "DOUBLE") {
		t.Errorf("Output missing custom instruction name\nOutput:\n%s", output)
	}
}

func TestDisassemblerOptions(t *testing.T) {
	builder := NewProgramBuilder()
	program, err := builder.
		Label("START").
		Push(1).
		Push(2).
		Add().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	t.Run("With addresses", func(t *testing.T) {
		disasm := NewDisassemblerWithOptions(DisassemblerOptions{
			IncludeAddresses:   true,
			IncludeMetadata:    false,
			IndentInstructions: false,
		})

		output, err := disasm.Disassemble(program)
		if err != nil {
			t.Fatalf("Disassemble() failed: %v", err)
		}

		// Should contain address markers
		if !strings.Contains(output, "[") {
			t.Error("Output missing address markers")
		}
	})

	t.Run("Without indentation", func(t *testing.T) {
		disasm := NewDisassemblerWithOptions(DisassemblerOptions{
			IncludeAddresses:   false,
			IncludeMetadata:    false,
			IndentInstructions: false,
		})

		output, err := disasm.Disassemble(program)
		if err != nil {
			t.Fatalf("Disassemble() failed: %v", err)
		}

		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "    ") && !strings.HasSuffix(line, ":") {
				t.Error("Found indented instruction when indentation disabled")
				break
			}
		}
	})
}
