package stackvm

import "fmt"

// ProgramBuilder provides a fluent API for constructing programs.
type ProgramBuilder struct {
	instructions []Instruction
	labels       map[string]int  // label name -> instruction index
	references   []labelRef      // unresolved label references
	metadata     ProgramMetadata
}

// labelRef tracks an unresolved label reference.
type labelRef struct {
	labelName string
	instIndex int // index of instruction that references the label
}

// NewProgramBuilder creates a new ProgramBuilder.
func NewProgramBuilder() *ProgramBuilder {
	return &ProgramBuilder{
		instructions: make([]Instruction, 0),
		labels:       make(map[string]int),
		references:   make([]labelRef, 0),
	}
}

// Stack Operations

// Push adds a PUSH instruction (push float value).
func (b *ProgramBuilder) Push(v float64) *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpPUSH, int32(v)))
	return b
}

// PushInt adds a PUSHI instruction (push int value).
func (b *ProgramBuilder) PushInt(v int64) *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpPUSHI, int32(v)))
	return b
}

// Pop adds a POP instruction.
func (b *ProgramBuilder) Pop() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpPOP, 0))
	return b
}

// Dup adds a DUP instruction.
func (b *ProgramBuilder) Dup() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpDUP, 0))
	return b
}

// Swap adds a SWAP instruction.
func (b *ProgramBuilder) Swap() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSWAP, 0))
	return b
}

// Over adds an OVER instruction.
func (b *ProgramBuilder) Over() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpOVER, 0))
	return b
}

// Rot adds a ROT instruction.
func (b *ProgramBuilder) Rot() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpROT, 0))
	return b
}

// Arithmetic Operations

// Add adds an ADD instruction.
func (b *ProgramBuilder) Add() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpADD, 0))
	return b
}

// Sub adds a SUB instruction.
func (b *ProgramBuilder) Sub() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSUB, 0))
	return b
}

// Mul adds a MUL instruction.
func (b *ProgramBuilder) Mul() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpMUL, 0))
	return b
}

// Div adds a DIV instruction.
func (b *ProgramBuilder) Div() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpDIV, 0))
	return b
}

// Mod adds a MOD instruction.
func (b *ProgramBuilder) Mod() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpMOD, 0))
	return b
}

// Neg adds a NEG instruction.
func (b *ProgramBuilder) Neg() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpNEG, 0))
	return b
}

// Abs adds an ABS instruction.
func (b *ProgramBuilder) Abs() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpABS, 0))
	return b
}

// Inc adds an INC instruction.
func (b *ProgramBuilder) Inc() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpINC, 0))
	return b
}

// Dec adds a DEC instruction.
func (b *ProgramBuilder) Dec() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpDEC, 0))
	return b
}

// Logic Operations

// And adds an AND instruction.
func (b *ProgramBuilder) And() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpAND, 0))
	return b
}

// Or adds an OR instruction.
func (b *ProgramBuilder) Or() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpOR, 0))
	return b
}

// Not adds a NOT instruction.
func (b *ProgramBuilder) Not() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpNOT, 0))
	return b
}

// Xor adds an XOR instruction.
func (b *ProgramBuilder) Xor() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpXOR, 0))
	return b
}

// Comparison Operations

// Eq adds an EQ instruction.
func (b *ProgramBuilder) Eq() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpEQ, 0))
	return b
}

// Ne adds a NE instruction.
func (b *ProgramBuilder) Ne() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpNE, 0))
	return b
}

// Gt adds a GT instruction.
func (b *ProgramBuilder) Gt() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpGT, 0))
	return b
}

// Lt adds a LT instruction.
func (b *ProgramBuilder) Lt() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpLT, 0))
	return b
}

// Ge adds a GE instruction.
func (b *ProgramBuilder) Ge() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpGE, 0))
	return b
}

// Le adds a LE instruction.
func (b *ProgramBuilder) Le() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpLE, 0))
	return b
}

// Memory Operations

// Load adds a LOAD instruction.
func (b *ProgramBuilder) Load(index int) *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpLOAD, int32(index)))
	return b
}

// Store adds a STORE instruction.
func (b *ProgramBuilder) Store(index int) *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSTORE, int32(index)))
	return b
}

// LoadD adds a LOADD instruction (load dynamic).
func (b *ProgramBuilder) LoadD() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpLOADD, 0))
	return b
}

// StoreD adds a STORED instruction (store dynamic).
func (b *ProgramBuilder) StoreD() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSTORED, 0))
	return b
}

// Control Flow Operations

// Label defines a label at the current position.
func (b *ProgramBuilder) Label(name string) *ProgramBuilder {
	b.labels[name] = len(b.instructions)
	return b
}

// Jmp adds a JMP instruction to the specified label.
func (b *ProgramBuilder) Jmp(label string) *ProgramBuilder {
	instIndex := len(b.instructions)
	b.instructions = append(b.instructions, NewInstruction(OpJMP, 0)) // Will be resolved later
	b.references = append(b.references, labelRef{label, instIndex})
	return b
}

// JmpZ adds a JMPZ instruction to the specified label.
func (b *ProgramBuilder) JmpZ(label string) *ProgramBuilder {
	instIndex := len(b.instructions)
	b.instructions = append(b.instructions, NewInstruction(OpJMPZ, 0))
	b.references = append(b.references, labelRef{label, instIndex})
	return b
}

// JmpNZ adds a JMPNZ instruction to the specified label.
func (b *ProgramBuilder) JmpNZ(label string) *ProgramBuilder {
	instIndex := len(b.instructions)
	b.instructions = append(b.instructions, NewInstruction(OpJMPNZ, 0))
	b.references = append(b.references, labelRef{label, instIndex})
	return b
}

// Call adds a CALL instruction to the specified label.
func (b *ProgramBuilder) Call(label string) *ProgramBuilder {
	instIndex := len(b.instructions)
	b.instructions = append(b.instructions, NewInstruction(OpCALL, 0))
	b.references = append(b.references, labelRef{label, instIndex})
	return b
}

// Ret adds a RET instruction.
func (b *ProgramBuilder) Ret() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpRET, 0))
	return b
}

// Halt adds a HALT instruction.
func (b *ProgramBuilder) Halt() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpHALT, 0))
	return b
}

// Nop adds a NOP instruction.
func (b *ProgramBuilder) Nop() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpNOP, 0))
	return b
}

// Math Functions

// Sqrt adds a SQRT instruction.
func (b *ProgramBuilder) Sqrt() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSQRT, 0))
	return b
}

// Sin adds a SIN instruction.
func (b *ProgramBuilder) Sin() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpSIN, 0))
	return b
}

// Cos adds a COS instruction.
func (b *ProgramBuilder) Cos() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpCOS, 0))
	return b
}

// Tan adds a TAN instruction.
func (b *ProgramBuilder) Tan() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpTAN, 0))
	return b
}

// Min adds a MIN instruction.
func (b *ProgramBuilder) Min() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpMIN, 0))
	return b
}

// Max adds a MAX instruction.
func (b *ProgramBuilder) Max() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpMAX, 0))
	return b
}

// Floor adds a FLOOR instruction.
func (b *ProgramBuilder) Floor() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpFLOOR, 0))
	return b
}

// Ceil adds a CEIL instruction.
func (b *ProgramBuilder) Ceil() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpCEIL, 0))
	return b
}

// Round adds a ROUND instruction.
func (b *ProgramBuilder) Round() *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(OpROUND, 0))
	return b
}

// Custom Operations

// Custom adds a custom instruction with the specified opcode and operand.
func (b *ProgramBuilder) Custom(opcode Opcode, operand int32) *ProgramBuilder {
	b.instructions = append(b.instructions, NewInstruction(opcode, operand))
	return b
}

// Metadata Operations

// SetMetadata sets the program metadata.
func (b *ProgramBuilder) SetMetadata(metadata ProgramMetadata) *ProgramBuilder {
	b.metadata = metadata
	return b
}

// Build constructs the final Program.
// Returns an error if there are unresolved label references.
func (b *ProgramBuilder) Build() (Program, error) {
	// Resolve label references
	for _, ref := range b.references {
		targetAddr, exists := b.labels[ref.labelName]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrUnresolvedLabel, ref.labelName)
		}
		// Update the instruction's operand with the target address
		b.instructions[ref.instIndex].Operand = int32(targetAddr)
	}

	// Create symbol table from labels
	symbols := make(map[int]string)
	for name, addr := range b.labels {
		symbols[addr] = name
	}

	program := NewProgramWithMetadata(b.instructions, b.metadata)
	program.SetSymbolTable(symbols)

	return program, nil
}
