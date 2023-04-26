{
  "Factory": {
    "fields": {
    }
  },
  "Pool": {
    "fields": {
      "token1": { "association": true },
      "token0": { "association": true }
    }
  },
  "Token": {
    "fields": {
      "whitelistPools": { "association": true, "array": true }
    }
  },
  "Tick": {
    "fields": {
     "pool": { "association": true }
    }
  },
  "Mint": {
    "fields": {
      "pool": { "association": true },
      "token1": { "association": true },
      "token0": { "association": true },
      "transaction": { "association": true }
    }
  },
  "Burn": {
      "fields": {
        "pool": { "association": true },
        "token1": { "association": true },
        "token0": { "association": true },
        "transaction": { "association": true }
      }
    },
  "Swap": {
    "fields": {
      "pool": { "association": true },
      "token1": { "association": true },
      "token0": { "association": true },
      "transaction": { "association": true }
    }
  },
  "PoolDayData": {
    "fields": {
      "pool": { "association": true }
    }
  },
  "PoolHourData": {
    "fields": {
      "pool": { "association": true }
    }
  },
 "TokenDayData": {
    "fields": {
      "token": { "association": true }
    }
  },
  "TokenHourData": {
    "fields": {
      "token": { "association": true }
    }
  },
  "Position": {
    "fields": {
      "pool": { "association": true },
      "token0": { "association": true },
      "token1": { "association": true },
      "transaction": { "association": true },
      "tickLower": { "association": true },
      "tickUpper": { "association": true }
    }
  },
  "PositionSnapshot": {
    "fields": {
      "pool": { "association": true },
      "position": { "association": true },
      "token0": { "association": true },
      "token1": { "association": true },
    }
  }
}