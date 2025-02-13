package engine

import (
	"math"
)

type Value struct {
	Data  float64
	Label string
	Op    string
	prev  []Value
}

func NewValue(data float64, label string) Value {
	return Value{data, label, "", []Value{}}
}

func (v Value) Add(value Value, label string) Value {
	return Value{v.Data + value.Data, label, "+", []Value{v, value}}
}

func (v Value) Sub(value Value, label string) Value {
	return Value{v.Data - value.Data, label, "-", []Value{v, value}}
}

func (v Value) Mul(value Value, label string) Value {
	return Value{v.Data * value.Data, label, "*", []Value{v, value}}
}

func (v Value) Div(value Value, label string) Value {
	return Value{v.Data / value.Data, label, "/", []Value{v, value}}
}

func (v Value) Pow(value Value, label string) Value {
	return Value{math.Pow(v.Data, value.Data), label, "**", []Value{v, value}}
}
