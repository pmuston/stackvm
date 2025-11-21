package stackvm

import "time"

// Program represents a sequence of instructions that can be executed by the VM.
type Program interface {
	// Instructions returns the instruction sequence.
	Instructions() []Instruction

	// SymbolTable returns the address to label mapping for debugging.
	// May return nil if no debug information is available.
	SymbolTable() map[int]string

	// Metadata returns program information.
	Metadata() ProgramMetadata
}

// ProgramMetadata contains information about a program.
type ProgramMetadata struct {
	Name        string
	Version     string
	Author      string
	Description string
	Created     time.Time
}

// SimpleProgram is a basic implementation of the Program interface.
type SimpleProgram struct {
	instructions []Instruction
	symbols      map[int]string
	metadata     ProgramMetadata
}

// NewProgram creates a new SimpleProgram with the given instructions.
func NewProgram(instructions []Instruction) *SimpleProgram {
	return &SimpleProgram{
		instructions: instructions,
		symbols:      nil,
		metadata:     ProgramMetadata{},
	}
}

// NewProgramWithMetadata creates a new SimpleProgram with instructions and metadata.
func NewProgramWithMetadata(instructions []Instruction, metadata ProgramMetadata) *SimpleProgram {
	return &SimpleProgram{
		instructions: instructions,
		symbols:      nil,
		metadata:     metadata,
	}
}

// Instructions returns the instruction sequence.
func (p *SimpleProgram) Instructions() []Instruction {
	return p.instructions
}

// SymbolTable returns the address to label mapping.
func (p *SimpleProgram) SymbolTable() map[int]string {
	return p.symbols
}

// Metadata returns program information.
func (p *SimpleProgram) Metadata() ProgramMetadata {
	return p.metadata
}

// SetSymbolTable sets the symbol table for the program.
func (p *SimpleProgram) SetSymbolTable(symbols map[int]string) {
	p.symbols = symbols
}

// AddSymbol adds a single symbol to the symbol table.
func (p *SimpleProgram) AddSymbol(address int, label string) {
	if p.symbols == nil {
		p.symbols = make(map[int]string)
	}
	p.symbols[address] = label
}
