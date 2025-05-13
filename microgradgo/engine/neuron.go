package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func NewInputs(values []float64) []*Value {
	inputSlice := make([]*Value, 0)
	for _, value := range values {
		inputSlice = append(inputSlice, NewValue(value, "input"))
	}
	return inputSlice
}

type Parameters []*Value

type Neuron struct {
	Weights []*Value `json:"weights"`
	Bias    *Value   `json:"bias"`
}

func ParametersNeuron(neuron Neuron) Parameters {
	return append(neuron.Weights, neuron.Bias)
}

func NewNeuron(numInputs int) Neuron {
	weights := make([]*Value, numInputs)
	for index, _ := range weights {
		weights[index] = NewRandomValueWithLabelPrefix("w")
	}
	return Neuron{weights, NewRandomValueWithLabelPrefix("b")}
}

func Forward(neuron Neuron, inputs []*Value) (*Value, error) {
	if len(neuron.Weights) != len(inputs) {
		return nil, errors.New(fmt.Sprintf("The neuron is set to take input of length %v while the input has length %v)", len(neuron.Weights), len(inputs)))
	}
	result := NewValue(0.0, "")
	for idx, _ := range neuron.Weights {
		product := Mul(neuron.Weights[idx], inputs[idx], "")
		result = Add(result, product, "")
	}
	result = Add(result, neuron.Bias, "")
	result = Tanh(result, "")

	return result, nil
}

type Layer struct {
	Neurons []Neuron `json:"neurons"`
}

func ParametersLayer(layer Layer) Parameters {
	params := make([]*Value, 0)
	for _, neuron := range layer.Neurons {
		params = append(params, ParametersNeuron(neuron)...)
	}
	return params
}

func NewLayer(numInputs, numNeurons int) Layer {
	neurons := make([]Neuron, numNeurons)
	for index, _ := range neurons {
		neurons[index] = NewNeuron(numInputs)
	}
	return Layer{neurons}
}

func ForwardLayer(layer Layer, inputs []*Value) ([]*Value, error) {
	outputs := make([]*Value, len(layer.Neurons))
	for idx, neuron := range layer.Neurons {
		output, err := Forward(neuron, inputs)
		if err != nil {
			return nil, err
		} else {
			outputs[idx] = output
		}
	}
	return outputs, nil
}

type MLP struct {
	NumInputs  int     `json:"inputs"`
	NumOutputs int     `json:"outputs"`
	Layers     []Layer `json:"layers"`
}

func ParametersMLP(mlp MLP) Parameters {
	params := make([]*Value, 0)
	for _, layer := range mlp.Layers {
		params = append(params, ParametersLayer(layer)...)
	}
	return params
}

func NewMLP(numInputs, numOutputs int, hiddenLayerSizes []int) MLP {
	layers := make([]Layer, len(hiddenLayerSizes)+1)
	inputSize := numInputs
	for index, size := range hiddenLayerSizes {
		layers[index] = NewLayer(inputSize, size)
		inputSize = size
	}
	layers[len(hiddenLayerSizes)] = NewLayer(hiddenLayerSizes[len(hiddenLayerSizes)-1], numOutputs)
	return MLP{numInputs, numOutputs, layers}
}

func ToJSON(mlp MLP) ([]byte, error) {
	jsonData, err := json.Marshal(mlp)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func FromJSON(jsonData []byte) (MLP, error) {
	var mlp MLP
	err := json.Unmarshal(jsonData, &mlp)
	if err != nil {
		return MLP{}, err
	}
	return mlp, nil
}

// forward pass
func ForwardMLP(mlp MLP, inputs []*Value) ([]*Value, error) {
	outputs := inputs
	var err error
	for _, layer := range mlp.Layers {
		outputs, err = ForwardLayer(layer, outputs)
		if err != nil {
			return nil, err
		}
	}
	return outputs, nil
}

// forward pass
func ForwardMLPMultipleInputs(mlp MLP, inputs [][]*Value) ([][]*Value, error) {
	// tried a go routine here which ended up being twice as slow
	outputs := make([][]*Value, len(inputs))
	for idx, input := range inputs {
		output, _ := ForwardMLP(mlp, input)
		outputs[idx] = output
	}

	return outputs, nil
}

// zero the grad before re-doing back propagation
func ZeroGrad(mlp MLP) {
	for _, value := range ParametersMLP(mlp) {
		value.GradientData = 0.0
	}
}

// adjust parameters with simple stochastic gradient descent updates
// learning rate can too high and you will overstep (explode the loss function), too low and it will take too long to converge
func AdjustParametersSimpleGradientDescent(mlp MLP, learningRate float64) {
	for _, value := range ParametersMLP(mlp) {
		value.Data += learningRate * value.GradientData
	}
}

// simple squared difference error function
func LossSquaredError(outputs, targets []*Value) (*Value, error) {
	if len(outputs) != len(targets) {
		return nil, errors.New(fmt.Sprintf("The length of the outputs array  %v should be the same length as the length of targets array %v", len(outputs), len(targets)))
	}
	loss := NewValue(0.0, "loss")
	for idx, _ := range outputs {
		tmp := Sub(outputs[idx], targets[idx], "")
		tmp = Mul(tmp, tmp, "")
		loss = Add(loss, tmp, "")
	}
	return loss, nil
}

// sum squared difference error function across 2 dimensional array
// TODO - calculate the square difference rather than the avg?
func LossSquaredErrorMultiple(outputs, targets [][]*Value) (*Value, error) {
	if len(outputs) != len(targets) {
		return nil, errors.New(fmt.Sprintf("The length of the 2 dimensional outputs array  %v should be the same length as the length of 2 dimensional targets array %v", len(outputs), len(targets)))
	}
	loss := NewValue(0.0, "loss")
	for idx, _ := range outputs {
		tmp, err := LossSquaredError(outputs[idx], targets[idx])
		if err != nil {
			return nil, err
		}
		loss = Add(loss, tmp, "")
	}
	return Div(loss, NewValue(float64(len(outputs)), ""), ""), nil
}

// run forward pass against input, calculate loss, do backward propagation and adjust the parameters
// return loss, outputs
func TrainSingle(mlp MLP, inputs []*Value, targetOutputs []*Value, learningRate float64) (*Value, []*Value, error) {
	inputsMultiple := [][]*Value{inputs}
	targetOutputsMultiple := [][]*Value{targetOutputs}

	loss, outputs, err := TrainMultiple(mlp, inputsMultiple, targetOutputsMultiple, learningRate)
	return loss, outputs[0], err
}

// run forward pass against multiple inputs, calculate loss, do backward propagation and adjust the parameters
// return loss, outputs
func TrainMultiple(mlp MLP, inputsMultiple [][]*Value, targetOutputsMultiple [][]*Value, learningRate float64) (*Value, [][]*Value, error) {
	var loss *Value

	verbose := GetConfig().Verbose

	// forward pass
	if verbose {
		log.Println("Starting training with batch of inputs of length " + strconv.Itoa(len(inputsMultiple)))
	}
	outputs, err := ForwardMLPMultipleInputs(mlp, inputsMultiple)
	if err != nil {
		return nil, nil, err
	}
	if verbose {
		log.Println("Completed forward pass")
	}

	// calculate loss
	loss, err = LossSquaredErrorMultiple(outputs, targetOutputsMultiple)
	if err != nil {
		return nil, nil, err
	}
	if verbose {
		log.Println("Loss calculated")
	}

	// zero the grad before re-doing back propagation
	ZeroGrad(mlp)
	if verbose {
		log.Println("Zeroed gradients/parameters")
	}

	// back propagation
	loss.BackPropagate()
	if verbose {
		log.Println("Back propagation completed")
	}

	AdjustParametersSimpleGradientDescent(mlp, learningRate)
	if verbose {
		log.Println("Adjusted parameters")
	}

	return loss, outputs, nil
}
