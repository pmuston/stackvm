# StackVM Example Programs

This directory contains example programs demonstrating various features of the StackVM assembly language.

## Running Programs

```bash
# Build the CLI tool first
go build ./cmd/stackvm-asm

# Run a program
./stackvm-asm -r <program.asm>

# Run with statistics
./stackvm-asm -r -s <program.asm>

# Disassemble a program
./stackvm-asm -d <program.asm>
```

## Basic Examples

### simple_add.asm
**Difficulty:** Beginner
**Concepts:** PUSH, ADD, HALT

Adds two numbers (10 + 5).

```bash
./stackvm-asm -r testdata/programs/simple_add.asm
# Result: 15
```

### math.asm
**Difficulty:** Beginner
**Concepts:** Stack operations, SQRT, Pythagorean theorem

Calculates √(3² + 4²) = 5 using the Pythagorean theorem.

```bash
./stackvm-asm -r testdata/programs/math.asm
# Result: 5.0
```

### temperature.asm
**Difficulty:** Beginner
**Concepts:** Arithmetic operations, formula implementation

Converts 25°C to Fahrenheit (77°F) using the formula F = (C × 9/5) + 32.

```bash
./stackvm-asm -r testdata/programs/temperature.asm
# Result: 77.0
```

## Memory Examples

### memory.asm
**Difficulty:** Beginner
**Concepts:** LOAD, STORE, memory operations

Demonstrates basic memory operations with a simple calculation: c = a + b.

```bash
./stackvm-asm -r testdata/programs/memory.asm
# Result: 15
```

### array_sum.asm
**Difficulty:** Intermediate
**Concepts:** Arrays, loops, dynamic memory access

Sums an array of 5 numbers [10, 20, 30, 40, 50] = 150.

```bash
./stackvm-asm -r testdata/programs/array_sum.asm
# Result: 150
```

## Control Flow Examples

### conditional.asm
**Difficulty:** Intermediate
**Concepts:** Conditional jumps, if-then-else pattern, comparison

Demonstrates if-then-else logic: if (15 > 10) then 1 else 0.

```bash
./stackvm-asm -r testdata/programs/conditional.asm
# Result: 1
```

### loop.asm
**Difficulty:** Intermediate
**Concepts:** Loops, counters, JMPNZ

Counts from 1 to 5 using a while loop.

```bash
./stackvm-asm -r testdata/programs/loop.asm
# Result: 6 (counter value after loop exit)
```

## Advanced Examples

### factorial.asm
**Difficulty:** Intermediate
**Concepts:** Loops, stack manipulation, mathematical computation

Calculates 5! = 5 × 4 × 3 × 2 × 1 = 120.

```bash
./stackvm-asm -r -s testdata/programs/factorial.asm
# Result: 120
# Shows 49 instructions executed
```

### fibonacci_simple.asm
**Difficulty:** Advanced
**Concepts:** Dynamic programming, array building, memory indexing

Calculates the first 10 Fibonacci numbers and stores them in memory.
Returns the 10th Fibonacci number (fib(9) = 34).

Sequence: 0, 1, 1, 2, 3, 5, 8, 13, 21, 34

```bash
./stackvm-asm -r testdata/programs/fibonacci_simple.asm
# Result: 34
```

## Learning Path

1. **Start Here:**
   - simple_add.asm
   - temperature.asm
   - math.asm

2. **Memory Basics:**
   - memory.asm
   - array_sum.asm

3. **Control Flow:**
   - conditional.asm
   - loop.asm

4. **Challenge Yourself:**
   - factorial.asm
   - fibonacci_simple.asm

## Program Statistics

| Program | Instructions | Executed | Result |
|---------|-------------|----------|--------|
| simple_add.asm | 4 | 4 | 15 |
| temperature.asm | 8 | 8 | 77.0 |
| math.asm | 9 | 9 | 5.0 |
| memory.asm | 11 | 11 | 15 |
| conditional.asm | 10 | 7 | 1 |
| loop.asm | 9 | 33 | 6 |
| array_sum.asm | 27 | 51 | 150 |
| factorial.asm | 15 | 49 | 120 |
| fibonacci_simple.asm | 26 | ~140 | 34 |

## Tips

- Use `-s` flag to see execution statistics
- Use `-d` flag to see how your code is assembled
- All programs HALT at the end
- Stack depth of 1 means one value left on stack (the result)

## Next Steps

After trying these examples:

1. Modify the programs (change input values)
2. Combine concepts from different programs
3. Write your own programs
4. See `docs/GETTING_STARTED.md` for a complete programming guide
5. Check `docs/SPEC.md` for the full opcode reference

---

**Happy Learning!** 🎓
