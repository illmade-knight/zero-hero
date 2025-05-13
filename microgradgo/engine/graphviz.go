package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// generates graph from Value recursively examining parent nodes and returns Dot graph rendering of this graph
func DotGraph(value *Value) (string, error) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return "", err
	}
	graph, err := g.Graph()
	graph.SetRankDir(cgraph.LRRank)
	graph.SetConcentrate(true)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	nodeLabelStr, err := createValueString(value)
	if err != nil {
		return "", err
	}

	leaf, err := graph.CreateNodeByName(value.Label)
	leaf.SetShape("record")
	leaf.SetLabel(nodeLabelStr)
	if err != nil {
		return "", err
	}

	err = recurseValuesGraphiz(graph, value, leaf)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func recurseValuesGraphiz(graph *graphviz.Graph, value *Value, node *cgraph.Node) error {
	if (len(value.Parents) > 0) && value.Op != "" {
		// "merge" edges by creating tiny opnode and edges from childNodes to opNode and from opNode to parentNode
		// bit of an effort to "hide" the opNode. Will produce graphviz dot warnings
		opNode, err := graph.CreateNodeByName("op" + value.UUID.String())
		if err != nil {
			return err
		}
		opNode.SetHeight(0)
		opNode.SetWidth(0)
		opNode.SetFixedSize(true)
		opNode.SetStyle(cgraph.NodeStyle("invis"))

		opEdge, err := graph.CreateEdgeByName(value.Op, opNode, node)
		if err != nil {
			return err
		}
		opEdge.SetLabel(value.Op)

		for _, parentValue := range value.Parents {
			nodeLabelStr, err := createValueString(parentValue)
			if err != nil {
				return err
			}

			parentNode, err := graph.CreateNodeByName(parentValue.UUID.String())
			parentNode.SetShape("record")
			parentNode.SetLabel(nodeLabelStr)
			if err != nil {
				return err
			}

			e, err := graph.CreateEdgeByName("", parentNode, opNode)
			if err != nil {
				return err
			}
			e.SetLabel("")
			e.SetDir(cgraph.NoneDir)
			recurseValuesGraphiz(graph, parentValue, parentNode)
		}
	}
	return nil
}

func createValueString(value *Value) (string, error) {
	var nodeLabelStr string
	nodeLabelStr = fmt.Sprintf("%s | %.4f | %.4f", value.Label, value.Data, value.GradientData)
	switch len(value.Derivatives) {
	case 0:

	case 1:
		keys := slices.Collect(maps.Keys(value.Derivatives))
		nodeLabelStr = fmt.Sprintf("%s | d%s/d%s  %.4f", nodeLabelStr,
			value.Label, value.Derivatives[keys[0]].VariableLabel, value.Derivatives[keys[0]].GradientData)
	case 2:
		keys := slices.Collect(maps.Keys(value.Derivatives))
		nodeLabelStr = fmt.Sprintf("%s | d%s/d%s  %.4f | d%s/d%s  %.4f", nodeLabelStr,
			value.Label, value.Derivatives[keys[0]].VariableLabel, value.Derivatives[keys[0]].GradientData,
			value.Label, value.Derivatives[keys[1]].VariableLabel, value.Derivatives[keys[1]].GradientData)
	default:
		return "", errors.New(fmt.Sprintf("The node %v does not have 0,1 or 2 parents (%v)", value, len(value.Derivatives)))
	}

	var leafGradientsStr string
	if len(value.AncestorDerivatives) > 0 {
		// sorting based on leafGradient.Label
		leafGradients := make(map[string]float64)
		keys := make([]string, 0)
		for _, leafGradient := range value.AncestorDerivatives {
			keys = append(keys, leafGradient.VariableLabel)
			leafGradients[leafGradient.VariableLabel] = leafGradient.GradientData
		}
		slices.Sort(keys)

		for _, leafGradientLabel := range keys {
			leafGradientsStr = leafGradientsStr + fmt.Sprintf("{d%s/d%s  %.4f}", value.Label, leafGradientLabel, leafGradients[leafGradientLabel])
			leafGradientsStr = leafGradientsStr + "|"
		}
		leafGradientsStr = leafGradientsStr[:len(leafGradientsStr)-1] // remove last separator
		leafGradientsStr = "|{" + leafGradientsStr + "}"
	}
	return "{" + nodeLabelStr + leafGradientsStr + "}", nil
}
