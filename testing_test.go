package stackvm

import (
	"testing"
)

func TestNewTestRunner(t *testing.T) {
	runner := NewTestRunner(t)
	if runner == nil {
		t.Fatal("NewTestRunner() returned nil")
	}
}

func TestTestRunnerAssembleAndRun(t *testing.T) {
	runner := NewTestRunner(t)

	source := `
		PUSH 10
		PUSH 5
		ADD
		HALT
	`

	result := runner.AssembleAndRun(source)

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}

	if !result.Halted {
		t.Error("Program should have halted")
	}
}

func TestTestRunnerExpectStackDepth(t *testing.T) {
	runner := NewTestRunner(t)

	source := `
		PUSH 1
		PUSH 2
		PUSH 3
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectStackDepth(result, 3)
	runner.ExpectHalted(result)
}

func TestTestRunnerMemoryOperations(t *testing.T) {
	runner := NewTestRunner(t)

	source := `
		PUSH 42
		STORE 0
		PUSH 99
		STORE 1
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectHalted(result)
	runner.ExpectMemoryValue(0, 42)
	runner.ExpectMemoryValue(1, 99)
}

func TestTestRunnerMemoryInt(t *testing.T) {
	runner := NewTestRunner(t)

	source := `
		PUSHI 42
		STORE 0
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectHalted(result)
	runner.ExpectMemoryInt(0, 42)
}

func TestTestRunnerReset(t *testing.T) {
	runner := NewTestRunner(t)

	source1 := `
		PUSH 10
		STORE 0
		HALT
	`

	runner.AssembleAndRun(source1)
	runner.ExpectMemoryValue(0, 10)

	runner.Reset()

	source2 := `
		PUSH 20
		STORE 0
		HALT
	`

	runner.AssembleAndRun(source2)
	runner.ExpectMemoryValue(0, 20)
}

func TestRunProgramTests(t *testing.T) {
	tests := []ProgramTest{
		{
			Name: "Simple addition",
			Source: `
				PUSH 10
				PUSH 5
				ADD
				HALT
			`,
			ExpectedStack: 1,
		},
		{
			Name: "Memory test",
			Source: `
				PUSH 42
				STORE 5
				LOAD 5
				HALT
			`,
			ExpectedStack: 1,
			ExpectedMemory: map[int]float64{
				5: 42,
			},
		},
		{
			Name: "Arithmetic chain",
			Source: `
				PUSH 10
				PUSH 5
				ADD
				PUSH 2
				MUL
				HALT
			`,
			ExpectedStack: 1,
		},
	}

	RunProgramTests(t, tests)
}

func TestAssertProgramOutput(t *testing.T) {
	source := `
		PUSH 1
		PUSH 2
		PUSH 3
		HALT
	`

	AssertProgramOutput(t, source, 3)
}

func TestMustAssemble(t *testing.T) {
	source := `
		PUSH 10
		HALT
	`

	program := MustAssemble(source)
	if program == nil {
		t.Fatal("MustAssemble() returned nil")
	}

	instructions := program.Instructions()
	if len(instructions) != 2 {
		t.Errorf("Expected 2 instructions, got %d", len(instructions))
	}
}

func TestMustAssemblePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAssemble() should panic on invalid source")
		}
	}()

	// Invalid source should panic
	MustAssemble("INVALID OPCODE")
}

func TestMustAssembleFile(t *testing.T) {
	program := MustAssembleFile("testdata/programs/simple_add.asm")
	if program == nil {
		t.Fatal("MustAssembleFile() returned nil")
	}

	instructions := program.Instructions()
	if len(instructions) != 4 {
		t.Errorf("Expected 4 instructions, got %d", len(instructions))
	}
}

func TestBenchmarkProgram(t *testing.T) {
	source := `
		PUSH 10
		PUSH 5
		ADD
		HALT
	`

	// Just verify it doesn't panic
	result := testing.Benchmark(func(b *testing.B) {
		BenchmarkProgram(b, source)
	})

	if result.N == 0 {
		t.Error("Benchmark didn't run")
	}
}

func TestTestRunnerWithCustomInstructions(t *testing.T) {
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

	runner := NewTestRunner(t)
	runner.SetRegistry(registry)

	source := `
		PUSH 5
		DOUBLE
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectStackDepth(result, 1)
	runner.ExpectHalted(result)
}

func TestTestRunnerSetMemory(t *testing.T) {
	runner := NewTestRunner(t)

	// Create custom memory with initial values
	memory := NewSimpleMemory(10)
	memory.Store(0, FloatValue(100))

	runner.SetMemory(memory)

	source := `
		LOAD 0
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectStackDepth(result, 1)
	runner.ExpectMemoryValue(0, 100)
}

func TestNewTestRunnerWithConfig(t *testing.T) {
	config := Config{
		StackSize: 1024,
	}

	runner := NewTestRunnerWithConfig(t, config)
	if runner == nil {
		t.Fatal("NewTestRunnerWithConfig() returned nil")
	}

	// Test with larger stack
	source := `
		PUSH 1
		PUSH 2
		PUSH 3
		PUSH 4
		PUSH 5
		HALT
	`

	result := runner.AssembleAndRun(source)
	runner.ExpectStackDepth(result, 5)
}
