# StackVM Assembly Language Specification

**Version:** 1.0
**Status:** Stable
**Last Updated:** 2025

---

## Table of Contents

1. [Introduction](#introduction)
2. [Lexical Structure](#lexical-structure)
3. [Syntax](#syntax)
4. [Type System](#type-system)
5. [Memory Model](#memory-model)
6. [Execution Model](#execution-model)
7. [Instruction Set Reference](#instruction-set-reference)
8. [Assembler Directives](#assembler-directives)
9. [Error Conditions](#error-conditions)
10. [Conformance](#conformance)

---

## 1. Introduction

### 1.1 Purpose

This document formally specifies the StackVM assembly language. It defines:
- Lexical structure (tokens, identifiers, literals)
- Syntactic structure (statements, labels, instructions)
- Semantic behavior (instruction effects, type conversions)
- Memory and execution models

### 1.2 Scope

This specification covers the assembly language only. For VM implementation details, see `SPEC.md`.

### 1.3 Notation

In this document:
- `UPPERCASE` indicates keywords or instruction names
- `lowercase` indicates meta-syntactic variables
- `[optional]` indicates optional elements
- `{repeated}` indicates zero or more repetitions
- `|` indicates alternatives
- Stack notation: `a b → c` means "pops b, pops a, pushes c"

### 1.4 Design Principles

1. **Simplicity**: Minimal syntax, easy to parse
2. **Stack-oriented**: All operations work with the stack
3. **Case-insensitive**: `PUSH`, `push`, and `Push` are equivalent
4. **Whitespace-agnostic**: Flexible formatting
5. **Self-documenting**: Labels provide structure

---

## 2. Lexical Structure

### 2.1 Character Set

StackVM assembly uses ASCII text encoding. The following characters are significant:

- **Letters**: `a-z`, `A-Z`
- **Digits**: `0-9`
- **Operators**: `+`, `-`, `.`
- **Delimiters**: `:`, whitespace, newline
- **Comments**: `;`, `#`

### 2.2 Whitespace

Whitespace includes:
- Space (U+0020)
- Tab (U+0009)
- Carriage return (U+000D)

Whitespace is used to separate tokens but is otherwise ignored.

### 2.3 Line Terminators

Newline (U+000A) terminates statements. Multiple newlines are collapsed to one.

### 2.4 Comments

Two comment styles are supported:

```assembly
; Semicolon comment (to end of line)
# Hash comment (to end of line)
```

Comments are treated as whitespace and ignored during assembly.

**Example:**
```assembly
PUSH 10     ; Push the value 10
# This is also a comment
```

### 2.5 Identifiers

Identifiers are used for labels and (potentially) custom instruction names.

**Syntax:**
```
identifier ::= letter { letter | digit | '_' }
letter     ::= 'A'..'Z' | 'a'..'z'
digit      ::= '0'..'9'
```

**Rules:**
- Must start with a letter
- May contain letters, digits, and underscores
- Case-insensitive for opcodes
- Case-sensitive for labels (recommended to be case-insensitive in practice)
- No length limit (practical limit: 255 characters)

**Valid identifiers:**
```
START
loop_1
myFunction
PROCESS_DATA
```

**Invalid identifiers:**
```
1start      ; Cannot start with digit
my-label    ; Hyphen not allowed
loop.end    ; Period not allowed
```

### 2.6 Numeric Literals

#### 2.6.1 Integer Literals

Integer literals represent whole numbers.

**Syntax:**
```
integer ::= ['-'] digit {digit}
```

**Range:** -2,147,483,648 to 2,147,483,647 (32-bit signed)

**Examples:**
```
0
42
-10
2147483647
```

#### 2.6.2 Floating-Point Literals

Floating-point literals represent real numbers.

**Syntax:**
```
float ::= ['-'] digit {digit} '.' {digit}
```

**Precision:** 64-bit IEEE 754 double precision

**Examples:**
```
3.14
-0.5
10.0
2.718281828
```

#### 2.6.3 Literal Type Inference

The assembler infers type from syntax:
- Contains decimal point → float
- No decimal point → integer

### 2.7 Keywords

The following are reserved instruction names (case-insensitive):

**Stack Operations:**
`PUSH`, `PUSHI`, `POP`, `DUP`, `SWAP`, `OVER`, `ROT`

**Arithmetic:**
`ADD`, `SUB`, `MUL`, `DIV`, `MOD`, `NEG`, `ABS`, `INC`, `DEC`

**Logic:**
`AND`, `OR`, `NOT`, `XOR`

**Comparison:**
`EQ`, `NE`, `GT`, `LT`, `GE`, `LE`

**Memory:**
`LOAD`, `STORE`, `LOADD`, `STORED`

**Control Flow:**
`JMP`, `JMPZ`, `JMPNZ`, `CALL`, `RET`, `HALT`, `NOP`

**Math Functions:**
`SQRT`, `SIN`, `COS`, `TAN`, `ASIN`, `ACOS`, `ATAN`, `ATAN2`,
`LOG`, `LOG10`, `EXP`, `POW`, `MIN`, `MAX`, `FLOOR`, `CEIL`, `ROUND`, `TRUNC`

---

## 3. Syntax

### 3.1 Program Structure

A program consists of zero or more statements separated by newlines.

**Grammar:**
```
program    ::= {statement}
statement  ::= label_def | instruction | empty_line
empty_line ::= [comment] newline
```

### 3.2 Label Definitions

Labels mark positions in the program for jump targets.

**Syntax:**
```
label_def ::= identifier ':' [comment] newline
```

**Semantics:**
- Labels define the address of the next instruction
- Multiple labels can refer to the same address
- Labels must be unique within a program
- Labels are resolved at assembly time

**Examples:**
```assembly
START:
LOOP:
    PUSH 1
END:
    HALT
```

### 3.3 Instructions

Instructions perform operations on the stack or control program flow.

**Syntax:**
```
instruction ::= opcode [operand] [comment] newline
opcode      ::= identifier
operand     ::= number | identifier
number      ::= integer | float
```

**Semantics:**
- Opcodes are case-insensitive
- Operand type and presence depends on opcode
- Unknown opcodes cause assembly error
- Invalid operands cause assembly error

**Examples:**
```assembly
PUSH 10
PUSHI 42
ADD
STORE 5
JMP LOOP
HALT
```

### 3.4 Operand Types

Instructions use three operand types:

#### 3.4.1 No Operand

Most instructions take no operand.

**Examples:**
```assembly
ADD
HALT
DUP
```

#### 3.4.2 Numeric Operand

Used for immediate values and memory addresses.

**Instructions:**
- `PUSH value` - immediate float value
- `PUSHI value` - immediate integer value
- `LOAD address` - memory address
- `STORE address` - memory address

**Examples:**
```assembly
PUSH 3.14
PUSHI 42
LOAD 0
STORE 10
```

#### 3.4.3 Label Operand

Used for control flow targets.

**Instructions:**
- `JMP label`
- `JMPZ label`
- `JMPNZ label`
- `CALL label`

**Examples:**
```assembly
JMP START
JMPZ END
CALL FUNCTION
```

---

## 4. Type System

### 4.1 Value Types

StackVM supports five value types:

| Type | Description | Size | Range/Precision |
|------|-------------|------|-----------------|
| `Nil` | Absence of value | - | N/A |
| `Int` | Signed integer | 64-bit | -2^63 to 2^63-1 |
| `Float` | Floating-point | 64-bit | IEEE 754 double |
| `Bool` | Boolean | - | true (1) or false (0) |
| `String` | Text (reserved) | - | Future use |

### 4.2 Type Coercion

Automatic type conversions occur in arithmetic operations:

**Rules:**
1. `Int` + `Int` → `Int`
2. `Float` + `Float` → `Float`
3. `Int` + `Float` → `Float` (int promoted to float)
4. Comparison operations always return `Bool`
5. Logic operations treat non-zero as true

**Examples:**
```assembly
PUSHI 10        ; Int 10
PUSHI 5         ; Int 5
ADD             ; Result: Int 15

PUSH 3.14       ; Float 3.14
PUSHI 2         ; Int 2
MUL             ; Result: Float 6.28 (int promoted)
```

### 4.3 Truthiness

For boolean contexts (logic operations, conditional jumps):

- `0`, `0.0`, `Nil` → false
- All other values → true

**Example:**
```assembly
PUSHI 0
JMPZ TRUE_BRANCH    ; Jumps (0 is false)

PUSH 3.14
JMPZ SKIP           ; Does not jump (3.14 is true)
```

---

## 5. Memory Model

### 5.1 Stack

The stack is a Last-In-First-Out (LIFO) structure.

**Properties:**
- Maximum depth: Implementation-defined (default 256)
- Grows upward (toward higher addresses conceptually)
- Overflow/underflow cause runtime errors

**Stack Operations:**
- `PUSH` - adds value to top
- `POP` - removes value from top
- Instructions operate on top values

### 5.2 Memory

Memory is a linear array of values indexed from 0.

**Properties:**
- Size: Implementation-defined (default 256)
- Indexed from 0 to size-1
- Initially all values are `Nil`
- Access out of bounds causes runtime error

**Access Modes:**

1. **Static** - address known at assembly time:
   ```assembly
   LOAD 0      ; Load from memory[0]
   STORE 5     ; Store to memory[5]
   ```

2. **Dynamic** - address computed at runtime:
   ```assembly
   PUSHI 3     ; Address
   LOADD       ; Load from memory[3]

   PUSH 42     ; Value
   PUSHI 7     ; Address
   STORED      ; Store to memory[7]
   ```

### 5.3 Program Counter (PC)

- Points to the next instruction to execute
- Automatically incremented after each instruction
- Modified by jump instructions
- Not directly accessible to programs

### 5.4 Call Stack

**Note:** Current implementation has limited call stack support.

- `CALL` jumps to a label
- `RET` should return to caller (simplified in current version)
- Full call stack with return addresses is implementation-defined

---

## 6. Execution Model

### 6.1 Instruction Cycle

Each instruction executes in this sequence:

1. **Fetch**: Load instruction at PC
2. **Decode**: Determine operation and operand
3. **Execute**: Perform operation
4. **Advance**: Increment PC (unless jump occurred)
5. **Check**: Test for halt, errors, limits

### 6.2 Program Termination

Programs terminate when:

1. **HALT** instruction executes (normal)
2. **Error** occurs (abnormal)
3. **Instruction limit** reached (abnormal)
4. **Timeout** occurs (abnormal)
5. **Context cancelled** (abnormal)

### 6.3 Error Handling

Errors halt execution immediately and return an error code.

**Common errors:**
- Stack overflow/underflow
- Division by zero
- Invalid memory access
- Invalid instruction
- Type mismatch

### 6.4 Execution Limits

Configurable limits prevent runaway programs:

- **MaxInstructions**: Maximum instructions to execute
- **MaxStackDepth**: Maximum stack depth
- **Timeout**: Maximum execution time
- **Memory size**: Maximum memory cells

---

## 7. Instruction Set Reference

### 7.1 Instruction Format

Each instruction is documented with:

- **Name**: Instruction mnemonic
- **Opcode**: Numeric opcode value
- **Operand**: Type and meaning
- **Stack Effect**: Before → After
- **Description**: What the instruction does
- **Errors**: Possible error conditions

### 7.2 Stack Operations (Opcodes 0-15)

#### PUSH value

| Property | Value |
|----------|-------|
| Opcode | 0 |
| Operand | Float value |
| Stack | → a |
| Description | Push immediate float value onto stack |

**Example:**
```assembly
PUSH 3.14       ; Stack: [3.14]
```

---

#### PUSHI value

| Property | Value |
|----------|-------|
| Opcode | 1 |
| Operand | Integer value |
| Stack | → a |
| Description | Push immediate integer value onto stack |

**Example:**
```assembly
PUSHI 42        ; Stack: [42]
```

---

#### POP

| Property | Value |
|----------|-------|
| Opcode | 2 |
| Operand | None |
| Stack | a → |
| Description | Remove and discard top value |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 10
POP             ; Stack: []
```

---

#### DUP

| Property | Value |
|----------|-------|
| Opcode | 3 |
| Operand | None |
| Stack | a → a a |
| Description | Duplicate top value |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 5
DUP             ; Stack: [5, 5]
```

---

#### SWAP

| Property | Value |
|----------|-------|
| Opcode | 4 |
| Operand | None |
| Stack | a b → b a |
| Description | Exchange top two values |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 1
PUSH 2
SWAP            ; Stack: [2, 1]
```

---

#### OVER

| Property | Value |
|----------|-------|
| Opcode | 5 |
| Operand | None |
| Stack | a b → a b a |
| Description | Copy second value to top |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 10
PUSH 20
OVER            ; Stack: [10, 20, 10]
```

---

#### ROT

| Property | Value |
|----------|-------|
| Opcode | 6 |
| Operand | None |
| Stack | a b c → b c a |
| Description | Rotate top three values |
| Errors | Stack underflow if fewer than 3 values |

**Example:**
```assembly
PUSH 1
PUSH 2
PUSH 3
ROT             ; Stack: [2, 3, 1]
```

---

### 7.3 Arithmetic Operations (Opcodes 16-31)

#### ADD

| Property | Value |
|----------|-------|
| Opcode | 16 |
| Operand | None |
| Stack | a b → (a+b) |
| Description | Add top two values |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 10
PUSH 5
ADD             ; Result: 15
```

---

#### SUB

| Property | Value |
|----------|-------|
| Opcode | 17 |
| Operand | None |
| Stack | a b → (a-b) |
| Description | Subtract b from a |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 10
PUSH 3
SUB             ; Result: 7
```

---

#### MUL

| Property | Value |
|----------|-------|
| Opcode | 18 |
| Operand | None |
| Stack | a b → (a*b) |
| Description | Multiply top two values |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 6
PUSH 7
MUL             ; Result: 42
```

---

#### DIV

| Property | Value |
|----------|-------|
| Opcode | 19 |
| Operand | None |
| Stack | a b → (a/b) |
| Description | Divide a by b |
| Errors | Stack underflow, division by zero |

**Example:**
```assembly
PUSH 20
PUSH 4
DIV             ; Result: 5.0
```

---

#### MOD

| Property | Value |
|----------|-------|
| Opcode | 20 |
| Operand | None |
| Stack | a b → (a%b) |
| Description | Remainder of a divided by b |
| Errors | Stack underflow, division by zero |

**Example:**
```assembly
PUSH 17
PUSH 5
MOD             ; Result: 2
```

---

#### NEG

| Property | Value |
|----------|-------|
| Opcode | 21 |
| Operand | None |
| Stack | a → (-a) |
| Description | Negate top value |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 42
NEG             ; Result: -42
```

---

#### ABS

| Property | Value |
|----------|-------|
| Opcode | 22 |
| Operand | None |
| Stack | a → \|a\| |
| Description | Absolute value |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH -10
ABS             ; Result: 10
```

---

#### INC

| Property | Value |
|----------|-------|
| Opcode | 23 |
| Operand | None |
| Stack | a → (a+1) |
| Description | Increment by 1 |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 5
INC             ; Result: 6
```

---

#### DEC

| Property | Value |
|----------|-------|
| Opcode | 24 |
| Operand | None |
| Stack | a → (a-1) |
| Description | Decrement by 1 |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 5
DEC             ; Result: 4
```

---

### 7.4 Logic Operations (Opcodes 32-39)

Logic operations treat 0 as false and non-zero as true. Results are 1 (true) or 0 (false).

#### AND

| Property | Value |
|----------|-------|
| Opcode | 32 |
| Operand | None |
| Stack | a b → (a && b) |
| Description | Logical AND |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 1
PUSH 1
AND             ; Result: 1 (true)
```

---

#### OR

| Property | Value |
|----------|-------|
| Opcode | 33 |
| Operand | None |
| Stack | a b → (a \|\| b) |
| Description | Logical OR |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 1
PUSH 0
OR              ; Result: 1 (true)
```

---

#### NOT

| Property | Value |
|----------|-------|
| Opcode | 34 |
| Operand | None |
| Stack | a → (!a) |
| Description | Logical NOT |
| Errors | Stack underflow if empty |

**Example:**
```assembly
PUSH 0
NOT             ; Result: 1 (true)
```

---

#### XOR

| Property | Value |
|----------|-------|
| Opcode | 35 |
| Operand | None |
| Stack | a b → (a XOR b) |
| Description | Logical XOR (exclusive OR) |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 1
PUSH 1
XOR             ; Result: 0 (false - both true)
```

---

### 7.5 Comparison Operations (Opcodes 40-47)

Comparison operations return 1 (true) or 0 (false).

#### EQ

| Property | Value |
|----------|-------|
| Opcode | 40 |
| Operand | None |
| Stack | a b → (a == b) |
| Description | Test equality |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 5
PUSH 5
EQ              ; Result: 1 (true)
```

---

#### NE

| Property | Value |
|----------|-------|
| Opcode | 41 |
| Operand | None |
| Stack | a b → (a != b) |
| Description | Test inequality |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 5
PUSH 3
NE              ; Result: 1 (true)
```

---

#### GT

| Property | Value |
|----------|-------|
| Opcode | 42 |
| Operand | None |
| Stack | a b → (a > b) |
| Description | Greater than |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 10
PUSH 5
GT              ; Result: 1 (true, 10 > 5)
```

---

#### LT

| Property | Value |
|----------|-------|
| Opcode | 43 |
| Operand | None |
| Stack | a b → (a < b) |
| Description | Less than |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 3
PUSH 8
LT              ; Result: 1 (true, 3 < 8)
```

---

#### GE

| Property | Value |
|----------|-------|
| Opcode | 44 |
| Operand | None |
| Stack | a b → (a >= b) |
| Description | Greater than or equal |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 5
PUSH 5
GE              ; Result: 1 (true, 5 >= 5)
```

---

#### LE

| Property | Value |
|----------|-------|
| Opcode | 45 |
| Operand | None |
| Stack | a b → (a <= b) |
| Description | Less than or equal |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 3
PUSH 5
LE              ; Result: 1 (true, 3 <= 5)
```

---

### 7.6 Memory Operations (Opcodes 48-55)

#### LOAD address

| Property | Value |
|----------|-------|
| Opcode | 48 |
| Operand | Memory address (integer) |
| Stack | → value |
| Description | Load value from memory[address] |
| Errors | Invalid address |

**Example:**
```assembly
LOAD 0          ; Load memory[0] onto stack
```

---

#### STORE address

| Property | Value |
|----------|-------|
| Opcode | 49 |
| Operand | Memory address (integer) |
| Stack | value → |
| Description | Store top value to memory[address] |
| Errors | Stack underflow, invalid address |

**Example:**
```assembly
PUSH 42
STORE 5         ; memory[5] = 42
```

---

#### LOADD

| Property | Value |
|----------|-------|
| Opcode | 50 |
| Operand | None |
| Stack | address → value |
| Description | Load from memory[address] where address is on stack |
| Errors | Stack underflow, invalid address |

**Example:**
```assembly
PUSHI 3
LOADD           ; Load memory[3]
```

---

#### STORED

| Property | Value |
|----------|-------|
| Opcode | 51 |
| Operand | None |
| Stack | value address → |
| Description | Store value to memory[address] |
| Errors | Stack underflow, invalid address |

**Example:**
```assembly
PUSH 42
PUSHI 7
STORED          ; memory[7] = 42
```

---

### 7.7 Control Flow Operations (Opcodes 56-63)

#### JMP label

| Property | Value |
|----------|-------|
| Opcode | 56 |
| Operand | Label name or address |
| Stack | - |
| Description | Unconditional jump to label |
| Errors | Unresolved label |

**Example:**
```assembly
JMP START
```

---

#### JMPZ label

| Property | Value |
|----------|-------|
| Opcode | 57 |
| Operand | Label name or address |
| Stack | a → |
| Description | Jump if a is zero/false |
| Errors | Stack underflow, unresolved label |

**Example:**
```assembly
PUSH 0
JMPZ BRANCH     ; Jumps because 0 is false
```

---

#### JMPNZ label

| Property | Value |
|----------|-------|
| Opcode | 58 |
| Operand | Label name or address |
| Stack | a → |
| Description | Jump if a is non-zero/true |
| Errors | Stack underflow, unresolved label |

**Example:**
```assembly
PUSH 1
JMPNZ BRANCH    ; Jumps because 1 is true
```

---

#### CALL label

| Property | Value |
|----------|-------|
| Opcode | 59 |
| Operand | Label name or address |
| Stack | - |
| Description | Call subroutine at label |
| Errors | Unresolved label |

**Note:** Current implementation provides simplified CALL without full return address stack.

**Example:**
```assembly
CALL FUNCTION
```

---

#### RET

| Property | Value |
|----------|-------|
| Opcode | 60 |
| Operand | None |
| Stack | - |
| Description | Return from subroutine |

**Note:** Current implementation treats RET as HALT.

**Example:**
```assembly
FUNCTION:
    PUSH 42
    RET
```

---

#### HALT

| Property | Value |
|----------|-------|
| Opcode | 61 |
| Operand | None |
| Stack | - |
| Description | Stop program execution |

**Example:**
```assembly
PUSH 42
HALT
```

---

#### NOP

| Property | Value |
|----------|-------|
| Opcode | 62 |
| Operand | None |
| Stack | - |
| Description | No operation (do nothing) |

**Example:**
```assembly
NOP             ; Does nothing
```

---

### 7.8 Math Functions (Opcodes 64-79)

#### SQRT

| Property | Value |
|----------|-------|
| Opcode | 64 |
| Operand | None |
| Stack | a → √a |
| Description | Square root |
| Errors | Stack underflow, negative value |

**Example:**
```assembly
PUSH 16
SQRT            ; Result: 4.0
```

---

#### SIN

| Property | Value |
|----------|-------|
| Opcode | 65 |
| Operand | None |
| Stack | a → sin(a) |
| Description | Sine (radians) |
| Errors | Stack underflow |

**Example:**
```assembly
PUSH 0
SIN             ; Result: 0.0
```

---

#### COS

| Property | Value |
|----------|-------|
| Opcode | 66 |
| Operand | None |
| Stack | a → cos(a) |
| Description | Cosine (radians) |
| Errors | Stack underflow |

**Example:**
```assembly
PUSH 0
COS             ; Result: 1.0
```

---

#### TAN

| Property | Value |
|----------|-------|
| Opcode | 67 |
| Operand | None |
| Stack | a → tan(a) |
| Description | Tangent (radians) |
| Errors | Stack underflow |

---

#### MIN

| Property | Value |
|----------|-------|
| Opcode | 73 |
| Operand | None |
| Stack | a b → min(a, b) |
| Description | Minimum of two values |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 5
PUSH 10
MIN             ; Result: 5
```

---

#### MAX

| Property | Value |
|----------|-------|
| Opcode | 74 |
| Operand | None |
| Stack | a b → max(a, b) |
| Description | Maximum of two values |
| Errors | Stack underflow if fewer than 2 values |

**Example:**
```assembly
PUSH 5
PUSH 10
MAX             ; Result: 10
```

---

#### FLOOR

| Property | Value |
|----------|-------|
| Opcode | 75 |
| Operand | None |
| Stack | a → ⌊a⌋ |
| Description | Round down to integer |
| Errors | Stack underflow |

**Example:**
```assembly
PUSH 3.7
FLOOR           ; Result: 3.0
```

---

#### CEIL

| Property | Value |
|----------|-------|
| Opcode | 76 |
| Operand | None |
| Stack | a → ⌈a⌉ |
| Description | Round up to integer |
| Errors | Stack underflow |

**Example:**
```assembly
PUSH 3.2
CEIL            ; Result: 4.0
```

---

#### ROUND

| Property | Value |
|----------|-------|
| Opcode | 77 |
| Operand | None |
| Stack | a → round(a) |
| Description | Round to nearest integer |
| Errors | Stack underflow |

**Example:**
```assembly
PUSH 3.5
ROUND           ; Result: 4.0
```

---

### 7.9 Custom Instructions (Opcodes 128-255)

Opcodes 128-255 are reserved for custom, user-defined instructions.

**Properties:**
- Must be registered via InstructionRegistry
- Can have custom semantics
- May or may not use operands
- Behavior defined by implementation

**Example:**
```assembly
; Assuming DOUBLE (opcode 128) is registered
PUSH 5
DOUBLE          ; Result: 10 (if DOUBLE multiplies by 2)
```

---

## 8. Assembler Directives

**Note:** Current version has no assembler directives. This section is reserved for future extensions.

Potential future directives:
- `.data` - Data section
- `.org` - Set origin address
- `.align` - Alignment
- `.include` - File inclusion

---

## 9. Error Conditions

### 9.1 Assembly-Time Errors

Errors detected during assembly:

| Error | Cause | Example |
|-------|-------|---------|
| Unknown opcode | Invalid instruction name | `BADOP` |
| Unresolved label | Jump to undefined label | `JMP MISSING` |
| Invalid operand | Wrong operand type | `PUSH LABEL` |
| Invalid number | Malformed literal | `PUSH 3.14.15` |
| Duplicate label | Label defined twice | Two `START:` |
| Syntax error | Invalid syntax | `PUSH` (missing operand) |

### 9.2 Runtime Errors

Errors detected during execution:

| Error | Cause | Recovery |
|-------|-------|----------|
| Stack overflow | Stack exceeds max depth | Increase stack size |
| Stack underflow | Pop from empty stack | Fix program logic |
| Division by zero | Divide/mod by 0 | Check divisor |
| Invalid address | Memory out of bounds | Check address |
| Instruction limit | Too many instructions | Increase limit or fix infinite loop |
| Timeout | Execution too long | Increase timeout |
| Type mismatch | Invalid type operation | Ensure correct types |

---

## 10. Conformance

### 10.1 Conforming Assembler

A conforming assembler must:

1. Accept all syntactically valid programs
2. Reject invalid programs with clear error messages
3. Generate correct bytecode for valid programs
4. Resolve labels correctly
5. Support all standard opcodes (0-79)
6. Be case-insensitive for opcodes
7. Support both comment styles

### 10.2 Conforming Implementation

A conforming VM implementation must:

1. Execute all standard opcodes correctly
2. Enforce stack and memory bounds
3. Detect and report errors
4. Support configurable limits
5. Provide execution statistics
6. Handle all value types correctly

### 10.3 Extensions

Implementations may provide extensions:

1. Custom opcodes (128-255)
2. Additional value types
3. Assembler directives
4. Optimization passes
5. Debugging features

Extensions must not conflict with standard behavior.

---

## Appendix A: Grammar Summary

Complete EBNF grammar:

```ebnf
program        ::= {statement}
statement      ::= label_def | instruction | empty_line
label_def      ::= identifier ':' [comment] newline
instruction    ::= opcode [operand] [comment] newline
empty_line     ::= [comment] newline

opcode         ::= identifier
operand        ::= number | identifier
number         ::= integer | float

identifier     ::= letter {letter | digit | '_'}
integer        ::= ['-'] digit {digit}
float          ::= ['-'] digit {digit} '.' {digit}

letter         ::= 'A'..'Z' | 'a'..'z'
digit          ::= '0'..'9'

comment        ::= ';' {any_char} newline
                 | '#' {any_char} newline
newline        ::= '\n'
```

---

## Appendix B: Opcode Quick Reference

| Range | Category | Opcodes |
|-------|----------|---------|
| 0-15 | Stack | PUSH, PUSHI, POP, DUP, SWAP, OVER, ROT |
| 16-31 | Arithmetic | ADD, SUB, MUL, DIV, MOD, NEG, ABS, INC, DEC |
| 32-39 | Logic | AND, OR, NOT, XOR |
| 40-47 | Comparison | EQ, NE, GT, LT, GE, LE |
| 48-55 | Memory | LOAD, STORE, LOADD, STORED |
| 56-63 | Control | JMP, JMPZ, JMPNZ, CALL, RET, HALT, NOP |
| 64-79 | Math | SQRT, SIN, COS, TAN, MIN, MAX, FLOOR, CEIL, ROUND, etc. |
| 128-255 | Custom | User-defined |

---

## Appendix C: Example Programs

### C.1 Hello World (Stack VM Style)

```assembly
; Push a number and halt
PUSH 42
HALT
```

### C.2 Arithmetic Expression

```assembly
; (10 + 5) * 2 = 30
PUSH 10
PUSH 5
ADD
PUSH 2
MUL
HALT
```

### C.3 Conditional Branch

```assembly
; If x > 10 then y = 1 else y = 0
PUSH 15         ; x
PUSH 10
GT
JMPZ ELSE

THEN:
    PUSH 1
    JMP END

ELSE:
    PUSH 0

END:
    HALT
```

### C.4 Loop

```assembly
; Count from 0 to 4
PUSHI 0

LOOP:
    DUP
    PUSHI 5
    GE
    JMPNZ END

    INC
    JMP LOOP

END:
    HALT
```

### C.5 Function Call

```assembly
MAIN:
    PUSH 5
    CALL SQUARE
    HALT

SQUARE:
    DUP
    MUL
    RET
```

---

## Appendix D: Version History

### Version 1.0 (2025)
- Initial specification
- All standard opcodes (0-79)
- Basic assembler syntax
- Memory model
- Type system

---

## Appendix E: References

1. **StackVM VM Specification** - `SPEC.md`
2. **Getting Started Guide** - `GETTING_STARTED.md`
3. **Example Programs** - `testdata/programs/`

---

**End of Specification**
