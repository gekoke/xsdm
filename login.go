package main

import (
	"os"
	"syscall"
)

func login() error {
	// FIXME: drop to user privileges

	env := os.Environ()

	if err := syscall.Exec("/bin/sh", []string{"-sh"}, env); err != nil {
		return err
	}

	panic("unreachable")
}
