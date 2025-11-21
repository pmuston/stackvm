# StackVM - Standalone Virtual Machine Specification

## Implementation Specification for Claude Code

---

## 1. Overview

### 1.1 Purpose

StackVM is a standalone, reusable stack-based virtual machine package designed to:
- Execute bytecode programs safely and efficiently
- Be embedded in multiple host systems
- Allow extension with custom instructions
- Support various memory/storage backends
- Provide tooling for program development

### 1.2 Design Principles

- **Independence**: No knowledge of host system concepts
- **Simplicity**: Clean, minimal interfaces
- **Extensibility**: Pluggable instructions, memory providers, value types
- **Safety**: Enforced limits, comprehensive error handling
- **Performance**: Zero-allocation hot paths, pooling support
- **Testability**: Easy to test in isolation

### 1.3 Technology Constraints

- Go 1.21+
- Zero external dependencies (standard library only)
- Pure Go implementation

---

## 2. Package Structure

```
stackvm/
├── pkg/
│   └── stackvm/                  # Public API
│       ├── vm.go                 # Core VM interface and implementation
│       ├── instruction.go        # Instruction and opcode definitions
│       ├── value.go              # Value type system
│       ├── memory.go             # Memory interface
│       ├── program.go            # Program interface and builders
│       ├── context.go            # Execution context
│       ├── registry.go           # Instruction registry
│       ├── pool.go               # VM pooling
│       ├── errors.go             # Error types
│       ├── assembler.go          # Assembler interface
│       ├── disassembler.go       # Disassembler interface
│       ├── encoder.go            # Binary encoding
│       ├── decoder.go            # Binary decoding
│       └── testing.go            # Test utilities
│
├── internal/
│   ├── executor/                 # VM execution implementation
│   │   └── executor.go
│   │
│   ├── opcodes/                  # Standard instruction handlers
│   │   ├── stack.go              # PUSH, POP, DUP, SWAP
│   │   ├── arithmetic.go         # ADD, SUB, MUL, DIV, etc.
│   │   ├── logic.go              # AND, OR, NOT, XOR
│   │   ├── comparison.go         # EQ, NE, GT, LT, etc.
│   │   ├── control.go            # JMP, JMPZ, JMPNZ, HALT
│   │   └── math.go               # SQRT, SIN, COS, etc.
│   │
│   ├── asm/                      # Assembler implementation
│   │   ├── lexer.go
│   │   ├── parser.go
│   │   └── codegen.go
│   │
│   └── encoding/                 # Binary format implementation
│       ├── encoder.go
│       └── decoder.go
│
├── cmd/
│   └── stackvm-asm/              # Standalone assembler tool
│       └── main.go
│
├── examples/
│   ├── simple/                   # Basic usage example
│   │   └── main.go
│   ├── custom_instructions/      # Extending instruction set
│   │   └── main.go
│   └── custom_memory/            # Custom memory provider
│       └── main.go
│
├── testdata/
│   ├── programs/                 # Test programs
│   │   ├── arithmetic.asm
│   │   ├── control_flow.asm
│   │   └── math_functions.asm
│   └── expected/                 # Expected outputs
│
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 3. Value System

### 3.1 Value Types

StackVM supports a fixed set of value types that can be extended by host systems.

**Core Types:**

| Type | Go Type | Description |
|------|---------|-------------|
| Nil | nil | Absence of value |
| Float | float64 | IEEE 754 double precision |
| Int | int64 | Signed 64-bit integer |
| Bool | bool | True or false |
| String | string | UTF-8 text |

**Reserved for Extension:**

| Type | Range | Description |
|------|-------|-------------|
| Custom | 128-255 | Host-defined types |

### 3.2 Value Structure

```
Value:
  Type: ValueType (uint8)
  Data: interface{} (underlying Go value)
```

### 3.3 Value Constructors

```
FloatValue(v float64) Value
IntValue(v int64) Value
BoolValue(v bool) Value
StringValue(v string) Value
NilValue() Value
CustomValue(typ ValueType, data interface{}) Value
```

### 3.4 Value Accessors

```
(v Value) AsFloat() (float64, error)
(v Value) AsInt() (int64, error)
(v Value) AsBool() (bool, error)
(v Value) AsString() (string, error)
(v Value) IsNil() bool
```

### 3.5 Value Utilities

```
(v Value) IsNumeric() bool
  - Returns true for Float and Int types
  
(v Value) IsTruthy() bool
  - Float: true if != 0.0
  - Int: true if != 0
  - Bool: the value itself
  - String: true if not empty
  - Nil: false
  
(v Value) String() string
  - Human-readable representation
  
(v Value) Equal(other Value) bool
  - Type-aware equality comparison
```

### 3.6 Type Conversion

```
ValueConverter interface:
  Convert(value Value, targetType ValueType) (Value, error)
```

**Default Conversions:**
- Int → Float: Allowed (may lose precision)
- Float → Int: Allowed (truncates)
- Numeric → Bool: 0 = false, else true
- Bool → Numeric: false = 0, true = 1
- Any → String: String representation

---

## 4. Memory Interface

### 4.1 Purpose

Memory abstraction allows host systems to provide storage for VM programs. The VM reads and writes values by index.

### 4.2 Memory Interface

```
Memory interface:
  Load(index int) (Value, error)
    - Retrieve value at index
    - Returns ErrInvalidMemoryAddress if out of bounds
    
  Store(index int, value Value) error
    - Save value at index
    - Returns ErrInvalidMemoryAddress if out of bounds
    - Returns ErrReadOnlyMemory if not writable
    
  Size() int
    - Number of addressable locations
```

### 4.3 ReadOnlyMemory Interface

```
ReadOnlyMemory interface:
  Memory
  IsReadOnly() bool
```

### 4.4 SimpleMemory Implementation

StackVM provides a basic implementation for testing and simple use cases.

```
SimpleMemory:
  Constructor:
    NewSimpleMemory(size int) *SimpleMemory
    
  Additional Methods:
    Values() []Value
      - Returns copy of all values
      
    SetValues(values []Value)
      - Bulk set all values
      
    Reset()
      - Clear all values to Nil
```

### 4.5 Memory Usage Notes

- Index 0 is the first memory location
- Indices must be non-negative
- Size is fixed at creation (for SimpleMemory)
- Host systems can implement dynamic sizing

---

## 5. Instruction Set

### 5.1 Instruction Format

```
Instruction:
  Opcode: Opcode (uint8)
  Operand: int32 (signed 32-bit)
```

**Binary Encoding:** 5 bytes per instruction
- Byte 0: Opcode
- Bytes 1-4: Operand (big-endian, signed)

### 5.2 Opcode Ranges

| Range | Purpose |
|-------|---------|
| 0-63 | Core operations |
| 64-127 | Reserved for future standard ops |
| 128-255 | Custom/host-defined operations |

### 5.3 Stack Operations (0-15)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 0 | PUSH | value | → a | Push immediate value (as float) |
| 1 | PUSHI | value | → a | Push immediate value (as int) |
| 2 | POP | - | a → | Remove top of stack |
| 3 | DUP | - | a → a a | Duplicate top |
| 4 | SWAP | - | a b → b a | Exchange top two |
| 5 | OVER | - | a b → a b a | Copy second to top |
| 6 | ROT | - | a b c → b c a | Rotate top three |

### 5.4 Arithmetic Operations (16-31)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 16 | ADD | - | a b → (a+b) | Addition |
| 17 | SUB | - | a b → (a-b) | Subtraction |
| 18 | MUL | - | a b → (a*b) | Multiplication |
| 19 | DIV | - | a b → (a/b) | Division (error if b=0) |
| 20 | MOD | - | a b → (a%b) | Modulo |
| 21 | NEG | - | a → (-a) | Negate |
| 22 | ABS | - | a → |a| | Absolute value |
| 23 | INC | - | a → (a+1) | Increment |
| 24 | DEC | - | a → (a-1) | Decrement |

### 5.5 Logic Operations (32-39)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 32 | AND | - | a b → (a && b) | Logical AND |
| 33 | OR | - | a b → (a \|\| b) | Logical OR |
| 34 | NOT | - | a → (!a) | Logical NOT |
| 35 | XOR | - | a b → (a xor b) | Logical XOR |

### 5.6 Comparison Operations (40-47)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 40 | EQ | - | a b → (a == b) | Equal |
| 41 | NE | - | a b → (a != b) | Not equal |
| 42 | GT | - | a b → (a > b) | Greater than |
| 43 | LT | - | a b → (a < b) | Less than |
| 44 | GE | - | a b → (a >= b) | Greater or equal |
| 45 | LE | - | a b → (a <= b) | Less or equal |

### 5.7 Memory Operations (48-55)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 48 | LOAD | index | → a | Load from memory[index] |
| 49 | STORE | index | a → | Store to memory[index] |
| 50 | LOADD | - | i → a | Load from memory[pop()] |
| 51 | STORED | - | a i → | Store to memory[pop()] |

### 5.8 Control Flow Operations (56-63)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 56 | JMP | offset | - | Jump to offset |
| 57 | JMPZ | offset | a → | Jump if zero/false |
| 58 | JMPNZ | offset | a → | Jump if non-zero/true |
| 59 | CALL | offset | - | Call subroutine (push return address) |
| 60 | RET | - | - | Return from subroutine |
| 61 | HALT | - | - | Stop execution |
| 62 | NOP | - | - | No operation |

### 5.9 Math Functions (64-79)

| Opcode | Name | Operand | Stack Effect | Description |
|--------|------|---------|--------------|-------------|
| 64 | SQRT | - | a → √a | Square root |
| 65 | SIN | - | a → sin(a) | Sine (radians) |
| 66 | COS | - | a → cos(a) | Cosine (radians) |
| 67 | TAN | - | a → tan(a) | Tangent (radians) |
| 68 | ASIN | - | a → asin(a) | Arc sine |
| 69 | ACOS | - | a → acos(a) | Arc cosine |
| 70 | ATAN | - | a → atan(a) | Arc tangent |
| 71 | ATAN2 | - | y x → atan2(y,x) | Two-argument arc tangent |
| 72 | LOG | - | a → ln(a) | Natural logarithm |
| 73 | LOG10 | - | a → log10(a) | Base-10 logarithm |
| 74 | EXP | - | a → e^a | Exponential |
| 75 | POW | - | a b → a^b | Power |
| 76 | MIN | - | a b → min(a,b) | Minimum |
| 77 | MAX | - | a b → max(a,b) | Maximum |
| 78 | FLOOR | - | a → floor(a) | Floor |
| 79 | CEIL | - | a → ceil(a) | Ceiling |
| 80 | ROUND | - | a → round(a) | Round to nearest |
| 81 | TRUNC | - | a → trunc(a) | Truncate toward zero |

### 5.10 Custom Operations (128-255)

Reserved for host system extensions. Host systems register handlers via InstructionRegistry.

---

## 6. VM Interface

### 6.1 VM Interface Definition

```
VM interface:
  Execute(program Program, memory Memory, opts ExecuteOptions) (*Result, error)
    - Run program with memory provider
    - Returns result with stats or error
    
  Reset()
    - Clear VM state for reuse
```

### 6.2 ExecuteOptions

```
ExecuteOptions:
  MaxInstructions: uint32
    - Limit on instructions executed (0 = unlimited)
    - Returns ErrInstructionLimit if exceeded
    
  MaxStackDepth: int
    - Stack size limit (0 = default 256)
    - Returns ErrStackOverflow if exceeded
    
  Timeout: time.Duration
    - Wall-clock timeout (0 = no timeout)
    - Returns ErrTimeout if exceeded
    
  Context: context.Context
    - Cancellation context (nil = no cancellation)
    - Returns context error if cancelled
```

### 6.3 Result

```
Result:
  InstructionCount: uint32
    - Number of instructions executed
    
  StackDepth: int
    - Final stack depth
    
  ExecutionTime: time.Duration
    - Total execution time
    
  Halted: bool
    - True if HALT instruction reached
    
  Error: error
    - Execution error (nil if successful)
```

### 6.4 VM Constructor

```
New() VM
  - Create VM with default configuration
  
NewWithConfig(config Config) VM
  - Create VM with custom configuration
```

### 6.5 Config

```
Config:
  StackSize: int
    - Initial stack capacity (default 256)
    
  DefaultInstrLimit: uint32
    - Default instruction limit (0 = unlimited)
    
  InstructionRegistry: InstructionRegistry
    - Custom instruction handlers (nil = standard only)
    
  ValueConverter: ValueConverter
    - Custom type conversions (nil = defaults)
```

---

## 7. Execution Context

### 7.1 Purpose

ExecutionContext provides access to VM state during instruction execution. Used by custom instruction handlers.

### 7.2 ExecutionContext Interface

```
ExecutionContext interface:
  Stack Operations:
    Push(value Value) error
      - Push value onto stack
      - Returns ErrStackOverflow if full
      
    Pop() (Value, error)
      - Remove and return top value
      - Returns ErrStackUnderflow if empty
      
    Peek() (Value, error)
      - Return top value without removing
      - Returns ErrStackUnderflow if empty
      
    PeekN(n int) (Value, error)
      - Return nth value from top (0 = top)
      
    StackDepth() int
      - Current number of values on stack
      
  Program Counter:
    PC() int
      - Current program counter
      
    SetPC(pc int)
      - Set program counter directly
      
    Jump(offset int)
      - Set PC to offset
      
  Memory:
    Memory() Memory
      - Access memory provider
      
  Execution Control:
    InstructionCount() uint32
      - Instructions executed so far
      
    IncrementInstructionCount()
      - Add to instruction counter
      
    Halt()
      - Stop execution
      
    IsHalted() bool
      - Check if halted
```

---

## 8. Instruction Registry

### 8.1 Purpose

Allows host systems to register custom instruction handlers for opcodes 128-255.

### 8.2 InstructionHandler Interface

```
InstructionHandler interface:
  Execute(ctx ExecutionContext, operand int32) error
    - Perform the instruction operation
    - Access stack/memory via context
    - Return error on failure
    
  Name() string
    - Mnemonic for assembler/disassembler
```

### 8.3 InstructionRegistry Interface

```
InstructionRegistry interface:
  Register(opcode Opcode, handler InstructionHandler) error
    - Add handler for opcode
    - Returns error if opcode < 128 (reserved)
    - Returns error if opcode already registered
    
  Unregister(opcode Opcode) error
    - Remove handler for opcode
    - Returns error if not registered
    
  Get(opcode Opcode) (InstructionHandler, bool)
    - Retrieve handler for opcode
    - Returns false if not registered
    
  List() []Opcode
    - Return all registered custom opcodes
    
  Names() map[Opcode]string
    - Return opcode → name mapping
```

### 8.4 Registry Constructor

```
NewInstructionRegistry() InstructionRegistry
  - Create empty registry
```

### 8.5 Example Custom Instruction

```
TimestampHandler:
  Execute(ctx ExecutionContext, operand int32) error:
    - Push current Unix timestamp (nanoseconds) as Int
    
  Name() string:
    - Returns "TIMESTAMP"

Usage:
  registry := stackvm.NewInstructionRegistry()
  registry.Register(128, &TimestampHandler{})
```

---

## 9. Program Interface

### 9.1 Program Interface

```
Program interface:
  Instructions() []Instruction
    - Return instruction sequence
    
  SymbolTable() map[int]string
    - Return address → label mapping (for debugging)
    - May be nil if no debug info
    
  Metadata() ProgramMetadata
    - Return program information
```

### 9.2 ProgramMetadata

```
ProgramMetadata:
  Name: string
  Version: string
  Author: string
  Description: string
  Created: time.Time
```

### 9.3 SimpleProgram

Basic program implementation provided by StackVM.

```
SimpleProgram:
  Constructor:
    NewProgram(instructions []Instruction) *SimpleProgram
    NewProgramWithMetadata(instructions []Instruction, metadata ProgramMetadata) *SimpleProgram
    
  Methods:
    SetSymbolTable(symbols map[int]string)
    AddSymbol(address int, label string)
```

### 9.4 ProgramBuilder

Fluent builder for constructing programs.

```
ProgramBuilder:
  Constructor:
    NewProgramBuilder() *ProgramBuilder
    
  Stack Operations:
    Push(v float64) *ProgramBuilder
    PushInt(v int64) *ProgramBuilder
    Pop() *ProgramBuilder
    Dup() *ProgramBuilder
    Swap() *ProgramBuilder
    
  Arithmetic:
    Add() *ProgramBuilder
    Sub() *ProgramBuilder
    Mul() *ProgramBuilder
    Div() *ProgramBuilder
    Neg() *ProgramBuilder
    
  Logic:
    And() *ProgramBuilder
    Or() *ProgramBuilder
    Not() *ProgramBuilder
    
  Comparison:
    Eq() *ProgramBuilder
    Ne() *ProgramBuilder
    Gt() *ProgramBuilder
    Lt() *ProgramBuilder
    Ge() *ProgramBuilder
    Le() *ProgramBuilder
    
  Memory:
    Load(index int) *ProgramBuilder
    Store(index int) *ProgramBuilder
    
  Control Flow:
    Label(name string) *ProgramBuilder
    Jmp(label string) *ProgramBuilder
    JmpZ(label string) *ProgramBuilder
    JmpNZ(label string) *ProgramBuilder
    Halt() *ProgramBuilder
    
  Math:
    Sqrt() *ProgramBuilder
    Sin() *ProgramBuilder
    Cos() *ProgramBuilder
    Min() *ProgramBuilder
    Max() *ProgramBuilder
    
  Custom:
    Custom(opcode Opcode, operand int32) *ProgramBuilder
    
  Build:
    Build() (Program, error)
      - Resolve labels and return program
      - Returns error if unresolved labels
```

---

## 10. VM Pool

### 10.1 Purpose

Pool VMs for reuse to avoid allocation overhead in high-frequency execution scenarios.

### 10.2 Pool Interface

```
Pool interface:
  Get() VM
    - Retrieve VM from pool
    - Creates new VM if pool empty
    
  Put(vm VM)
    - Return VM to pool
    - VM is reset before pooling
    
  Size() int
    - Current number of pooled VMs
    
  Stats() PoolStats
    - Pool statistics
```

### 10.3 PoolStats

```
PoolStats:
  Created: uint64
    - Total VMs created
    
  Reused: uint64
    - Times VM was reused from pool
    
  CurrentSize: int
    - Current pool size
    
  MaxSize: int
    - Maximum pool size reached
```

### 10.4 PoolConfig

```
PoolConfig:
  InitialSize: int
    - Pre-create this many VMs (default 0)
    
  MaxSize: int
    - Maximum VMs to keep pooled (default 100)
    - Excess VMs discarded on Put
    
  VMConfig: Config
    - Configuration for created VMs
```

### 10.5 Pool Constructor

```
NewPool(config PoolConfig) Pool
```

---

## 11. Assembler

### 11.1 Purpose

Convert human-readable assembly source to bytecode programs.

### 11.2 Assembly Syntax

```
; Comment (to end of line)
# Also a comment

LABEL:              ; Label definition (address marker)
    OPCODE          ; Instruction with no operand
    OPCODE operand  ; Instruction with operand
    OPCODE LABEL    ; Instruction with label reference
```

### 11.3 Assembler Interface

```
Assembler interface:
  Assemble(source string) (Program, error)
    - Parse and compile source to program
    - Returns error with line number on failure
    
  AssembleFile(path string) (Program, error)
    - Read file and assemble
    
  SetRegistry(registry InstructionRegistry)
    - Enable custom instruction names
```

### 11.4 AssemblerError

```
AssemblerError:
  Line: int
  Column: int
  Message: string
  Source: string (the problematic line)
```

### 11.5 Assembler Constructor

```
NewAssembler() Assembler
```

### 11.6 Assembly Example

```assembly
; Calculate: c = a + b
; Memory: 0=a, 1=b, 2=c

START:
    LOAD 0          ; Load a
    LOAD 1          ; Load b
    ADD             ; a + b
    STORE 2         ; Store to c
    HALT

; Conditional example
CHECK:
    LOAD 0          ; Load value
    PUSH 10         ; Compare against 10
    GT              ; value > 10?
    JMPZ BELOW      ; Jump if not greater
    PUSH 1          ; Result = 1
    JMP DONE
BELOW:
    PUSH 0          ; Result = 0
DONE:
    STORE 1         ; Store result
    HALT
```

---

## 12. Disassembler

### 12.1 Purpose

Convert bytecode programs back to human-readable assembly.

### 12.2 Disassembler Interface

```
Disassembler interface:
  Disassemble(program Program) (string, error)
    - Convert program to assembly source
    
  DisassembleWithOptions(program Program, opts DisassembleOptions) (string, error)
    - Convert with formatting options
```

### 12.3 DisassembleOptions

```
DisassembleOptions:
  ShowAddresses: bool
    - Include address prefix (e.g., "0000:")
    
  ShowHex: bool
    - Include hex encoding
    
  ShowComments: bool
    - Include opcode descriptions
    
  UseSymbols: bool
    - Use symbol table labels if available
    
  Registry: InstructionRegistry
    - For custom instruction names
```

### 12.4 Disassembler Constructor

```
NewDisassembler() Disassembler
```

### 12.5 Output Example

```assembly
; With ShowAddresses and ShowComments
0000: LOAD 0          ; Load from memory[0]
0001: LOAD 1          ; Load from memory[1]
0002: ADD             ; Pop b, pop a, push a+b
0003: STORE 2         ; Store to memory[2]
0004: HALT            ; Stop execution
```

---

## 13. Binary Encoding

### 13.1 Instruction Encoding

Each instruction is 5 bytes:

```
Byte 0: Opcode (uint8)
Bytes 1-4: Operand (int32, big-endian)
```

### 13.2 Program Encoding

Programs are simply concatenated instructions with optional header.

**Simple Format (no header):**
```
[Instruction 0][Instruction 1]...[Instruction N]
```

**With Header (for metadata):**
```
[Magic: 4 bytes "SVMP"]
[Version: 1 byte]
[Flags: 1 byte]
[Instruction Count: 4 bytes, big-endian]
[Metadata Length: 4 bytes, big-endian]
[Metadata: JSON encoded ProgramMetadata]
[Instructions...]
[Symbol Table Length: 4 bytes] (if flag set)
[Symbol Table: JSON encoded map[int]string] (if flag set)
```

### 13.3 Encoder Interface

```
Encoder interface:
  Encode(program Program, w io.Writer) error
    - Write program in simple format
    
  EncodeWithHeader(program Program, w io.Writer) error
    - Write program with header and metadata
    
  EncodeInstruction(inst Instruction) []byte
    - Encode single instruction
```

### 13.4 Decoder Interface

```
Decoder interface:
  Decode(r io.Reader) (Program, error)
    - Read program (auto-detects format)
    
  DecodeInstruction(data []byte) (Instruction, error)
    - Decode single instruction
```

### 13.5 Encoder/Decoder Constructors

```
NewEncoder() Encoder
NewDecoder() Decoder
```

---

## 14. Error System

### 14.1 Standard Errors

```
var (
    ErrStackOverflow        = errors.New("stack overflow")
    ErrStackUnderflow       = errors.New("stack underflow")
    ErrInvalidMemoryAddress = errors.New("invalid memory address")
    ErrReadOnlyMemory       = errors.New("memory is read-only")
    ErrInvalidInstruction   = errors.New("invalid instruction")
    ErrInvalidOpcode        = errors.New("invalid opcode")
    ErrInstructionLimit     = errors.New("instruction limit exceeded")
    ErrDivisionByZero       = errors.New("division by zero")
    ErrTypeMismatch         = errors.New("type mismatch")
    ErrTimeout              = errors.New("execution timeout")
    ErrInvalidOperand       = errors.New("invalid operand")
    ErrInvalidProgram       = errors.New("invalid program")
    ErrUnresolvedLabel      = errors.New("unresolved label")
)
```

### 14.2 VMError

Wraps errors with execution context.

```
VMError:
  Err: error
    - Underlying error
    
  PC: int
    - Program counter at failure
    
  InstructionCount: uint32
    - Instructions executed before failure
    
  StackDepth: int
    - Stack depth at failure
    
  Opcode: Opcode
    - Instruction that failed
    
  Message: string
    - Additional context

Methods:
  Error() string
  Unwrap() error
  Is(target error) bool
```

### 14.3 Error Checking

```
IsStackError(err error) bool
  - Returns true for stack overflow/underflow
  
IsMemoryError(err error) bool
  - Returns true for memory errors
  
IsLimitError(err error) bool
  - Returns true for instruction limit/timeout
```

---

## 15. Testing Utilities

### 15.1 TestMemory

Simple memory for testing.

```
TestMemory:
  Constructor:
    NewTestMemory(size int) *TestMemory
    NewTestMemoryWithValues(values []Value) *TestMemory
    
  Additional Methods:
    Set(index int, value Value)
      - Direct set (no error checking)
      
    Get(index int) Value
      - Direct get (no error checking)
      
    All() []Value
      - Return all values
```

### 15.2 TestProgramBuilder

Simplified builder for tests.

```
TestProgramBuilder:
  Constructor:
    NewTestProgram() *TestProgramBuilder
    
  All ProgramBuilder methods plus:
    MustBuild() Program
      - Build or panic (for tests)
```

### 15.3 Assertions

```
AssertStackDepth(t *testing.T, vm VM, expected int)
AssertMemoryValue(t *testing.T, memory Memory, index int, expected Value)
AssertExecutionSuccess(t *testing.T, result *Result, err error)
AssertExecutionError(t *testing.T, err error, expectedErr error)
```

### 15.4 Test Example

```go
func TestAddition(t *testing.T) {
    vm := stackvm.New()
    memory := stackvm.NewTestMemory(3)
    memory.Set(0, stackvm.FloatValue(10))
    memory.Set(1, stackvm.FloatValue(20))
    
    program := stackvm.NewTestProgram().
        Load(0).
        Load(1).
        Add().
        Store(2).
        Halt().
        MustBuild()
    
    result, err := vm.Execute(program, memory, stackvm.ExecuteOptions{})
    
    stackvm.AssertExecutionSuccess(t, result, err)
    stackvm.AssertMemoryValue(t, memory, 2, stackvm.FloatValue(30))
}
```

---

## 16. CLI Tool: stackvm-asm

### 16.1 Commands

**Compile:**
```bash
stackvm-asm compile <input.asm> -o <output.bin>
stackvm-asm compile <input.asm> --stdout  # Output to stdout
```

**Disassemble:**
```bash
stackvm-asm disasm <input.bin>
stackvm-asm disasm <input.bin> -o <output.asm>
stackvm-asm disasm <input.bin> --show-addresses --show-hex
```

**Validate:**
```bash
stackvm-asm validate <input.asm>
```

**Run (for testing):**
```bash
stackvm-asm run <input.bin> --memory 10 --set 0=5.0 --set 1=3.0
stackvm-asm run <input.asm> --memory 10  # Auto-compile
```

**Info:**
```bash
stackvm-asm info <input.bin>  # Show metadata, instruction count
```

### 16.2 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Compilation error |
| 2 | Runtime error |
| 3 | File I/O error |
| 4 | Invalid arguments |

---

## 17. Implementation Requirements

### 17.1 VM Execution

- Single-threaded execution per VM instance
- No heap allocations in main execution loop
- Stack implemented as pre-allocated slice
- Instruction dispatch via switch or jump table
- Bounds checking on all memory access
- Type checking on operations (where required)

### 17.2 Numeric Operations

- Float operations use Go float64
- Int operations use Go int64
- Mixed operations: convert to float
- Division by zero returns error (not panic)
- Math functions use Go math package

### 17.3 Control Flow

- CALL pushes return address to call stack (separate from data stack)
- RET pops from call stack
- Call stack depth limit: 64 (configurable)
- JMP/JMPZ/JMPNZ use absolute addresses

### 17.4 Thread Safety

- VM instances are NOT thread-safe
- Pool is thread-safe
- InstructionRegistry is thread-safe after initial setup
- Memory implementations should document thread safety

---

## 18. Testing Requirements

### 18.1 Unit Tests

**Value System:**
- All constructors
- All accessors with type checking
- Type conversions
- Equality comparisons
- Truthy evaluation

**Instructions:**
- Each opcode individually
- Edge cases (overflow, underflow, divide by zero)
- Type coercion behavior

**Assembler:**
- Valid programs
- All opcodes
- Labels and jumps
- Syntax errors (line numbers)
- Custom instructions

**Disassembler:**
- Round-trip (assemble → disassemble → assemble)
- All formatting options
- Symbol table handling

**Encoder/Decoder:**
- Simple format
- Header format
- Corruption detection

### 18.2 Integration Tests

- Complete program execution
- Memory provider integration
- Custom instruction registration
- Pool behavior under load

### 18.3 Benchmark Tests

- Instruction dispatch overhead
- Memory access patterns
- Pool performance
- Large program execution

---

## 19. Documentation Requirements

### 19.1 Package Documentation

- Comprehensive godoc for all public types
- Usage examples in doc comments
- Package-level overview

### 19.2 README.md

- Quick start guide
- Installation
- Basic usage example
- Link to full documentation

### 19.3 Instruction Reference

- Complete opcode listing
- Stack effects
- Error conditions
- Examples

### 19.4 Integration Guide

- Implementing Memory interface
- Adding custom instructions
- Error handling patterns
- Performance best practices

---

## 20. Example Programs

### 20.1 Arithmetic

```assembly
; Calculate: result = (a + b) * c
; Memory: 0=a, 1=b, 2=c, 3=result

    LOAD 0          ; a
    LOAD 1          ; b
    ADD             ; a + b
    LOAD 2          ; c
    MUL             ; (a + b) * c
    STORE 3         ; result
    HALT
```

### 20.2 Conditional

```assembly
; Clamp value to range [min, max]
; Memory: 0=value, 1=min, 2=max, 3=result

    LOAD 0          ; value
    LOAD 1          ; min
    MAX             ; max(value, min)
    LOAD 2          ; max
    MIN             ; min(max(value, min), max)
    STORE 3         ; result
    HALT
```

### 20.3 Loop

```assembly
; Sum numbers from 1 to n
; Memory: 0=n, 1=sum, 2=counter

    PUSH 0
    STORE 1         ; sum = 0
    PUSH 1
    STORE 2         ; counter = 1

LOOP:
    LOAD 2          ; counter
    LOAD 0          ; n
    GT              ; counter > n?
    JMPNZ DONE      ; Exit if true
    
    LOAD 1          ; sum
    LOAD 2          ; counter
    ADD             ; sum + counter
    STORE 1         ; sum = sum + counter
    
    LOAD 2          ; counter
    INC             ; counter + 1
    STORE 2         ; counter = counter + 1
    
    JMP LOOP

DONE:
    HALT
```

### 20.4 Function Call

```assembly
; Calculate distance: sqrt(x*x + y*y)
; Memory: 0=x, 1=y, 2=result

    LOAD 0          ; x
    CALL SQUARE
    LOAD 1          ; y
    CALL SQUARE
    ADD             ; x² + y²
    SQRT            ; sqrt(x² + y²)
    STORE 2
    HALT

SQUARE:             ; Square top of stack
    DUP
    MUL
    RET
```

---

## 21. Versioning

### 21.1 Version Format

```
MAJOR.MINOR.PATCH
```

### 21.2 Compatibility Guarantees

**Within Major Version (1.x.x):**
- Public interfaces stable
- Standard opcodes (0-127) stable
- Binary format backward compatible
- New features via minor versions
- Bug fixes via patch versions

**Breaking Changes (2.0.0, etc.):**
- Interface changes
- Opcode reassignment
- Binary format changes

### 21.3 Version Constants

```
const (
    VersionMajor = 1
    VersionMinor = 0
    VersionPatch = 0
)

func Version() string
  - Returns "1.0.0" format
  
func VersionInfo() VersionInfo
  - Returns structured version info
```

---

## 22. Success Criteria

Implementation is complete when:

1. **All opcodes implemented** and tested
2. **Memory interface** works with custom implementations
3. **Custom instructions** can be registered and executed
4. **Assembler** compiles valid programs and reports clear errors
5. **Disassembler** produces correct output
6. **Encoder/Decoder** handles all formats
7. **Pool** manages VMs efficiently
8. **CLI tool** provides all commands
9. **Documentation** is comprehensive
10. **Tests** have good coverage
11. **Examples** are working and instructive
12. **Zero external dependencies**

---

## 23. Non-Requirements (Out of Scope)

- Garbage collection (host manages memory)
- Concurrency primitives (single-threaded execution)
- File I/O instructions
- Network instructions
- Graphics/display
- Debugging protocol (IDE integration)
- JIT compilation
- Optimization passes

These may be added in future versions or provided by host systems.
