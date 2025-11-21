package stackvm

import (
	"fmt"
	"testing"
)

// TestRunner provides utilities for testing VM programs.
type TestRunner struct {
	vm       VM
	memory   Memory
	t        *testing.T
	registry InstructionRegistry
}

// NewTestRunner creates a new test runner.
func NewTestRunner(t *testing.T) *TestRunner {
	return &TestRunner{
		vm:     New(),
		memory: NewSimpleMemory(256),
		t:      t,
	}
}

// NewTestRunnerWithConfig creates a test runner with custom configuration.
func NewTestRunnerWithConfig(t *testing.T, config Config) *TestRunner {
	return &TestRunner{
		vm:     NewWithConfig(config),
		memory: NewSimpleMemory(256),
		t:      t,
	}
}

// SetRegistry sets the instruction registry for custom opcodes.
func (tr *TestRunner) SetRegistry(registry InstructionRegistry) {
	tr.registry = registry
	tr.vm = NewWithConfig(Config{
		StackSize:           256,
		InstructionRegistry: registry,
	})
}

// SetMemory sets the memory for the test runner.
func (tr *TestRunner) SetMemory(memory Memory) {
	tr.memory = memory
}

// AssembleAndRun assembles source code and executes it.
// Returns the result or fails the test.
func (tr *TestRunner) AssembleAndRun(source string, opts ...ExecuteOptions) *Result {
	tr.t.Helper()

	asm := NewAssembler()
	if tr.registry != nil {
		asm.SetRegistry(tr.registry)
	}

	program, err := asm.Assemble(source)
	if err != nil {
		tr.t.Fatalf("Failed to assemble: %v", err)
	}

	return tr.Run(program, opts...)
}

// Run executes a program and returns the result or fails the test.
func (tr *TestRunner) Run(program Program, opts ...ExecuteOptions) *Result {
	tr.t.Helper()

	var executeOpts ExecuteOptions
	if len(opts) > 0 {
		executeOpts = opts[0]
	}

	// Set default max instructions if not specified
	if executeOpts.MaxInstructions == 0 {
		executeOpts.MaxInstructions = 10000
	}

	result, err := tr.vm.Execute(program, tr.memory, executeOpts)
	if err != nil {
		tr.t.Fatalf("Execution failed: %v", err)
	}

	return result
}

// ExpectStackDepth verifies the stack depth.
func (tr *TestRunner) ExpectStackDepth(result *Result, expected int) {
	tr.t.Helper()
	if result.StackDepth != expected {
		tr.t.Errorf("Stack depth = %d, want %d", result.StackDepth, expected)
	}
}

// ExpectHalted verifies the program halted.
func (tr *TestRunner) ExpectHalted(result *Result) {
	tr.t.Helper()
	if !result.Halted {
		tr.t.Error("Expected program to halt")
	}
}

// ExpectMemoryValue verifies a memory value.
func (tr *TestRunner) ExpectMemoryValue(index int, expectedValue float64) {
	tr.t.Helper()

	val, err := tr.memory.Load(index)
	if err != nil {
		tr.t.Fatalf("Failed to load memory[%d]: %v", index, err)
	}

	f, err := val.AsFloat()
	if err != nil {
		tr.t.Fatalf("Memory[%d] is not a float: %v", index, err)
	}

	if f != expectedValue {
		tr.t.Errorf("Memory[%d] = %f, want %f", index, f, expectedValue)
	}
}

// ExpectMemoryInt verifies an integer memory value.
func (tr *TestRunner) ExpectMemoryInt(index int, expectedValue int64) {
	tr.t.Helper()

	val, err := tr.memory.Load(index)
	if err != nil {
		tr.t.Fatalf("Failed to load memory[%d]: %v", index, err)
	}

	i, err := val.AsInt()
	if err != nil {
		tr.t.Fatalf("Memory[%d] is not an int: %v", index, err)
	}

	if i != expectedValue {
		tr.t.Errorf("Memory[%d] = %d, want %d", index, i, expectedValue)
	}
}

// Reset resets the VM and memory for the next test.
func (tr *TestRunner) Reset() {
	tr.vm.Reset()
	if sm, ok := tr.memory.(*SimpleMemory); ok {
		sm.Reset()
	}
}

// ProgramTest represents a test case for a program.
type ProgramTest struct {
	Name           string
	Source         string
	ExpectedStack  int
	ExpectedMemory map[int]float64
	ExpectedError  bool
	Options        ExecuteOptions
}

// RunProgramTests runs a suite of program tests.
func RunProgramTests(t *testing.T, tests []ProgramTest) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runner := NewTestRunner(t)

			// Assemble
			asm := NewAssembler()
			program, err := asm.Assemble(tt.Source)

			if tt.ExpectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to assemble: %v", err)
			}

			// Execute
			opts := tt.Options
			if opts.MaxInstructions == 0 {
				opts.MaxInstructions = 10000
			}

			result, err := runner.vm.Execute(program, runner.memory, opts)
			if err != nil {
				t.Fatalf("Execution failed: %v", err)
			}

			// Verify stack depth
			if result.StackDepth != tt.ExpectedStack {
				t.Errorf("Stack depth = %d, want %d", result.StackDepth, tt.ExpectedStack)
			}

			// Verify memory values
			for index, expected := range tt.ExpectedMemory {
				runner.ExpectMemoryValue(index, expected)
			}

			// Verify halted
			if !result.Halted {
				t.Error("Expected program to halt")
			}
		})
	}
}

// BenchmarkProgram benchmarks program execution.
func BenchmarkProgram(b *testing.B, source string) {
	asm := NewAssembler()
	program, err := asm.Assemble(source)
	if err != nil {
		b.Fatalf("Failed to assemble: %v", err)
	}

	pool := NewDefaultVMPool()
	memory := NewSimpleMemory(256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pool.Execute(program, memory, ExecuteOptions{
			MaxInstructions: 10000,
		})
		if err != nil {
			b.Fatalf("Execution failed: %v", err)
		}
	}
}

// AssertProgramOutput is a helper for quick program testing.
func AssertProgramOutput(t *testing.T, source string, expectedStackDepth int) {
	t.Helper()

	runner := NewTestRunner(t)
	result := runner.AssembleAndRun(source)

	runner.ExpectStackDepth(result, expectedStackDepth)
	runner.ExpectHalted(result)
}

// MustAssemble assembles source or panics. Useful for test setup.
func MustAssemble(source string) Program {
	asm := NewAssembler()
	program, err := asm.Assemble(source)
	if err != nil {
		panic(fmt.Sprintf("Failed to assemble: %v", err))
	}
	return program
}

// MustAssembleFile assembles a file or panics. Useful for test setup.
func MustAssembleFile(path string) Program {
	asm := NewAssembler()
	program, err := asm.AssembleFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to assemble file: %v", err))
	}
	return program
}
