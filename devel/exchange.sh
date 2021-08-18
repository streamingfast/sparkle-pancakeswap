#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
SUBGRAPH_NAME="pancakeswap/exchange-v2"
SUBGRAPH_VERSION="pedantic-mogwai"
SUBGRAPH="${SUBGRAPH_NAME}@${SUBGRAPH_VERSION}"
PSQL_DSN="${PSQL_DSN:-"postgresql://postgres:${PG_PASSWORD}@127.0.0.1:5432/graph-node?enable_incremental_sort=off&sslmode=disable"}"
main() {
  pushd "$ROOT" &> /dev/null

  go install -v ./cmd/exchange

#   subgraph \
#    deploy \
#    --psql-dsn="${PSQL_DSN}"\
#    "$SUBGRAPH_ID" \
#    "$SCHEMA"subgraph \
#    deploy \
#    --psql-dsn="${PSQL_DSN}"\
#    "$SUBGRAPH_ID" \
#    "$SCHEMA"

  INFO=".*" exchange \
    index \
    "$SUBGRAPH"  \
    --psql-dsn="${PSQL_DSN}" \
    --enable-poi \
    --network-name="ethereum/mainnet" \
    --rpc-endpoint=http://localhost:8545 \
    --sf-api-key="${DFUSE_API_KEY}" \
    --sf-endpoint=${DFUSE_SF_ENDPOINT} \
    --start-block-num=6851000
    "$SUBGRAPH" \
    $@
  popd &> /dev/null
}

main $@
