# StreamingFast Pancake Generated
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/pancake-generated-priv)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Usage

Build with:

```bash

go build -o pancakeswap-exchange ./cmd/exchange
```

Call the `graph-node` to create the deployment:

```
./pancakeswap-exchange create namespace/target_name    # create  a row in `subgraph` table (current_version = nil, previsou_version = nil)
```

Deploy it, as you would with normal subgraphs:

```
./pancakeswap-exchange deploy namespace/target_name    # create  a row in `subgraph_deployment` &`subgraph_version` & IPS upload & `deployment_schemas` & Update `subgraph` table current_version, previous_version (MAYBE)
```

Start the linear indexer:

``` 
./pancakeswap-exchange index namespace/target_name@VERSION
```

You will find the `VERSION` printed when you `deploy` the subgraph.



## Updating `sparkle`

Generate using `sparkle` (see https://github.com/streamingfast/sparkle)

```shell
sparkle codegen ./subgraph/exchange.yaml github.com/streamingfast/sparkle-pancakeswap
go mod tidy
```

To init a database
```shell
go run ./cmd/exchange -- deploy \
   --psql-dsn="postgresql://postgres:@localhost:5432/YOUR_DATABASE?enable_incremental_sort=off&sslmode=disable" \
   project/subgraph
```

## Contributing

**Issues and PR in this repo related strictly to Pancake Generated.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

## License

[Apache 2.0](LICENSE)

