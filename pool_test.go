package stackvm

import (
	"sync"
	"testing"
)

func TestNewVMPool(t *testing.T) {
	config := Config{
		StackSize: 512,
	}
	pool := NewVMPool(config)

	if pool == nil {
		t.Fatal("NewVMPool() returned nil")
	}
}

func TestNewDefaultVMPool(t *testing.T) {
	pool := NewDefaultVMPool()

	if pool == nil {
		t.Fatal("NewDefaultVMPool() returned nil")
	}
}

func TestVMPoolGetPut(t *testing.T) {
	pool := NewDefaultVMPool()

	// Get a VM
	vm := pool.Get()
	if vm == nil {
		t.Fatal("Get() returned nil")
	}

	// Return it
	pool.Put(vm)

	// Get another one (might be the same, might be new)
	vm2 := pool.Get()
	if vm2 == nil {
		t.Fatal("Get() returned nil on second call")
	}

	pool.Put(vm2)
}

func TestVMPoolExecute(t *testing.T) {
	pool := NewDefaultVMPool()

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

	memory := NewSimpleMemory(0)
	result, err := pool.Execute(program, memory, ExecuteOptions{})

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

func TestVMPoolExecuteFunc(t *testing.T) {
	pool := NewDefaultVMPool()

	// Create a simple program
	builder := NewProgramBuilder()
	program, err := builder.
		Push(42).
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	memory := NewSimpleMemory(0)

	var result *Result
	err = pool.ExecuteFunc(func(vm VM) error {
		var execErr error
		result, execErr = vm.Execute(program, memory, ExecuteOptions{})
		return execErr
	})

	if err != nil {
		t.Fatalf("ExecuteFunc() failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1", result.StackDepth)
	}
}

func TestVMPoolConcurrency(t *testing.T) {
	pool := NewDefaultVMPool()

	// Create a simple program
	builder := NewProgramBuilder()
	program, err := builder.
		Push(1).
		Push(1).
		Add().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Execute concurrently from multiple goroutines
	const goroutines = 100
	const execsPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, goroutines*execsPerGoroutine)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < execsPerGoroutine; j++ {
				memory := NewSimpleMemory(0)
				result, err := pool.Execute(program, memory, ExecuteOptions{})

				if err != nil {
					errors <- err
					return
				}

				if result.StackDepth != 1 {
					errors <- ErrStackUnderflow
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent execution error: %v", err)
	}
}

func TestVMPoolReset(t *testing.T) {
	pool := NewDefaultVMPool()

	// Get a VM and use it
	vm := pool.Get()

	builder := NewProgramBuilder()
	program, err := builder.
		Push(99).
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	memory := NewSimpleMemory(0)
	_, err = vm.Execute(program, memory, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Return it to pool
	pool.Put(vm)

	// Get it again (might be the same VM)
	vm2 := pool.Get()

	// Execute a different program - should work if reset properly
	builder2 := NewProgramBuilder()
	program2, err := builder2.
		Push(1).
		Push(2).
		Add().
		Halt().
		Build()

	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	memory2 := NewSimpleMemory(0)
	result, err := vm2.Execute(program2, memory2, ExecuteOptions{})
	if err != nil {
		t.Fatalf("Execute() failed after reuse: %v", err)
	}

	if result.StackDepth != 1 {
		t.Errorf("StackDepth = %d, want 1 (VM should be reset)", result.StackDepth)
	}

	pool.Put(vm2)
}

func TestVMPoolPutNil(t *testing.T) {
	pool := NewDefaultVMPool()

	// Putting nil should not panic
	pool.Put(nil)
}

func BenchmarkVMPoolGet(b *testing.B) {
	pool := NewDefaultVMPool()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm := pool.Get()
		pool.Put(vm)
	}
}

func BenchmarkVMPoolExecute(b *testing.B) {
	pool := NewDefaultVMPool()

	builder := NewProgramBuilder()
	program, err := builder.
		Push(10).
		Push(5).
		Add().
		Halt().
		Build()

	if err != nil {
		b.Fatalf("Build() failed: %v", err)
	}

	memory := NewSimpleMemory(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pool.Execute(program, memory, ExecuteOptions{})
		if err != nil {
			b.Fatalf("Execute() failed: %v", err)
		}
	}
}

func BenchmarkVMPoolParallel(b *testing.B) {
	pool := NewDefaultVMPool()

	builder := NewProgramBuilder()
	program, err := builder.
		Push(10).
		Push(5).
		Add().
		Halt().
		Build()

	if err != nil {
		b.Fatalf("Build() failed: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		memory := NewSimpleMemory(0)
		for pb.Next() {
			_, err := pool.Execute(program, memory, ExecuteOptions{})
			if err != nil {
				b.Fatalf("Execute() failed: %v", err)
			}
		}
	})
}
