package main

import (
	"fmt"
	"github.com/ilius/glopher/glopher"
	_ "github.com/ilius/glopher/glopher/plugins/tabfile"
	"os"
)

func main() {
	fmt.Println("Supported formats:", glopher.PluginNames())
	if len(os.Args) == 3 {
		plug := glopher.PluginByName("tabfile")
		if plug == nil {
			panic("tabfile plugin was not found")
		}
		reader, err := plug.Read(os.Args[1])
		if err != nil {
			panic(err)
		}
		err = plug.Write(os.Args[2], reader, nil, nil)
		if err != nil {
			panic(err)
		}
	}
}
