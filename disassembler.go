package stackvm

import (
	"fmt"
	"strings"
)

// Disassembler converts bytecode programs back to assembly source.
type Disassembler interface {
	// Disassemble converts a program to assembly source.
	Disassemble(program Program) (string, error)

	// SetRegistry enables custom instruction names.
	SetRegistry(registry InstructionRegistry)
}

// DisassemblerOptions configures disassembler output.
type DisassemblerOptions struct {
	// IncludeAddresses adds instruction addresses as comments
	IncludeAddresses bool

	// IncludeMetadata adds program metadata as comments
	IncludeMetadata bool

	// IndentInstructions indents instructions under labels
	IndentInstructions bool
}

// disassembler implements the Disassembler interface.
type disassembler struct {
	registry InstructionRegistry
	options  DisassemblerOptions
}

// NewDisassembler creates a new disassembler with default options.
func NewDisassembler() Disassembler {
	return NewDisassemblerWithOptions(DisassemblerOptions{
		IncludeAddresses:   false,
		IncludeMetadata:    true,
		IndentInstructions: true,
	})
}

// NewDisassemblerWithOptions creates a disassembler with custom options.
func NewDisassemblerWithOptions(opts DisassemblerOptions) Disassembler {
	return &disassembler{
		options: opts,
	}
}

// SetRegistry sets the instruction registry for custom opcodes.
func (d *disassembler) SetRegistry(registry InstructionRegistry) {
	d.registry = registry
}

// Disassemble converts a program to assembly source.
func (d *disassembler) Disassemble(program Program) (string, error) {
	var sb strings.Builder

	// Add metadata if requested
	if d.options.IncludeMetadata {
		metadata := program.Metadata()
		if metadata.Name != "" || metadata.Version != "" || metadata.Author != "" {
			sb.WriteString("; Program Metadata\n")
			if metadata.Name != "" {
				sb.WriteString(fmt.Sprintf("; Name: %s\n", metadata.Name))
			}
			if metadata.Version != "" {
				sb.WriteString(fmt.Sprintf("; Version: %s\n", metadata.Version))
			}
			if metadata.Author != "" {
				sb.WriteString(fmt.Sprintf("; Author: %s\n", metadata.Author))
			}
			if metadata.Description != "" {
				sb.WriteString(fmt.Sprintf("; Description: %s\n", metadata.Description))
			}
			sb.WriteString("\n")
		}
	}

	// Build opcode name map
	opcodeNames := d.makeOpcodeNameMap()

	// Get custom opcode names if registry is set
	if d.registry != nil {
		customNames := d.registry.Names()
		for opcode, name := range customNames {
			opcodeNames[opcode] = name
		}
	}

	// Get symbol table for labels
	symbols := program.SymbolTable()

	// Disassemble instructions
	instructions := program.Instructions()
	for i, inst := range instructions {
		// Check if there's a label at this address
		if label, exists := symbols[i]; exists {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("%s:\n", label))
		}

		// Add address comment if requested
		if d.options.IncludeAddresses {
			sb.WriteString(fmt.Sprintf("; [%04d] ", i))
		}

		// Add indentation if requested
		if d.options.IndentInstructions {
			sb.WriteString("    ")
		}

		// Disassemble instruction
		line, err := d.disassembleInstruction(inst, opcodeNames)
		if err != nil {
			return "", fmt.Errorf("error at instruction %d: %w", i, err)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (d *disassembler) disassembleInstruction(inst Instruction, opcodeNames map[Opcode]string) (string, error) {
	opcodeName, exists := opcodeNames[inst.Opcode]
	if !exists {
		return "", fmt.Errorf("unknown opcode %d", inst.Opcode)
	}

	// Instructions that don't use operands
	if d.hasNoOperand(inst.Opcode) {
		return opcodeName, nil
	}

	// Instructions with numeric operands
	if d.hasNumericOperand(inst.Opcode) {
		return fmt.Sprintf("%s %d", opcodeName, inst.Operand), nil
	}

	// Instructions with label operands (control flow)
	// For disassembly, we just show the address
	// A smarter version would look up the label name from symbol table
	return fmt.Sprintf("%s %d", opcodeName, inst.Operand), nil
}

func (d *disassembler) hasNoOperand(opcode Opcode) bool {
	noOperandOps := []Opcode{
		// Stack
		OpPOP, OpDUP, OpSWAP, OpOVER, OpROT,
		// Arithmetic
		OpADD, OpSUB, OpMUL, OpDIV, OpMOD, OpNEG, OpABS, OpINC, OpDEC,
		// Logic
		OpAND, OpOR, OpNOT, OpXOR,
		// Comparison
		OpEQ, OpNE, OpGT, OpLT, OpGE, OpLE,
		// Memory (dynamic)
		OpLOADD, OpSTORED,
		// Control
		OpRET, OpHALT, OpNOP,
		// Math
		OpSQRT, OpSIN, OpCOS, OpTAN, OpASIN, OpACOS, OpATAN, OpATAN2,
		OpLOG, OpLOG10, OpEXP, OpPOW,
		OpMIN, OpMAX, OpFLOOR, OpCEIL, OpROUND, OpTRUNC,
	}

	for _, op := range noOperandOps {
		if opcode == op {
			return true
		}
	}

	return false
}

func (d *disassembler) hasNumericOperand(opcode Opcode) bool {
	// PUSH, PUSHI, LOAD, STORE, and custom instructions use numeric operands
	return opcode == OpPUSH || opcode == OpPUSHI || opcode == OpLOAD || opcode == OpSTORE || opcode >= 128
}

// makeOpcodeNameMap creates a reverse mapping from opcode to name.
func (d *disassembler) makeOpcodeNameMap() map[Opcode]string {
	return map[Opcode]string{
		// Stack operations
		OpPUSH:  "PUSH",
		OpPUSHI: "PUSHI",
		OpPOP:   "POP",
		OpDUP:   "DUP",
		OpSWAP:  "SWAP",
		OpOVER:  "OVER",
		OpROT:   "ROT",

		// Arithmetic
		OpADD: "ADD",
		OpSUB: "SUB",
		OpMUL: "MUL",
		OpDIV: "DIV",
		OpMOD: "MOD",
		OpNEG: "NEG",
		OpABS: "ABS",
		OpINC: "INC",
		OpDEC: "DEC",

		// Logic
		OpAND: "AND",
		OpOR:  "OR",
		OpNOT: "NOT",
		OpXOR: "XOR",

		// Comparison
		OpEQ: "EQ",
		OpNE: "NE",
		OpGT: "GT",
		OpLT: "LT",
		OpGE: "GE",
		OpLE: "LE",

		// Memory
		OpLOAD:   "LOAD",
		OpSTORE:  "STORE",
		OpLOADD:  "LOADD",
		OpSTORED: "STORED",

		// Control flow
		OpJMP:   "JMP",
		OpJMPZ:  "JMPZ",
		OpJMPNZ: "JMPNZ",
		OpCALL:  "CALL",
		OpRET:   "RET",
		OpHALT:  "HALT",
		OpNOP:   "NOP",

		// Math functions
		OpSQRT:  "SQRT",
		OpSIN:   "SIN",
		OpCOS:   "COS",
		OpTAN:   "TAN",
		OpASIN:  "ASIN",
		OpACOS:  "ACOS",
		OpATAN:  "ATAN",
		OpATAN2: "ATAN2",
		OpLOG:   "LOG",
		OpLOG10: "LOG10",
		OpEXP:   "EXP",
		OpPOW:   "POW",
		OpMIN:   "MIN",
		OpMAX:   "MAX",
		OpFLOOR: "FLOOR",
		OpCEIL:  "CEIL",
		OpROUND: "ROUND",
		OpTRUNC: "TRUNC",
	}
}
