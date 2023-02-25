package main

import (
	"fmt"

	"github.com/ilius/glopher/glopher"
)

func FormatEntry(entry *glopher.Entry) string {
	s := fmt.Sprintf(">>> %s\n", entry.Word)
	for _, alt := range entry.AltWord {
		s += fmt.Sprintf("Alt: %s\n", alt)
	}
	s += fmt.Sprintf("\n%s\n", entry.Defi)
	for _, defi := range entry.AltDefi {
		s += fmt.Sprintf("\n%s\n", defi)
	}
	return s
}
