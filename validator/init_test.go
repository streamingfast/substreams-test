package validator

import (
	"github.com/streamingfast/logging"
)

var zlog, tracer = logging.RootLogger("substreams-test", "github.com/streamingfast/substreams-test/test")

func init() {
	//logging.WithDefaultLevel(zapcore.DebugLevel)
	logging.InstantiateLoggers()
}
