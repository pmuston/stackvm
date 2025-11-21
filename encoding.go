package stackvm

import (
	"encoding/binary"
	"fmt"
)

// EncodeProgram encodes a Program to binary bytecode.
// Format:
//   - Header: 4 bytes instruction count (little-endian uint32)
//   - Body: For each instruction:
//       - 1 byte: opcode
//       - 4 bytes: operand (little-endian int32)
// Returns the encoded bytecode or an error.
func EncodeProgram(program Program) ([]byte, error) {
	if program == nil {
		return nil, fmt.Errorf("%w: program is nil", ErrInvalidProgram)
	}

	instructions := program.Instructions()
	if instructions == nil {
		instructions = []Instruction{}
	}

	// Calculate total size: 4 bytes header + (1 + 4) bytes per instruction
	totalSize := 4 + (len(instructions) * 5)
	bytecode := make([]byte, totalSize)

	// Write instruction count (4 bytes, little-endian)
	binary.LittleEndian.PutUint32(bytecode[0:4], uint32(len(instructions)))

	// Write each instruction
	offset := 4
	for _, instr := range instructions {
		// Write opcode (1 byte)
		bytecode[offset] = byte(instr.Opcode)
		offset++

		// Write operand (4 bytes, little-endian)
		binary.LittleEndian.PutUint32(bytecode[offset:offset+4], uint32(instr.Operand))
		offset += 4
	}

	return bytecode, nil
}

// DecodeProgram decodes binary bytecode to a Program.
// Validates the bytecode format and returns a Program or an error.
// Returns ErrInvalidProgram if the bytecode is malformed.
func DecodeProgram(data []byte) (Program, error) {
	// Minimum valid bytecode is 4 bytes (header with 0 instructions)
	if len(data) < 4 {
		return nil, fmt.Errorf("%w: bytecode too short (minimum 4 bytes required)", ErrInvalidProgram)
	}

	// Read instruction count from header
	instrCount := binary.LittleEndian.Uint32(data[0:4])

	// Validate bytecode length matches instruction count
	expectedSize := 4 + (instrCount * 5)
	if uint32(len(data)) != expectedSize {
		return nil, fmt.Errorf("%w: bytecode length mismatch (expected %d bytes, got %d bytes)",
			ErrInvalidProgram, expectedSize, len(data))
	}

	// Decode instructions
	instructions := make([]Instruction, instrCount)
	offset := 4
	for i := uint32(0); i < instrCount; i++ {
		// Read opcode (1 byte)
		opcode := Opcode(data[offset])
		offset++

		// Read operand (4 bytes, little-endian)
		operand := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
		offset += 4

		instructions[i] = Instruction{
			Opcode:  opcode,
			Operand: operand,
		}
	}

	return NewProgram(instructions), nil
}
