package main

import (
	"fmt"
	"askeladden/internal/commands"
)

func main() {
	helpText := commands.GetHelpText()
	fmt.Println(helpText)
}
