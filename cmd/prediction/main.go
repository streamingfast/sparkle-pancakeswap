package main

import (
	"github.com/streamingfast/sparkle/cli"
	"github.com/streamingfast/sparkle/subgraph"
	"github.com/streamingfast/sparkle-pancakeswap/prediction"
)

func main() {
	subgraph.MainSubgraphDef = prediction.Definition
	cli.Execute()
}
