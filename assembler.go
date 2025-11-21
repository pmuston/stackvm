package stackvm

import (
	"fmt"
	"os"
	"strings"

	"github.com/pmuston/stackvm/internal/asm"
)

// Assembler converts assembly source code to bytecode programs.
type Assembler interface {
	// Assemble parses and compiles source to a program.
	// Returns an error with line number on failure.
	Assemble(source string) (Program, error)

	// AssembleFile reads a file and assembles it.
	AssembleFile(path string) (Program, error)

	// SetRegistry enables custom instruction names.
	SetRegistry(registry InstructionRegistry)
}

// AssemblerError represents an error during assembly.
type AssemblerError struct {
	Line    int
	Column  int
	Message string
	Source  string // The problematic line
}

func (e *AssemblerError) Error() string {
	if e.Source != "" {
		return fmt.Sprintf("assembler error at %d:%d: %s\n%s", e.Line, e.Column, e.Message, e.Source)
	}
	return fmt.Sprintf("assembler error at %d:%d: %s", e.Line, e.Column, e.Message)
}

// assembler implements the Assembler interface.
type assembler struct {
	registry InstructionRegistry
}

// NewAssembler creates a new assembler.
func NewAssembler() Assembler {
	return &assembler{}
}

// SetRegistry sets the instruction registry for custom opcodes.
func (a *assembler) SetRegistry(registry InstructionRegistry) {
	a.registry = registry
}

// Assemble parses and compiles source to a program.
func (a *assembler) Assemble(source string) (Program, error) {
	// Lexical analysis
	lexer := asm.NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, a.wrapError(err, source)
	}

	// Parsing
	parser := asm.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return nil, a.wrapError(err, source)
	}

	// Code generation
	program, err := a.generate(statements)
	if err != nil {
		return nil, a.wrapError(err, source)
	}

	return program, nil
}

// AssembleFile reads a file and assembles it.
func (a *assembler) AssembleFile(path string) (Program, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	program, err := a.Assemble(string(data))
	if err != nil {
		// Add file path to error message
		if asmErr, ok := err.(*AssemblerError); ok {
			asmErr.Message = fmt.Sprintf("%s (in file %s)", asmErr.Message, path)
			return nil, asmErr
		}
		return nil, fmt.Errorf("failed to assemble %s: %w", path, err)
	}

	return program, nil
}

// generate generates a program from parsed statements.
func (a *assembler) generate(statements []asm.Statement) (Program, error) {
	builder := NewProgramBuilder()
	opcodeMap := makeOpcodeMap()
	customMap := make(map[string]Opcode)

	// Build custom opcode map if registry is set
	if a.registry != nil {
		names := a.registry.Names()
		for opcode, name := range names {
			customMap[strings.ToUpper(name)] = opcode
		}
	}

	// Process statements
	for _, stmt := range statements {
		if stmt.Type == asm.StmtLabel {
			builder.Label(stmt.Label)
		} else if stmt.Type == asm.StmtInstruction {
			if err := a.emitInstruction(builder, stmt, opcodeMap, customMap); err != nil {
				return nil, fmt.Errorf("line %d: %w", stmt.Line, err)
			}
		}
	}

	// Build the program (resolves label references)
	program, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return program, nil
}

func (a *assembler) emitInstruction(builder *ProgramBuilder, stmt asm.Statement, opcodeMap, customMap map[string]Opcode) error {
	opcodeName := strings.ToUpper(stmt.Opcode)

	// Check for standard opcode
	opcode, exists := opcodeMap[opcodeName]
	if !exists {
		// Check for custom opcode
		opcode, exists = customMap[opcodeName]
		if !exists {
			return fmt.Errorf("unknown opcode '%s'", stmt.Opcode)
		}
	}

	// Emit instruction based on opcode and operand
	if stmt.Operand == nil {
		return a.emitNoOperand(builder, opcode)
	} else {
		return a.emitWithOperand(builder, opcode, stmt.Operand)
	}
}

func (a *assembler) emitNoOperand(builder *ProgramBuilder, opcode Opcode) error {
	switch opcode {
	// Stack operations
	case OpPOP:
		builder.Pop()
	case OpDUP:
		builder.Dup()
	case OpSWAP:
		builder.Swap()
	case OpOVER:
		builder.Over()
	case OpROT:
		builder.Rot()

	// Arithmetic
	case OpADD:
		builder.Add()
	case OpSUB:
		builder.Sub()
	case OpMUL:
		builder.Mul()
	case OpDIV:
		builder.Div()
	case OpMOD:
		builder.Mod()
	case OpNEG:
		builder.Neg()
	case OpABS:
		builder.Abs()
	case OpINC:
		builder.Inc()
	case OpDEC:
		builder.Dec()

	// Logic
	case OpAND:
		builder.And()
	case OpOR:
		builder.Or()
	case OpNOT:
		builder.Not()
	case OpXOR:
		builder.Xor()

	// Comparison
	case OpEQ:
		builder.Eq()
	case OpNE:
		builder.Ne()
	case OpGT:
		builder.Gt()
	case OpLT:
		builder.Lt()
	case OpGE:
		builder.Ge()
	case OpLE:
		builder.Le()

	// Memory (dynamic)
	case OpLOADD:
		builder.LoadD()
	case OpSTORED:
		builder.StoreD()

	// Control flow
	case OpRET:
		builder.Ret()
	case OpHALT:
		builder.Halt()
	case OpNOP:
		builder.Nop()

	// Math
	case OpSQRT:
		builder.Sqrt()
	case OpSIN:
		builder.Sin()
	case OpCOS:
		builder.Cos()
	case OpTAN:
		builder.Tan()
	case OpASIN, OpACOS, OpATAN, OpATAN2:
		// These require special handling
		return fmt.Errorf("opcode %d not yet implemented", opcode)
	case OpLOG, OpLOG10, OpEXP, OpPOW:
		// These require special handling
		return fmt.Errorf("opcode %d not yet implemented", opcode)
	case OpMIN:
		builder.Min()
	case OpMAX:
		builder.Max()
	case OpFLOOR:
		builder.Floor()
	case OpCEIL:
		builder.Ceil()
	case OpROUND:
		builder.Round()
	case OpTRUNC:
		// TRUNC not in builder yet
		return fmt.Errorf("opcode TRUNC not yet implemented")

	default:
		// For custom instructions without operands, use operand 0
		if opcode >= 128 {
			builder.Custom(opcode, 0)
		} else {
			return fmt.Errorf("opcode %d requires an operand", opcode)
		}
	}

	return nil
}

func (a *assembler) emitWithOperand(builder *ProgramBuilder, opcode Opcode, operand *asm.Operand) error {
	switch opcode {
	// Stack operations with operands
	case OpPUSH:
		if operand.Type != asm.OperandNumber {
			return fmt.Errorf("PUSH requires a numeric operand")
		}
		if operand.IsFloat {
			builder.Push(operand.FloatValue)
		} else {
			builder.Push(float64(operand.Number))
		}

	case OpPUSHI:
		if operand.Type != asm.OperandNumber {
			return fmt.Errorf("PUSHI requires a numeric operand")
		}
		builder.PushInt(operand.Number)

	// Memory operations with static address
	case OpLOAD:
		if operand.Type != asm.OperandNumber {
			return fmt.Errorf("LOAD requires a numeric operand")
		}
		builder.Load(int(operand.Number))

	case OpSTORE:
		if operand.Type != asm.OperandNumber {
			return fmt.Errorf("STORE requires a numeric operand")
		}
		builder.Store(int(operand.Number))

	// Control flow with labels
	case OpJMP:
		if operand.Type != asm.OperandLabel {
			return fmt.Errorf("JMP requires a label operand")
		}
		builder.Jmp(operand.Label)

	case OpJMPZ:
		if operand.Type != asm.OperandLabel {
			return fmt.Errorf("JMPZ requires a label operand")
		}
		builder.JmpZ(operand.Label)

	case OpJMPNZ:
		if operand.Type != asm.OperandLabel {
			return fmt.Errorf("JMPNZ requires a label operand")
		}
		builder.JmpNZ(operand.Label)

	case OpCALL:
		if operand.Type != asm.OperandLabel {
			return fmt.Errorf("CALL requires a label operand")
		}
		builder.Call(operand.Label)

	default:
		// For custom instructions, use the Custom method
		if opcode >= 128 {
			if operand.Type != asm.OperandNumber {
				return fmt.Errorf("custom instruction requires a numeric operand")
			}
			builder.Custom(opcode, int32(operand.Number))
		} else {
			return fmt.Errorf("opcode %d does not accept operands", opcode)
		}
	}

	return nil
}

// wrapError wraps an error in an AssemblerError if possible.
func (a *assembler) wrapError(err error, source string) error {
	if err == nil {
		return nil
	}

	// Try to extract line information from error message
	// Errors from the lexer/parser/codegen should include line numbers
	// For now, just wrap in a generic AssemblerError
	return &AssemblerError{
		Line:    0,
		Column:  0,
		Message: err.Error(),
		Source:  "",
	}
}

// makeOpcodeMap creates a map of opcode names to opcode values.
func makeOpcodeMap() map[string]Opcode {
	return map[string]Opcode{
		// Stack operations
		"PUSH":   OpPUSH,
		"PUSHI":  OpPUSHI,
		"POP":    OpPOP,
		"DUP":    OpDUP,
		"SWAP":   OpSWAP,
		"OVER":   OpOVER,
		"ROT":    OpROT,

		// Arithmetic
		"ADD": OpADD,
		"SUB": OpSUB,
		"MUL": OpMUL,
		"DIV": OpDIV,
		"MOD": OpMOD,
		"NEG": OpNEG,
		"ABS": OpABS,
		"INC": OpINC,
		"DEC": OpDEC,

		// Logic
		"AND": OpAND,
		"OR":  OpOR,
		"NOT": OpNOT,
		"XOR": OpXOR,

		// Comparison
		"EQ": OpEQ,
		"NE": OpNE,
		"GT": OpGT,
		"LT": OpLT,
		"GE": OpGE,
		"LE": OpLE,

		// Memory
		"LOAD":   OpLOAD,
		"STORE":  OpSTORE,
		"LOADD":  OpLOADD,
		"STORED": OpSTORED,

		// Control flow
		"JMP":   OpJMP,
		"JMPZ":  OpJMPZ,
		"JMPNZ": OpJMPNZ,
		"CALL":  OpCALL,
		"RET":   OpRET,
		"HALT":  OpHALT,
		"NOP":   OpNOP,

		// Math functions
		"SQRT":  OpSQRT,
		"SIN":   OpSIN,
		"COS":   OpCOS,
		"TAN":   OpTAN,
		"ASIN":  OpASIN,
		"ACOS":  OpACOS,
		"ATAN":  OpATAN,
		"ATAN2": OpATAN2,
		"LOG":   OpLOG,
		"LOG10": OpLOG10,
		"EXP":   OpEXP,
		"POW":   OpPOW,
		"MIN":   OpMIN,
		"MAX":   OpMAX,
		"FLOOR": OpFLOOR,
		"CEIL":  OpCEIL,
		"ROUND": OpROUND,
		"TRUNC": OpTRUNC,
	}
}
