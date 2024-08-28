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
	return []key.Binding{k.Prev, k.Next, k.Submit, k.ToggleShowPassword, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var keys = keyMap{
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("↑tab", "previous"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "submit"),
	),
	ToggleShowPassword: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl-r", "toggle show password"),
	),
}
