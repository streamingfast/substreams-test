{
  "Factory": {
    "fields": {
      "id": { "ignore": true },
    }
  },
  "Pool": {
    "fields": {
      "id": { "ignore": true },
      "token1": { "association": true },
      "token0": { "association": true }
    }
  },
  "Token": {
    "fields": {
      "id": { "ignore": true },
      "whitelistPools": { "association": true, "array": true }
    }
  },
  "Tick": {
    "fields": {
      "id": { "ignore": true },
    }
  },
  "Transaction": {
    "fields": {
      "id": { "ignore": true },
    }
  },
  "Mint": {
    "fields": {
      "id": { "ignore": true },
      "pool": { "association": true },
      "token1": { "association": true },
      "token0": { "association": true },
      "transaction": { "association": true }
    }
  },
  "Burn": {
      "fields": {
        "id": { "ignore": true },
        "pool": { "association": true },
        "token1": { "association": true },
        "token0": { "association": true },
        "transaction": { "association": true }
      }
    },
  "Swap": {
    "fields": {
      "id": { "ignore": true },
      "pool": { "association": true },
      "token1": { "association": true },
      "token0": { "association": true },
      "transaction": { "association": true }
    }
  },
  "PoolDayData": {
    "fields": {
      "id": { "ignore": true },
      "pool": { "association": true }
    }
  },
  "PoolHourData": {
    "fields": {
      "id": { "ignore": true },
      "pool": { "association": true }
    }
  },
 "TokenDayData": {
    "fields": {
      "id": { "ignore": true },
      "token": { "association": true }
    }
  },
  "TokenHourData": {
    "fields": {
      "id": { "ignore": true },
      "token": { "association": true }
    }
  }
}