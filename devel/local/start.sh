#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

stest="$ROOT/../substreams-test"

finish() {
    kill -s TERM $active_pid &> /dev/null || true
}

main() {
  trap "finish" EXIT
  pushd "$ROOT" &> /dev/null

    set -e

    export INFO=".*"
    $stest generate "config.json"\
      16371050 \
      5 \
      "$@"

  popd &> /dev/null
}

main "$@"
