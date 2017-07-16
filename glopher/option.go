package glopher

import (
	"reflect"
)

type Option struct {
	Name  string
	Value interface{}
}

type OptionType struct {
	Name      string
	ValueType reflect.Type
}
