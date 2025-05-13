package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	// "strconv"

	"microgradgo/engine"
)

// making some Values, combining them, generating a Graphiz dot graph
func basic() {
	a := engine.NewValue(3.0, "a")
	b := engine.NewValue(-2.0, "b")
	c := engine.NewValue(4.0, "c")
	d := engine.NewValue(23.12, "d")
	e := engine.NewValue(-423.56, "e")
	f := engine.NewValue(36.9, "f")
	g := engine.NewValue(23232.3467, "g")
	h := engine.NewValue(-727.82, "h")
	i := engine.Pow(a, b, "i")
	j := engine.Sub(c, d, "j")
	k := engine.Mul(e, f, "k")
	l := engine.Div(g, h, "l")
	m := engine.Mul(i, j, "m")
	n := engine.Div(k, l, "n")
	O := engine.Add(m, n, "o")

	O.BackPropagate()
	dotGraph, err := engine.DotGraph(O)
	if err != nil {
		log.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
	if err != nil {
		log.Println("Failed to write to file:", err)
	}
}

// repeatedly using same Value
func basicWithRepeatedValues() {
	a := engine.NewValue(-2.0, "a")
	b := engine.NewValue(3.0, "b")
	d := engine.Mul(a, b, "d")
	e := engine.Add(a, b, "e")
	f := engine.Mul(d, e, "f")

	f.BackPropagate()
	dotGraph, err := engine.DotGraph(f)
	if err != nil {
		log.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
	if err != nil {
		log.Println("Failed to write to file:", err)
	}
}

// simple neural net
func neuralNet() {
	// input values
	x1 := engine.NewValue(2.0, "x1")
	x2 := engine.NewValue(0.0, "x2")
	// weights
	w1 := engine.NewValue(-3.0, "w1")
	w2 := engine.NewValue(1.0, "w2")
	// bias of neuron
	b := engine.NewValue(6.8813735870195432, "b")

	x1w1 := engine.Mul(x1, w1, "x1*w1")
	x2w2 := engine.Mul(x2, w2, "x2*w2")
	x1w1_x2w2 := engine.Add(x1w1, x2w2, "x1*w1 + x2*w2")
	n := engine.Add(x1w1_x2w2, b, "n")

	// activation function
	o := engine.Tanh(n, "o")

	o.BackPropagate()
	dotGraph, err := engine.DotGraph(o)
	if err != nil {
		log.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
	if err != nil {
		log.Println("Failed to write to file:", err)
	}
}

// many layers
func mlp() {
	inputs := engine.NewInputs([]float64{1.0, -2.0, 3.5})
	mlp := engine.NewMLP(3, 1, []int{10, 11, 12, 11, 10})
	outputs, err := engine.ForwardMLP(mlp, inputs)
	output := outputs[0]

	output.BackPropagate()

	dotGraph, err := engine.DotGraph(output)
	if err != nil {
		log.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
	if err != nil {
		log.Println("Failed to write to file:", err)
	}
}

func mlpTrain() {
	// input, target outputs
	multipleInputs := [][]*engine.Value{
		engine.NewInputs([]float64{1.0, 3.0, -1.0}),
		engine.NewInputs([]float64{3.0, -1.0, 0.5}),
		engine.NewInputs([]float64{0.5, 1.0, 1.0}),
		engine.NewInputs([]float64{1.0, 1.0, -1.0}),
	}
	targetOutputs := [][]*engine.Value{
		engine.NewInputs([]float64{1.0}),
		engine.NewInputs([]float64{-1.0}),
		engine.NewInputs([]float64{-1.0}),
		engine.NewInputs([]float64{1.0})}

	// set up network
	mlp := engine.NewMLP(3, 1, []int{4, 4})

	var loss *engine.Value
	var outputs [][]*engine.Value
	var err error
	for range 100 {
		loss, outputs, err = engine.TrainMultiple(mlp, multipleInputs, targetOutputs, -0.05)
		if err != nil {
			log.Println("Failed to train:", err)
		}
		log.Println("loss is ", loss.Data)
	}

	// summary
	log.Println("outputs are")
	for _, subrange := range outputs {
		for _, value := range subrange {
			fmt.Printf("%v", value.Data)
		}
		fmt.Printf("\n")
	}
	log.Println("loss is ", loss.Data)

	// todo training run for different epochs (* 100) and plot epochs vs loss avg/stn dev
	// graph
	dotGraph, err := engine.DotGraph(loss)
	if err != nil {
		log.Println("Failed to generate graph:", err)
	}

	err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
	if err != nil {
		log.Println("Failed to write to file:", err)
	}
}

// will load
func mnistTrain() {

	// network
	// 784 inputs - 28 x 28 pixels, each pixel is an input
	// 100 hidden nodes - arbitrary
	// 10 outputs - digits 0 to 9
	var mlp engine.MLP
	modelParams, err := os.ReadFile(*modelFile)
	if err != nil {
		log.Println("No model file exists", *modelFile)
		mlp = engine.NewMLP(784, 10, []int{100})
		log.Println("New model created")
	} else {
		mlp, err = engine.FromJSON(modelParams)
		if err != nil {
			log.Fatal("Unable to unmarshal from json ", err)
		}
		log.Println("Loaded model", *modelFile)
	}

	// load random batch from 6000 line training file, iterate for 6000/batchSize iterations, then iterate for epochs
	// save model every batch
	batchSize := 64
	numEpochs := 5
	learningRate := 0.025
	var loss *engine.Value
	var outputs [][]*engine.Value
	for epoch := 0; epoch < numEpochs; epoch++ {
		var batchInputs [][]*engine.Value
		var batchTargets [][]*engine.Value

		content, err := os.ReadFile(*trainingFile)
		if err != nil {
			log.Fatal("Unable to read input file "+*trainingFile, err)
		}
		log.Println("Loaded training file", *trainingFile)
		lines := strings.Split(string(content), "\n")

		for batchIdx := range len(lines) / batchSize {
			rand.Shuffle(len(lines), func(i, j int) {
				lines[i], lines[j] = lines[j], lines[i]
			})
			batchLines := lines[:batchSize]

			for _, line := range batchLines {
				reader := csv.NewReader(strings.NewReader(line))

				record, err := reader.Read()
				if err == io.EOF {
					break
				}

				// TODO - why in gonn is there a slight offset (multiply 0.999 and add 0.001)
				inputs := make([]float64, 784)
				for i := range inputs {
					x, _ := strconv.ParseFloat(record[i], 64)
					// inputs[i] = (x / 255.0 * 0.999) + 0.001
					inputs[i] = x / 255.0
				}
				targets := make([]float64, 10)
				for i := range targets {
					// targets[i] = 0.001
					targets[i] = 0.0
				}
				x, _ := strconv.Atoi(record[0])
				// targets[x] = 0.999
				targets[x] = 1.0

				batchInputs = append(batchInputs, engine.NewInputs(inputs))
				batchTargets = append(batchTargets, engine.NewInputs(targets))
			}

			// TODO can we use a goroutine to run this in parallel
			loss, outputs, err = engine.TrainMultiple(mlp, batchInputs, batchTargets, -1.0*learningRate)
			if err != nil {
				log.Println("Failed to train:", err)
			}
			log.Println("loss after batch is", loss.Data)
			batchInputs = nil
			batchTargets = nil

			log.Println("Saving model file", *modelFile)
			modelBytes, err := engine.ToJSON(mlp)
			if err != nil {
				log.Fatal("Unable to marshal model to json ", err)
			}
			err = os.WriteFile(*modelFile, modelBytes, fs.ModePerm)
			if err != nil {
				log.Fatal("Unable to write model parameters file "+*modelFile, err)
			}
			log.Println("completed epoch", epoch, "batch", batchIdx)
		}
		log.Println("completed epoch", epoch)
	}

	// summary
	log.Println("outputs are")
	for _, subrange := range outputs {
		for _, value := range subrange {
			fmt.Printf(", %v", value.Data)
		}
		fmt.Printf("\n")
	}

	log.Println("\nloss is ", loss.Data)
}

func mnistTest() {

	modelParams, err := os.ReadFile(*modelFile)
	if err != nil {
		log.Fatal("Unable to read model parameters file "+*modelFile, err)
	}
	mlp, err := engine.FromJSON(modelParams)
	if err != nil {
		log.Fatal("Unable to unmarshal from json ", err)
	}
	log.Println("Loaded model", *modelFile)

	content, err := os.ReadFile(*testFile)
	if err != nil {
		log.Fatal("Unable to read input file "+*testFile, err)
	}
	log.Println("Loaded test file", *testFile)
	lines := strings.Split(string(content), "\n")

	var outputs []*engine.Value
	score := 0
	for _, line := range lines {
		reader := csv.NewReader(strings.NewReader(line))

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		inputs := make([]float64, 784)
		for i := range inputs {
			x, _ := strconv.ParseFloat(record[i], 64)
			// inputs[i] = (x / 255.0 * 0.999) + 0.001
			inputs[i] = x / 255.0
		}
		target, _ := strconv.Atoi(record[0])

		outputs, err = engine.ForwardMLP(mlp, engine.NewInputs(inputs))
		if err != nil {
			log.Println("Failed to infer:", err)
		}

		// summary
		max := 0.0
		predict := 0
		log.Println("outputs are")
		for idx, value := range outputs {
			if value.Data >= max {
				max = value.Data
				predict = idx
			}
		}
		if target == predict {
			score++
			log.Println("Predicted correctly:", max)
		} else {
			log.Println("Predicted incorrectly:", max)
		}

		// graph - TODO - graph takes too long to be generated 
		// dotGraph, err := engine.DotGraph(outputs[0])
		// if err != nil {
		// 	log.Println("Failed to generate graph:", err)
		// }

		// err = os.WriteFile(*graphFile, []byte(dotGraph), 0644)
		// if err != nil {
		// 	log.Println("Failed to write to file:", err)
		// }

	}
	log.Println("Predicted ", score, " out of ", len(lines))
}

func debug() {

}

func parseLogs() {
	content, err := os.ReadFile(*testFile)
	if err != nil {
		log.Fatal("Unable to read input file "+*testFile, err)
	}
	log.Println("Loaded test file", *testFile)
	lines := strings.Split(string(content), "\n")

	var numIters, timeIters int64
	startStr := " Starting training with batch of inputs"
	endStr := " Completed forward pass"
	var startTime, endTime time.Time
	startLoop := false
	endLoop := false

	layout := "15:04:05"
	for _, line := range lines {
		reader := csv.NewReader(strings.NewReader(line))

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !startLoop && !endLoop && strings.HasPrefix(record[1], startStr) {
			startLoop = true
			startTime, err = time.Parse(layout, record[0])

			if err != nil {
				log.Println(err)
			}
			// log.Println("Start record", record, "time", startTime)
		}
		if startLoop && strings.HasPrefix(record[1], endStr) {
			endLoop = true
			endTime, err = time.Parse(layout, record[0])
			if err != nil {
				log.Println(err)
			}
			// log.Println("End record", record, "time", endTime)
		}

		if startLoop && endLoop {
			numIters++
			diff := endTime.Sub(startTime).Milliseconds()
			timeIters += diff
			startLoop = false
			endLoop = false
		}
	}
	log.Println("Iters ", numIters, " total time ", timeIters, " average ", timeIters/numIters)
}

var modelFile, graphFile, testFile, trainingFile *string

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	propagationForward := flag.Bool("propagationforward", false, "Enable propagation forward")
	action := flag.String("action", "debug", "Which action to execute")
	graphFile = flag.String("graphfile", "graph.dot", "Path to store the graphviz dot file")
	trainingFile = flag.String("trainingfile", "mnist_dataset/mnist_train_100.csv", "Path to mnist training set")
	testFile = flag.String("testfile", "mnist_dataset/mnist_test_100.csv", "Path to mnist test set")
	modelFile = flag.String("modelfile", "model.json", "Path to model file")

	flag.Parse()

	// logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	// logger.

	engine.Init(*propagationForward, *verbose)
	switch *action {
	case "basic":
		basic()
	case "debug":
		debug()
	case "basicrepeat":
		basicWithRepeatedValues()
	case "neuralnet":
		neuralNet()
	case "mlp":
		mlp()
	case "mlptrain":
		mlpTrain()
	case "mnisttrain":
		mnistTrain()
	case "mnisttest":
		mnistTest()
	case "logs":
		parseLogs()
	default:
		log.Println("Unknown action:", *action)
	}
}
