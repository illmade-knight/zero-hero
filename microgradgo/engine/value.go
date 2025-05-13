package engine

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

type Config struct {
	PropagationForward bool
	Verbose            bool
}

var (
	config     *Config   // Package-level "static" config
	configOnce sync.Once // Ensures initialization happens once
)

// forward propagation will be calculated as well (if true)
func Init(propagationForward, verbose bool) {
	configOnce.Do(func() {
		config = &Config{
			PropagationForward: propagationForward,
			Verbose:            verbose,
		}
	})
}

// GetConfig returns the initialized config, will default to false
func GetConfig() *Config {
	if config == nil {
		config = &Config{
			PropagationForward: false,
			Verbose:            false,
		}
	}
	return config
}

type BackwardFunctionType func(value *Value)

// implementing custom marshaling/unmarshaling for BackwardFunctionType (so loaded model has DefaultBackward)
func (backwardFunctionType BackwardFunctionType) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (backwardFunctionType *BackwardFunctionType) UnmarshalJSON(input []byte) error {
	if string(input) == "null" {
		*backwardFunctionType = DefaultBackward
		return nil
	} else {
		return errors.New(fmt.Sprintf("Can't unmarshal Backward function: %v", input))
	}
}

func DefaultBackward(value *Value) {}

// gradient/derivative is of node wrt to another node (identified by UUID of the associated Value)
type Derivative struct {
	VariableValueUUID uuid.UUID
	VariableLabel     string
	GradientData      float64
}

// as Values are operated on (with themselves or other values), if config.PropagationForward is true
// then the gradients/derivatives (wrt parents) and ancestor gradients/derivatives (wrt) ancestral nodes are calculated.
type Value struct {
	UUID                uuid.UUID                `json:"uuid"`
	Label               string                   `json:"label"` // mutable, prob should be unique
	Data                float64                  `json:"data"`
	Parents             []*Value                 `json:"-"` // 0, 1 or 2 parents (more parents? TODO)
	Op                  string                   `json:"-"` // op used on parents (if they exist) to produce current Value
	Backward            BackwardFunctionType     `json:"backward"`
	GradientData        float64                  `json:"-"` // gradient/derivative of a node
	Derivatives         map[uuid.UUID]Derivative `json:"-"` // gradients/derivatives of a node wrt parents (if they exist)
	AncestorDerivatives map[uuid.UUID]Derivative `json:"-"` // gradients/derivatives of a node wrt ancestral nodes
}

func NewValue(data float64, label string) *Value {
	return &Value{uuid.New(), label, data, []*Value{}, "", DefaultBackward, 0.0, make(map[uuid.UUID]Derivative), make(map[uuid.UUID]Derivative)}
}

func NewRandomValue() *Value {
	return NewValue(rand.Float64()*2-1, RandStringRunes(4))
}

func NewRandomValueWithLabelPrefix(prefix string) *Value {
	return NewValue(rand.Float64()*2-1, prefix+"-"+RandStringRunes(4))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Add(v, value *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent1 := value.Parents[0]
		parent2 := value.Parents[1]
		if reflect.DeepEqual(parent1, parent2) {
			parent1.GradientData += 2.0 * value.GradientData
		} else {
			(parent1).GradientData += value.GradientData
			(parent2).GradientData += value.GradientData
		}
	}
	if GetConfig().PropagationForward {
		if reflect.DeepEqual(v, value) {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 2.0}
		} else {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 1.0}
			derivatives[value.UUID] = Derivative{value.UUID, value.Label, 1.0}
		}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v, value}, derivatives)
	}

	return &Value{uuid.New(), label, v.Data + value.Data, []*Value{v, value}, "+", backward, 0.0, derivatives, ancestorDerivatives}
}

func Sub(v, value *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent1 := value.Parents[0]
		parent2 := value.Parents[1]
		if reflect.DeepEqual(parent1, parent2) {
			parent1.GradientData += 2.0 * value.GradientData
		} else {
			parent1.GradientData += value.GradientData
			parent2.GradientData += -1.0 * value.GradientData
		}
	}
	if GetConfig().PropagationForward {
		if reflect.DeepEqual(v, value) {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 0.0}
		} else {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 1.0}
			derivatives[value.UUID] = Derivative{value.UUID, value.Label, -1.0}
		}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v, value}, derivatives)
	}

	return &Value{uuid.New(), label, v.Data - value.Data, []*Value{v, value}, "-", backward, 0.0, derivatives, ancestorDerivatives}
}

func Mul(v, value *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent1 := value.Parents[0]
		parent2 := value.Parents[1]
		if reflect.DeepEqual(parent1, parent2) {
			parent1.GradientData += 2.0 * parent1.Data * value.GradientData
		} else {
			parent1.GradientData += parent2.Data * value.GradientData
			parent2.GradientData += parent1.Data * value.GradientData
		}
	}
	if GetConfig().PropagationForward {
		if reflect.DeepEqual(v, value) {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 2.0 * value.Data}
		} else {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, value.Data}
			derivatives[value.UUID] = Derivative{value.UUID, value.Label, v.Data}
		}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v, value}, derivatives)
	}

	return &Value{uuid.New(), label, v.Data * value.Data, []*Value{v, value}, "*", backward, 0.0, derivatives, ancestorDerivatives}
}

func Div(v, value *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent1 := value.Parents[0]
		parent2 := value.Parents[1]
		if reflect.DeepEqual(parent1, parent2) {
			parent1.GradientData += 0.0 * parent1.Data * value.GradientData
		} else {
			parent1.GradientData += (1.0 / parent2.Data) * value.GradientData
			parent2.GradientData += (-1.0 * parent1.Data / math.Pow(parent2.Data, 2)) * value.GradientData
		}
	}
	if GetConfig().PropagationForward {
		if reflect.DeepEqual(v, value) {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 0.0}
		} else {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, 1.0 / value.Data}
			derivatives[value.UUID] = Derivative{value.UUID, value.Label, -1 * v.Data / math.Pow(value.Data, 2)}
		}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v, value}, derivatives)
	}
	return &Value{uuid.New(), label, v.Data / value.Data, []*Value{v, value}, "/", backward, 0.0, derivatives, ancestorDerivatives}
}

func Pow(v, value *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent1 := value.Parents[0]
		parent2 := value.Parents[1]
		if reflect.DeepEqual(parent1, parent2) {
			parent1.GradientData += math.Pow(parent1.Data, 2) * (math.Log(parent1.Data) + 1) * value.GradientData
		} else {
			parent1.GradientData += parent2.Data * math.Pow(parent1.Data, parent2.Data-1) * value.GradientData
			parent2.GradientData += math.Log(parent1.Data) * math.Pow(parent1.Data, parent2.Data) * value.GradientData
		}
	}
	if GetConfig().PropagationForward {
		if reflect.DeepEqual(v, value) {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, math.Pow(v.Data, 2) * (math.Log(v.Data) + 1)}
		} else {
			derivatives[v.UUID] = Derivative{v.UUID, v.Label, value.Data * math.Pow(v.Data, value.Data-1)}
			derivatives[value.UUID] = Derivative{value.UUID, value.Label, math.Log(v.Data) * math.Pow(v.Data, value.Data)}
		}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v, value}, derivatives)
	}
	return &Value{uuid.New(), label, math.Pow(v.Data, value.Data), []*Value{v, value}, "**", backward, 0.0, derivatives, ancestorDerivatives}
}

func Exp(v *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent := value.Parents[0]
		parent.GradientData += math.Exp(parent.Data) * value.GradientData
	}
	if GetConfig().PropagationForward {
		derivatives[v.UUID] = Derivative{v.UUID, v.Label, math.Exp(v.Data)}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v}, derivatives)
	}

	return &Value{uuid.New(), label, math.Exp(v.Data), []*Value{v}, "exp", backward, 0.0, derivatives, ancestorDerivatives}
}

func Tanh(v *Value, label string) *Value {
	derivatives := make(map[uuid.UUID]Derivative)
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	backward := func(value *Value) {
		parent := value.Parents[0]
		parent.GradientData += (1 - math.Pow(math.Tanh(parent.Data), 2)) * value.GradientData
	}
	if GetConfig().PropagationForward {
		derivatives[v.UUID] = Derivative{v.UUID, v.Label, 1 - math.Pow(math.Tanh(v.Data), 2)}
		ancestorDerivatives = calculateAncestorDerivatives([]*Value{v}, derivatives)
	}

	return &Value{uuid.New(), label, math.Tanh(v.Data), []*Value{v}, "tanh", backward, 0.0, derivatives, ancestorDerivatives}
}

func (v *Value) BackPropagate() {
	v.GradientData = 1.0
	v.backwardDFS(make(map[uuid.UUID]bool))
}

// visit all parent nodes (using DFS) and execute Backward func
func (v *Value) backwardDFS(visited map[uuid.UUID]bool) {
	if visited[v.UUID] {
		return
	}
	visited[v.UUID] = true
	v.Backward(v)

	for _, parent := range v.Parents {
		parent.backwardDFS(visited)
	}
}

func calculateAncestorDerivatives(parents []*Value, derivatives map[uuid.UUID]Derivative) map[uuid.UUID]Derivative {
	ancestorDerivatives := make(map[uuid.UUID]Derivative)
	for _, parent := range parents {
		localGradient := derivatives[parent.UUID]
		if len(parent.AncestorDerivatives) > 0 {
			for _, parentAncestorDerivative := range parent.AncestorDerivatives {
				gradientValue := parentAncestorDerivative.GradientData * localGradient.GradientData
				if _, ok := ancestorDerivatives[parentAncestorDerivative.VariableValueUUID]; ok {
					gradientValue += ancestorDerivatives[parentAncestorDerivative.VariableValueUUID].GradientData
				}
				ancestorDerivatives[parentAncestorDerivative.VariableValueUUID] = Derivative{parentAncestorDerivative.VariableValueUUID, parentAncestorDerivative.VariableLabel, gradientValue}
			}
		} else {
			for _, parentDerivative := range parent.Derivatives {
				gradientValue := parentDerivative.GradientData * localGradient.GradientData
				if _, ok := ancestorDerivatives[parentDerivative.VariableValueUUID]; ok {
					gradientValue += ancestorDerivatives[parentDerivative.VariableValueUUID].GradientData
				}
				ancestorDerivatives[parentDerivative.VariableValueUUID] = Derivative{parentDerivative.VariableValueUUID, parentDerivative.VariableLabel, gradientValue}
			}
		}
	}
	return ancestorDerivatives
}
