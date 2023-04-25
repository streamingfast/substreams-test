local base = {
  "Factory": {
      "fields": {
          "totalFeesETH": { "ignore": true },
          "totalVolumeETH": { "rename": "volumeETH" },
          "totalVolumeUSD": { "rename": "volumeUSD" },
          "untrackedVolumeUSD": { "rename": "volumeUSDUntracked" },
          "totalFeesUSD": { "rename": "feesUSD" },
      }
  },
  "Pool": {
      "fields": {
        "untrackedVolumeUSD": { "rename": "volumeUSDUntracked" },
        "token1Price": { "opt": {"tolerance": 0.013}}
      }
  },
  "Transaction": {
    "fields": {
      "gasUsed": { "opt": {"error": 15}}
    }
  },
  "Token": {
      "fields": {
        "totalValueLockedUSDUntracked": { "ignore": true },
        "untrackedVolumeUSD": { "ignore": true }
      }
  },
  "PositionSnapshot": { "ignore": true },
  "Position": { "ignore": true },
  "Tick": {
    "fields": {
      "volumeToken0": { "ignore": true},
      "volumeToken1": { "ignore": true},
      "volumeUSD": { "ignore": true},
      "untrackedVolumeUSD": { "ignore": true},
      "feesUSD": { "ignore": true},
      "collectedFeesToken0": { "ignore": true},
      "collectedFeesToken1": { "ignore": true},
      "collectedFeesUSD": { "ignore": true},
      "liquidityProviderCount": { "ignore": true},
      "feeGrowthOutside0X128": { "ignore": true},
      "feeGrowthOutside1X128": { "ignore": true},
    }
  }
};


local to_fix = {
  "TokenHourData": {
    "fields": {
      "tokenPrice": {
        "ignore": true
      }
    }
  },
  "TokenDayData": {
    "fields": {
      "tokenPrice": {
        "ignore": true
      }
    }
  }
};

std.mergePatch(import 'base.libsonnet',std.mergePatch(base, to_fix))