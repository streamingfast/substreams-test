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
    jsonnet minimal.jsonnet > run_config.json
    $stest test substream \
      ../../../substreams-uniswap-v3/substreams.yaml \
      https://api.thegraph.com/subgraphs/name/ianlapham/v3-minimal \
      ./run_config.json \
      12369621:12469621\
      "$@"

  popd &> /dev/null
}

main "$@"
