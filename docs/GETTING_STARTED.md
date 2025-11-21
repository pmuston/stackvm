# Getting Started with StackVM Assembly

A practical guide to programming the Stack Virtual Machine.

## Table of Contents

1. [Introduction](#introduction)
2. [Your First Program](#your-first-program)
3. [Understanding the Stack](#understanding-the-stack)
4. [Basic Operations](#basic-operations)
5. [Control Flow](#control-flow)
6. [Working with Memory](#working-with-memory)
7. [Complete Examples](#complete-examples)
8. [Best Practices](#best-practices)

---

## Introduction

StackVM is a stack-based virtual machine, similar to the Java Virtual Machine or Forth. Instead of using registers like traditional CPUs, all operations work with a stack of values.

### Key Concepts

- **Stack**: A Last-In-First-Out (LIFO) data structure where operations take place
- **Instructions**: Simple commands that manipulate the stack
- **Memory**: Separate storage for persistent values
- **Labels**: Named positions in your program for jumps and loops

### Running Programs

```bash
# Assemble and run a program
stackvm-asm -r myprogram.asm

# Run with statistics
stackvm-asm -r -s myprogram.asm

# Disassemble to see structure
stackvm-asm -d myprogram.asm
```

---

## Your First Program

Let's start with the classic "Hello World" equivalent - adding two numbers:

```assembly
; my_first_program.asm
; Adds 10 + 5

PUSH 10         ; Push 10 onto the stack
PUSH 5          ; Push 5 onto the stack
ADD             ; Pop two values, add them, push result
HALT            ; Stop execution
```

**What happens:**
1. Stack after `PUSH 10`: `[10]`
2. Stack after `PUSH 5`: `[10, 5]` (5 is on top)
3. Stack after `ADD`: `[15]` (10 + 5)
4. `HALT` stops the program

### Comments

Use `;` or `#` for comments:

```assembly
; This is a comment
# This is also a comment
PUSH 42         ; Inline comment
```

---

## Understanding the Stack

The stack is the core of StackVM. All arithmetic and logic operations work with stack values.

### Stack Visualization

Think of the stack as a pile of plates. You can only add (push) or remove (pop) from the top:

```
    [Top]
    -----
     15     <- Most recent value (top of stack)
     42
     10
    -----
   [Bottom]
```

### Basic Stack Operations

#### PUSH - Add a Value

```assembly
PUSH 42         ; Push float 42.0 onto stack
PUSHI 42        ; Push integer 42 onto stack
```

**Stack effect:** `→ value`

#### POP - Remove a Value

```assembly
PUSH 10
PUSH 20
POP             ; Remove 20, stack now has just [10]
```

**Stack effect:** `value →`

#### DUP - Duplicate Top Value

```assembly
PUSH 5
DUP             ; Stack: [5, 5]
```

**Stack effect:** `a → a a`

#### SWAP - Exchange Top Two Values

```assembly
PUSH 10
PUSH 20
SWAP            ; Stack: [10, 20] → [20, 10]
```

**Stack effect:** `a b → b a`

#### OVER - Copy Second Value to Top

```assembly
PUSH 10
PUSH 20
OVER            ; Stack: [10, 20, 10]
```

**Stack effect:** `a b → a b a`

#### ROT - Rotate Top Three Values

```assembly
PUSH 1
PUSH 2
PUSH 3
ROT             ; Stack: [1, 2, 3] → [2, 3, 1]
```

**Stack effect:** `a b c → b c a`

### Example: Using Stack Operations

```assembly
; Calculate (10 + 5) * 2

PUSH 10         ; [10]
PUSH 5          ; [10, 5]
ADD             ; [15]
PUSH 2          ; [15, 2]
MUL             ; [30]
HALT
```

---

## Basic Operations

### Arithmetic

All arithmetic operations pop two values, perform the operation, and push the result.

```assembly
; Addition
PUSH 10
PUSH 5
ADD             ; Result: 15

; Subtraction (10 - 5)
PUSH 10
PUSH 5
SUB             ; Result: 5

; Multiplication
PUSH 6
PUSH 7
MUL             ; Result: 42

; Division
PUSH 20
PUSH 4
DIV             ; Result: 5.0

; Modulo (remainder)
PUSH 17
PUSH 5
MOD             ; Result: 2
```

**Stack effects:**
- `ADD, SUB, MUL, DIV, MOD`: `a b → (a op b)`

### Unary Operations

These operations work on a single value:

```assembly
; Negate
PUSH 42
NEG             ; Result: -42

; Absolute value
PUSH -10
ABS             ; Result: 10

; Increment
PUSH 5
INC             ; Result: 6

; Decrement
PUSH 5
DEC             ; Result: 4
```

**Stack effects:**
- `NEG, ABS`: `a → (op a)`
- `INC, DEC`: `a → (a±1)`

### Example: Complex Expression

Calculate: `(a * b) + (c / d)` where a=10, b=5, c=20, d=4

```assembly
; First part: a * b
PUSH 10         ; a
PUSH 5          ; b
MUL             ; a*b = 50

; Second part: c / d
PUSH 20         ; c
PUSH 4          ; d
DIV             ; c/d = 5.0

; Combine
ADD             ; 50 + 5 = 55
HALT
```

### Logic Operations

Boolean operations treat 0 as false, non-zero as true.

```assembly
; AND (both must be true)
PUSH 1          ; true
PUSH 1          ; true
AND             ; Result: 1 (true)

; OR (at least one must be true)
PUSH 1          ; true
PUSH 0          ; false
OR              ; Result: 1 (true)

; NOT (invert)
PUSH 0          ; false
NOT             ; Result: 1 (true)

; XOR (exactly one must be true)
PUSH 1          ; true
PUSH 1          ; true
XOR             ; Result: 0 (false)
```

**Stack effects:**
- `AND, OR, XOR`: `a b → (a op b)`
- `NOT`: `a → (!a)`

### Comparison Operations

These operations compare two values and push 1 (true) or 0 (false).

```assembly
; Equal
PUSH 5
PUSH 5
EQ              ; Result: 1 (true)

; Not equal
PUSH 5
PUSH 3
NE              ; Result: 1 (true)

; Greater than
PUSH 10
PUSH 5
GT              ; Result: 1 (true, because 10 > 5)

; Less than
PUSH 3
PUSH 8
LT              ; Result: 1 (true, because 3 < 8)

; Greater or equal
PUSH 5
PUSH 5
GE              ; Result: 1 (true)

; Less or equal
PUSH 3
PUSH 5
LE              ; Result: 1 (true)
```

**Stack effects:** `a b → (a cmp b)`

---

## Control Flow

### Labels

Labels mark positions in your code:

```assembly
START:          ; Define a label
    PUSH 1
    HALT
```

### Unconditional Jump

```assembly
    PUSH 1
    JMP SKIP    ; Always jump to SKIP

    PUSH 999    ; This is skipped

SKIP:
    PUSH 2
    HALT
```

### Conditional Jumps

#### JMPZ - Jump if Zero (false)

```assembly
    PUSH 0      ; false condition
    JMPZ TRUE_BRANCH

    PUSH 100    ; Skipped because condition was 0

TRUE_BRANCH:
    PUSH 200
    HALT
```

#### JMPNZ - Jump if Not Zero (true)

```assembly
    PUSH 1      ; true condition
    JMPNZ TRUE_BRANCH

    PUSH 100    ; Skipped because condition was non-zero

TRUE_BRANCH:
    PUSH 200
    HALT
```

### If-Then-Else Pattern

```assembly
; if (x > 10) then y = 100 else y = 200

    PUSH 15     ; x = 15
    PUSH 10
    GT          ; x > 10? → 1 (true)
    JMPZ ELSE

THEN:
    PUSH 100    ; y = 100
    JMP END

ELSE:
    PUSH 200    ; y = 200

END:
    HALT
```

### While Loop Pattern

```assembly
; while (counter < 5) { counter++ }

    PUSHI 0     ; counter = 0

LOOP:
    DUP         ; Copy counter
    PUSHI 5     ; Push limit
    LT          ; counter < 5?
    JMPZ END    ; Exit if false

    INC         ; counter++
    JMP LOOP    ; Continue loop

END:
    HALT
```

### For Loop Pattern

```assembly
; for (i = 1; i <= 5; i++) { sum += i }

    PUSHI 0     ; sum = 0
    PUSHI 1     ; i = 1

LOOP:
    OVER        ; Copy i to top
    PUSHI 5
    GT          ; i > 5?
    JMPNZ END   ; Exit if true

    OVER        ; Copy i
    ADD         ; sum += i

    SWAP        ; Get i to top
    INC         ; i++
    SWAP        ; Put sum back on top

    JMP LOOP

END:
    SWAP        ; Get sum to top
    POP         ; Remove i
    HALT
```

---

## Working with Memory

Memory provides persistent storage separate from the stack. Memory is indexed starting at 0.

### STORE - Save to Memory

```assembly
PUSH 42
STORE 0         ; Store 42 at memory[0]
```

**Stack effect:** `value →`

### LOAD - Load from Memory

```assembly
LOAD 0          ; Load value from memory[0]
```

**Stack effect:** `→ value`

### Example: Using Variables

```assembly
; Calculate: c = a + b
; Where a=10, b=5, c is stored in memory[2]

    ; Initialize variables
    PUSH 10
    STORE 0     ; a = 10

    PUSH 5
    STORE 1     ; b = 5

    ; Calculate c = a + b
    LOAD 0      ; Load a
    LOAD 1      ; Load b
    ADD         ; a + b

    STORE 2     ; c = result

    ; Show result
    LOAD 2
    HALT
```

### Dynamic Memory Access

Use `LOADD` and `STORED` when the address is on the stack:

```assembly
; Load from memory[index] where index is calculated

    PUSHI 5     ; index = 5
    LOADD       ; Load memory[5]

; Store to memory[index]
    PUSH 42     ; value to store
    PUSHI 3     ; index
    STORED      ; memory[3] = 42
```

**Stack effects:**
- `LOADD`: `index → value`
- `STORED`: `value index →`

### Example: Array Sum

```assembly
; Sum array: memory[0..4]

    ; Initialize array
    PUSH 10
    STORE 0
    PUSH 20
    STORE 1
    PUSH 30
    STORE 2
    PUSH 40
    STORE 3
    PUSH 50
    STORE 4

    ; sum = 0, i = 0
    PUSHI 0     ; sum
    PUSHI 0     ; i

LOOP:
    DUP         ; Copy i
    PUSHI 5     ; array length
    GE          ; i >= 5?
    JMPNZ END

    DUP         ; Copy i
    LOADD       ; Load memory[i]

    ; Add to sum (sum is at position 2 on stack)
    SWAP        ; Move i down
    SWAP        ; Get sum to top
    ADD         ; sum += value
    SWAP        ; Put i back on top

    INC         ; i++
    JMP LOOP

END:
    POP         ; Remove i
    HALT        ; sum is on top of stack
```

---

## Complete Examples

### Example 1: Factorial

Calculate n! (factorial) for n=5:

```assembly
; factorial.asm
; Calculate 5! = 5 * 4 * 3 * 2 * 1 = 120

    PUSHI 5     ; n = 5
    PUSHI 1     ; result = 1

LOOP:
    OVER        ; Copy n
    PUSHI 1
    LE          ; n <= 1?
    JMPNZ DONE

    OVER        ; Copy n
    MUL         ; result *= n

    SWAP        ; Get n to top
    DEC         ; n--
    SWAP        ; Put result back

    JMP LOOP

DONE:
    SWAP        ; Get result to top
    POP         ; Remove n
    HALT        ; Result: 120
```

### Example 2: Fibonacci

Calculate the nth Fibonacci number:

```assembly
; fibonacci.asm
; Calculate Fib(10)

    PUSHI 10    ; n = 10
    PUSHI 0     ; fib(0) = 0
    PUSHI 1     ; fib(1) = 1
    PUSHI 2     ; i = 2

LOOP:
    ; Check if i > n
    DUP         ; Copy i
    PUSH 3
    LOADD       ; Load n (3 positions down)
    GT
    JMPNZ DONE

    ; Calculate next fibonacci: fib(i) = fib(i-1) + fib(i-2)
    OVER        ; Copy fib(i-1)
    PUSH 2
    LOADD       ; Get fib(i-2)
    ADD         ; fib(i) = fib(i-1) + fib(i-2)

    ; Shift values: fib(i-2) = fib(i-1), fib(i-1) = fib(i)
    SWAP        ; Get old fib(i-1) to top
    POP         ; Remove it
    SWAP        ; Position new value

    ; Increment i
    SWAP
    INC
    SWAP

    JMP LOOP

DONE:
    ; Clean up and return result
    POP         ; Remove i
    POP         ; Remove fib(i-2)
    ; fib(i-1) is the result
    HALT
```

### Example 3: Prime Number Check

Check if a number is prime:

```assembly
; isprime.asm
; Check if 17 is prime

    PUSHI 17    ; Number to test
    PUSHI 2     ; Divisor starts at 2

CHECK_LOOP:
    ; Check if divisor * divisor > n
    DUP         ; Copy divisor
    DUP         ; Copy divisor again
    MUL         ; divisor * divisor

    PUSH 2      ; Access n
    LOADD
    GT          ; divisor*divisor > n?
    JMPNZ IS_PRIME

    ; Check if n % divisor == 0
    PUSH 2      ; Get n
    LOADD
    OVER        ; Copy divisor
    MOD         ; n % divisor
    PUSHI 0
    EQ          ; == 0?
    JMPNZ NOT_PRIME

    ; Try next divisor
    INC
    JMP CHECK_LOOP

IS_PRIME:
    POP         ; Remove divisor
    POP         ; Remove n
    PUSHI 1     ; Return true
    HALT

NOT_PRIME:
    POP         ; Remove divisor
    POP         ; Remove n
    PUSHI 0     ; Return false
    HALT
```

### Example 4: Mathematical Functions

Using built-in math functions:

```assembly
; math_demo.asm
; Pythagorean theorem: c = sqrt(a² + b²)

    ; Calculate 3² + 4² = 9 + 16 = 25
    PUSH 3
    DUP
    MUL         ; 3² = 9

    PUSH 4
    DUP
    MUL         ; 4² = 16

    ADD         ; 9 + 16 = 25
    SQRT        ; sqrt(25) = 5.0

    HALT
```

### Example 5: Min/Max Functions

```assembly
; minmax.asm
; Find minimum and maximum of three numbers

    PUSH 42
    PUSH 17
    PUSH 99

    ; Find maximum
    OVER        ; Get second number
    OVER        ; Get third number
    MAX         ; max(17, 99) = 99
    PUSH 3
    LOADD       ; Get first number (42)
    MAX         ; max(99, 42) = 99

    ; Result: 99 is on stack
    HALT
```

### Example 6: Temperature Converter

Convert Celsius to Fahrenheit: F = (C × 9/5) + 32

```assembly
; temp_convert.asm
; Convert 25°C to Fahrenheit

    PUSH 25     ; Celsius temperature

    ; Multiply by 9
    PUSH 9
    MUL         ; 25 * 9 = 225

    ; Divide by 5
    PUSH 5
    DIV         ; 225 / 5 = 45

    ; Add 32
    PUSH 32
    ADD         ; 45 + 32 = 77°F

    HALT
```

---

## Best Practices

### 1. Use Comments Liberally

```assembly
; BAD: No comments
PUSH 10
PUSH 5
GT
JMPZ ELSE
PUSH 100
JMP END
ELSE:
PUSH 200
END:
HALT

; GOOD: Clear comments
; Check if age > 5
PUSH 10         ; age
PUSH 5          ; minimum age
GT              ; age > 5?
JMPZ ELSE       ; If not, go to else

THEN:
    PUSH 100    ; adult_price
    JMP END

ELSE:
    PUSH 200    ; child_price

END:
    HALT
```

### 2. Use Meaningful Label Names

```assembly
; BAD
L1:
    PUSH 1
    JMPZ L2
    JMP L3

; GOOD
VALIDATE_INPUT:
    PUSH 1
    JMPZ INVALID_INPUT
    JMP PROCESS_DATA
```

### 3. Document Stack State

```assembly
; Calculate area of rectangle
; Input: width and height in memory[0] and memory[1]
; Output: area on stack

    LOAD 0      ; Stack: [width]
    LOAD 1      ; Stack: [width, height]
    MUL         ; Stack: [area]
    HALT
```

### 4. Keep Functions Small

Break complex programs into logical sections:

```assembly
; Main program
MAIN:
    CALL INITIALIZE
    CALL PROCESS
    CALL CLEANUP
    HALT

INITIALIZE:
    ; Setup code here
    RET

PROCESS:
    ; Main logic here
    RET

CLEANUP:
    ; Cleanup code here
    RET
```

### 5. Use Memory for Complex Data

```assembly
; Store multiple values for later use
; Memory layout:
; [0] = width
; [1] = height
; [2] = depth
; [3] = result

    PUSH 10
    STORE 0     ; width = 10

    PUSH 20
    STORE 1     ; height = 20

    PUSH 5
    STORE 2     ; depth = 5

    ; Calculate volume = width * height * depth
    LOAD 0
    LOAD 1
    MUL
    LOAD 2
    MUL
    STORE 3     ; result = volume

    LOAD 3
    HALT
```

### 6. Validate Stack Assumptions

```assembly
; Make sure stack has expected values
; Expected: [base, height] on stack

CALCULATE_AREA:
    ; Could add assertions here
    ; In production, validate inputs

    MUL         ; area = base * height
    PUSH 2
    DIV         ; Triangle area = (base * height) / 2
    RET
```

### 7. Plan Your Stack Usage

Before writing complex operations, sketch out stack states:

```
; Goal: Calculate (a + b) * (c + d)
;
; Step 1: PUSH a       -> [a]
; Step 2: PUSH b       -> [a, b]
; Step 3: ADD          -> [a+b]
; Step 4: PUSH c       -> [a+b, c]
; Step 5: PUSH d       -> [a+b, c, d]
; Step 6: ADD          -> [a+b, c+d]
; Step 7: MUL          -> [(a+b)*(c+d)]
```

---

## Common Patterns

### Pattern: Swap Variables

```assembly
; Swap values in memory[0] and memory[1]
LOAD 0
LOAD 1
SWAP
STORE 0
SWAP
STORE 1
```

### Pattern: Conditional Assignment

```assembly
; x = (condition) ? valueA : valueB
PUSH condition
JMPZ ELSE
    PUSH valueA
    JMP END
ELSE:
    PUSH valueB
END:
STORE 0     ; x = result
```

### Pattern: Loop with Counter

```assembly
; for (i = 0; i < N; i++)
PUSHI 0         ; i = 0
LOOP:
    DUP         ; Copy i
    PUSHI N
    GE          ; i >= N?
    JMPNZ END

    ; Loop body here

    INC         ; i++
    JMP LOOP
END:
    POP         ; Remove counter
```

### Pattern: Array Initialization

```assembly
; Initialize array at memory[10..14] with value 0
PUSHI 10        ; Start index
INIT_LOOP:
    DUP
    PUSHI 15    ; End index (exclusive)
    GE
    JMPNZ INIT_DONE

    DUP         ; Copy index
    PUSHI 0     ; Value to store
    SWAP
    STORED      ; memory[index] = 0

    INC
    JMP INIT_LOOP
INIT_DONE:
    POP
```

---

## Tips for Debugging

### 1. Test Small Pieces

Build your program incrementally and test each piece:

```assembly
; Test just the addition first
PUSH 10
PUSH 5
ADD
HALT

; Then add more complexity
```

### 2. Use Strategic HALTs

Insert `HALT` statements to inspect stack at different points:

```assembly
PUSH 10
PUSH 5
ADD
HALT        ; Check: stack should be [15]

PUSH 2
MUL
HALT        ; Check: stack should be [30]
```

### 3. Run with Statistics

```bash
stackvm-asm -r -s myprogram.asm
```

This shows:
- Instructions executed
- Final stack depth
- Execution time
- Instructions per second

### 4. Disassemble Complex Programs

```bash
stackvm-asm -d myprogram.asm
```

See how the assembler interpreted your code.

---

## Next Steps

1. **Try the examples** - Run each example program to see how they work
2. **Modify examples** - Change values and see what happens
3. **Write your own** - Start with simple calculations, then build up
4. **Read the spec** - See `SPEC.md` for complete opcode reference
5. **Explore math functions** - Try `SIN`, `COS`, `LOG`, `POW`, etc.

## Additional Resources

- `SPEC.md` - Complete technical specification
- `testdata/programs/` - More example programs
- `stackvm-asm -h` - CLI tool help

---

**Happy coding with StackVM!** 🚀
