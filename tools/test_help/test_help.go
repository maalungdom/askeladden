package main

import (
	"askeladden/internal/commands"
	"fmt"
)

func main() {
	helpText := commands.GetHelpText()
	fmt.Println(helpText)
}
