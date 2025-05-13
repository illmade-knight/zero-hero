package engine

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		value2 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{1.0, "a"}, fields{1.0, "b"}, "c", NewValue(2.0, "c")},
		{"simple2", fields{-1.0, "a"}, fields{1.0, "b"}, "c", NewValue(0.0, "c")},
		{"simple3", fields{-110.12121, "a"}, fields{1234.1212099090, "b"}, "c", NewValue(1123.999999909, "c")},
		{"simple4", fields{1234567890, "a"}, fields{1234567890, "b"}, "c", NewValue(2469135780, "c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			value2 := NewValue(tt.value2.Data, tt.value2.Label)
			fmt.Println()
			if got := Add(value1, value2, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Add() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestSub(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		value2 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{1.0, "a"}, fields{1.0, "b"}, "c", NewValue(0.0, "c")},
		{"simple2", fields{-1.0, "a"}, fields{1.0, "b"}, "c", NewValue(-2.0, "c")},
		{"simple3", fields{-110.12121, "a"}, fields{1234.1212099090, "b"}, "c", NewValue(-1344.242419909, "c")},
		{"simple4", fields{1234567891, "a"}, fields{1234567890, "b"}, "c", NewValue(1, "c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			value2 := NewValue(tt.value2.Data, tt.value2.Label)
			fmt.Println()
			if got := Sub(value1, value2, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Sub() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestMul(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		value2 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{1.0, "a"}, fields{1.0, "b"}, "c", NewValue(1.0, "c")},
		{"simple2", fields{-1.0, "a"}, fields{1.0, "b"}, "c", NewValue(-1.0, "c")},
		{"simple3", fields{-2.5, "a"}, fields{-2.5, "b"}, "c", NewValue(6.25, "c")},
		{"simple4", fields{123456, "a"}, fields{7890, "b"}, "c", NewValue(974067840, "c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			value2 := NewValue(tt.value2.Data, tt.value2.Label)
			fmt.Println()
			if got := Mul(value1, value2, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Mul() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		value2 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{1.0, "a"}, fields{1.0, "b"}, "c", NewValue(1.0, "c")},
		{"simple2", fields{-1.0, "a"}, fields{1.0, "b"}, "c", NewValue(-1.0, "c")},
		{"simple3", fields{10.0, "a"}, fields{2.0, "b"}, "c", NewValue(5.0, "c")},
		{"simple4", fields{123456, "a"}, fields{7890, "b"}, "c", NewValue(15.647148288973384, "c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			value2 := NewValue(tt.value2.Data, tt.value2.Label)
			fmt.Println()
			if got := Div(value1, value2, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Div() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestPow(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		value2 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{1.0, "a"}, fields{1.0, "b"}, "c", NewValue(1.0, "c")},
		{"simple2", fields{10.0, "a"}, fields{0.0, "b"}, "c", NewValue(1.0, "c")},
		{"simple3", fields{10.0, "a"}, fields{2.0, "b"}, "c", NewValue(100.0, "c")},
		{"simple4", fields{2.5, "a"}, fields{2.5, "b"}, "c", NewValue(9.882117688026186, "c")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			value2 := NewValue(tt.value2.Data, tt.value2.Label)
			fmt.Println()
			if got := Pow(value1, value2, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Pow() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestTanh(t *testing.T) {
	type fields struct {
		Data  float64
		Label string
	}
	tests := []struct {
		name   string
		value1 fields
		label  string
		want   *Value
	}{
		{"simple1", fields{0.0, "a"}, "b", NewValue(0.0, "b")},
		{"simple2", fields{1.0, "a"}, "b", NewValue(0.76159415595576488811945828260479359041276859725793655159681050012195324457, "b")},
		{"simple3", fields{2.0, "a"}, "b", NewValue(0.9640275800758169, "b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value1 := NewValue(tt.value1.Data, tt.value1.Label)
			fmt.Println()
			if got := Tanh(value1, tt.label); got.Data != tt.want.Data {
				t.Errorf("Value.Tanh() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestCalculateAncestorGradient(t *testing.T) {
	Init(true, false)

	// run thru some basic calculations
	a := NewValue(3.0, "a")
	b := NewValue(-2.0, "b")
	c := NewValue(4.0, "c")
	d := NewValue(23.12, "d")
	e := NewValue(-423.56, "e")
	f := NewValue(36.9, "f")
	g := NewValue(23232.3467, "g")
	h := NewValue(-727.82, "h")
	i := Pow(a, b, "i")
	j := Sub(c, d, "j")
	k := Mul(e, f, "k")
	l := Div(g, h, "l")
	m := Mul(i, j, "m")
	n := Div(k, l, "n")
	O := Add(m, n, "o")

	// simple neural net
	x1 := NewValue(2.0, "x1")
	x2 := NewValue(0.0, "x2")
	w1 := NewValue(-3.0, "w1")
	w2 := NewValue(1.0, "w2")
	bias := NewValue(6.8813735870195432, "bias")

	x1w1 := Mul(x1, w1, "x1*w1")
	x2w2 := Mul(x2, w2, "x2*w2")
	x1w1_x2w2 := Add(x1w1, x2w2, "x1*w1 + x2*w2")
	neural := Add(x1w1_x2w2, bias, "neural")

	neuralOutput := Tanh(neural, "neuralOutput")

	// multi-variable case
	multia := NewValue(-2.0, "multia")
	multib := NewValue(3.0, "multib")
	multid := Mul(multia, multib, "multid")
	multie := Add(multia, multib, "multie")
	multif := Mul(multid, multie, "multif")

	tests := []struct {
		name     string
		output   *Value
		ancestor *Value
		want     float64
	}{
		{"simple a", O, a, 1.4162962962962964},
		{"simple b", O, b, -2.3339407732593616},
		{"simple c", O, c, 0.1111111111111111},
		{"simple d", O, d, -0.1111111111111111},
		{"simple e", O, e, -1.1559985027255126},
		{"simple f", O, f, 13.26923376190835},
		{"simple g", O, g, -0.02107556038729476},
		{"simple h", O, h, -0.6727415100085434},
		{"neural x1", neuralOutput, x1, -1.4999999999999996},
		{"neural x1", neuralOutput, x2, 0.4999999999999999},
		{"neural x1", neuralOutput, w1, 0.9999999999999998},
		{"neural x1", neuralOutput, w2, 0.0},
		{"multi-variable a", multif, multia, -3.0},
		{"multi-variable b", multif, multib, -8.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.output.AncestorDerivatives[tt.ancestor.UUID].GradientData; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateAncestorGradient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAncestorGradientSelf(t *testing.T) {
	Init(true, false)

	// operating on self case
	selfa := NewValue(17.53, "selfa")
	selfOutputAdd := Add(selfa, selfa, "selfOutputAdd")
	selfOutputSub := Sub(selfa, selfa, "selfOutputSub")
	selfOutputMul := Mul(selfa, selfa, "selfOutputMul")
	selfOutputDiv := Div(selfa, selfa, "selfOutputDiv")
	selfOutputPow := Pow(selfa, selfa, "selfOutputPow")

	tests := []struct {
		name   string
		output *Value
		parent *Value
		want   float64
	}{
		{"self operating add", selfOutputAdd, selfa, 2},
		{"self operating sub", selfOutputSub, selfa, 0},
		{"self operating mul", selfOutputMul, selfa, 2 * 17.53},
		{"self operating div", selfOutputDiv, selfa, 0},
		{"self operating pow", selfOutputPow, selfa, 1187.38415720448392384136203281126852494785640925855515736143890026596058},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.output.Derivatives[tt.parent.UUID].GradientData; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateAncestorGradient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeuralNetBack(t *testing.T) {

	// input values
	x1 := NewValue(2.0, "x1")
	x2 := NewValue(0.0, "x2")
	// weights
	w1 := NewValue(-3.0, "w1")
	w2 := NewValue(1.0, "w2")
	// bias of neuron
	b := NewValue(6.8813735870195432, "b")

	x1w1 := Mul(x1, w1, "x1*w1")
	x2w2 := Mul(x2, w2, "x2*w2")
	x1w1_x2w2 := Add(x1w1, x2w2, "x1*w1 + x2*w2")
	n := Add(x1w1_x2w2, b, "n")

	// activation function
	o := Tanh(n, "o")

	o.BackPropagate()

	tests := []struct {
		name  string
		value *Value
		want  float64
	}{
		{"self operating add", o, 1.0},
		{"self operating add", n, 0.5},
		{"self operating add", b, 0.5},
		{"self operating add", x1w1_x2w2, 0.5},
		{"self operating add", x2w2, 0.5},
		{"self operating add", x1w1, 0.5},
		{"self operating add", x1, -1.5},
		{"self operating add", x2, 0.5},
		{"self operating add", w1, 1.0},
		{"self operating add", w2, 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.GradientData; !almostEqual(got, tt.want) {
				t.Errorf("TestNeuralNetBack() = %v, want %v", got, tt.want)
			}
		})
	}
}

const float64EqualityThreshold = 1e-15

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
