package engine

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

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

	leaf, err := graph.CreateNodeByName(value.Label + " | " + strconv.FormatFloat(value.Data, 'f', -1, 64))
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

func recurseValuesGraph(graph *graphviz.Graph, childValue Value, childNode *cgraph.Node) error {
	if (len(childValue.prev) > 0) && childValue.Op != "" {
		// "merge" edges by creating tiny opnode and edges from childNodes to opNode and from opNode to parentNode
		// bit of an effort to "hide" the opNode. Will produce graphviz dot warnings
		opNode, err := graph.CreateNodeByName("op" + childValue.Label)
		if err != nil {
			return err
		}
		opNode.SetHeight(0)
		opNode.SetWidth(0)
		opNode.SetFixedSize(true)
		opNode.SetStyle(cgraph.NodeStyle("invis"))

		opEdge, err := graph.CreateEdgeByName(childValue.Op, opNode, childNode)
		if err != nil {
			return err
		}
		opEdge.SetLabel(childValue.Op)
		
		for _, parentValue := range childValue.prev {
			parentNode, err := graph.CreateNodeByName(fmt.Sprintf("%s | %.4f", parentValue.Label, parentValue.Data))
			if err != nil {
				return err
			}
	
			e, err := graph.CreateEdgeByName("", parentNode, opNode)
			if err != nil {
				return err
			}
			e.SetLabel("")
			e.SetDir(cgraph.NoneDir)
			recurseValuesGraph(graph, parentValue, parentNode)
		}
	}
	return nil
}