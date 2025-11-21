# stackvm - Stack-Based Virtual Machine

## Project Overview

This is a standalone, reusable stack-based virtual machine package for Go.
See `docs/SPEC.md` for the complete technical specification.

## Key Principles

- Zero external dependencies (standard library only)
- All public API in `pkg/stackvm/`
- Internal implementation in `internal/`
- Comprehensive tests for all components
- No knowledge of host systems - pure VM implementation

## Build Commands
```bash
go build ./...
go test ./...
go test -race ./...
go test -bench=. ./...
```

## Implementation Order

1. Value system (`pkg/stackvm/value.go`)
2. Instructions and opcodes (`pkg/stackvm/instruction.go`)
3. Memory interface (`pkg/stackvm/memory.go`)
4. Errors (`pkg/stackvm/errors.go`)
5. Execution context (`pkg/stackvm/context.go`)
6. Core VM (`pkg/stackvm/vm.go`, `internal/executor/`)
7. Standard instruction handlers (`internal/opcodes/`)
8. Instruction registry (`pkg/stackvm/registry.go`)
9. Program and builder (`pkg/stackvm/program.go`)
10. Encoder/decoder (`pkg/stackvm/encoder.go`, `pkg/stackvm/decoder.go`)
11. Assembler (`pkg/stackvm/assembler.go`, `internal/asm/`)
12. Disassembler (`pkg/stackvm/disassembler.go`)
13. VM Pool (`pkg/stackvm/pool.go`)
14. Testing utilities (`pkg/stackvm/testing.go`)
15. CLI tool (`cmd/stackvm-asm/`)

## Code Style

- Use standard Go formatting (gofmt)
- Comprehensive godoc comments on all public types
- Table-driven tests
- Avoid global state
- Return errors, don't panic

## Testing Requirements

- Unit test each opcode
- Test error conditions (stack overflow, underflow, etc.)
- Benchmark critical paths
- Test assembler round-trips