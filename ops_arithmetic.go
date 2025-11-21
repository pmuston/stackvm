package stackvm

// opAdd pops two values, adds them, and pushes the result.
func opAdd(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	result, err := numericOp(a, b, func(x, y float64) float64 { return x + y })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opSub pops two values, subtracts them, and pushes the result.
func opSub(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	result, err := numericOp(a, b, func(x, y float64) float64 { return x - y })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opMul pops two values, multiplies them, and pushes the result.
func opMul(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	result, err := numericOp(a, b, func(x, y float64) float64 { return x * y })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opDiv pops two values, divides them, and pushes the result.
func opDiv(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	bVal, err := toFloat64(b)
	if err != nil {
		return stack, err
	}
	if bVal == 0 {
		return stack, ErrDivisionByZero
	}

	result, err := numericOp(a, b, func(x, y float64) float64 { return x / y })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opMod pops two values, computes modulo, and pushes the result.
func opMod(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	aVal, err := toInt64(a)
	if err != nil {
		return stack, err
	}
	bVal, err := toInt64(b)
	if err != nil {
		return stack, err
	}
	if bVal == 0 {
		return stack, ErrDivisionByZero
	}

	result := IntValue(aVal % bVal)
	return append(stack, result), nil
}

// opNeg pops a value, negates it, and pushes the result.
func opNeg(stack []Value) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	result, err := unaryOp(a, func(x float64) float64 { return -x })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opAbs pops a value, computes absolute value, and pushes the result.
func opAbs(stack []Value) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}

	result := aVal
	if result < 0 {
		result = -result
	}

	return append(stack, FloatValue(result)), nil
}

// opInc pops a value, increments it, and pushes the result.
func opInc(stack []Value) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	result, err := unaryOp(a, func(x float64) float64 { return x + 1 })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// opDec pops a value, decrements it, and pushes the result.
func opDec(stack []Value) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	result, err := unaryOp(a, func(x float64) float64 { return x - 1 })
	if err != nil {
		return stack, err
	}

	return append(stack, result), nil
}

// Helper function for unary operations
func unaryOp(v Value, op func(float64) float64) (Value, error) {
	val, err := toFloat64(v)
	if err != nil {
		return NilValue(), err
	}
	result := op(val)
	return FloatValue(result), nil
}
