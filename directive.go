package main

type DirectiveToken struct {
	Name   Action
	Repeat int
}

func NewDirectiveToken(name Action, repeat int) *DirectiveToken {
	return &DirectiveToken{
		Name:   name,
		Repeat: repeat,
	}
}
