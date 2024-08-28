package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Next               key.Binding
	Prev               key.Binding
	Submit             key.Binding
	ToggleShowPassword key.Binding
	Quit               key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit, k.ToggleShowPassword}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var keys = keyMap{
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	ToggleShowPassword: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl-r", "toggle show password"),
	),
}
