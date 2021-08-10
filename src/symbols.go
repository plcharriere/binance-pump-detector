package main

import "strings"

type Symbols struct {
	FirstSymbol  string
	SecondSymbol string
}

func NewSymbols(firstSymbol, secondSymbol string) Symbols {
	return Symbols{
		FirstSymbol:  strings.ToUpper(firstSymbol),
		SecondSymbol: strings.ToUpper(secondSymbol),
	}
}

func (symbols *Symbols) ToString() string {
	return symbols.FirstSymbol + symbols.SecondSymbol
}

func (symbols *Symbols) ToStringSeparated() string {
	return symbols.FirstSymbol + "/" + symbols.SecondSymbol
}
