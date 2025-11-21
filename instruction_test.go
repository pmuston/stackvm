package stackvm

import (
	"testing"
)

func TestNewInstruction(t *testing.T) {
	tests := []struct {
		name    string
		opcode  Opcode
		operand int32
	}{
		{"PUSH with operand", OpPUSH, 42},
		{"POP no operand", OpPOP, 0},
		{"JMP with offset", OpJMP, -10},
		{"LOAD with index", OpLOAD, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := NewInstruction(tt.opcode, tt.operand)
			if inst.Opcode != tt.opcode {
				t.Errorf("Opcode = %v, want %v", inst.Opcode, tt.opcode)
			}
			if inst.Operand != tt.operand {
				t.Errorf("Operand = %v, want %v", inst.Operand, tt.operand)
			}
		})
	}
}

func TestOpcodeString(t *testing.T) {
	tests := []struct {
		name   string
		opcode Opcode
		want   string
	}{
		// Stack operations
		{"PUSH", OpPUSH, "PUSH"},
		{"PUSHI", OpPUSHI, "PUSHI"},
		{"POP", OpPOP, "POP"},
		{"DUP", OpDUP, "DUP"},
		{"SWAP", OpSWAP, "SWAP"},
		{"OVER", OpOVER, "OVER"},
		{"ROT", OpROT, "ROT"},

		// Arithmetic operations
		{"ADD", OpADD, "ADD"},
		{"SUB", OpSUB, "SUB"},
		{"MUL", OpMUL, "MUL"},
		{"DIV", OpDIV, "DIV"},
		{"MOD", OpMOD, "MOD"},
		{"NEG", OpNEG, "NEG"},
		{"ABS", OpABS, "ABS"},
		{"INC", OpINC, "INC"},
		{"DEC", OpDEC, "DEC"},

		// Logic operations
		{"AND", OpAND, "AND"},
		{"OR", OpOR, "OR"},
		{"NOT", OpNOT, "NOT"},
		{"XOR", OpXOR, "XOR"},

		// Comparison operations
		{"EQ", OpEQ, "EQ"},
		{"NE", OpNE, "NE"},
		{"GT", OpGT, "GT"},
		{"LT", OpLT, "LT"},
		{"GE", OpGE, "GE"},
		{"LE", OpLE, "LE"},

		// Memory operations
		{"LOAD", OpLOAD, "LOAD"},
		{"STORE", OpSTORE, "STORE"},
		{"LOADD", OpLOADD, "LOADD"},
		{"STORED", OpSTORED, "STORED"},

		// Control flow operations
		{"JMP", OpJMP, "JMP"},
		{"JMPZ", OpJMPZ, "JMPZ"},
		{"JMPNZ", OpJMPNZ, "JMPNZ"},
		{"CALL", OpCALL, "CALL"},
		{"RET", OpRET, "RET"},
		{"HALT", OpHALT, "HALT"},
		{"NOP", OpNOP, "NOP"},

		// Math functions
		{"SQRT", OpSQRT, "SQRT"},
		{"SIN", OpSIN, "SIN"},
		{"COS", OpCOS, "COS"},
		{"TAN", OpTAN, "TAN"},
		{"ASIN", OpASIN, "ASIN"},
		{"ACOS", OpACOS, "ACOS"},
		{"ATAN", OpATAN, "ATAN"},
		{"ATAN2", OpATAN2, "ATAN2"},
		{"LOG", OpLOG, "LOG"},
		{"LOG10", OpLOG10, "LOG10"},
		{"EXP", OpEXP, "EXP"},
		{"POW", OpPOW, "POW"},
		{"MIN", OpMIN, "MIN"},
		{"MAX", OpMAX, "MAX"},
		{"FLOOR", OpFLOOR, "FLOOR"},
		{"CEIL", OpCEIL, "CEIL"},
		{"ROUND", OpROUND, "ROUND"},
		{"TRUNC", OpTRUNC, "TRUNC"},

		// Custom opcodes
		{"Custom 128", Opcode(128), "CUSTOM_128"},
		{"Custom 200", Opcode(200), "CUSTOM_200"},
		{"Custom 255", Opcode(255), "CUSTOM_255"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opcode.String(); got != tt.want {
				t.Errorf("Opcode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstructionString(t *testing.T) {
	tests := []struct {
		name string
		inst Instruction
		want string
	}{
		{"PUSH with operand", NewInstruction(OpPUSH, 42), "PUSH 42"},
		{"POP no operand", NewInstruction(OpPOP, 0), "POP"},
		{"JMP with offset", NewInstruction(OpJMP, -10), "JMP -10"},
		{"LOAD with index", NewInstruction(OpLOAD, 5), "LOAD 5"},
		{"HALT no operand", NewInstruction(OpHALT, 0), "HALT"},
		{"ADD no operand", NewInstruction(OpADD, 0), "ADD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.inst.String(); got != tt.want {
				t.Errorf("Instruction.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpcodeIsStandardOpcode(t *testing.T) {
	tests := []struct {
		name   string
		opcode Opcode
		want   bool
	}{
		{"PUSH is standard", OpPUSH, true},
		{"ADD is standard", OpADD, true},
		{"HALT is standard", OpHALT, true},
		{"SQRT is standard", OpSQRT, true},
		{"Opcode 127 is standard", Opcode(127), true},
		{"Opcode 128 is not standard", Opcode(128), false},
		{"Opcode 200 is not standard", Opcode(200), false},
		{"Opcode 255 is not standard", Opcode(255), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opcode.IsStandardOpcode(); got != tt.want {
				t.Errorf("IsStandardOpcode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpcodeIsCustomOpcode(t *testing.T) {
	tests := []struct {
		name   string
		opcode Opcode
		want   bool
	}{
		{"PUSH is not custom", OpPUSH, false},
		{"ADD is not custom", OpADD, false},
		{"HALT is not custom", OpHALT, false},
		{"SQRT is not custom", OpSQRT, false},
		{"Opcode 127 is not custom", Opcode(127), false},
		{"Opcode 128 is custom", Opcode(128), true},
		{"Opcode 200 is custom", Opcode(200), true},
		{"Opcode 255 is custom", Opcode(255), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opcode.IsCustomOpcode(); got != tt.want {
				t.Errorf("IsCustomOpcode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpcodeRanges(t *testing.T) {
	t.Run("Stack operations are 0-15", func(t *testing.T) {
		stackOps := []Opcode{OpPUSH, OpPUSHI, OpPOP, OpDUP, OpSWAP, OpOVER, OpROT}
		for _, op := range stackOps {
			if op < 0 || op > 15 {
				t.Errorf("Stack operation %v (%d) is not in range 0-15", op, op)
			}
		}
	})

	t.Run("Arithmetic operations are 16-31", func(t *testing.T) {
		arithOps := []Opcode{OpADD, OpSUB, OpMUL, OpDIV, OpMOD, OpNEG, OpABS, OpINC, OpDEC}
		for _, op := range arithOps {
			if op < 16 || op > 31 {
				t.Errorf("Arithmetic operation %v (%d) is not in range 16-31", op, op)
			}
		}
	})

	t.Run("Logic operations are 32-39", func(t *testing.T) {
		logicOps := []Opcode{OpAND, OpOR, OpNOT, OpXOR}
		for _, op := range logicOps {
			if op < 32 || op > 39 {
				t.Errorf("Logic operation %v (%d) is not in range 32-39", op, op)
			}
		}
	})

	t.Run("Comparison operations are 40-47", func(t *testing.T) {
		cmpOps := []Opcode{OpEQ, OpNE, OpGT, OpLT, OpGE, OpLE}
		for _, op := range cmpOps {
			if op < 40 || op > 47 {
				t.Errorf("Comparison operation %v (%d) is not in range 40-47", op, op)
			}
		}
	})

	t.Run("Memory operations are 48-55", func(t *testing.T) {
		memOps := []Opcode{OpLOAD, OpSTORE, OpLOADD, OpSTORED}
		for _, op := range memOps {
			if op < 48 || op > 55 {
				t.Errorf("Memory operation %v (%d) is not in range 48-55", op, op)
			}
		}
	})

	t.Run("Control flow operations are 56-63", func(t *testing.T) {
		ctrlOps := []Opcode{OpJMP, OpJMPZ, OpJMPNZ, OpCALL, OpRET, OpHALT, OpNOP}
		for _, op := range ctrlOps {
			if op < 56 || op > 63 {
				t.Errorf("Control flow operation %v (%d) is not in range 56-63", op, op)
			}
		}
	})

	t.Run("Math functions are 64-81", func(t *testing.T) {
		mathOps := []Opcode{
			OpSQRT, OpSIN, OpCOS, OpTAN, OpASIN, OpACOS, OpATAN, OpATAN2,
			OpLOG, OpLOG10, OpEXP, OpPOW, OpMIN, OpMAX, OpFLOOR, OpCEIL, OpROUND, OpTRUNC,
		}
		for _, op := range mathOps {
			if op < 64 || op > 81 {
				t.Errorf("Math function %v (%d) is not in range 64-81", op, op)
			}
		}
	})
}

func TestOpcodeValues(t *testing.T) {
	// Test specific opcode values as defined in the spec
	tests := []struct {
		name   string
		opcode Opcode
		want   uint8
	}{
		{"PUSH", OpPUSH, 0},
		{"PUSHI", OpPUSHI, 1},
		{"POP", OpPOP, 2},
		{"ADD", OpADD, 16},
		{"SUB", OpSUB, 17},
		{"AND", OpAND, 32},
		{"EQ", OpEQ, 40},
		{"LOAD", OpLOAD, 48},
		{"JMP", OpJMP, 56},
		{"SQRT", OpSQRT, 64},
		{"HALT", OpHALT, 61},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint8(tt.opcode) != tt.want {
				t.Errorf("Opcode %v has value %d, want %d", tt.name, tt.opcode, tt.want)
			}
		})
	}
}
