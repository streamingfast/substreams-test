package main

import (
	"fmt"
	"time"

	"github.com/streamingfast/dstore"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli/sflags"
	sink "github.com/streamingfast/substreams-sink"

	"github.com/streamingfast/substreams-test/thegraph"
	"github.com/streamingfast/substreams-test/validator"
	"github.com/streamingfast/substreams-test/validator/config"

	"go.uber.org/zap"
)

var testSubstreamCmd = &cobra.Command{
	Use:          "substream <spkg> <graph-url> <configpath> [<start>:<stop>]",
	Short:        "Tests the output graph_out of a substreams ",
	RunE:         runCmdE,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
}

func init() {
	testCmd.AddCommand(testSubstreamCmd)
	testSubstreamCmd.Flags().String("output-module", "graph_out", "Output module under test, it needs to be of output type proto:substreams.entity.v1.EntityChanges")
	testSubstreamCmd.Flags().StringP("endpoint", "e", "mainnet.eth.streamingfast.io:443", "Substreams endpoint")
	testSubstreamCmd.Flags().String("graphql-cache-store", "./localdata/.cache", "TheGraph GraphQL response caches store")
	testSubstreamCmd.Flags().Bool("only-error", false, "Only shows error")
	sink.AddFlagsToSet(testSubstreamCmd.Flags())
}

func runCmdE(cmd *cobra.Command, args []string) (err error) {
	ctx := cmd.Context()

	endpoint := sflags.MustGetString(cmd, "endpoint")
	outputModuleName := sflags.MustGetString(cmd, "output-module")
	graphqlCacheStore := sflags.MustGetString(cmd, "graphql-cache-store")

	onlyError := sflags.MustGetBool(cmd, "only-error")

	manifestPath := args[0]
	graphURL := args[1]
	configPath := args[2]

	blockRange := ""
	if len(args) > 3 {
		blockRange = args[3]
	}

	zlog.Info("running sustreams subgraph test",
		zap.String("endpoint", endpoint),
		zap.String("output_module_name", outputModuleName),
		zap.String("graph_url", graphURL),
		zap.String("manifest_path", manifestPath),
		zap.String("block_range", blockRange),
		zap.String("graphql_cache_store", graphqlCacheStore),
		zap.String("config_path", configPath),
		zap.Bool("only_error", onlyError),
	)

	config, err := config.ReadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read confile file %q: %w", configPath, err)
	}

	s, err := sink.NewFromViper(
		cmd,
		"proto:substreams.entity.v1.EntityChanges",
		endpoint,
		manifestPath,
		outputModuleName,
		blockRange,
		zlog,
		tracer,
	)
	if err != nil {
		return fmt.Errorf("failed to setup sink: %w", err)
	}

	opts := []thegraph.Option{
		thegraph.WithLogger(zlog),
	}
	if graphqlCacheStore != "" {
		cacheStore, err := dstore.NewStore(graphqlCacheStore, "", "", true)
		if err != nil {
			return fmt.Errorf("unable to create graphql cache store: %w", err)
		}
		opts = append(opts, thegraph.WithCache(cacheStore))
	}

	graphClient := thegraph.New(graphURL, opts...)

	var valOpts []validator.Option
	if onlyError {
		valOpts = append(valOpts, validator.WithOnlyError())
	}
	testRunner := validator.New(config, graphClient, zlog, valOpts...)
	t0 := time.Now()
	s.Run(ctx, nil, testRunner)

	if err := s.Err(); err != nil {
		fmt.Println("*************************************************************")
		fmt.Println("Error: ", err)
		fmt.Println("*************************************************************")
	}

	stats := testRunner.GetStats()
	fmt.Println()
	fmt.Println("Summary Stats")
	fmt.Printf(" > Duration %s\n", time.Since(t0).String())
	fmt.Printf(" > Block %d blocks [%d-%d]\n", testRunner.CurrentBlock-testRunner.FirstBlock, testRunner.FirstBlock, testRunner.CurrentBlock)
	fmt.Printf(" > Perform %d GraphQL Queries\n", graphClient.GetTotalQueries())
	fmt.Printf("     > %d cache hit\n", graphClient.GetCacheHitCount())
	fmt.Printf("     > %d cache missed\n", graphClient.GetCacheMissCount())

	fmt.Println("")
	fmt.Println(stats.Print())
	return
}
