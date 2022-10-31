package glopher

import (
	"reflect"
)

type Option struct {
	Value interface{}
	Name  string
}

type OptionType struct {
	ValueType reflect.Type
	Name      string
}
