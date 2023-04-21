#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

stest="$ROOT/../stest"

finish() {
    kill -s TERM $active_pid &> /dev/null || true
}

main() {
  trap "finish" EXIT
  pushd "$ROOT" &> /dev/null

    set -e

    export INFO=".*"
    $stest test substream \
      ../../../substreams-uniswap-v3/substreams.yaml \
      https://api.thegraph.com/subgraphs/name/ianlapham/v3-minimal \
      ./config.json \
      12369621:12469621\
      "$@"

  popd &> /dev/null
}

main "$@"
