// Package stackvm provides a standalone, reusable stack-based virtual machine.
package stackvm

import (
	"fmt"
	"strconv"
)

// ValueType represents the type of a Value in the VM.
type ValueType uint8

// Core value types supported by the VM.
const (
	TypeNil    ValueType = 0
	TypeFloat  ValueType = 1
	TypeInt    ValueType = 2
	TypeBool   ValueType = 3
	TypeString ValueType = 4
	// TypeCustom range: 128-255 reserved for host-defined types
)

// Value represents a typed value in the VM.
// It consists of a type tag and the underlying Go value.
type Value struct {
	Type ValueType
	Data interface{}
}

// NilValue returns a new nil Value.
func NilValue() Value {
	return Value{Type: TypeNil, Data: nil}
}

// FloatValue returns a new float Value.
func FloatValue(v float64) Value {
	return Value{Type: TypeFloat, Data: v}
}

// IntValue returns a new integer Value.
func IntValue(v int64) Value {
	return Value{Type: TypeInt, Data: v}
}

// BoolValue returns a new boolean Value.
func BoolValue(v bool) Value {
	return Value{Type: TypeBool, Data: v}
}

// StringValue returns a new string Value.
func StringValue(v string) Value {
	return Value{Type: TypeString, Data: v}
}

// CustomValue returns a new custom-typed Value.
// The type must be in the range 128-255.
func CustomValue(typ ValueType, data interface{}) Value {
	return Value{Type: typ, Data: data}
}

// IsNil returns true if the Value is nil.
func (v Value) IsNil() bool {
	return v.Type == TypeNil
}

// AsFloat returns the Value as a float64.
// Returns an error if the Value is not a float.
func (v Value) AsFloat() (float64, error) {
	if v.Type != TypeFloat {
		return 0, ErrTypeMismatch
	}
	f, ok := v.Data.(float64)
	if !ok {
		return 0, ErrTypeMismatch
	}
	return f, nil
}

// AsInt returns the Value as an int64.
// Returns an error if the Value is not an integer.
func (v Value) AsInt() (int64, error) {
	if v.Type != TypeInt {
		return 0, ErrTypeMismatch
	}
	i, ok := v.Data.(int64)
	if !ok {
		return 0, ErrTypeMismatch
	}
	return i, nil
}

// AsBool returns the Value as a bool.
// Returns an error if the Value is not a boolean.
func (v Value) AsBool() (bool, error) {
	if v.Type != TypeBool {
		return false, ErrTypeMismatch
	}
	b, ok := v.Data.(bool)
	if !ok {
		return false, ErrTypeMismatch
	}
	return b, nil
}

// AsString returns the Value as a string.
// Returns an error if the Value is not a string.
func (v Value) AsString() (string, error) {
	if v.Type != TypeString {
		return "", ErrTypeMismatch
	}
	s, ok := v.Data.(string)
	if !ok {
		return "", ErrTypeMismatch
	}
	return s, nil
}

// IsNumeric returns true if the Value is a numeric type (Float or Int).
func (v Value) IsNumeric() bool {
	return v.Type == TypeFloat || v.Type == TypeInt
}

// IsTruthy returns the truthiness of the Value.
// - Float: true if != 0.0
// - Int: true if != 0
// - Bool: the value itself
// - String: true if not empty
// - Nil: false
// - Custom: false (default)
func (v Value) IsTruthy() bool {
	switch v.Type {
	case TypeNil:
		return false
	case TypeFloat:
		f, _ := v.AsFloat()
		return f != 0.0
	case TypeInt:
		i, _ := v.AsInt()
		return i != 0
	case TypeBool:
		b, _ := v.AsBool()
		return b
	case TypeString:
		s, _ := v.AsString()
		return s != ""
	default:
		// Custom types default to false
		return false
	}
}

// String returns a human-readable representation of the Value.
func (v Value) String() string {
	switch v.Type {
	case TypeNil:
		return "nil"
	case TypeFloat:
		f, _ := v.AsFloat()
		return strconv.FormatFloat(f, 'g', -1, 64)
	case TypeInt:
		i, _ := v.AsInt()
		return strconv.FormatInt(i, 10)
	case TypeBool:
		b, _ := v.AsBool()
		return strconv.FormatBool(b)
	case TypeString:
		s, _ := v.AsString()
		return s
	default:
		// Custom types
		return fmt.Sprintf("<custom:%d:%v>", v.Type, v.Data)
	}
}

// Equal performs type-aware equality comparison.
func (v Value) Equal(other Value) bool {
	// Different types are never equal
	if v.Type != other.Type {
		return false
	}

	switch v.Type {
	case TypeNil:
		return true // All nils are equal
	case TypeFloat:
		f1, _ := v.AsFloat()
		f2, _ := other.AsFloat()
		return f1 == f2
	case TypeInt:
		i1, _ := v.AsInt()
		i2, _ := other.AsInt()
		return i1 == i2
	case TypeBool:
		b1, _ := v.AsBool()
		b2, _ := other.AsBool()
		return b1 == b2
	case TypeString:
		s1, _ := v.AsString()
		s2, _ := other.AsString()
		return s1 == s2
	default:
		// Custom types - compare underlying data
		return v.Data == other.Data
	}
}
