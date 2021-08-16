#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
main() {
    if uname |grep -q Darwin; then
        echo "Transpiling from OSX"
    elif ldd --version |grep -q 2.31; then
	echo "Compiling locally on linux: glibc 2.31 found"
    else
        echo "Cannot compile locally.. wrong version of glibc (needs ubuntu 20.04 for glibc2.31) Use workflow.. sorry..."
	exit 1
    fi

    cd "$ROOT"
    go test -v ./... || exit 1

    echo "Building binary (GOOS=linux GOARCH=amd64)"
    GOOS=linux GOARCH=amd64 go build -o ./deploy/exchange ./cmd/exchange

    TAG=gcr.io/eoscanada-shared-services/pancake-exchange
    REV=$(git describe --long --always --dirty)$1
    IMG=$TAG:$REV

    cd deploy
    docker build . -t $IMG
    docker push $IMG

    echo
    echo Image pushed: $IMG
    echo
}

main $@
