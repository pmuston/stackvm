package stackvm

import "fmt"

// Opcode represents a VM instruction opcode.
type Opcode uint8

// Stack operations (0-15)
const (
	OpPUSH  Opcode = 0  // Push immediate value (as float)
	OpPUSHI Opcode = 1  // Push immediate value (as int)
	OpPOP   Opcode = 2  // Remove top of stack
	OpDUP   Opcode = 3  // Duplicate top
	OpSWAP  Opcode = 4  // Exchange top two
	OpOVER  Opcode = 5  // Copy second to top
	OpROT   Opcode = 6  // Rotate top three
)

// Arithmetic operations (16-31)
const (
	OpADD Opcode = 16 // Addition
	OpSUB Opcode = 17 // Subtraction
	OpMUL Opcode = 18 // Multiplication
	OpDIV Opcode = 19 // Division
	OpMOD Opcode = 20 // Modulo
	OpNEG Opcode = 21 // Negate
	OpABS Opcode = 22 // Absolute value
	OpINC Opcode = 23 // Increment
	OpDEC Opcode = 24 // Decrement
)

// Logic operations (32-39)
const (
	OpAND Opcode = 32 // Logical AND
	OpOR  Opcode = 33 // Logical OR
	OpNOT Opcode = 34 // Logical NOT
	OpXOR Opcode = 35 // Logical XOR
)

// Comparison operations (40-47)
const (
	OpEQ Opcode = 40 // Equal
	OpNE Opcode = 41 // Not equal
	OpGT Opcode = 42 // Greater than
	OpLT Opcode = 43 // Less than
	OpGE Opcode = 44 // Greater or equal
	OpLE Opcode = 45 // Less or equal
)

// Memory operations (48-55)
const (
	OpLOAD   Opcode = 48 // Load from memory[index]
	OpSTORE  Opcode = 49 // Store to memory[index]
	OpLOADD  Opcode = 50 // Load from memory[pop()]
	OpSTORED Opcode = 51 // Store to memory[pop()]
)

// Control flow operations (56-63)
const (
	OpJMP   Opcode = 56 // Jump to offset
	OpJMPZ  Opcode = 57 // Jump if zero/false
	OpJMPNZ Opcode = 58 // Jump if non-zero/true
	OpCALL  Opcode = 59 // Call subroutine
	OpRET   Opcode = 60 // Return from subroutine
	OpHALT  Opcode = 61 // Stop execution
	OpNOP   Opcode = 62 // No operation
)

// Math functions (64-81)
const (
	OpSQRT   Opcode = 64 // Square root
	OpSIN    Opcode = 65 // Sine (radians)
	OpCOS    Opcode = 66 // Cosine (radians)
	OpTAN    Opcode = 67 // Tangent (radians)
	OpASIN   Opcode = 68 // Arc sine
	OpACOS   Opcode = 69 // Arc cosine
	OpATAN   Opcode = 70 // Arc tangent
	OpATAN2  Opcode = 71 // Two-argument arc tangent
	OpLOG    Opcode = 72 // Natural logarithm
	OpLOG10  Opcode = 73 // Base-10 logarithm
	OpEXP    Opcode = 74 // Exponential
	OpPOW    Opcode = 75 // Power
	OpMIN    Opcode = 76 // Minimum
	OpMAX    Opcode = 77 // Maximum
	OpFLOOR  Opcode = 78 // Floor
	OpCEIL   Opcode = 79 // Ceiling
	OpROUND  Opcode = 80 // Round to nearest
	OpTRUNC  Opcode = 81 // Truncate toward zero
)

// Custom operations (128-255) are reserved for host-defined extensions.

// Instruction represents a VM instruction with an opcode and operand.
type Instruction struct {
	Opcode  Opcode
	Operand int32
}

// NewInstruction creates a new Instruction.
func NewInstruction(opcode Opcode, operand int32) Instruction {
	return Instruction{
		Opcode:  opcode,
		Operand: operand,
	}
}

// String returns a human-readable representation of the instruction.
func (i Instruction) String() string {
	name := i.Opcode.String()
	if i.Operand != 0 {
		return fmt.Sprintf("%s %d", name, i.Operand)
	}
	return name
}

// String returns the mnemonic name of the opcode.
func (op Opcode) String() string {
	switch op {
	// Stack operations
	case OpPUSH:
		return "PUSH"
	case OpPUSHI:
		return "PUSHI"
	case OpPOP:
		return "POP"
	case OpDUP:
		return "DUP"
	case OpSWAP:
		return "SWAP"
	case OpOVER:
		return "OVER"
	case OpROT:
		return "ROT"

	// Arithmetic operations
	case OpADD:
		return "ADD"
	case OpSUB:
		return "SUB"
	case OpMUL:
		return "MUL"
	case OpDIV:
		return "DIV"
	case OpMOD:
		return "MOD"
	case OpNEG:
		return "NEG"
	case OpABS:
		return "ABS"
	case OpINC:
		return "INC"
	case OpDEC:
		return "DEC"

	// Logic operations
	case OpAND:
		return "AND"
	case OpOR:
		return "OR"
	case OpNOT:
		return "NOT"
	case OpXOR:
		return "XOR"

	// Comparison operations
	case OpEQ:
		return "EQ"
	case OpNE:
		return "NE"
	case OpGT:
		return "GT"
	case OpLT:
		return "LT"
	case OpGE:
		return "GE"
	case OpLE:
		return "LE"

	// Memory operations
	case OpLOAD:
		return "LOAD"
	case OpSTORE:
		return "STORE"
	case OpLOADD:
		return "LOADD"
	case OpSTORED:
		return "STORED"

	// Control flow operations
	case OpJMP:
		return "JMP"
	case OpJMPZ:
		return "JMPZ"
	case OpJMPNZ:
		return "JMPNZ"
	case OpCALL:
		return "CALL"
	case OpRET:
		return "RET"
	case OpHALT:
		return "HALT"
	case OpNOP:
		return "NOP"

	// Math functions
	case OpSQRT:
		return "SQRT"
	case OpSIN:
		return "SIN"
	case OpCOS:
		return "COS"
	case OpTAN:
		return "TAN"
	case OpASIN:
		return "ASIN"
	case OpACOS:
		return "ACOS"
	case OpATAN:
		return "ATAN"
	case OpATAN2:
		return "ATAN2"
	case OpLOG:
		return "LOG"
	case OpLOG10:
		return "LOG10"
	case OpEXP:
		return "EXP"
	case OpPOW:
		return "POW"
	case OpMIN:
		return "MIN"
	case OpMAX:
		return "MAX"
	case OpFLOOR:
		return "FLOOR"
	case OpCEIL:
		return "CEIL"
	case OpROUND:
		return "ROUND"
	case OpTRUNC:
		return "TRUNC"

	default:
		// Custom opcodes (128-255) or unknown
		if op >= 128 {
			return fmt.Sprintf("CUSTOM_%d", op)
		}
		return fmt.Sprintf("UNKNOWN_%d", op)
	}
}

// IsStandardOpcode returns true if the opcode is a standard (non-custom) opcode.
func (op Opcode) IsStandardOpcode() bool {
	return op < 128
}

// IsCustomOpcode returns true if the opcode is in the custom range (128-255).
func (op Opcode) IsCustomOpcode() bool {
	return op >= 128
}
