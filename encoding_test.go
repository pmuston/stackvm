package stackvm

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodeProgram(t *testing.T) {
	tests := []struct {
		name        string
		program     Program
		wantErr     bool
		errContains string
	}{
		{
			name: "simple program with two instructions",
			program: NewProgram([]Instruction{
				{Opcode: OpPUSH, Operand: 42},
				{Opcode: OpPOP, Operand: 0},
			}),
			wantErr: false,
		},
		{
			name:    "empty program",
			program: NewProgram([]Instruction{}),
			wantErr: false,
		},
		{
			name: "program with all standard opcodes",
			program: NewProgram([]Instruction{
				{Opcode: OpPUSH, Operand: 10},
				{Opcode: OpADD, Operand: 0},
				{Opcode: OpSUB, Operand: 0},
				{Opcode: OpMUL, Operand: 0},
				{Opcode: OpDIV, Operand: 0},
			}),
			wantErr: false,
		},
		{
			name: "program with custom opcodes",
			program: NewProgram([]Instruction{
				{Opcode: 128, Operand: 100}, // Custom opcode
				{Opcode: 200, Operand: 200}, // Custom opcode
				{Opcode: 255, Operand: 300}, // Custom opcode
			}),
			wantErr: false,
		},
		{
			name: "program with negative operands",
			program: NewProgram([]Instruction{
				{Opcode: OpPUSH, Operand: -42},
				{Opcode: OpJMP, Operand: -10},
			}),
			wantErr: false,
		},
		{
			name: "program with large operands",
			program: NewProgram([]Instruction{
				{Opcode: OpPUSH, Operand: 2147483647},  // Max int32
				{Opcode: OpPUSH, Operand: -2147483648}, // Min int32
			}),
			wantErr: false,
		},
		{
			name:        "nil program",
			program:     nil,
			wantErr:     true,
			errContains: "program is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytecode, err := EncodeProgram(tt.program)

			if tt.wantErr {
				if err == nil {
					t.Errorf("EncodeProgram() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("EncodeProgram() error = %v, want error containing %q", err, tt.errContains)
				}
				if !errors.Is(err, ErrInvalidProgram) {
					t.Errorf("EncodeProgram() error should wrap ErrInvalidProgram, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("EncodeProgram() unexpected error = %v", err)
				return
			}

			if bytecode == nil {
				t.Error("EncodeProgram() returned nil bytecode")
				return
			}

			// Validate bytecode structure
			instructions := tt.program.Instructions()
			expectedSize := 4 + (len(instructions) * 5)
			if len(bytecode) != expectedSize {
				t.Errorf("EncodeProgram() bytecode length = %d, want %d", len(bytecode), expectedSize)
			}
		})
	}
}

func TestDecodeProgram(t *testing.T) {
	tests := []struct {
		name        string
		bytecode    []byte
		wantErr     bool
		errContains string
		wantInstrs  []Instruction
	}{
		{
			name: "valid simple program",
			bytecode: []byte{
				0x02, 0x00, 0x00, 0x00, // 2 instructions
				0x00, 0x2A, 0x00, 0x00, 0x00, // PUSH 42
				0x02, 0x00, 0x00, 0x00, 0x00, // POP 0
			},
			wantErr: false,
			wantInstrs: []Instruction{
				{Opcode: OpPUSH, Operand: 42},
				{Opcode: OpPOP, Operand: 0},
			},
		},
		{
			name: "empty program",
			bytecode: []byte{
				0x00, 0x00, 0x00, 0x00, // 0 instructions
			},
			wantErr:    false,
			wantInstrs: []Instruction{},
		},
		{
			name: "program with custom opcodes",
			bytecode: []byte{
				0x02, 0x00, 0x00, 0x00, // 2 instructions
				0x80, 0x64, 0x00, 0x00, 0x00, // CUSTOM_128 100
				0xFF, 0xC8, 0x00, 0x00, 0x00, // CUSTOM_255 200
			},
			wantErr: false,
			wantInstrs: []Instruction{
				{Opcode: 128, Operand: 100},
				{Opcode: 255, Operand: 200},
			},
		},
		{
			name: "program with negative operands",
			bytecode: []byte{
				0x01, 0x00, 0x00, 0x00, // 1 instruction
				0x00, 0xD6, 0xFF, 0xFF, 0xFF, // PUSH -42 (little-endian two's complement)
			},
			wantErr: false,
			wantInstrs: []Instruction{
				{Opcode: OpPUSH, Operand: -42},
			},
		},
		{
			name:        "bytecode too short (empty)",
			bytecode:    []byte{},
			wantErr:     true,
			errContains: "too short",
		},
		{
			name:        "bytecode too short (partial header)",
			bytecode:    []byte{0x01, 0x00},
			wantErr:     true,
			errContains: "too short",
		},
		{
			name: "bytecode length mismatch (truncated)",
			bytecode: []byte{
				0x02, 0x00, 0x00, 0x00, // 2 instructions
				0x00, 0x2A, 0x00, 0x00, 0x00, // PUSH 42
				// Missing second instruction
			},
			wantErr:     true,
			errContains: "length mismatch",
		},
		{
			name: "bytecode length mismatch (extra bytes)",
			bytecode: []byte{
				0x01, 0x00, 0x00, 0x00, // 1 instruction
				0x00, 0x2A, 0x00, 0x00, 0x00, // PUSH 42
				0xFF, // Extra byte
			},
			wantErr:     true,
			errContains: "length mismatch",
		},
		{
			name: "bytecode with instruction count but no body",
			bytecode: []byte{
				0x05, 0x00, 0x00, 0x00, // 5 instructions (but no data follows)
			},
			wantErr:     true,
			errContains: "length mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := DecodeProgram(tt.bytecode)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DecodeProgram() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("DecodeProgram() error = %v, want error containing %q", err, tt.errContains)
				}
				if !errors.Is(err, ErrInvalidProgram) {
					t.Errorf("DecodeProgram() error should wrap ErrInvalidProgram, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeProgram() unexpected error = %v", err)
				return
			}

			if program == nil {
				t.Error("DecodeProgram() returned nil program")
				return
			}

			// Validate decoded instructions
			instructions := program.Instructions()
			if len(instructions) != len(tt.wantInstrs) {
				t.Errorf("DecodeProgram() instruction count = %d, want %d", len(instructions), len(tt.wantInstrs))
				return
			}

			for i, instr := range instructions {
				if instr.Opcode != tt.wantInstrs[i].Opcode {
					t.Errorf("DecodeProgram() instruction[%d].Opcode = %d, want %d",
						i, instr.Opcode, tt.wantInstrs[i].Opcode)
				}
				if instr.Operand != tt.wantInstrs[i].Operand {
					t.Errorf("DecodeProgram() instruction[%d].Operand = %d, want %d",
						i, instr.Operand, tt.wantInstrs[i].Operand)
				}
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name         string
		instructions []Instruction
	}{
		{
			name: "simple program",
			instructions: []Instruction{
				{Opcode: OpPUSH, Operand: 42},
				{Opcode: OpPUSHI, Operand: 100},
				{Opcode: OpADD, Operand: 0},
			},
		},
		{
			name:         "empty program",
			instructions: []Instruction{},
		},
		{
			name: "program with all opcode types",
			instructions: []Instruction{
				// Stack operations
				{Opcode: OpPUSH, Operand: 10},
				{Opcode: OpDUP, Operand: 0},
				{Opcode: OpSWAP, Operand: 0},
				// Arithmetic
				{Opcode: OpADD, Operand: 0},
				{Opcode: OpSUB, Operand: 0},
				{Opcode: OpMUL, Operand: 0},
				// Logic
				{Opcode: OpAND, Operand: 0},
				{Opcode: OpOR, Operand: 0},
				// Comparison
				{Opcode: OpEQ, Operand: 0},
				{Opcode: OpGT, Operand: 0},
				// Memory
				{Opcode: OpLOAD, Operand: 5},
				{Opcode: OpSTORE, Operand: 10},
				// Control flow
				{Opcode: OpJMP, Operand: 100},
				{Opcode: OpJMPZ, Operand: 50},
				{Opcode: OpHALT, Operand: 0},
				// Math functions
				{Opcode: OpSQRT, Operand: 0},
				{Opcode: OpSIN, Operand: 0},
				{Opcode: OpPOW, Operand: 0},
				// Custom opcodes
				{Opcode: 128, Operand: 999},
				{Opcode: 255, Operand: -123},
			},
		},
		{
			name: "program with extreme operand values",
			instructions: []Instruction{
				{Opcode: OpPUSH, Operand: 2147483647},  // Max int32
				{Opcode: OpPUSH, Operand: -2147483648}, // Min int32
				{Opcode: OpPUSH, Operand: 0},
				{Opcode: OpPUSH, Operand: 1},
				{Opcode: OpPUSH, Operand: -1},
			},
		},
		{
			name: "large program",
			instructions: func() []Instruction {
				instrs := make([]Instruction, 1000)
				for i := 0; i < 1000; i++ {
					instrs[i] = Instruction{
						Opcode:  Opcode(i % 256),
						Operand: int32(i),
					}
				}
				return instrs
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original program
			original := NewProgram(tt.instructions)

			// Encode
			bytecode, err := EncodeProgram(original)
			if err != nil {
				t.Fatalf("EncodeProgram() error = %v", err)
			}

			// Decode
			decoded, err := DecodeProgram(bytecode)
			if err != nil {
				t.Fatalf("DecodeProgram() error = %v", err)
			}

			// Compare instructions
			originalInstrs := original.Instructions()
			decodedInstrs := decoded.Instructions()

			if len(decodedInstrs) != len(originalInstrs) {
				t.Fatalf("Round-trip instruction count mismatch: got %d, want %d",
					len(decodedInstrs), len(originalInstrs))
			}

			for i := 0; i < len(originalInstrs); i++ {
				if decodedInstrs[i].Opcode != originalInstrs[i].Opcode {
					t.Errorf("Round-trip instruction[%d].Opcode mismatch: got %d, want %d",
						i, decodedInstrs[i].Opcode, originalInstrs[i].Opcode)
				}
				if decodedInstrs[i].Operand != originalInstrs[i].Operand {
					t.Errorf("Round-trip instruction[%d].Operand mismatch: got %d, want %d",
						i, decodedInstrs[i].Operand, originalInstrs[i].Operand)
				}
			}
		})
	}
}

func TestEncodeProgramFormat(t *testing.T) {
	// Test that encoding produces the exact expected byte format
	program := NewProgram([]Instruction{
		{Opcode: OpPUSH, Operand: 42},
		{Opcode: OpADD, Operand: 0},
	})

	bytecode, err := EncodeProgram(program)
	if err != nil {
		t.Fatalf("EncodeProgram() error = %v", err)
	}

	expected := []byte{
		0x02, 0x00, 0x00, 0x00, // 2 instructions (little-endian)
		0x00, 0x2A, 0x00, 0x00, 0x00, // PUSH (0) 42 (little-endian)
		0x10, 0x00, 0x00, 0x00, 0x00, // ADD (16) 0
	}

	if !bytes.Equal(bytecode, expected) {
		t.Errorf("EncodeProgram() bytecode format mismatch\ngot:  %v\nwant: %v", bytecode, expected)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
