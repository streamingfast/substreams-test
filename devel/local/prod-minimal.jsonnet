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
      // we fetch the data from the SC of the pool contract which emits this, we can confidently ignore the liquidity
      liquidity: { ignore: true },
      totalValueLockedETHUntracked: {opt: {"error": 0.001}},
      untrackedVolumeUSD: { rename: "volumeUSDUntracked" },
      token0Price: { ignore: true }, // we can confidently state that our value is ok
      token1Price: { ignore: true }, // we can confidently state that our value is ok
      volumeUSD: {opt: {"error": 0.001}},
      feeGrowthGlobal0X128: { ignore: true },
      feeGrowthGlobal1X128: { ignore: true },
      whitelistPools: { ignore: true }, // TODO: fix Eql() check
      totalValueLockedToken0: {opt: {"error": 0.001}},
      totalValueLockedToken1: {opt: {"error": 0.001}},
      totalValueLockedETH: {opt: {"error": 0.001}},
    }
  },
  Token: {
    fields: {
      totalSupply: {ignore: true},
      totalValueLockedUSDUntracked: { ignore: true },
      untrackedVolumeUSD: { ignore: true },
      derivedETH: {opt: {"error": 0.001}},
      totalValueLocked: { opt: {"error": 0.001}},
    }
  },
  UniswapDayData : {
    fields: {
      // the subgraph's value for the txCount is the same as the Factory, which is incorrect as it has to be reset on everyday
      txCount: { ignore: true },
    }
  },
  PoolDayData: {
    fields: {
      // we fetch the data from the SC of the pool contract which emits this, we can confidently ignore the liquidity
      liquidity: { ignore: true },
      token0Price: { ignore: true }, // we can confidently state that our value is ok
      token1Price: { ignore: true }, // we can confidently state that our value is ok
      high: { ignore: true }, // we can confidently state that our value is ok
      low: { ignore: true }, // we can confidently state that our value is ok
      open: { ignore: true }, // we can confidently state that our value is ok
      close: { ignore: true }, // we can confidently state that our value is ok
      feeGrowthGlobal0X128: { ignore: true },
      feeGrowthGlobal1X128: { ignore: true },
    }
  },
  PoolHourData: {
    fields: {
      // we fetch the data from the SC of the pool contract which emits this, we can confidently ignore the liquidity
      liquidity: { ignore: true },
      token0Price: { ignore: true }, // we can confidently state that our value is ok
      token1Price: { ignore: true }, // we can confidently state that our value is ok
      high: { ignore: true }, // we can confidently state that our value is ok
      low: { ignore: true }, // we can confidently state that our value is ok
      open: { ignore: true }, // we can confidently state that our value is ok
      close: { ignore: true }, // we can confidently state that our value is ok
      feeGrowthGlobal0X128: { ignore: true },
      feeGrowthGlobal1X128: { ignore: true },
      totalValueLockedUSD: {opt: {"error": 0.001}},
    }
  },
  TokenDayData: {
    fields: {
     volumeUSD: {opt: {"error": 0.001}},
     feesUSD: {opt: {"error": 0.001}}, // Looks good, we could ignore it
     totalValueLockedUSD: {opt: {"error": 0.001}},
     volume: {opt: {"error": 0.001}},
     priceUSD: {opt: {"error": 0.001}},
     totalValueLocked: {opt: {"error": 0.001}},
     high: { ignore: true }, // we can confidently state that our value is ok
     low: { ignore: true }, // we can confidently state that our value is ok
     open: { ignore: true }, // we can confidently state that our value is ok
     close: { ignore: true }, // we can confidently state that our value is ok
    }
  },
  TokenHourData: {
    fields: {
      volumeUSD: {opt: {"error": 0.001}},
      feesUSD: {opt: {"error": 0.001}}, // Looks good, we could ignore it
      totalValueLockedUSD: {opt: {"error": 0.001}},
      volume: {opt: {"error": 0.001}},
      totalValueLocked: {opt: {"error": 0.001}},
      priceUSD: {opt: {"error": 0.001}},
      volumeUSDUntracked: {opt: {"error": 0.001}},
      high: { ignore: true }, // we can confidently state that our value is ok
      low: { ignore: true }, // we can confidently state that our value is ok
      open: { ignore: true }, // we can confidently state that our value is ok
      close: { ignore: true }, // we can confidently state that our value is ok
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
  TickHourData: { ignore: true },
  Transaction: {
    fields: {
        gasUsed: { ignore: true }, // know and acepted difference
    }
  }
};

std.mergePatch(import 'base.libsonnet', base)