package main

import (
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
)

func main() {
	// Finn alle yaml-filer i commands-mappa
	files, err := filepath.Glob("internal/commands/*.yaml")
	if err != nil {
		fmt.Printf("Feil ved søking etter filer: %v\n", err)
		return
	}

	fmt.Println("UUID-ar for kommandoar:")
	for _, file := range files {
		// Hent filnamnet utan mappe og ending
		baseName := filepath.Base(file)
		cmdName := strings.TrimSuffix(baseName, ".yaml")

		// Generer UUID
		newUUID := uuid.New()
		fmt.Printf("%s: %s\n", cmdName, newUUID.String())
	}
}
