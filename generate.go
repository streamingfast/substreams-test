package subcmp

import (
	"context"
	"fmt"
	"github.com/streamingfast/substreams-test/thegraph"
	"go.uber.org/zap"
)

type Generate struct {
	StartBlockNum uint64
	StopBlockNum  uint64
	Config        *Config
	Logger        *zap.Logger
}

func (g *Generate) Run(ctx context.Context) error {

	graphClient := thegraph.New(g.Config.GraphURL, thegraph.WithLogger(g.Logger))

	var generators []*testGenerator
	for graphqlQueryPath, config := range g.Config.Tests {
		gen, err := createTestGenerator(ctx, graphqlQueryPath, config, graphClient, g.Logger)
		if err != nil {
			return fmt.Errorf("failed setting up test %q: %w", graphqlQueryPath, err)
		}
		generators = append(generators, gen)
	}

	g.Logger.Info("setup test generators", zap.Int("count", len(generators)))

	g.Logger.Info("setting test writer")
	testWriter, err := newWriter(g.Config.TestOutputPath)
	if err != nil {
		return fmt.Errorf("failed to open test writer: %w", err)
	}

	go func() {
		if err := testWriter.Write(); err != nil {
			g.Logger.Info("writer completed", zap.Error(err))
		}
	}()

	for _, gen := range generators {
		go gen.run(ctx, g.StartBlockNum, g.StopBlockNum, testWriter.in)
	}

	for _, gen := range generators {
		gen.wait()
	}

	testWriter.closeTest()
	g.Logger.Info("waiting for writer to complete")
	testWriter.wait()
	return nil

}
