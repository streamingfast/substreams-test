package fields

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/streamingfast/eth-go"

	"github.com/streamingfast/substreams-test/validator/config"
)

type Bytes struct {
	v []byte
}

func newBytes(v []byte, opt config.Options) *Bytes {
	return &Bytes{v: v}
}

func newBytesFromStr(v string) (Comparable, error) {
	var data []byte
	var err error
	if strings.HasPrefix(v, "0x") {
		if data, err = hex.DecodeString(eth.SanitizeHex(v)); err != nil {
			return nil, fmt.Errorf("failed to convert %s to byte array: %w", v, err)
		}
	}
	return &Bytes{v: data}, nil
}

func (f *Bytes) Eql(v Comparable) bool {
	return bytes.Compare(f.v, v.(*Bytes).v) == 0
}

func (f *Bytes) String() string {
	return hex.EncodeToString(f.v)
}
