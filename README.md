# Substreams Subgraph Test

This tools allows you to compare the values from a Substreams graph_out module with a Subgraph.

## Installing

```bash
go install ./cmd/substreams-graph-test
```

## Running

```bash
substreams-test generate <path-to-config-file.json> <start-block-num> <block_count>
substreams-graph-test test substream <substream-manifes-path> <subgraph-api-rul> <config_file> [<start:stop>]

```

## Configuration file 

The configuration files allows you to ignore Entity or Fields during testing and allows you to specify tolerances when comparing values.

```json
{
  "Factory": {
    "ignore": true
  },
  "PositionSnapshot": {
    "ignore": true
  },
  "Position": {
    "ignore": true
  },
  "Token": {
    "fields": {
      "untrackedVolumeUSD": {
        "rename": "volumeUSDUntracked"
      }
    }
  },
  "TokenDayData": {
    "fields": {
      "tokenPrice": {
        "ignore": true
      }
    }
  },
  "TokenHourData": {
    "fields": {
      "tokenPrice": {
        "ignore": true
      }
    }
  }
}```