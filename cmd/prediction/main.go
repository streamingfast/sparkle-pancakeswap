package main

import (
	"github.com/streamingfast/sparkle/cli"
	"github.com/streamingfast/sparkle/subgraph"
	"github.com/streamingfast/pancake-generated-priv/prediction"
)

func main() {
	subgraph.MainSubgraphDef = prediction.Definition
	cli.Execute()
}
