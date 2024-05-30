package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	return sym, ok
}
