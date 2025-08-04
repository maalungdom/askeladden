package main

import (
	"fmt"
	"roersla.no/askeladden/internal/commands"
)

func main() {
	helpText := commands.GetHelpText()
	fmt.Println(helpText)
}
