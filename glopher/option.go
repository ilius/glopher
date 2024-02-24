package glopher

import (
	"reflect"
)

type Option struct {
	Value any
	Name  string
}

type OptionType struct {
	ValueType reflect.Type
	Name      string
}
