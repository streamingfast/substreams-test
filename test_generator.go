package subcmp

import (
	"context"
	"fmt"
	"github.com/streamingfast/dstore"
	"github.com/streamingfast/substreams-test/thegraph"
	subtest "github.com/streamingfast/substreams/tools/test"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

const BLOCK_NUM_VAR = "block"
const CONCURRENT_GRAPHQL_WORKER = 5

type testGenerator struct {
	config *TestConfig

	graphClient  *thegraph.Graph
	graphqlQuery string

	queries chan *thegraph.Query
	done    chan bool

	logger *zap.Logger
}

func createTestGenerator(ctx context.Context, queryPath string, config *TestConfig, graphClient *thegraph.Graph, logger *zap.Logger) (*testGenerator, error) {
	cnt, err := dstore.ReadObject(ctx, queryPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read graphql graphqlQuery: %w", err)
	}
	graphqlQuery := string(cnt)

	graphqlQueryVars, err := extractGraphqlVariables(graphqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to extract graphql vars from graphqlQuery: %w", err)
	}

	seenBlock := false
	requiredVars := map[string]bool{}
	for _, graphqlVar := range graphqlQueryVars {
		if graphqlVar == BLOCK_NUM_VAR {
			seenBlock = true
			continue
		}
		requiredVars[graphqlVar] = true
	}
	if !seenBlock {
		return nil, fmt.Errorf("graphql graphqlQuery needs to define a '$%s' varaible", BLOCK_NUM_VAR)
	}

	for idx, varMap := range config.Vars {
		for varName, _ := range requiredVars {
			if _, found := varMap[varName]; !found {
				return nil, fmt.Errorf("vars.%d > missing variable %q required by graphl query", idx, varName)
			}
		}

	}

	for idx, path := range config.Paths {
		subVariables := extractSubstreamsPathVariables(path.Substreams)
		if len(subVariables) != len(requiredVars) {
			return nil, fmt.Errorf("substreams path.%d %q > variable count mismatch substreams path contains %d variables while graphql contains  %d (w/out counting the block variable)", idx, path.Substreams, len(subVariables), len(requiredVars))
		}
		for _, subVariable := range subVariables {
			if _, found := requiredVars[subVariable]; !found {
				return nil, fmt.Errorf("substreams path.%d %q > defines variable %q not in graphql", idx, path.Substreams, subVariable)
			}
		}
	}

	return &testGenerator{
		config: config,

		graphClient:  graphClient,
		graphqlQuery: graphqlQuery,

		queries: make(chan *thegraph.Query),
		done:    make(chan bool),

		logger: logger.Named(queryPath),
	}, nil
}

func (t *testGenerator) run(ctx context.Context, startBlockNum, stopBlockNum uint64, writerChan chan *subtest.TestConfig) {
	t.logger.Info("launching test generator")
	var workers []*worker
	for i := 0; i < CONCURRENT_GRAPHQL_WORKER; i++ {
		w := &worker{
			graphClient: t.graphClient,
			done:        make(chan bool),
			logger:      t.logger.With(zap.Int("worker", i)),
		}
		go w.work(ctx, t.queries, t.config.Paths, writerChan)
		workers = append(workers, w)
	}

	queryCount := 0
	for blockNum := startBlockNum; blockNum < stopBlockNum; blockNum++ {
		for _, vars := range t.config.Vars {

			q := &thegraph.Query{
				Query: t.graphqlQuery,
				Variables: map[string]interface{}{
					BLOCK_NUM_VAR: blockNum,
				},
			}
			for varName, varValue := range vars {
				q.Variables[varName] = varValue
			}
			queryCount++
			t.queries <- q
		}
	}
	t.logger.Info("queued graphql queries", zap.Int("count", queryCount))
	close(t.queries)

	for _, w := range workers {
		w.Wait()
	}
	t.done <- true

}

func (t *testGenerator) wait() {
	t.logger.Info("waiting for test generator")
	<-t.done
}

type worker struct {
	queryHandleCount uint64
	done             chan bool
	graphClient      *thegraph.Graph
	startTime        time.Time
	logger           *zap.Logger
}

func (w *worker) Wait() {
	<-w.done
	w.logger.Info("test worker completed",
		zap.Uint64("query_count", w.queryHandleCount),
		zap.Duration("duration", time.Since(w.startTime)),
	)
}

func (w *worker) work(ctx context.Context, queries chan *thegraph.Query, paths []*Path, writerChan chan *subtest.TestConfig) {
	w.startTime = time.Now()
	for query := range queries {
		w.queryHandleCount++
		w.logger.Debug("handling query", zap.Reflect("query", query))

		resp, err := w.graphClient.Fetch(ctx, query)
		if err != nil {
			w.logger.Warn("failed to perform graphql call", zap.Error(err))
		}

		blockNum := getBlockNum(query)

		for _, path := range paths {

			expectVal := gjson.GetBytes(resp, path.Graph).String()
			if expectVal == "" {
				w.logger.Warn("expected value is blank skipping",
					zap.Uint64("block_num", blockNum),
					zap.String("query", query.Query),
				)
				continue
			}
			writerChan <- &subtest.TestConfig{
				Module: "module_name",
				Block:  blockNum,
				Path:   replaceVarsSubstreamPath(getVars(query), path.Substreams),
				Expect: expectVal,
				//Op:     "",
				//Args:   "",
			}

		}
	}
	w.done <- true
}

func getBlockNum(query *thegraph.Query) uint64 {
	if blockNum, found := query.Variables[BLOCK_NUM_VAR]; found {
		return blockNum.(uint64)
	}
	panic("expected block num in query variables")
}

func getVars(query *thegraph.Query) map[string]string {
	out := map[string]string{}
	for key, value := range query.Variables {
		if key == BLOCK_NUM_VAR {
			continue
		}
		out[key] = value.(string)
	}
	return out
}
