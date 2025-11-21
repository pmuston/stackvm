package stackvm

import (
	"fmt"
	"sync"
)

// instructionRegistry implements the InstructionRegistry interface.
type instructionRegistry struct {
	mu       sync.RWMutex
	handlers map[Opcode]InstructionHandler
}

// NewInstructionRegistry creates a new instruction registry.
func NewInstructionRegistry() InstructionRegistry {
	return &instructionRegistry{
		handlers: make(map[Opcode]InstructionHandler),
	}
}

// Register adds a handler for a custom opcode (128-255).
// Returns an error if the opcode is in the standard range (0-127) or already registered.
func (r *instructionRegistry) Register(opcode Opcode, handler InstructionHandler) error {
	if opcode < 128 {
		return fmt.Errorf("cannot register standard opcode %d: reserved for built-in instructions", opcode)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[opcode]; exists {
		return fmt.Errorf("opcode %d already registered", opcode)
	}

	r.handlers[opcode] = handler
	return nil
}

// Unregister removes a handler for an opcode.
// Returns an error if the opcode is not registered.
func (r *instructionRegistry) Unregister(opcode Opcode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[opcode]; !exists {
		return fmt.Errorf("opcode %d not registered", opcode)
	}

	delete(r.handlers, opcode)
	return nil
}

// Get retrieves a handler for an opcode.
// Returns false if the opcode is not registered.
func (r *instructionRegistry) Get(opcode Opcode) (InstructionHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[opcode]
	return handler, exists
}

// List returns all registered custom opcodes.
func (r *instructionRegistry) List() []Opcode {
	r.mu.RLock()
	defer r.mu.RUnlock()

	opcodes := make([]Opcode, 0, len(r.handlers))
	for opcode := range r.handlers {
		opcodes = append(opcodes, opcode)
	}
	return opcodes
}

// Names returns a mapping of opcodes to their names.
func (r *instructionRegistry) Names() map[Opcode]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make(map[Opcode]string, len(r.handlers))
	for opcode, handler := range r.handlers {
		names[opcode] = handler.Name()
	}
	return names
}

// executionContextImpl implements the ExecutionContext interface.
// This is used by custom instruction handlers to interact with the VM.
type executionContextImpl struct {
	vm       *executor
	memory   Memory
	userData map[string]interface{}
}

// newExecutionContext creates a new execution context.
func newExecutionContext(vm *executor, memory Memory) *executionContextImpl {
	return &executionContextImpl{
		vm:       vm,
		memory:   memory,
		userData: make(map[string]interface{}),
	}
}

// Push adds a value to the top of the stack.
func (ctx *executionContextImpl) Push(value Value) error {
	// Get the max stack depth from the VM config
	maxDepth := ctx.vm.config.StackSize
	if len(ctx.vm.stack) >= maxDepth {
		return ErrStackOverflow
	}
	ctx.vm.stack = append(ctx.vm.stack, value)
	return nil
}

// Pop removes and returns the value from the top of the stack.
func (ctx *executionContextImpl) Pop() (Value, error) {
	if len(ctx.vm.stack) == 0 {
		return NilValue(), ErrStackUnderflow
	}
	val := ctx.vm.stack[len(ctx.vm.stack)-1]
	ctx.vm.stack = ctx.vm.stack[:len(ctx.vm.stack)-1]
	return val, nil
}

// Peek returns the value at the top of the stack without removing it.
func (ctx *executionContextImpl) Peek() (Value, error) {
	if len(ctx.vm.stack) == 0 {
		return NilValue(), ErrStackUnderflow
	}
	return ctx.vm.stack[len(ctx.vm.stack)-1], nil
}

// PeekN returns the nth value from the top of the stack (0 = top).
func (ctx *executionContextImpl) PeekN(n int) (Value, error) {
	if n < 0 || n >= len(ctx.vm.stack) {
		return NilValue(), ErrStackUnderflow
	}
	return ctx.vm.stack[len(ctx.vm.stack)-1-n], nil
}

// StackDepth returns the current number of values on the stack.
func (ctx *executionContextImpl) StackDepth() int {
	return len(ctx.vm.stack)
}

// PC returns the current program counter value.
func (ctx *executionContextImpl) PC() int {
	return ctx.vm.pc
}

// SetPC sets the program counter to the specified value.
func (ctx *executionContextImpl) SetPC(pc int) {
	ctx.vm.pc = pc
}

// Jump sets the program counter to the specified offset.
func (ctx *executionContextImpl) Jump(offset int) {
	ctx.vm.pc = offset
}

// Memory returns the memory provider associated with this execution.
func (ctx *executionContextImpl) Memory() Memory {
	return ctx.memory
}

// InstructionCount returns the number of instructions executed so far.
func (ctx *executionContextImpl) InstructionCount() uint32 {
	return ctx.vm.instrCount
}

// IncrementInstructionCount increments the instruction counter by one.
func (ctx *executionContextImpl) IncrementInstructionCount() {
	ctx.vm.instrCount++
}

// Halt stops execution.
func (ctx *executionContextImpl) Halt() {
	ctx.vm.halted = true
}

// IsHalted returns true if execution has been halted.
func (ctx *executionContextImpl) IsHalted() bool {
	return ctx.vm.halted
}

// UserData returns a map for storing and retrieving custom execution context data.
func (ctx *executionContextImpl) UserData() map[string]interface{} {
	return ctx.userData
}
