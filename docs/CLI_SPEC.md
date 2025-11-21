# StackVM Assembler CLI Tool Specification

**Version:** 1.0.0
**Tool:** `stackvm-asm`

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Command Syntax](#command-syntax)
4. [Flags Reference](#flags-reference)
5. [Operation Modes](#operation-modes)
6. [Memory Features](#memory-features)
7. [Output Format](#output-format)
8. [Error Handling](#error-handling)
9. [Examples](#examples)
10. [Exit Codes](#exit-codes)

---

## Overview

`stackvm-asm` is the command-line assembler and runtime tool for StackVM. It provides:

- **Assembly** - Converts `.asm` files to executable programs
- **Execution** - Runs assembled programs with configurable limits
- **Disassembly** - Converts programs back to readable assembly
- **Memory Management** - Initialize and inspect memory during execution
- **Statistics** - Performance metrics and execution analysis

## Installation

### Build from Source

```bash
# From the repository root
go build ./cmd/stackvm-asm

# The binary will be created as 'stackvm-asm'
```

### Install to GOPATH

```bash
go install github.com/pmuston/stackvm/cmd/stackvm-asm@latest
```

## Command Syntax

```
stackvm-asm [OPTIONS] <input-file>
```

**Arguments:**
- `<input-file>` - Path to assembly (.asm) file (required)

**Options:** See [Flags Reference](#flags-reference)

## Flags Reference

### Basic Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-h` | bool | false | Show help message and exit |
| `-v` | bool | false | Show version and exit |
| `-q` | bool | false | Quiet mode (suppress non-error output) |
| `-o` | string | "" | Output file for assembled/disassembled code |

### Execution Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-r` | bool | false | Run the program after assembling |
| `-s` | bool | false | Show execution statistics (requires `-r`) |
| `-m` | string | "" | Show memory values after execution (comma-separated indices) |
| `-M` | string | "" | Set memory values before execution (index=value pairs) |

### Configuration Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--max-instr` | uint | 100000 | Maximum instructions to execute |
| `--stack-size` | int | 256 | Stack size for execution |
| `--memory-size` | int | 256 | Memory size (number of addressable locations) |

### Mode Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-d` | bool | false | Disassemble mode (reads assembly, outputs disassembled form) |

## Operation Modes

### 1. Assemble Only (Default)

Assembles the input file and reports success/failure.

```bash
stackvm-asm program.asm
```

**Output:**
```
Assembly successful: 15 instructions
Use -o to save output or -r to run
```

### 2. Assemble and Run

Assembles and executes the program.

```bash
stackvm-asm -r program.asm
```

**Output:**
```
Assembly successful: 15 instructions

=== Executing Program ===

=== Execution Complete ===
Status: HALTED
```

### 3. Assemble and Save

Assembles and saves the disassembled output to a file.

```bash
stackvm-asm -o output.asm program.asm
```

**Output:**
```
Assembly successful: 15 instructions
Output written to output.asm
```

### 4. Disassemble

Reads assembly, assembles it, then outputs the disassembled form.

```bash
stackvm-asm -d program.asm
```

**Output:**
```
PUSH 10
PUSH 5
ADD
HALT
```

### 5. Run with Statistics

Executes and displays performance metrics.

```bash
stackvm-asm -r -s program.asm
```

**Output:**
```
Assembly successful: 15 instructions

=== Executing Program ===

=== Execution Complete ===
Status: HALTED

=== Statistics ===
Instructions executed: 49
Final stack depth: 1
Execution time: 19.542µs
Instructions/sec: 2507420
```

## Memory Features

### Memory Initialization (`-M`)

Pre-initialize memory locations before program execution.

**Syntax:**
```bash
-M "index=value,index=value,..."
```

**Type Detection:**
- Values without decimal point → Integer
- Values with decimal point → Float

**Examples:**

```bash
# Set memory[0]=42 and memory[1]=100
stackvm-asm -r -M "0=42,1=100" program.asm

# Mix integers and floats
stackvm-asm -r -M "0=10,1=3.14,2=100" program.asm

# Spaces are allowed (quote the argument)
stackvm-asm -r -M "0 = 42, 1 = 100" program.asm
```

**Use Case:**
```asm
; Program expects inputs in memory[0] and memory[1]
    LOAD 0
    LOAD 1
    ADD
    STORE 2
    HALT
```

```bash
# Run with different inputs
stackvm-asm -r -M "0=10,1=5" program.asm
stackvm-asm -r -M "0=100,1=50" program.asm
```

**Error Handling:**
- Invalid format → Warning, skip entry
- Out-of-bounds index → Warning, skip entry
- Invalid value → Warning, skip entry
- Execution continues with valid entries

### Memory Inspection (`-m`)

Display memory values after program execution.

**Syntax:**
```bash
-m "index,index,..."
```

**Examples:**

```bash
# Show memory[0]
stackvm-asm -r -m 0 program.asm

# Show multiple locations
stackvm-asm -r -m 0,1,2,3,4 program.asm

# Spaces allowed
stackvm-asm -r -m "0, 1, 2" program.asm
```

**Output Format:**
```
=== Memory ===
memory[0] = 42 (int)
memory[1] = 3.14 (float)
memory[2] = nil
memory[3] = true (bool)
```

**Error Handling:**
- Invalid index → Warning, skip
- Out-of-bounds → Warning, skip
- Other indices still displayed

### Combined Memory Operations

Initialize and inspect memory together:

```bash
stackvm-asm -r -M "0=10,1=5" -m 0,1,2 program.asm
```

**Output:**
```
Assembly successful: 5 instructions

=== Executing Program ===

=== Execution Complete ===
Status: HALTED

=== Memory ===
memory[0] = 10 (int)
memory[1] = 5 (int)
memory[2] = 15 (float)
```

## Output Format

### Standard Output (stdout)

- Disassembled code (when using `-d` without `-o`)
- Help text (when using `-h`)
- Version info (when using `-v`)

### Standard Error (stderr)

All other output goes to stderr:
- Assembly status messages
- Execution status
- Statistics
- Memory inspection
- Warnings and errors

This allows filtering:
```bash
# Capture only disassembled code
stackvm-asm -d program.asm > output.asm

# Capture statistics to file
stackvm-asm -r -s program.asm 2> stats.txt
```

### Quiet Mode

With `-q` flag, suppresses:
- Assembly success messages
- Execution status messages
- "Use -o or -r" hints

Still shows:
- Errors (always displayed)
- Statistics (if `-s` specified)
- Memory values (if `-m` specified)

## Error Handling

### Assembly Errors

```bash
$ stackvm-asm invalid.asm
Error: assembly failed: line 5: unknown opcode INVALID
```

**Exit code:** 1

### Execution Errors

```bash
$ stackvm-asm -r program.asm
Assembly successful: 5 instructions

=== Executing Program ===

Error: execution failed: stack underflow
```

**Exit code:** 1

### Warnings

Warnings are non-fatal and allow execution to continue:

```bash
$ stackvm-asm -r -M "0=10,invalid,999=42" program.asm
Assembly successful: 5 instructions

=== Executing Program ===

Warning: invalid memory init format 'invalid' (expected index=value)
Warning: cannot initialize memory[999]: invalid memory address

=== Execution Complete ===
Status: HALTED
```

**Exit code:** 0 (warnings don't cause failure)

### File Errors

```bash
$ stackvm-asm nonexistent.asm
Error: failed to read input file: open nonexistent.asm: no such file or directory
```

**Exit code:** 1

## Examples

### Basic Examples

**Simple execution:**
```bash
stackvm-asm -r simple_add.asm
```

**With statistics:**
```bash
stackvm-asm -r -s factorial.asm
```

**Quiet mode:**
```bash
stackvm-asm -r -q program.asm
```

### Memory Examples

**Initialize input values:**
```bash
stackvm-asm -r -M "0=100" factorial_from_memory.asm
```

**Inspect results:**
```bash
stackvm-asm -r -m 0,1,2 array_processor.asm
```

**Full workflow:**
```bash
stackvm-asm -r -s -M "0=10,1=5" -m 0,1,2 calculator.asm
```

### Configuration Examples

**Large programs:**
```bash
stackvm-asm -r --max-instr 1000000 large_program.asm
```

**Large stack:**
```bash
stackvm-asm -r --stack-size 1024 deep_recursion.asm
```

**Large memory:**
```bash
stackvm-asm -r --memory-size 1024 -M "0=10" data_processing.asm
```

### Pipeline Examples

**Save assembled output:**
```bash
stackvm-asm -o output.asm program.asm
```

**Capture disassembled to file:**
```bash
stackvm-asm -d program.asm > disassembled.asm
```

**Capture statistics:**
```bash
stackvm-asm -r -s program.asm 2> stats.txt
```

**Combine with other tools:**
```bash
# Run multiple tests
for i in testdata/programs/*.asm; do
    echo "Testing $i"
    stackvm-asm -r -q "$i" || echo "Failed: $i"
done
```

### Advanced Examples

**Benchmarking:**
```bash
# Run with stats, capture output
stackvm-asm -r -s -M "0=1000" benchmark.asm 2>&1 | grep "Instructions/sec"
```

**Testing different inputs:**
```bash
# Test factorial with different values
for n in 5 10 15 20; do
    echo "Factorial($n):"
    stackvm-asm -r -M "0=$n" -m 1 factorial.asm 2>&1 | grep "memory\[1\]"
done
```

**Memory debugging:**
```bash
# Initialize, run, inspect all memory
stackvm-asm -r -M "0=1,1=2,2=3,3=4,4=5" \
    -m 0,1,2,3,4,5,6,7,8,9 \
    array_processor.asm
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (assembly, execution, file I/O, etc.) |

**Note:** Warnings do not cause non-zero exit codes.

### Checking Exit Codes

```bash
# Shell script example
if stackvm-asm -r program.asm; then
    echo "Success"
else
    echo "Failed"
fi
```

```bash
# Capture exit code
stackvm-asm -r program.asm
EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "Program executed successfully"
fi
```

## Appendix: Complete Flag Summary

### Alphabetical Reference

```
-M string
    Set memory values before execution (index=value pairs, e.g., 0=42,1=3.14)

-d
    Disassemble mode

-h
    Show help

-m string
    Show memory values (comma-separated indices, e.g., 0,1,2)

-max-instr uint
    Maximum instructions to execute (default 100000)

-memory-size int
    Memory size for execution (default 256)

-o string
    Output file (default: stdout)

-q
    Quiet mode (suppress non-error output)

-r
    Run the program after assembling

-s
    Show execution statistics

-stack-size int
    Stack size for execution (default 256)

-v
    Show version
```

## Version History

### 1.0.0 (Current)

- Initial release
- Assembly and execution
- Disassembly support
- Memory initialization (`-M`)
- Memory inspection (`-m`)
- Execution statistics
- Configurable limits
- Comprehensive error handling

---

**For language specification, see:** [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md)
**For programming guide, see:** [GETTING_STARTED.md](GETTING_STARTED.md)
**For VM specification, see:** [SPEC.md](SPEC.md)
