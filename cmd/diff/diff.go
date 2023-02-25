package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ilius/glopher/glopher"
)

var entrySep = strings.Repeat("_", 20)

func main_diff() {
	if len(os.Args) < 3 {
		panic("not enough arguments")
	}
	path1 := os.Args[1]
	path2 := os.Args[2]
	format1 := ""
	format2 := ""
	if len(os.Args) > 4 {
		format1 = os.Args[3]
		format2 = os.Args[4]
	}
	glos1 := glopher.NewGlossary()
	glos2 := glopher.NewGlossary()
	{
		err := glos1.Read(path1, format1)
		if err != nil {
			panic(err)
		}
	}
	{
		err := glos2.Read(path2, format2)
		if err != nil {
			panic(err)
		}
	}
	for entry := range glos1.Iter() {
		fmt.Printf("%v\n%s\n\n", entry, entrySep)
	}
	for entry := range glos2.Iter() {
		fmt.Printf("%v\n%s\n\n", entry, entrySep)
	}
}

func main_word_split_file() {
	textB, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	text := string(textB)
	idx := xmlWordSplit(text)
	// fmt.Printf("%#v\n", idx)
	// fmt.Printf("%#v\n", text[:idx[0]])
	// lastNewline := 0
	for i := 0; i < len(idx)-1; i++ {
		start := idx[i]
		end := idx[i+1]
		str := text[start:end]
		if text[end-1] == '\n' {
			fmt.Printf("%#v\n\n", str)
			// fmt.Printf("%v\n\n", text[lastNewline:end])
			// lastNewline = end
			continue
		}
		fmt.Printf("%#v%s,%s ", str, red, reset)
	}
}

func main_word_diff_files() {
	if len(os.Args) < 3 {
		panic("not enough arguments")
	}
	a_bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	b_bytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	a_str := string(a_bytes)
	b_str := string(b_bytes)
	diff_str := xmlFormattedWordDiff(a_str, b_str)
	fmt.Println(diff_str)
}

func main() {
	main_word_diff_files()
}
