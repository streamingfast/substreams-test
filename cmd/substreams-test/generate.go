package main

import (
	"fmt"
	"github.com/spf13/cobra"
	subcmp "github.com/streamingfast/substreams-test"
	"go.uber.org/zap"
	"strconv"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

// runCmd represents the command to run substreams remotely
var generateCmd = &cobra.Command{
	Use:          "generate <config.yml> <start-block-num> <block-count>",
	Short:        "Generate a test file from GraphQL",
	RunE:         generateRunE,
	Args:         cobra.ExactArgs(3),
	SilenceUsage: true,
}

func generateRunE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	configPath := args[0]

	startBlockNum, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse start block num %q: %w", args[1], err)
	}

	blockCount, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse block count %q: %w", args[2], err)
	}

	stopBlockNum := startBlockNum + blockCount

	zlog.Info("running generation",
		zap.String("config_path", configPath),
		zap.Uint64("start_block_num", startBlockNum),
		zap.Uint64("block_count", blockCount),
		zap.Uint64("stop_block_num", stopBlockNum),
	)

	config, err := subcmp.ReadConfigFromFile(configPath)
	if err != nil {
		return fmt.Errorf("unable to read config %q: %w", configPath, err)
	}

	gen := subcmp.Generate{
		StartBlockNum: startBlockNum,
		StopBlockNum:  stopBlockNum,
		Config:        config,
		Logger:        zlog,
	}

	return gen.Run(ctx)
}
