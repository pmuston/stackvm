package stackvm

// opEq pops two values, compares for equality, and pushes the result.
func opEq(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	result := a.Equal(b)
	return append(stack, BoolValue(result)), nil
}

// opNe pops two values, compares for inequality, and pushes the result.
func opNe(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	result := !a.Equal(b)
	return append(stack, BoolValue(result)), nil
}

// opGt pops two values, checks if first > second, and pushes the result.
func opGt(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return stack, err
	}
	result := aVal > bVal
	return append(stack, BoolValue(result)), nil
}

// opLt pops two values, checks if first < second, and pushes the result.
func opLt(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return stack, err
	}
	result := aVal < bVal
	return append(stack, BoolValue(result)), nil
}

// opGe pops two values, checks if first >= second, and pushes the result.
func opGe(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return stack, err
	}
	result := aVal >= bVal
	return append(stack, BoolValue(result)), nil
}

// opLe pops two values, checks if first <= second, and pushes the result.
func opLe(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	b := stack[len(stack)-1]
	a := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}
	bVal, err := toFloat64(b)
	if err != nil {
		return stack, err
	}
	result := aVal <= bVal
	return append(stack, BoolValue(result)), nil
}
