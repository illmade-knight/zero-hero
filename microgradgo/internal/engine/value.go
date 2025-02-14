package engine

import (
	"errors"
	"fmt"
	"math"
)

type ParentGradient struct {
	Label    string
	Gradient float64
}
type ValueGradient struct {
	ParentA ParentGradient
	ParentB ParentGradient
}
type Value struct {
	Data     float64
	Gradient ValueGradient
	Label    string
	Op       string
	parents  []*Value
}

func NewValue(data float64, label string) Value {
	return Value{data, ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "", []*Value{}}
}

func (v Value) Add(value Value, label string) Value {
	return Value{v.Data + value.Data, ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "+", []*Value{&v, &value}}
}

func (v Value) Sub(value Value, label string) Value {
	return Value{v.Data - value.Data, ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "-", []*Value{&v, &value}}
}

func (v Value) Mul(value Value, label string) Value {
	return Value{v.Data * value.Data, ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "*", []*Value{&v, &value}}
}

func (v Value) Div(value Value, label string) Value {
	return Value{v.Data / value.Data, ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "/", []*Value{&v, &value}}
}

func (v Value) Pow(value Value, label string) Value {
	return Value{math.Pow(v.Data, value.Data), ValueGradient{
		ParentGradient{label, 1.0},
		ParentGradient{label, 1.0},
	}, label, "**", []*Value{&v, &value}}
}

func RecurseValuesGradient(node *Value) error {
	if (len(node.parents) > 0) && node.Op != "" {
		if len(node.parents) != 2 {
			return errors.New(fmt.Sprintf("A node should have zero or 2 parents. This node has %s. Node is %v", len(node.parents), node))
		}
		parentA := node.parents[0]
		parentB := node.parents[1]

		switch node.Op {
		case "+":
			node.Gradient.ParentA.Gradient = 1.0
			node.Gradient.ParentB.Gradient = 1.0
		case "-":
			node.Gradient.ParentA.Gradient = 1.0
			node.Gradient.ParentB.Gradient = 1.0
		case "*":
			node.Gradient.ParentA.Gradient = parentB.Data
			node.Gradient.ParentB.Gradient = parentA.Data
		case "/":
			node.Gradient.ParentA.Gradient = 1 / parentB.Data
			node.Gradient.ParentB.Gradient = -1 * parentA.Data / (parentB.Data * parentB.Data)
		case "**":
			node.Gradient.ParentA.Gradient = parentB.Data * math.Pow(parentA.Data, parentB.Data-1)
			node.Gradient.ParentB.Gradient = math.Log(parentA.Data) * math.Pow(parentA.Data, parentB.Data)
		default:
			return errors.New(fmt.Sprintf("Unknown op %s for node  %v", node.Op, node))
		}
		node.Gradient.ParentA.Label = parentA.Label
		node.Gradient.ParentB.Label = parentB.Label

		RecurseValuesGradient(parentA)
		RecurseValuesGradient(parentB)
	}
	return nil
}
