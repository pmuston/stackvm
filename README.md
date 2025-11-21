# StackVM

A standalone, reusable stack-based virtual machine implementation in Go. StackVM provides a complete bytecode execution environment with zero external dependencies, designed to be embedded in other applications or used standalone via the included assembler CLI tool.

## Features

- **Stack-Based Architecture**: Simple, powerful execution model similar to JVM or Forth
- **Zero Dependencies**: Pure Go implementation using only the standard library
- **Extensible**: Support for custom instructions and memory providers
- **Complete Tooling**: Assembler, disassembler, and execution runtime
- **Well-Documented**: Comprehensive specification and getting started guides
- **Production Ready**: Full test coverage, benchmarks, and VM pooling support
- **Safe Execution**: Configurable limits, bounds checking, and comprehensive error handling

## Quick Start

### Installation

```bash
go get github.com/pmuston/stackvm
```

### Hello World

Create a file `hello.asm`:

```assembly
; Add two numbers
PUSH 10
PUSH 5
ADD
HALT
```

Assemble and run:

```bash
stackvm-asm -r hello.asm
```

### Using as a Library

```go
package main

import (
    "fmt"
    "github.com/pmuston/stackvm"
)

func main() {
    // Create a VM
    vm := stackvm.New()

    // Build a program using the fluent builder
    program := stackvm.NewProgramBuilder().
        Push(10).
        Push(5).
        Add().
        Halt().
        MustBuild()

    // Create memory
    memory := stackvm.NewSimpleMemory(10)

    // Execute
    result, err := vm.Execute(program, memory, stackvm.ExecuteOptions{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("Executed %d instructions in %v\n",
        result.InstructionCount, result.ExecutionTime)
}
```

## Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Tutorial with examples for learning StackVM assembly
- **[Language Specification](docs/LANGUAGE_SPEC.md)** - Complete assembly language reference
- **[Technical Specification](docs/SPEC.md)** - Complete VM architecture and API documentation
- **[CLI Tool Specification](docs/CLI_SPEC.md)** - Command-line tool usage

## Instruction Set

StackVM provides over 80 built-in instructions including:

### Stack Operations
`PUSH`, `PUSHI`, `POP`, `DUP`, `SWAP`, `OVER`, `ROT`

### Arithmetic
`ADD`, `SUB`, `MUL`, `DIV`, `MOD`, `NEG`, `ABS`, `INC`, `DEC`

### Logic & Comparison
`AND`, `OR`, `NOT`, `XOR`, `EQ`, `NE`, `GT`, `LT`, `GE`, `LE`

### Control Flow
`JMP`, `JMPZ`, `JMPNZ`, `CALL`, `RET`, `HALT`, `NOP`

### Memory
`LOAD`, `STORE`, `LOADD`, `STORED`

### Math Functions
`SQRT`, `SIN`, `COS`, `TAN`, `LOG`, `EXP`, `POW`, `MIN`, `MAX`, `FLOOR`, `CEIL`, `ROUND`, and more

### Custom Instructions
Opcodes 128-255 are reserved for host-defined custom instructions

## CLI Tool

The `stackvm-asm` command-line tool provides complete program development workflow:

```bash
# Assemble a program
stackvm-asm program.asm -o program.bin

# Run a program
stackvm-asm -r program.asm

# Run with statistics
stackvm-asm -r -s program.asm

# Run with memory initialization
stackvm-asm -r -M "0=10,1=5" program.asm

# Disassemble bytecode
stackvm-asm -d program.bin

# Show program information
stackvm-asm info program.bin
```

## Project Structure

```
stackvm/
├── pkg/stackvm/          # Public API (not used - flat structure)
├── internal/             # Internal implementation
│   ├── executor/         # VM execution engine
│   ├── opcodes/          # Standard instruction handlers
│   ├── asm/              # Assembler implementation
│   └── encoding/         # Binary format implementation
├── cmd/
│   └── stackvm-asm/      # CLI assembler tool
├── docs/                 # Documentation
├── examples/             # Usage examples
└── testdata/             # Test programs and data
```

## Building

### Build Everything

```bash
go build ./...
```

### Build CLI Tool

```bash
cd cmd/stackvm-asm
go build
```

Or from the root:

```bash
go build -o stackvm-asm ./cmd/stackvm-asm
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

## Advanced Features

### Custom Instructions

Extend the VM with your own instructions:

```go
// Define a custom instruction
type TimestampHandler struct{}

func (h *TimestampHandler) Execute(ctx stackvm.ExecutionContext, operand int32) error {
    timestamp := time.Now().UnixNano()
    return ctx.Push(stackvm.IntValue(timestamp))
}

func (h *TimestampHandler) Name() string {
    return "TIMESTAMP"
}

// Register and use
registry := stackvm.NewInstructionRegistry()
registry.Register(128, &TimestampHandler{})

config := stackvm.Config{
    InstructionRegistry: registry,
}
vm := stackvm.NewWithConfig(config)
```

### Custom Memory Providers

Implement your own memory backend:

```go
type MyMemory struct {
    // Your storage implementation
}

func (m *MyMemory) Load(index int) (stackvm.Value, error) {
    // Your load implementation
}

func (m *MyMemory) Store(index int, value stackvm.Value) error {
    // Your store implementation
}

func (m *MyMemory) Size() int {
    // Your size implementation
}
```

### VM Pooling

Reuse VM instances for high-performance scenarios:

```go
pool := stackvm.NewPool(stackvm.PoolConfig{
    InitialSize: 10,
    MaxSize:     100,
    VMConfig:    stackvm.Config{},
})

// Get VM from pool
vm := pool.Get()

// Use VM
result, err := vm.Execute(program, memory, opts)

// Return to pool
pool.Put(vm)
```

### Execution Limits

Control resource usage:

```go
result, err := vm.Execute(program, memory, stackvm.ExecuteOptions{
    MaxInstructions: 10000,           // Limit instruction count
    MaxStackDepth:   512,              // Limit stack size
    Timeout:         time.Second,      // Wall-clock timeout
    Context:         ctx,              // Cancellation support
})
```

## Example Programs

The `testdata/programs/` directory contains example programs demonstrating various features:

- **Arithmetic**: Basic calculations and expressions
- **Control Flow**: Loops, conditionals, and jumps
- **Memory Operations**: Variable storage and array manipulation
- **Math Functions**: Trigonometry and advanced math
- **Algorithms**: Factorial, Fibonacci, prime checking, etc.

## Design Principles

- **Independence**: No knowledge of host system concepts
- **Simplicity**: Clean, minimal interfaces
- **Extensibility**: Pluggable instructions, memory providers, value types
- **Safety**: Enforced limits, comprehensive error handling
- **Performance**: Zero-allocation hot paths, pooling support
- **Testability**: Easy to test in isolation

## Version

Current version: **1.0.0**

StackVM follows semantic versioning. Within a major version:
- Public interfaces remain stable
- Standard opcodes (0-127) remain stable
- Binary format is backward compatible
- New features added via minor versions
- Bug fixes via patch versions

## Requirements

- Go 1.21 or higher
- No external dependencies

## Contributing

This is a standalone VM implementation with a complete specification. When contributing:

- Follow the implementation order in CLAUDE.md
- Maintain zero external dependencies
- Add tests for all new features
- Update documentation
- Follow standard Go formatting (gofmt)

## License

See LICENSE file for details.

## Links

- [GitHub Repository](https://github.com/pmuston/stackvm)
- [Documentation](docs/)
- [Issue Tracker](https://github.com/pmuston/stackvm/issues)

## Acknowledgments

Inspired by classical stack-based architectures including Forth, PostScript, and the Java Virtual Machine.
