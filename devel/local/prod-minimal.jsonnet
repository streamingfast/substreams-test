local base = {
  Bundle: {
      fields: {
          ethPriceUSD: {opt: {"error": 0.001}},
      }
  },
  Factory: {
      fields: {
          totalFeesETH: { ignore: true },
          totalVolumeETH: { rename: "volumeETH" },
          totalVolumeUSD: { rename: "volumeUSD" },
          untrackedVolumeUSD: { rename: "volumeUSDUntracked" },
          totalFeesUSD: { rename: "feesUSD" },
      }
  },
   Pool: {
      fields: {
          totalValueLockedETHUntracked: {opt: {"error": 0.001}},
          untrackedVolumeUSD: { rename: "volumeUSDUntracked" },
          token0Price: {opt: {"error": 0.001}},
          token1Price: {opt: {"error": 0.001}},
          volumeUSD: {opt: {"error": 0.001}},
          feeGrowthGlobal0X128: { ignore: true }, // field not present in minimal
          feeGrowthGlobal1X128: { ignore: true},     // field not present in minimal
          whitelistPools: { ignore: true }, // TODO: fix Eql() check
      }
  },
  Token: {
      fields: {
        totalSupply: {ignore: true},
        totalValueLockedUSDUntracked: { ignore: true },
        untrackedVolumeUSD: { ignore: true },
        feeGrowthGlobal0X128: { ignore: true },    // field not present in minimal
        feeGrowthGlobal1X128: { ignore: true},     // field not present in minimal
        derivedETH: {opt: {"error": 0.001}},
      }
  },
  PoolDayData: {
    fields: {
      token0Price: {opt: {"error": 0.001}},
      token1Price: {opt: {"error": 0.001}},
      high: {opt: {"error": 0.001}},
      low: {opt: {"error": 0.001}},
      open: {opt: {"error": 0.001}},
      close: {opt: {"error": 0.001}},
      feeGrowthGlobal0X128: { ignore: true},     // field not present in minimal
      feeGrowthGlobal1X128: { ignore: true},     // field not present in minimal
    }
  },
  PoolHourData: {
    fields: {
      token0Price: {opt: {"error": 0.001}},
      token1Price: {opt: {"error": 0.001}},
      high: {opt: {"error": 0.001}},
      low: {opt: {"error": 0.001}},
      open: {opt: {"error": 0.001}},
      close: {opt: {"error": 0.001}},
    }
  },
  TokenDayData: {
    fields: {
     volumeUSD: {opt: {"error": 0.001}},
     feesUSD: {opt: {"error": 0.001}},
     totalValueLockedUSD: {opt: {"error": 0.001}},
     volume: {opt: {"error": 0.001}},
     priceUSD: {opt: {"error": 0.001}},
     totalValueLocked: {opt: {"error": 0.001}},
    }
  },
  TokenHourData: {
    fields: {
      volumeUSD: {opt: {"error": 0.001}},
      feesUSD: {opt: {"error": 0.001}},
      totalValueLockedUSD: {opt: {"error": 0.001}},
      volume: {opt: {"error": 0.001}},
      totalValueLocked: {opt: {"error": 0.001}},
      priceUSD: {opt: {"error": 0.001}},
      volumeUSDUntracked: {opt: {"error": 0.001}},
    }
  },
  Mint: {
    fields: {
      amountUSD: {opt: {"error": 0.001}},
    }
  },
  Tick: { ignore: true },
  PositionSnapshot: { ignore: true },
  Position: { ignore: true },
  TickDayData: { ignore: true },
  TickHourData: { ignore: true }
};

std.mergePatch(import 'base.libsonnet', base)