package stackvm

// Memory provides an abstraction for VM storage.
// Host systems can implement this interface to provide custom memory backends.
type Memory interface {
	// Load retrieves the value at the specified index.
	// Returns ErrInvalidMemoryAddress if the index is out of bounds.
	Load(index int) (Value, error)

	// Store saves the value at the specified index.
	// Returns ErrInvalidMemoryAddress if the index is out of bounds.
	// Returns ErrReadOnlyMemory if the memory is read-only.
	Store(index int, value Value) error

	// Size returns the number of addressable memory locations.
	Size() int
}

// ReadOnlyMemory extends Memory with read-only semantics.
type ReadOnlyMemory interface {
	Memory
	// IsReadOnly returns true if the memory cannot be written to.
	IsReadOnly() bool
}

// SimpleMemory is a basic memory implementation using a slice.
// It provides fixed-size, writable memory suitable for testing and simple use cases.
type SimpleMemory struct {
	data []Value
}

// NewSimpleMemory creates a new SimpleMemory with the specified size.
// All memory locations are initialized to NilValue().
func NewSimpleMemory(size int) *SimpleMemory {
	data := make([]Value, size)
	// Initialize all values to Nil
	for i := range data {
		data[i] = NilValue()
	}
	return &SimpleMemory{
		data: data,
	}
}

// Load retrieves the value at the specified index.
// Returns ErrInvalidMemoryAddress if the index is out of bounds or negative.
func (m *SimpleMemory) Load(index int) (Value, error) {
	if index < 0 || index >= len(m.data) {
		return NilValue(), ErrInvalidMemoryAddress
	}
	return m.data[index], nil
}

// Store saves the value at the specified index.
// Returns ErrInvalidMemoryAddress if the index is out of bounds or negative.
func (m *SimpleMemory) Store(index int, value Value) error {
	if index < 0 || index >= len(m.data) {
		return ErrInvalidMemoryAddress
	}
	m.data[index] = value
	return nil
}

// Size returns the number of addressable memory locations.
func (m *SimpleMemory) Size() int {
	return len(m.data)
}

// Values returns a copy of all memory values.
// This is useful for inspection and testing.
func (m *SimpleMemory) Values() []Value {
	result := make([]Value, len(m.data))
	copy(result, m.data)
	return result
}

// SetValues bulk sets all memory values.
// The values slice must match the memory size exactly.
// If the lengths don't match, this method does nothing.
func (m *SimpleMemory) SetValues(values []Value) {
	if len(values) != len(m.data) {
		return
	}
	copy(m.data, values)
}

// Reset clears all memory values back to NilValue().
func (m *SimpleMemory) Reset() {
	for i := range m.data {
		m.data[i] = NilValue()
	}
}
