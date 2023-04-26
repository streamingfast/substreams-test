#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

stest="$ROOT/../stest"

finish() {
    kill -s TERM $active_pid &> /dev/null || true
}

main() {
  trap "finish" EXIT
  pushd "$ROOT" &> /dev/null

    shift $(($OPTIND - 1))
    version=$1
    if ! echo "$version" |grep -qE "(^main|^prod-minimal$)"; then
      error "Invalid version: $version, expected minimal or prod-minimal"
    fi

    configJsonnet="main.jsonnet"
    graphUrl="https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"
    if [  "$version" == "prod-minimal" ]; then
      configJsonnet="prod-minimal.jsonnet"
      graphUrl="https://api.thegraph.com/subgraphs/name/ianlapham/v3-minimal"
    fi

    jsonnet "$configJsonnet" > run_config.json
    $stest test substream \
      ../../../substreams-uniswap-v3/substreams.yaml \
      "$graphUrl" \
      ./run_config.json \
      12369621:12379621 \
      "$@"

  popd &> /dev/null
}

error() {
  message="$1"
  exit_code="$2"
  printf "${RED} * $message * ${NC}\n"
  exit ${exit_code:-1}
}


main "$@"
