package stackvm

import (
	"testing"
)

// mockHandler is a mock implementation of InstructionHandler for testing.
type mockHandler struct {
	name string
	fn   func(ExecutionContext, int32) error
}

func (m *mockHandler) Execute(ctx ExecutionContext, operand int32) error {
	if m.fn != nil {
		return m.fn(ctx, operand)
	}
	return nil
}

func (m *mockHandler) Name() string {
	return m.name
}

func TestNewInstructionRegistry(t *testing.T) {
	registry := NewInstructionRegistry()
	if registry == nil {
		t.Fatal("NewInstructionRegistry() returned nil")
	}

	// Should start empty
	if len(registry.List()) != 0 {
		t.Errorf("New registry should be empty, got %d handlers", len(registry.List()))
	}
}

func TestRegisterCustomOpcode(t *testing.T) {
	registry := NewInstructionRegistry()
	handler := &mockHandler{name: "TEST"}

	// Register a custom opcode (128-255 range)
	err := registry.Register(128, handler)
	if err != nil {
		t.Errorf("Register(128) failed: %v", err)
	}

	// Verify it was registered
	retrieved, exists := registry.Get(128)
	if !exists {
		t.Error("Get(128) returned false, want true")
	}
	if retrieved.Name() != "TEST" {
		t.Errorf("Handler name = %s, want TEST", retrieved.Name())
	}
}

func TestRegisterStandardOpcodeError(t *testing.T) {
	registry := NewInstructionRegistry()
	handler := &mockHandler{name: "INVALID"}

	// Try to register a standard opcode (should fail)
	err := registry.Register(OpPUSH, handler)
	if err == nil {
		t.Error("Register(OpPUSH) should fail for standard opcodes")
	}

	// Try various standard opcodes
	standardOpcodes := []Opcode{0, 1, 50, 127}
	for _, opcode := range standardOpcodes {
		err := registry.Register(opcode, handler)
		if err == nil {
			t.Errorf("Register(%d) should fail for standard opcode", opcode)
		}
	}
}

func TestRegisterDuplicateOpcode(t *testing.T) {
	registry := NewInstructionRegistry()
	handler1 := &mockHandler{name: "FIRST"}
	handler2 := &mockHandler{name: "SECOND"}

	// Register first handler
	err := registry.Register(200, handler1)
	if err != nil {
		t.Fatalf("First Register(200) failed: %v", err)
	}

	// Try to register second handler for same opcode (should fail)
	err = registry.Register(200, handler2)
	if err == nil {
		t.Error("Register(200) should fail when opcode already registered")
	}

	// Verify first handler is still there
	retrieved, _ := registry.Get(200)
	if retrieved.Name() != "FIRST" {
		t.Errorf("Handler name = %s, want FIRST", retrieved.Name())
	}
}

func TestUnregister(t *testing.T) {
	registry := NewInstructionRegistry()
	handler := &mockHandler{name: "TEMP"}

	// Register and then unregister
	registry.Register(150, handler)
	err := registry.Unregister(150)
	if err != nil {
		t.Errorf("Unregister(150) failed: %v", err)
	}

	// Verify it's gone
	_, exists := registry.Get(150)
	if exists {
		t.Error("Get(150) returned true after unregister, want false")
	}
}

func TestUnregisterNonexistent(t *testing.T) {
	registry := NewInstructionRegistry()

	// Try to unregister an opcode that was never registered
	err := registry.Unregister(200)
	if err == nil {
		t.Error("Unregister(200) should fail for non-existent opcode")
	}
}

func TestList(t *testing.T) {
	registry := NewInstructionRegistry()

	// Register multiple handlers
	registry.Register(128, &mockHandler{name: "ONE"})
	registry.Register(129, &mockHandler{name: "TWO"})
	registry.Register(200, &mockHandler{name: "THREE"})

	opcodes := registry.List()
	if len(opcodes) != 3 {
		t.Errorf("List() returned %d opcodes, want 3", len(opcodes))
	}

	// Verify all opcodes are in the list
	opcodeSet := make(map[Opcode]bool)
	for _, op := range opcodes {
		opcodeSet[op] = true
	}

	for _, expected := range []Opcode{128, 129, 200} {
		if !opcodeSet[expected] {
			t.Errorf("List() missing opcode %d", expected)
		}
	}
}

func TestNames(t *testing.T) {
	registry := NewInstructionRegistry()

	// Register handlers with different names
	registry.Register(128, &mockHandler{name: "ALPHA"})
	registry.Register(129, &mockHandler{name: "BETA"})
	registry.Register(200, &mockHandler{name: "GAMMA"})

	names := registry.Names()
	if len(names) != 3 {
		t.Errorf("Names() returned %d entries, want 3", len(names))
	}

	// Verify names
	if names[128] != "ALPHA" {
		t.Errorf("names[128] = %s, want ALPHA", names[128])
	}
	if names[129] != "BETA" {
		t.Errorf("names[129] = %s, want BETA", names[129])
	}
	if names[200] != "GAMMA" {
		t.Errorf("names[200] = %s, want GAMMA", names[200])
	}
}

func TestCustomInstructionExecution(t *testing.T) {
	registry := NewInstructionRegistry()

	// Create a custom instruction that doubles the top value
	doubleHandler := &mockHandler{
		name: "DOUBLE",
		fn: func(ctx ExecutionContext, operand int32) error {
			val, err := ctx.Pop()
			if err != nil {
				return err
			}
			f, err := val.AsFloat()
			if err != nil {
				// Try int
				i, err := val.AsInt()
				if err != nil {
					return err
				}
				f = float64(i)
			}
			return ctx.Push(FloatValue(f * 2))
		},
	}

	err := registry.Register(128, doubleHandler)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Create VM with custom registry
	vmWithRegistry := NewWithConfig(Config{
		StackSize:           256,
		InstructionRegistry: registry,
	})

	// Test program: push 5, double it (custom instruction), halt
	program := NewProgram([]Instruction{
		NewInstruction(OpPUSH, 5),
		NewInstruction(128, 0), // Custom DOUBLE instruction
		NewInstruction(OpHALT, 0),
	})

	memory := NewSimpleMemory(0)
	result, err := vmWithRegistry.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
}

func TestCustomInstructionPushPop(t *testing.T) {
	registry := NewInstructionRegistry()

	// Create a custom instruction that pushes the operand value
	pushOperandHandler := &mockHandler{
		name: "PUSHOP",
		fn: func(ctx ExecutionContext, operand int32) error {
			return ctx.Push(IntValue(int64(operand)))
		},
	}

	registry.Register(200, pushOperandHandler)

	vmWithRegistry := NewWithConfig(Config{
		StackSize:           256,
		InstructionRegistry: registry,
	})

	// Test program: use custom instruction to push 42
	program := NewProgram([]Instruction{
		NewInstruction(200, 42), // Custom instruction with operand
		NewInstruction(OpHALT, 0),
	})

	memory := NewSimpleMemory(0)
	result, err := vmWithRegistry.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
}

func TestRegistryConcurrency(t *testing.T) {
	registry := NewInstructionRegistry()

	// Test concurrent registration
	done := make(chan bool)
	for i := 128; i < 138; i++ {
		go func(opcode Opcode) {
			handler := &mockHandler{name: "CONCURRENT"}
			registry.Register(opcode, handler)
			done <- true
		}(Opcode(i))
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all were registered
	opcodes := registry.List()
	if len(opcodes) != 10 {
		t.Errorf("List() returned %d opcodes, want 10", len(opcodes))
	}
}
