package validator

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/streamingfast/substreams-test/thegraph"

	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/test-go/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestValidator_handleEntityChanges(t *testing.T) {
	tests := []struct {
		blockNum     uint64
		expecteError bool
	}{
		{12370010, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.blockNum), func(t *testing.T) {
			filename := fmt.Sprintf("./testdata/%d.hex", test.blockNum)

			cnt, err := os.ReadFile(filename)
			require.NoError(t, err)

			data, err := hex.DecodeString(string(cnt))
			require.NoError(t, err)

			changes := &pbentities.EntityChanges{}
			require.NoError(t, proto.Unmarshal(data, changes))

			//cacheStore, err := dstore.NewStore("./testdata/.cache", "", "", false)
			require.NoError(t, err)
			v := &Validator{
				graphClient: thegraph.New("https://api.thegraph.com/subgraphs/name/ianlapham/v3-minimal", thegraph.WithLogger(zlog)),
				config: Config{
					"PositionSnapshot": &EntityConfig{
						Ignore: true,
					},
					"Position": &EntityConfig{
						Ignore: true,
					},
				},
				logger: zlog,
			}
			err = v.handleEntityChanges(context.Background(), test.blockNum, changes)
			require.NoError(t, err)
		})
	}
}
