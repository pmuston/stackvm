package stackvm

import "math"

// Math operations

func opSqrt(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Sqrt)
}

func opSin(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Sin)
}

func opCos(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Cos)
}

func opTan(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Tan)
}

func opAsin(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Asin)
}

func opAcos(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Acos)
}

func opAtan(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Atan)
}

func opAtan2(stack []Value) ([]Value, error) {
	if len(stack) < 2 {
		return stack, ErrStackUnderflow
	}
	x := stack[len(stack)-1]
	y := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	yVal, err := toFloat64(y)
	if err != nil {
		return stack, err
	}
	xVal, err := toFloat64(x)
	if err != nil {
		return stack, err
	}
	result := math.Atan2(yVal, xVal)
	return append(stack, FloatValue(result)), nil
}

func opLog(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Log)
}

func opLog10(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Log10)
}

func opExp(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Exp)
}

func opPow(stack []Value) ([]Value, error) {
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
	result := math.Pow(aVal, bVal)
	return append(stack, FloatValue(result)), nil
}

func opMin(stack []Value) ([]Value, error) {
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
	result := math.Min(aVal, bVal)
	return append(stack, FloatValue(result)), nil
}

func opMax(stack []Value) ([]Value, error) {
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
	result := math.Max(aVal, bVal)
	return append(stack, FloatValue(result)), nil
}

func opFloor(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Floor)
}

func opCeil(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Ceil)
}

func opRound(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Round)
}

func opTrunc(stack []Value) ([]Value, error) {
	return mathUnaryOp(stack, math.Trunc)
}

func mathUnaryOp(stack []Value, op func(float64) float64) ([]Value, error) {
	if len(stack) < 1 {
		return stack, ErrStackUnderflow
	}
	a := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	aVal, err := toFloat64(a)
	if err != nil {
		return stack, err
	}
	result := op(aVal)
	return append(stack, FloatValue(result)), nil
}
