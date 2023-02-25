package main

import (
	"fmt"
	"github.com/ilius/glopher/glopher"
	"os"
)

func main() {
	fmt.Println("Supported formats:", glopher.PluginNames())
	if len(os.Args) == 3 {
		inputPath := os.Args[1]
		outputPath := os.Args[2]
		glos := glopher.NewGlossary()
		glos.SetProgressBar(NewCmdProgressBar())
		err := glos.Read(inputPath, "")
		if err != nil {
			panic(err)
		}
		err = glos.Write(outputPath, "")
		if err != nil {
			panic(err)
		}
	}
}
