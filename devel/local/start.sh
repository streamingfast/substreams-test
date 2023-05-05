#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

stest="$ROOT/../stest"
clean=

finish() {
    kill -s TERM $active_pid &> /dev/null || true
}

main() {
  trap "finish" EXIT
  pushd "$ROOT" &> /dev/null

    while getopts "hc" opt; do
      case $opt in
        h) usage && exit 0;;
        c) clean=true;;
        \?) usage_error "Invalid option: -$OPTARG";;
      esac
    done
    shift $((OPTIND-1))
    [[ $1 = "--" ]] && shift

    set -e

    if [[ $clean == "true" ]]; then
      rm -rf localdata &> /dev/null || true
    fi

    version=$1
    if ! echo "$version" |grep -qE "(^main|^prod-minimal$)"; then
      error "Invalid version: $version, expected minimal or prod-minimal"
    fi

    configJsonnet="main.jsonnet"
    graphUrl="https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"
    if [ "$version" == "prod-minimal" ]; then
      configJsonnet="prod-minimal.jsonnet"
      graphUrl="https://api.thegraph.com/subgraphs/name/ianlapham/v3-minimal"
    fi

    substreamsEndpoint=$2
    if [ -z "$substreamsEndpoint" ]; then
      substreamsEndpoint="mainnet.eth.streamingfast.io:443"
    fi

    echo "Using substreams endpoint: $substreamsEndpoint"

    jsonnet "$configJsonnet" > run_config.json
    $stest test substream \
      ../../../substreams-uniswap-v3/substreams.spkg \
      "$graphUrl" \
      ./run_config.json \
      12369621:12375000 \
      --endpoint $substreamsEndpoint \
      "$@"

  popd &> /dev/null
}

error() {
  message="$1"
  exit_code="$2"
  printf "${RED} * $message * ${NC}\n"
  exit ${exit_code:-1}
}

usage() {
  echo "usage: start.sh [-c]"
  echo ""
  echo "Start substreams test."
  echo ""
  echo "Options"
  echo "    -c             Clean actual graphql cache"
}


main "$@"
