format:
	@go fmt ./...
.PHONY: format

build:
	@go build ./...
.PHONY : build

test:
	@go test ./...

testv:
	@go test -v -cover -count=1 ./...

DOT_GRAPH = graphs/graph.dot
DOT_GRAPH_IMAGE = graphs/graph.png

neuralnet:
	@go run cmd/main.go -action neuralnet -graphfile $(DOT_GRAPH) -propagationforward; dot -q -Tpng $(DOT_GRAPH) > $(DOT_GRAPH_IMAGE); open $(DOT_GRAPH_IMAGE)

mlp:
	@go run cmd/main.go -action mlp -graphfile $(DOT_GRAPH); dot -q -Tpng $(DOT_GRAPH) > $(DOT_GRAPH_IMAGE); open $(DOT_GRAPH_IMAGE)

mlptrain:
	@go run cmd/main.go -action mlptrain -graphfile $(DOT_GRAPH); dot -q -Tpng $(DOT_GRAPH) > $(DOT_GRAPH_IMAGE); open $(DOT_GRAPH_IMAGE)


TRAINING_FILE = mnist_dataset/mnist_train.csv
TEST_FILE = mnist_dataset/mnist_test.csv

mnistdata:
	@if [ ! -f $(TRAINING_FILE) ]; then curl -o $(TRAINING_FILE)  https://pjreddie.com/media/files/mnist_train.csv; fi
	@if [ ! -f mnist_dataset/mnist_test.csv ]; then curl -o mnist_dataset/mnist_test.csv  https://pjreddie.com/media/files/mnist_test.csv; fi

mnisttrain: mnistdata
	@go run cmd/main.go -action mnisttrain -trainingfile $(TRAINING_FILE) -modelfile mnist_dataset/model.json -verbose	

mnisttest:
	@go run cmd/main.go -action mnisttest -testfile mnist_dataset/mnist_test.csv -modelfile mnist_dataset/model-pretrained.json -verbose	

.PHONY: test-all
