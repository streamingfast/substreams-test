package main

import (
	"github.com/streamingfast/logging"
)

var zlog, tracer = logging.RootLogger("substreams-test", "github.com/streamingfast/substreams-test/cmd/substreams-test")

func init() {
	logging.InstantiateLoggers(logging.WithSwitcherServerAutoStart())
}
