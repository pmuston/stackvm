package stackvm

import (
	"sync"
)

// VMPool manages a pool of reusable VM instances.
// This is useful for high-throughput scenarios where creating new VMs
// for each execution would be expensive.
type VMPool struct {
	pool   sync.Pool
	config Config
}

// NewVMPool creates a new VM pool with the given configuration.
// All VMs in the pool will be created with this configuration.
func NewVMPool(config Config) *VMPool {
	return &VMPool{
		config: config,
		pool: sync.Pool{
			New: func() interface{} {
				return NewWithConfig(config)
			},
		},
	}
}

// NewDefaultVMPool creates a VM pool with default configuration.
func NewDefaultVMPool() *VMPool {
	return NewVMPool(Config{
		StackSize: 256,
	})
}

// Get retrieves a VM from the pool.
// The VM is reset before being returned.
// The caller must call Put() when done with the VM.
func (p *VMPool) Get() VM {
	vm := p.pool.Get().(VM)
	vm.Reset()
	return vm
}

// Put returns a VM to the pool.
// The VM is reset before being added back to the pool.
func (p *VMPool) Put(vm VM) {
	if vm == nil {
		return
	}
	vm.Reset()
	p.pool.Put(vm)
}

// Execute is a convenience method that gets a VM from the pool,
// executes the program, and returns the VM to the pool.
// This is safe for concurrent use.
func (p *VMPool) Execute(program Program, memory Memory, opts ExecuteOptions) (*Result, error) {
	vm := p.Get()
	defer p.Put(vm)
	return vm.Execute(program, memory, opts)
}

// ExecuteFunc executes a function with a VM from the pool.
// The VM is automatically returned to the pool when the function completes.
// This is useful for more complex execution scenarios.
func (p *VMPool) ExecuteFunc(fn func(VM) error) error {
	vm := p.Get()
	defer p.Put(vm)
	return fn(vm)
}
