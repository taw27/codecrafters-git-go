package models

import (
	"errors"
	"fmt"
)

type GitObject struct {
	Type    string
	Size    int
	Content string
}

func (g *GitObject) RunCommand(command string) error {
	switch command {
	case "-p":
		g.PrettyPrint()
	default:
		return errors.New(fmt.Sprintf("Error: command '%s' not recognized", command))
	}

	return nil
}

func (g *GitObject) PrettyPrint() {
	fmt.Printf(g.Content)
}
