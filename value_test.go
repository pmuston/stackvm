package stackvm

import (
	"testing"
)

func TestValueConstructors(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		wantType ValueType
		wantData interface{}
	}{
		{"Nil", NilValue(), TypeNil, nil},
		{"Float", FloatValue(3.14), TypeFloat, 3.14},
		{"Int", IntValue(42), TypeInt, int64(42)},
		{"Bool true", BoolValue(true), TypeBool, true},
		{"Bool false", BoolValue(false), TypeBool, false},
		{"String", StringValue("hello"), TypeString, "hello"},
		{"Custom", CustomValue(128, "custom"), ValueType(128), "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", tt.value.Type, tt.wantType)
			}
			if tt.value.Data != tt.wantData {
				t.Errorf("Data = %v, want %v", tt.value.Data, tt.wantData)
			}
		})
	}
}

func TestValueIsNil(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  bool
	}{
		{"Nil is nil", NilValue(), true},
		{"Float not nil", FloatValue(0), false},
		{"Int not nil", IntValue(0), false},
		{"Bool not nil", BoolValue(false), false},
		{"String not nil", StringValue(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueAsFloat(t *testing.T) {
	tests := []struct {
		name    string
		value   Value
		want    float64
		wantErr bool
	}{
		{"Valid float", FloatValue(3.14), 3.14, false},
		{"Zero float", FloatValue(0.0), 0.0, false},
		{"Negative float", FloatValue(-42.5), -42.5, false},
		{"Int returns error", IntValue(42), 0, true},
		{"Bool returns error", BoolValue(true), 0, true},
		{"Nil returns error", NilValue(), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.AsFloat()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("AsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueAsInt(t *testing.T) {
	tests := []struct {
		name    string
		value   Value
		want    int64
		wantErr bool
	}{
		{"Valid int", IntValue(42), 42, false},
		{"Zero int", IntValue(0), 0, false},
		{"Negative int", IntValue(-100), -100, false},
		{"Float returns error", FloatValue(3.14), 0, true},
		{"Bool returns error", BoolValue(true), 0, true},
		{"Nil returns error", NilValue(), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.AsInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueAsBool(t *testing.T) {
	tests := []struct {
		name    string
		value   Value
		want    bool
		wantErr bool
	}{
		{"Bool true", BoolValue(true), true, false},
		{"Bool false", BoolValue(false), false, false},
		{"Int returns error", IntValue(1), false, true},
		{"Float returns error", FloatValue(1.0), false, true},
		{"Nil returns error", NilValue(), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.AsBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueAsString(t *testing.T) {
	tests := []struct {
		name    string
		value   Value
		want    string
		wantErr bool
	}{
		{"Valid string", StringValue("hello"), "hello", false},
		{"Empty string", StringValue(""), "", false},
		{"Int returns error", IntValue(42), "", true},
		{"Float returns error", FloatValue(3.14), "", true},
		{"Nil returns error", NilValue(), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.AsString()
			if (err != nil) != tt.wantErr {
				t.Errorf("AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueIsNumeric(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  bool
	}{
		{"Float is numeric", FloatValue(3.14), true},
		{"Int is numeric", IntValue(42), true},
		{"Bool not numeric", BoolValue(true), false},
		{"String not numeric", StringValue("42"), false},
		{"Nil not numeric", NilValue(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsNumeric(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueIsTruthy(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  bool
	}{
		{"Nil is falsy", NilValue(), false},
		{"Float zero is falsy", FloatValue(0.0), false},
		{"Float non-zero is truthy", FloatValue(3.14), true},
		{"Float negative is truthy", FloatValue(-1.0), true},
		{"Int zero is falsy", IntValue(0), false},
		{"Int non-zero is truthy", IntValue(42), true},
		{"Int negative is truthy", IntValue(-1), true},
		{"Bool true is truthy", BoolValue(true), true},
		{"Bool false is falsy", BoolValue(false), false},
		{"Empty string is falsy", StringValue(""), false},
		{"Non-empty string is truthy", StringValue("hello"), true},
		{"Custom is falsy", CustomValue(128, "data"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsTruthy(); got != tt.want {
				t.Errorf("IsTruthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{"Nil", NilValue(), "nil"},
		{"Float", FloatValue(3.14), "3.14"},
		{"Float zero", FloatValue(0.0), "0"},
		{"Int", IntValue(42), "42"},
		{"Int negative", IntValue(-100), "-100"},
		{"Bool true", BoolValue(true), "true"},
		{"Bool false", BoolValue(false), "false"},
		{"String", StringValue("hello"), "hello"},
		{"Empty string", StringValue(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueEqual(t *testing.T) {
	tests := []struct {
		name  string
		v1    Value
		v2    Value
		want  bool
	}{
		{"Nil equals nil", NilValue(), NilValue(), true},
		{"Float equals", FloatValue(3.14), FloatValue(3.14), true},
		{"Float not equals", FloatValue(3.14), FloatValue(2.71), false},
		{"Int equals", IntValue(42), IntValue(42), true},
		{"Int not equals", IntValue(42), IntValue(100), false},
		{"Bool equals true", BoolValue(true), BoolValue(true), true},
		{"Bool equals false", BoolValue(false), BoolValue(false), true},
		{"Bool not equals", BoolValue(true), BoolValue(false), false},
		{"String equals", StringValue("hello"), StringValue("hello"), true},
		{"String not equals", StringValue("hello"), StringValue("world"), false},
		{"Different types not equal", IntValue(42), FloatValue(42.0), false},
		{"Int and bool not equal", IntValue(1), BoolValue(true), false},
		{"Nil and zero not equal", NilValue(), IntValue(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Equal(tt.v2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
			// Test symmetry
			if got := tt.v2.Equal(tt.v1); got != tt.want {
				t.Errorf("Equal() (reversed) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomValue(t *testing.T) {
	t.Run("Custom type 128", func(t *testing.T) {
		v := CustomValue(128, "custom data")
		if v.Type != 128 {
			t.Errorf("Type = %v, want 128", v.Type)
		}
		if v.Data != "custom data" {
			t.Errorf("Data = %v, want 'custom data'", v.Data)
		}
	})

	t.Run("Custom type 255", func(t *testing.T) {
		v := CustomValue(255, 12345)
		if v.Type != 255 {
			t.Errorf("Type = %v, want 255", v.Type)
		}
		if v.Data != 12345 {
			t.Errorf("Data = %v, want 12345", v.Data)
		}
	})

	t.Run("Custom equals", func(t *testing.T) {
		v1 := CustomValue(128, "data")
		v2 := CustomValue(128, "data")
		if !v1.Equal(v2) {
			t.Errorf("Custom values should be equal")
		}
	})

	t.Run("Custom not equals different data", func(t *testing.T) {
		v1 := CustomValue(128, "data1")
		v2 := CustomValue(128, "data2")
		if v1.Equal(v2) {
			t.Errorf("Custom values with different data should not be equal")
		}
	})

	t.Run("Custom not equals different type", func(t *testing.T) {
		v1 := CustomValue(128, "data")
		v2 := CustomValue(129, "data")
		if v1.Equal(v2) {
			t.Errorf("Custom values with different types should not be equal")
		}
	})
}
