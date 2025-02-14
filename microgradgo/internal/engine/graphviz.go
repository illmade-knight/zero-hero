package engine

import (
	"bytes"
	"context"
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// generates graph from Value recursively examining parent nodes and returns Dot graph rendering of this graph
func DotGraph(value Value) (string, error) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		return "", err
	}
	graph, err := g.Graph()
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

	leaf, err := graph.CreateNodeByName(fmt.Sprintf("%s | %.4f | d%s/d%s %.4f, d%s/d%s %.4f",
		value.Label, value.Data,
		value.Label, value.Gradient.ParentA.Label, value.Gradient.ParentA.Gradient,
		value.Label, value.Gradient.ParentB.Label, value.Gradient.ParentB.Gradient))
	if err != nil {
		return "", err
	}
	err = recurseValuesGraph(graph, value, leaf)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func recurseValuesGraph(graph *graphviz.Graph, value Value, node *cgraph.Node) error {
	if (len(value.parents) > 0) && value.Op != "" {
		// "merge" edges by creating tiny opnode and edges from childNodes to opNode and from opNode to parentNode
		// bit of an effort to "hide" the opNode. Will produce graphviz dot warnings
		opNode, err := graph.CreateNodeByName("op" + value.Label)
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

		for _, parentValue := range value.parents {
			parentNode, err := graph.CreateNodeByName(fmt.Sprintf("%s | %.4f | d%s/d%s %.4f, d%s/d%s %.4f",
				parentValue.Label, parentValue.Data,
				parentValue.Label, parentValue.Gradient.ParentA.Label, parentValue.Gradient.ParentA.Gradient,
				parentValue.Label, parentValue.Gradient.ParentB.Label, parentValue.Gradient.ParentB.Gradient))
			if err != nil {
				return err
			}

			e, err := graph.CreateEdgeByName("", parentNode, opNode)
			if err != nil {
				return err
			}
			e.SetLabel("")
			e.SetDir(cgraph.NoneDir)
			recurseValuesGraph(graph, *parentValue, parentNode)
		}
	}
	return nil
}
