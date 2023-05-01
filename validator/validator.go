package validator

import (
	"context"
	"fmt"

	config2 "github.com/streamingfast/substreams-test/validator/config"
	"github.com/streamingfast/substreams-test/validator/fields"

	"github.com/tidwall/gjson"

	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/thegraph"
)

type Validator struct {
	graphClient *thegraph.Client

	stats         *Stats
	config        config2.Config
	showOnlyError bool

	logger *zap.Logger

	FirstBlock   uint64
	CurrentBlock uint64
}

type Option func(v *Validator) *Validator

func WithOnlyError() Option {
	return func(v *Validator) *Validator {
		v.showOnlyError = true
		return v
	}
}
func New(config config2.Config, graphClient *thegraph.Client, logger *zap.Logger, opts ...Option) *Validator {
	v := &Validator{
		graphClient: graphClient,
		stats:       newStats(),
		config:      config,
		logger:      logger,
	}

	for _, opt := range opts {
		v = opt(v)
	}
	return v
}

func (v *Validator) GetStats() *Stats {
	return v.stats
}

func (v *Validator) HandleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	blockNum := cursor.Block().Num()
	if v.FirstBlock == 0 {
		v.FirstBlock = blockNum
	}
	v.CurrentBlock = blockNum

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
	return v.handleEntityChanges(ctx, blockNum, entityChanges)
}

func (v *Validator) HandleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("unimplemented")
}

func (v *Validator) handleEntityChanges(ctx context.Context, blockNum uint64, changes *pbentities.EntityChanges) error {
	v.logger.Debug("handling entity changes", zap.Uint64("block_num", blockNum), zap.Int("count", len(changes.EntityChanges)))

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

	var entityFields []*fields.Field
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

	resp, err := v.graphClient.Fetch(ctx, blockNum, query, vars)
	if err != nil {
		return fmt.Errorf("failed to query thegraph %s: %w", query, err)
	}

	if gjson.GetBytes(resp, fmt.Sprintf("data.%s", subgraphEntity)).String() == "" {
		fmt.Printf("❌ [%d] %s.%s unable to find entity [GRQLERR]\n", blockNum, subgraphEntity, change.Id)
		return nil
	}

	for _, field := range entityFields {
		prefix := fmt.Sprintf("[%d] %s.%s.%s", blockNum, subgraphEntity, change.Id, field.SubstreamsField)

		subgraphValueRes := gjson.GetBytes(resp, field.GraphqlJSONPath)
		if subgraphValueRes.Type == gjson.Null {
			fmt.Printf("❌ %s: sub: %s <-> grql: NULL\n", prefix, field.Obj.String())
			continue
		}

		actualValue, err := field.ObjFactory(subgraphValueRes.String())
		if err != nil {
			return fmt.Errorf("failed to convert subgraph value %s: %w", subgraphValueRes.String(), err)
		}

		if field.Obj.Eql(actualValue) {
			v.stats.Success(change.Entity, field.SubstreamsField)
			if !v.showOnlyError {
				fmt.Printf("✅ %-120s ✅ sub: %s <-> grql: %s\n", prefix, field.Obj.String(), subgraphValueRes.String())
			}
		} else {
			v.stats.Fail(change.Entity, field.SubstreamsField)
			fmt.Printf("❌ %-120s ❌ sub: %s <-> grql: %s\n", prefix, field.Obj.String(), subgraphValueRes.String())
		}
	}
	return nil

}
