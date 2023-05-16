package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()
	// TODO: Determine if the "err" variable is ErrNotFound
	fmt.Println(errors.Is(errors.Unwrap(err), ErrNotFound))

}

var ErrNotFound = errors.New("not found")

func A() error {
	return ErrNotFound
}

func B() error {
	err := A()
	if err != nil {
		return fmt.Errorf("b: %w", err)
	}

	return nil
}
