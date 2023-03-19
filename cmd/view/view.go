package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ilius/glopher/glopher"
)

var entrySep = strings.Repeat("_", 50)

func main() {
	if len(os.Args) < 2 {
		panic("not enough arguments")
	}
	inputPath := os.Args[1]
	inputFormat := ""
	if len(os.Args) > 2 {
		inputFormat = os.Args[2]
	}
	glos := glopher.NewGlossary()
	err := glos.Read(inputPath, inputFormat)
	if err != nil {
		panic(err)
	}
	index := 0
	for entry := range glos.Iter() {
		fmt.Printf(
			"#%d %s%s\n\n",
			index,
			FormatEntry(entry),
			entrySep,
		)
		index++
	}
}
