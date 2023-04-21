package validator

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"

	sink "github.com/streamingfast/substreams-sink"
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/thegraph"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Validator struct {
	graphClient *thegraph.Client

	stats  *Stats
	config Config

	logger *zap.Logger
}

func New(config Config, graphClient *thegraph.Client, logger *zap.Logger) *Validator {
	v := &Validator{
		graphClient: graphClient,
		stats: &Stats{
			successCount: 0,
			failedCount:  0,
		},
		config: config,
		logger: logger,
	}

	return v
}

func (v *Validator) GetStats() *Stats {
	return v.stats
}

func (v *Validator) HandleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	entityChanges := &pbentities.EntityChanges{}
	err := proto.Unmarshal(data.GetOutput().GetMapOutput().GetValue(), entityChanges)
	if err != nil {
		return fmt.Errorf("unmarshal database changes: %w", err)
	}
	v.logger.Debug("received blocked scoped data",
		zap.String("block_id", cursor.Block().ID()),
		zap.Uint64("block_num", cursor.Block().Num()),
		zap.Int("count", len(entityChanges.EntityChanges)),
	)

	if len(entityChanges.EntityChanges) == 0 {
		return nil
	}
	return v.handleEntityChanges(ctx, cursor.Block().Num(), entityChanges)
}

func (v *Validator) HandleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("unimplemented")
}

func (v *Validator) handleEntityChanges(ctx context.Context, blockNum uint64, changes *pbentities.EntityChanges) error {
	v.logger.Info("handling entity changes", zap.Uint64("block_num", blockNum), zap.Int("count", len(changes.EntityChanges)))

	for _, change := range changes.EntityChanges {
		if err := v.handleEntityChange(ctx, blockNum, change); err != nil {
			return fmt.Errorf("failed to handle entity change %q: %w", change.Entity, err)
		}
	}

	return nil
}

func (v *Validator) handleEntityChange(ctx context.Context, blockNum uint64, change *pbentities.EntityChange) error {
	logger := v.logger.With(zap.Uint64("block_num", blockNum), zap.String("entity", change.Entity))
	logger.Debug("entity_change", zap.Reflect("change", change))

	if v.shouldIgnoreEntity(change.Entity) {
		return nil
	}

	var entityFields []*Field
	for _, field := range change.Fields {
		if v.shouldIgnoreField(change.Entity, field.Name) {
			continue
		}
		entityFields = append(entityFields, v.newField(change.Entity, field))
	}

	subgraphEntity := normalizeEntityName(change.Entity)
	query := queryFromEntity(subgraphEntity, entityFields)

	vars := map[string]interface{}{
		"block": blockNum,
		"id":    change.Id,
	}

	logger.Debug("getting query for entity change",
		zap.String("query", query),
		zap.Reflect("vars", vars),
	)

	resp, err := v.graphClient.Fetch(ctx, query, vars)
	if err != nil {
		return fmt.Errorf("failed to query thegraph %s: %w", query, err)
	}

	if gjson.GetBytes(resp, fmt.Sprintf("data.%s", subgraphEntity)).String() == "" {
		fmt.Printf("❌ [%d] %s.%s unable to find entity [GRQLERR]\n", blockNum, subgraphEntity, change.Id)
		return nil
	}

	for _, field := range entityFields {
		prefix := fmt.Sprintf("[%d] %s.%s.%s", blockNum, subgraphEntity, change.Id, field.substreamsField)

		subgraphValue := gjson.GetBytes(resp, field.graphqlJSONPath).String()

		cmp, err := field.obj.eql(subgraphValue)
		if err != nil {
			return fmt.Errorf("failed to compare field %s: %w", field.substreamsField, err)
		}
		if cmp {
			v.stats.successCount++
			fmt.Printf("✅ %s\n", prefix)
		} else {
			v.stats.failedCount++
			fmt.Printf("❌ %s: sub: %s <-> grt: %s\n", prefix, field.obj.string(), subgraphValue)
		}
	}
	return nil

}

func normalizeEntityName(s string) string {
	if len(s) != 0 && (s[0] <= 90 && s[0] >= 65) {
		return string(s[0]+32) + string(s[1:])
	}
	return s
}
