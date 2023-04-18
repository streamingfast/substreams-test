# Substreams Test

Tool that helps you build a Substreams test file based on a deployed Subgraph

## Installing

```bash
go install ./cmd/substreams-test
```

## Running

```bash
substreams-test generate <path-to-config-file.json> <start-block-num> <block_count>
```

## Configuration file 

```json
{
  "graph_url": "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3",
  "test_output_path": "./test.jsonl",
  "substreams_module": "graph_out",

  "tests":{
   "./queries/pool.graphql": {
     "paths": [
       {"graph":  "data.pool.feeGrowthGlobal1X128", "substreams": ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"${pool}\") | .newValue.value" }
     ],
     "vars": [
       { "pool": "0x6c6bc977e13df9b0de53b251522280bb72383700" }
     ]
   }
  }
}
```