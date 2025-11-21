package stackvm

import (
	"testing"
)

func TestNewSimpleMemory(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"Size 0", 0},
		{"Size 1", 1},
		{"Size 10", 10},
		{"Size 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewSimpleMemory(tt.size)
			if mem == nil {
				t.Fatal("NewSimpleMemory returned nil")
			}
			if mem.Size() != tt.size {
				t.Errorf("Size() = %d, want %d", mem.Size(), tt.size)
			}

			// Verify all values are initialized to Nil
			for i := 0; i < tt.size; i++ {
				val, err := mem.Load(i)
				if err != nil {
					t.Errorf("Load(%d) returned error: %v", i, err)
				}
				if !val.IsNil() {
					t.Errorf("Load(%d) = %v, want Nil", i, val)
				}
			}
		})
	}
}

func TestSimpleMemoryLoad(t *testing.T) {
	mem := NewSimpleMemory(5)
	// Store some test values
	mem.Store(0, FloatValue(1.5))
	mem.Store(2, IntValue(42))
	mem.Store(4, StringValue("test"))

	tests := []struct {
		name    string
		index   int
		want    Value
		wantErr bool
	}{
		{"Valid index 0", 0, FloatValue(1.5), false},
		{"Valid index 1 (nil)", 1, NilValue(), false},
		{"Valid index 2", 2, IntValue(42), false},
		{"Valid index 4", 4, StringValue("test"), false},
		{"Negative index", -1, NilValue(), true},
		{"Out of bounds", 5, NilValue(), true},
		{"Out of bounds large", 100, NilValue(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mem.Load(tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && err != ErrInvalidMemoryAddress {
				t.Errorf("Load() error = %v, want ErrInvalidMemoryAddress", err)
			}
		})
	}
}

func TestSimpleMemoryStore(t *testing.T) {
	mem := NewSimpleMemory(5)

	tests := []struct {
		name    string
		index   int
		value   Value
		wantErr bool
	}{
		{"Store at 0", 0, FloatValue(3.14), false},
		{"Store at 4", 4, IntValue(100), false},
		{"Store nil", 2, NilValue(), false},
		{"Store bool", 1, BoolValue(true), false},
		{"Store string", 3, StringValue("hello"), false},
		{"Negative index", -1, FloatValue(1.0), true},
		{"Out of bounds", 5, FloatValue(1.0), true},
		{"Out of bounds large", 100, FloatValue(1.0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mem.Store(tt.index, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != ErrInvalidMemoryAddress {
				t.Errorf("Store() error = %v, want ErrInvalidMemoryAddress", err)
				return
			}

			// Verify the value was stored correctly
			if !tt.wantErr {
				got, err := mem.Load(tt.index)
				if err != nil {
					t.Errorf("Load() after Store() returned error: %v", err)
					return
				}
				if !got.Equal(tt.value) {
					t.Errorf("After Store(), Load() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestSimpleMemorySize(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"Empty memory", 0},
		{"Single location", 1},
		{"Multiple locations", 256},
		{"Large memory", 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewSimpleMemory(tt.size)
			if got := mem.Size(); got != tt.size {
				t.Errorf("Size() = %d, want %d", got, tt.size)
			}
		})
	}
}

func TestSimpleMemoryValues(t *testing.T) {
	mem := NewSimpleMemory(5)
	mem.Store(0, FloatValue(1.5))
	mem.Store(1, IntValue(42))
	mem.Store(2, BoolValue(true))
	mem.Store(3, StringValue("test"))
	// Index 4 remains nil

	values := mem.Values()

	// Check that we got all values
	if len(values) != 5 {
		t.Errorf("Values() returned %d values, want 5", len(values))
	}

	// Check specific values
	if !values[0].Equal(FloatValue(1.5)) {
		t.Errorf("values[0] = %v, want FloatValue(1.5)", values[0])
	}
	if !values[1].Equal(IntValue(42)) {
		t.Errorf("values[1] = %v, want IntValue(42)", values[1])
	}
	if !values[2].Equal(BoolValue(true)) {
		t.Errorf("values[2] = %v, want BoolValue(true)", values[2])
	}
	if !values[3].Equal(StringValue("test")) {
		t.Errorf("values[3] = %v, want StringValue(test)", values[3])
	}
	if !values[4].IsNil() {
		t.Errorf("values[4] = %v, want NilValue", values[4])
	}

	// Verify it's a copy (modifying returned slice doesn't affect memory)
	values[0] = IntValue(999)
	storedValue, _ := mem.Load(0)
	if !storedValue.Equal(FloatValue(1.5)) {
		t.Error("Modifying Values() result affected memory (should be a copy)")
	}
}

func TestSimpleMemorySetValues(t *testing.T) {
	t.Run("Set values with matching size", func(t *testing.T) {
		mem := NewSimpleMemory(3)
		newValues := []Value{
			FloatValue(1.1),
			IntValue(22),
			BoolValue(false),
		}

		mem.SetValues(newValues)

		// Verify all values were set
		for i, want := range newValues {
			got, err := mem.Load(i)
			if err != nil {
				t.Errorf("Load(%d) returned error: %v", i, err)
			}
			if !got.Equal(want) {
				t.Errorf("Load(%d) = %v, want %v", i, got, want)
			}
		}
	})

	t.Run("Set values with mismatched size does nothing", func(t *testing.T) {
		mem := NewSimpleMemory(3)
		mem.Store(0, FloatValue(1.5))

		// Try to set with wrong size
		newValues := []Value{IntValue(1), IntValue(2)}
		mem.SetValues(newValues)

		// Original value should be unchanged
		got, _ := mem.Load(0)
		if !got.Equal(FloatValue(1.5)) {
			t.Errorf("SetValues with wrong size modified memory")
		}
	})

	t.Run("SetValues creates a copy", func(t *testing.T) {
		mem := NewSimpleMemory(2)
		newValues := []Value{FloatValue(1.0), FloatValue(2.0)}
		mem.SetValues(newValues)

		// Modify the source slice
		newValues[0] = IntValue(999)

		// Memory should not be affected
		got, _ := mem.Load(0)
		if !got.Equal(FloatValue(1.0)) {
			t.Error("Modifying source slice after SetValues affected memory")
		}
	})
}

func TestSimpleMemoryReset(t *testing.T) {
	mem := NewSimpleMemory(5)

	// Store various values
	mem.Store(0, FloatValue(3.14))
	mem.Store(1, IntValue(42))
	mem.Store(2, BoolValue(true))
	mem.Store(3, StringValue("hello"))
	mem.Store(4, FloatValue(-1.0))

	// Reset memory
	mem.Reset()

	// Verify all values are now Nil
	for i := 0; i < mem.Size(); i++ {
		val, err := mem.Load(i)
		if err != nil {
			t.Errorf("Load(%d) after Reset() returned error: %v", i, err)
		}
		if !val.IsNil() {
			t.Errorf("Load(%d) after Reset() = %v, want NilValue", i, val)
		}
	}
}

func TestSimpleMemorySequentialOperations(t *testing.T) {
	mem := NewSimpleMemory(10)

	// Test sequential store and load
	for i := 0; i < 10; i++ {
		val := IntValue(int64(i * 10))
		if err := mem.Store(i, val); err != nil {
			t.Fatalf("Store(%d) failed: %v", i, err)
		}
	}

	for i := 0; i < 10; i++ {
		val, err := mem.Load(i)
		if err != nil {
			t.Fatalf("Load(%d) failed: %v", i, err)
		}
		expected := IntValue(int64(i * 10))
		if !val.Equal(expected) {
			t.Errorf("Load(%d) = %v, want %v", i, val, expected)
		}
	}
}

func TestSimpleMemoryOverwrite(t *testing.T) {
	mem := NewSimpleMemory(1)

	// Store initial value
	mem.Store(0, IntValue(100))

	// Overwrite multiple times
	values := []Value{
		FloatValue(1.5),
		BoolValue(true),
		StringValue("test"),
		NilValue(),
		IntValue(42),
	}

	for _, val := range values {
		if err := mem.Store(0, val); err != nil {
			t.Errorf("Store() failed: %v", err)
		}
		got, err := mem.Load(0)
		if err != nil {
			t.Errorf("Load() failed: %v", err)
		}
		if !got.Equal(val) {
			t.Errorf("After overwrite, Load() = %v, want %v", got, val)
		}
	}
}

func TestMemoryInterface(t *testing.T) {
	// Verify SimpleMemory implements Memory interface
	var _ Memory = (*SimpleMemory)(nil)

	// Test through interface
	var mem Memory = NewSimpleMemory(5)

	if mem.Size() != 5 {
		t.Errorf("Size() through interface = %d, want 5", mem.Size())
	}

	err := mem.Store(2, FloatValue(3.14))
	if err != nil {
		t.Errorf("Store() through interface failed: %v", err)
	}

	val, err := mem.Load(2)
	if err != nil {
		t.Errorf("Load() through interface failed: %v", err)
	}
	if !val.Equal(FloatValue(3.14)) {
		t.Errorf("Load() through interface = %v, want FloatValue(3.14)", val)
	}
}
