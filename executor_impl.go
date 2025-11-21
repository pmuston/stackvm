package stackvm

import (
	"math"
	"time"
)

// executor implements the VM interface.
type executor struct {
	config     Config
	stack      []Value
	pc         int
	halted     bool
	instrCount uint32
}

// newExecutor creates a new executor with the given configuration.
func newExecutor(config Config) *executor {
	if config.StackSize <= 0 {
		config.StackSize = 256
	}
	return &executor{
		config: config,
		stack:  make([]Value, 0, config.StackSize),
	}
}

// Execute runs the program with the given memory and options.
func (e *executor) Execute(program Program, memory Memory, opts ExecuteOptions) (*Result, error) {
	startTime := time.Now()

	// Reset state
	e.stack = e.stack[:0]
	e.pc = 0
	e.halted = false
	e.instrCount = 0

	// Apply options
	maxInstructions := opts.MaxInstructions
	if maxInstructions == 0 && e.config.DefaultInstrLimit > 0 {
		maxInstructions = e.config.DefaultInstrLimit
	}

	maxStackDepth := opts.MaxStackDepth
	if maxStackDepth <= 0 {
		maxStackDepth = e.config.StackSize
	}

	// Set up context for timeout/cancellation
	ctx := opts.Context
	var deadline time.Time
	if opts.Timeout > 0 {
		deadline = startTime.Add(opts.Timeout)
	}

	// Create execution context once for the entire execution
	// This ensures UserData persists across custom instructions
	execCtx := newExecutionContext(e, memory)

	instructions := program.Instructions()

	// Main execution loop
	for !e.halted && e.pc >= 0 && e.pc < len(instructions) {
		// Check instruction limit
		if maxInstructions > 0 && e.instrCount >= maxInstructions {
			return &Result{
				InstructionCount: e.instrCount,
				StackDepth:       len(e.stack),
				ExecutionTime:    time.Since(startTime),
				Halted:           false,
				Error:            ErrInstructionLimit,
			}, ErrInstructionLimit
		}

		// Check timeout
		if !deadline.IsZero() && time.Now().After(deadline) {
			return &Result{
				InstructionCount: e.instrCount,
				StackDepth:       len(e.stack),
				ExecutionTime:    time.Since(startTime),
				Halted:           false,
				Error:            ErrTimeout,
			}, ErrTimeout
		}

		// Check context cancellation
		if ctx != nil {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				return &Result{
					InstructionCount: e.instrCount,
					StackDepth:       len(e.stack),
					ExecutionTime:    time.Since(startTime),
					Halted:           false,
					Error:            err,
				}, err
			default:
			}
		}

		// Fetch instruction
		inst := instructions[e.pc]
		e.instrCount++

		// Execute instruction
		if err := e.executeInstruction(inst, memory, maxStackDepth, execCtx); err != nil {
			return &Result{
				InstructionCount: e.instrCount,
				StackDepth:       len(e.stack),
				ExecutionTime:    time.Since(startTime),
				Halted:           e.halted,
				Error:            err,
			}, err
		}

		// Move to next instruction (unless a jump occurred or halted)
		if !e.halted {
			e.pc++
		}
	}

	// Check if we ran out of instructions without halting
	if !e.halted && e.pc >= len(instructions) {
		// Reached end of program without HALT - this is allowed
		e.halted = true
	}

	return &Result{
		InstructionCount: e.instrCount,
		StackDepth:       len(e.stack),
		ExecutionTime:    time.Since(startTime),
		Halted:           e.halted,
		Error:            nil,
	}, nil
}

// Reset clears the VM state for reuse.
func (e *executor) Reset() {
	e.stack = e.stack[:0]
	e.pc = 0
	e.halted = false
	e.instrCount = 0
}

// executeInstruction executes a single instruction.
func (e *executor) executeInstruction(inst Instruction, memory Memory, maxStackDepth int, execCtx *executionContextImpl) error {
	var err error

	switch inst.Opcode {
	// Stack operations
	case OpPUSH:
		return e.push(FloatValue(float64(inst.Operand)), maxStackDepth)
	case OpPUSHI:
		return e.push(IntValue(int64(inst.Operand)), maxStackDepth)
	case OpPOP:
		_, err = e.pop()
		return err
	case OpDUP:
		val, err := e.peek()
		if err != nil {
			return err
		}
		return e.push(val, maxStackDepth)
	case OpSWAP:
		if len(e.stack) < 2 {
			return ErrStackUnderflow
		}
		top := len(e.stack) - 1
		e.stack[top], e.stack[top-1] = e.stack[top-1], e.stack[top]
		return nil
	case OpOVER:
		if len(e.stack) < 2 {
			return ErrStackUnderflow
		}
		val := e.stack[len(e.stack)-2]
		return e.push(val, maxStackDepth)
	case OpROT:
		if len(e.stack) < 3 {
			return ErrStackUnderflow
		}
		top := len(e.stack) - 1
		e.stack[top-2], e.stack[top-1], e.stack[top] = e.stack[top-1], e.stack[top], e.stack[top-2]
		return nil

	// Arithmetic operations
	case OpADD:
		e.stack, err = opAdd(e.stack)
	case OpSUB:
		e.stack, err = opSub(e.stack)
	case OpMUL:
		e.stack, err = opMul(e.stack)
	case OpDIV:
		e.stack, err = opDiv(e.stack)
	case OpMOD:
		e.stack, err = opMod(e.stack)
	case OpNEG:
		e.stack, err = opNeg(e.stack)
	case OpABS:
		e.stack, err = opAbs(e.stack)
	case OpINC:
		e.stack, err = opInc(e.stack)
	case OpDEC:
		e.stack, err = opDec(e.stack)

	// Logic operations
	case OpAND:
		e.stack, err = opAnd(e.stack)
	case OpOR:
		e.stack, err = opOr(e.stack)
	case OpNOT:
		e.stack, err = opNot(e.stack)
	case OpXOR:
		e.stack, err = opXor(e.stack)

	// Comparison operations
	case OpEQ:
		e.stack, err = opEq(e.stack)
	case OpNE:
		e.stack, err = opNe(e.stack)
	case OpGT:
		e.stack, err = opGt(e.stack)
	case OpLT:
		e.stack, err = opLt(e.stack)
	case OpGE:
		e.stack, err = opGe(e.stack)
	case OpLE:
		e.stack, err = opLe(e.stack)

	// Math functions
	case OpSQRT:
		e.stack, err = opSqrt(e.stack)
	case OpSIN:
		e.stack, err = opSin(e.stack)
	case OpCOS:
		e.stack, err = opCos(e.stack)
	case OpTAN:
		e.stack, err = opTan(e.stack)
	case OpASIN:
		e.stack, err = opAsin(e.stack)
	case OpACOS:
		e.stack, err = opAcos(e.stack)
	case OpATAN:
		e.stack, err = opAtan(e.stack)
	case OpATAN2:
		e.stack, err = opAtan2(e.stack)
	case OpLOG:
		e.stack, err = opLog(e.stack)
	case OpLOG10:
		e.stack, err = opLog10(e.stack)
	case OpEXP:
		e.stack, err = opExp(e.stack)
	case OpPOW:
		e.stack, err = opPow(e.stack)
	case OpMIN:
		e.stack, err = opMin(e.stack)
	case OpMAX:
		e.stack, err = opMax(e.stack)
	case OpFLOOR:
		e.stack, err = opFloor(e.stack)
	case OpCEIL:
		e.stack, err = opCeil(e.stack)
	case OpROUND:
		e.stack, err = opRound(e.stack)
	case OpTRUNC:
		e.stack, err = opTrunc(e.stack)

	// Memory operations
	case OpLOAD:
		val, err := memory.Load(int(inst.Operand))
		if err != nil {
			return err
		}
		return e.push(val, maxStackDepth)
	case OpSTORE:
		val, err := e.pop()
		if err != nil {
			return err
		}
		return memory.Store(int(inst.Operand), val)
	case OpLOADD:
		addr, err := e.pop()
		if err != nil {
			return err
		}
		addrInt, err := toInt64(addr)
		if err != nil {
			return err
		}
		val, err := memory.Load(int(addrInt))
		if err != nil {
			return err
		}
		return e.push(val, maxStackDepth)
	case OpSTORED:
		val, err := e.pop()
		if err != nil {
			return err
		}
		addr, err := e.pop()
		if err != nil {
			return err
		}
		addrInt, err := toInt64(addr)
		if err != nil {
			return err
		}
		return memory.Store(int(addrInt), val)

	// Control flow
	case OpJMP:
		// Set PC to target address (subtract 1 because main loop increments)
		e.pc = int(inst.Operand) - 1
		return nil
	case OpJMPZ:
		val, err := e.pop()
		if err != nil {
			return err
		}
		if !toBool(val) {
			e.pc = int(inst.Operand) - 1
		}
		return nil
	case OpJMPNZ:
		val, err := e.pop()
		if err != nil {
			return err
		}
		if toBool(val) {
			e.pc = int(inst.Operand) - 1
		}
		return nil
	case OpCALL:
		// TODO: Implement call stack for proper CALL/RET support
		// For now, just jump to the address
		e.pc = int(inst.Operand) - 1
		return nil
	case OpRET:
		// TODO: Implement call stack for proper CALL/RET support
		// For now, just halt
		e.halted = true
		return nil
	case OpHALT:
		e.halted = true
		return nil
	case OpNOP:
		// No operation
		return nil

	default:
		// Check for custom instructions
		if inst.Opcode >= 128 && e.config.InstructionRegistry != nil {
			handler, exists := e.config.InstructionRegistry.Get(inst.Opcode)
			if exists {
				// Reuse the execution context to maintain UserData across instructions
				return handler.Execute(execCtx, inst.Operand)
			}
		}
		return ErrInvalidOpcode
	}

	return err
}

// Stack operation helpers

func (e *executor) push(val Value, maxStackDepth int) error {
	if len(e.stack) >= maxStackDepth {
		return ErrStackOverflow
	}
	e.stack = append(e.stack, val)
	return nil
}

func (e *executor) pop() (Value, error) {
	if len(e.stack) == 0 {
		return NilValue(), ErrStackUnderflow
	}
	val := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]
	return val, nil
}

func (e *executor) peek() (Value, error) {
	if len(e.stack) == 0 {
		return NilValue(), ErrStackUnderflow
	}
	return e.stack[len(e.stack)-1], nil
}

func (e *executor) peekN(n int) (Value, error) {
	if n < 0 || n >= len(e.stack) {
		return NilValue(), ErrStackUnderflow
	}
	return e.stack[len(e.stack)-1-n], nil
}

// Conversion helpers for numeric operations (for future use)

func toFloat64(v Value) (float64, error) {
	switch v.Type {
	case TypeFloat:
		return v.AsFloat()
	case TypeInt:
		i, err := v.AsInt()
		if err != nil {
			return 0, err
		}
		return float64(i), nil
	default:
		return 0, ErrTypeMismatch
	}
}

func toInt64(v Value) (int64, error) {
	switch v.Type {
	case TypeInt:
		return v.AsInt()
	case TypeFloat:
		f, err := v.AsFloat()
		if err != nil {
			return 0, err
		}
		return int64(f), nil
	default:
		return 0, ErrTypeMismatch
	}
}

func toBool(v Value) bool {
	return v.IsTruthy()
}

func numericOp(a, b Value, op func(float64, float64) float64) (Value, error) {
	aVal, err := toFloat64(a)
	if err != nil {
		return NilValue(), err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return NilValue(), err
	}
	result := op(aVal, bVal)
	return FloatValue(result), nil
}

func compareOp(a, b Value, op func(float64, float64) bool) (Value, error) {
	aVal, err := toFloat64(a)
	if err != nil {
		return NilValue(), err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return NilValue(), err
	}
	result := op(aVal, bVal)
	return BoolValue(result), nil
}

func unaryMathOp(v Value, op func(float64) float64) (Value, error) {
	val, err := toFloat64(v)
	if err != nil {
		return NilValue(), err
	}
	result := op(val)
	return FloatValue(result), nil
}

// Helper to check for NaN and Inf
func isValidFloat(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}
