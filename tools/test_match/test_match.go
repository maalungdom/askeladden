package main

import (
	"askeladden/internal/commands"
	"fmt"
)

func main() {
	// Test matching
	fmt.Println("Testing MatchCommand:")

	// Test main command
	commands.MatchAndRunCommand("!hjelp", nil, nil, nil)

	// Test alias
	commands.MatchAndRunCommand("!help", nil, nil, nil)

	// Test another alias
	commands.MatchAndRunCommand("!h", nil, nil, nil)

	// Test non-existing command
	commands.MatchAndRunCommand("!blah", nil, nil, nil)
	fmt.Println("!blah not found (correct)")
}
