package stackvm

// opAnd pops two values, performs logical AND, and pushes the result.
func opAnd(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	result := a.IsTruthy() && b.IsTruthy()
	return append(stack, BoolValue(result)), nil
}

// opOr pops two values, performs logical OR, and pushes the result.
func opOr(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	result := a.IsTruthy() || b.IsTruthy()
	return append(stack, BoolValue(result)), nil
}

// opNot pops a value, performs logical NOT, and pushes the result.
func opNot(stack []Value) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	result := !a.IsTruthy()
	return append(stack, BoolValue(result)), nil
}

// opXor pops two values, performs logical XOR, and pushes the result.
func opXor(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	aTruthy := a.IsTruthy()
	bTruthy := b.IsTruthy()
	result := (aTruthy || bTruthy) && !(aTruthy && bTruthy)
	return append(stack, BoolValue(result)), nil
}
