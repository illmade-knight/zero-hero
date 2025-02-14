package main

import (
	"flag"
	"fmt"
	"os"

	"microgradgo/internal/engine"
)

func main() {
	// verbose := flag.Bool("verbose", false, "Enable verbose output")
	outputfile := flag.String("outputfile", "graph.dot", "Path to store the graphviz dot file")

	flag.Parse()

	// making some Values, combining them, generating a Graphiz dot graph
	a := engine.NewValue(3.0, "a")
	b := engine.NewValue(-11.0, "b")
	c := engine.NewValue(4.0, "c")
	fmt.Printf("a %v\n", a)
	fmt.Printf("b %v\n", b)
	fmt.Printf("c %v\n", c)
	d := a.Add(b, "d")
	e := a.Div(c, "e")
	f := b.Pow(c, "f")
	g := d.Mul(e, "g")
	h := e.Sub(f, "h")
	L := g.Add(h, "L")
	fmt.Printf("L %v\n", L)

	dotGraph, err := engine.DotGraph(L)
	if err != nil {
		fmt.Println("Failed to generate graph:", err)
	}

	err = engine.RecurseValuesGradient(&L)
	if err != nil {
		fmt.Println("Failed to recurse gradients:", err)
	}

	dotGraph, err = engine.DotGraph(L)
	if err != nil {
		fmt.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*outputfile, []byte(dotGraph), 0644)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
	}

}
