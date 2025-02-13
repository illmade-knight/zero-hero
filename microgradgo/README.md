- [Intro](#intro)
- [Usage](#usage)

## Intro
Following https://github.com/karpathy/nn-zero-to-hero and "implementing" it in Go as I go along


## Usage
Not much yet
Simple main func to generate a dot graph. To render you need to install graphiz

```
make build
go run cmd/main.go -outputfile graph.dot
dot -Tpng graph.dot > graph.png
open graph.png
```