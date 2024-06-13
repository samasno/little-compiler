package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
  LocalScope SymbolScope = "LOCAL"
)

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(st *SymbolTable) *SymbolTable {
  s := NewSymbolTable()
  s.Outer = st
  return s
}

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
  Outer *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}

  if s.Outer != nil {
    symbol.Scope = LocalScope
  }

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
  if !ok && s.Outer != nil {
    sym, ok = s.Outer.Resolve(name)
  }
	return sym, ok
}
