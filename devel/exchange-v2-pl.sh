#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
STOREURL=gs://dfuseio-global-blocks-us/eth-bsc-mainnet/v1
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

######
EOC
fi

function step1() {
    INFO=.* exchange parallel step -s 1 --output-path ./step1-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL &
}
function step2() {
    INFO=.* exchange parallel step -s 2 --input-path ./step1-v1 --output-path ./step2-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL &
}
function step3() {
    INFO=.* exchange parallel step -s 3 --input-path ./step2-v1 --output-path ./step3-v1 --start-block $1 --stop-block $2 --blocks-store-url $STOREURL &
}
function step4() {
    INFO=.* exchange parallel step -s 4 --flush-entities --store-snapshot=false --input-path ./step3-v1 --output-path ./step4-v1  --start-block $1 --stop-block $2  --blocks-store-url $STOREURL  &
}


main() {
  pushd "$ROOT" &> /dev/null
    go install -v ./cmd/exchange || exit 1

    if [ "$1" != "" ] && [ "$1" != 1 ]; then
      echo "SKIPPING STEP 1"
    else
      echo "LAUNCHING STEP 1"
      rm -rf ./step1-v1

      step1 6809700 6829699
      step1 6829700 6849699
      step1 6849700 6889699

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

      step2 6809700 6829699
      step2 6829700 6849699
      step2 6849700 6889699

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

      step3 6809700 6829699
      step3 6829700 6849699
      step3 6849700 6889699

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

      step4 6809737 6810736
      step4 6810737 6811736
#      step4  6811737 6812737

      for job in `jobs -p`; do
          echo "Waiting on $job"
          wait $job
      done
    fi

    if [ "$1" != "" ] && [ "$1" != csv ]; then
      echo "SKIPPING STEP CSV"
    else
      echo "Exporting to csv"
      INFO=.* exchange parallel to-csv --input-path ./step4-v1 --output-path ./stepcsvs --chunk-size 1000
    fi
  popd &> /dev/null
}

main $@
