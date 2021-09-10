#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
STOREURL=gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1
if [ "$RPCENDPOINT" == "" ]; then
	echo "you need to set RPCENDPOINT var to your 'https://.../privateendpoint/' URL"
	exit 1
fi
RPCCACHE=./localrpccache
if test -d ./localblocks; then
  echo "Using blocks from local store: ./localblocks"
    STOREURL=./localblocks
  else
    echo "Fetching blocks from remote store. You should copy them locally to make this faster..., ex:"
    cat <<EOC
######

mkdir ./localblocks
gsutil -m cp "gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1/0006809*" ./localblocks/
gsutil -m cp "gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1/0006810*" ./localblocks/
gsutil -m cp "gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1/0006811*" ./localblocks/
gsutil -m cp "gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1/0007000**" ./localblocks/

######
EOC
fi

function step1() {
    DEBUG=.* exchange parallel step -s 1 --enable-aggregate-snapshot-save --output-path ./step1-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL --rpc-cache-load-path $RPCCACHE --rpc-cache-save-path $RPCCACHE --rpc-endpoint $RPCENDPOINT &
}
function step2() {
    DEBUG=.* exchange parallel step -s 2 --enable-aggregate-snapshot-save --input-path ./step1-v1 --output-path ./step2-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL --rpc-cache-load-path $RPCCACHE --rpc-cache-save-path $RPCCACHE --rpc-endpoint $RPCENDPOINT &
}
function step3() {
    DEBUG=.* exchange parallel step -s 3 --enable-aggregate-snapshot-save --input-path ./step2-v1 --output-path ./step3-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL --rpc-cache-load-path $RPCCACHE --rpc-cache-save-path $RPCCACHE --rpc-endpoint $RPCENDPOINT &
}
function step4() {
    DEBUG=.* exchange parallel step -s 4 --enable-aggregate-snapshot-save --input-path ./step3-v1 --output-path ./step4-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL --rpc-cache-load-path $RPCCACHE --rpc-cache-save-path $RPCCACHE --rpc-endpoint $RPCENDPOINT &
}
function step5() {
    DEBUG=.* exchange parallel step -s 5 --enable-aggregate-snapshot-save --flush-entities --store-snapshot=false --input-path ./step4-v1 --output-path ./step5-v1  --start-block $1 --stop-block $2  --blocks-store-url $STOREURL  --enable-poi --rpc-cache-load-path $RPCCACHE --rpc-cache-save-path $RPCCACHE --rpc-endpoint $RPCENDPOINT &
}


main() {
  pushd "$ROOT" &> /dev/null
    go install -v ./cmd/exchange || exit 1

    if [ "$1" != "" ] && [ "$1" != 1 ]; then
      echo "SKIPPING STEP 1"
    else
      echo "LAUNCHING STEP 1"
      rm -rf ./step1-v1

      step1 7000000 7000009
      step1 7000010 7000019
      step1 7000020 7000029
      step1 7000030 7000039
      step1 7000040 7000049
      step1 7000050 7000059

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != 2 ]; then
      echo "SKIPPING STEP 2"
    else
      echo "LAUNCHING STEP 2"
      rm -rf ./step2-v1

      step2 7000000 7000009
      step2 7000010 7000019
      step2 7000020 7000029
      step2 7000030 7000039
      step2 7000040 7000049
      step2 7000050 7000059

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != 3 ]; then
      echo "SKIPPING STEP 3"
    else
      echo "LAUNCHING STEP 3"
      rm -rf ./step3-v1

#      step3 7000000 7000009
#      step3 7000010 7000019
#      step3 7000020 7000029
#      step3 7000030 7000039
      step3 7000040 7000049
#      step3 7000050 7000059

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != 4 ]; then
      echo "SKIPPING STEP 4"
    else
      echo "LAUNCHING STEP 4"
      rm -rf ./step4-v1

      step4 7000000 7000009
      step4 7000010 7000019
      step4 7000020 7000029
      step4 7000030 7000039
      step4 7000040 7000049
      step4 7000050 7000059

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != 5 ]; then
      echo "SKIPPING STEP 5"
    else
      echo "LAUNCHING STEP 5"
      rm -rf ./step5-v1

      step5 7000000 7000009
      step5 7000010 7000019
      step5 7000020 7000029
      step5 7000030 7000039
      step5 7000040 7000049
      step5 7000050 7000059

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != csv ]; then
      echo "SKIPPING STEP CSV"
    else
      echo "Exporting to csv"
      DEBUG=.* exchange parallel to-csv --only-tables "poi2\$" --input-path ./step4-v1 --output-path ./stepcsvs --chunk-size 1000
    fi
  popd &> /dev/null
}

main $@
