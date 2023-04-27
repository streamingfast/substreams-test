local base = {
  Factory: {
    fields: {
      totalFeesETH: { ignore: true },
      totalVolumeETH: { rename: 'volumeETH' },
      totalVolumeUSD: { rename: 'volumeUSD' },
      untrackedVolumeUSD: { rename: 'volumeUSDUntracked' },
      totalFeesUSD: { rename: 'feesUSD' },
    },
  },
  Pool: {
    fields: {
      untrackedVolumeUSD: { rename: 'volumeUSDUntracked' },
      token0Price: { opt: { round: 'shortest' } },
      token1Price: { opt: { round: 'shortest' } },
      volumeUSD: { opt: { round: 'shortest' } },
    },
  },
  Token: {
    fields: {
      totalSupply: { ignore: true },
      totalValueLockedUSDUntracked: { ignore: true },
      untrackedVolumeUSD: { ignore: true },
      feeGrowthGlobal0X128: { ignore: true },  // field not present in minimal
      feeGrowthGlobal1X128: { ignore: true },  // field not present in minimal
    },
  },
  PoolDayData: {
    fields: {
      token0Price: { opt: { "error": 0.00001 } },
      token1Price: { opt: { "error": 0.00001 } },
      high: { opt: { round: 'shortest' } },
    },
  },
  PoolHourData: {
    fields: {
      token0Price: { opt: { "error": 0.00001 } },
      token1Price: { opt: { "error": 0.00001 } },
    },
  },
  TokenDayData: {
    fields: {
    },
  },
  TokenHourData: {
    fields: {
    },
  },
  Mint: {
    fields: {
      amountUSD: { opt: { round: 'shortest' } },
    },
  },
  Tick: {
    fields: {
      volumeToken0: { ignore: true },
      volumeToken1: { ignore: true },
      volumeUSD: { ignore: true },
      untrackedVolumeUSD: { ignore: true },
      feesUSD: { ignore: true },
      collectedFeesToken0: { ignore: true },
      collectedFeesToken1: { ignore: true },
      collectedFeesUSD: { ignore: true },
      liquidityProviderCount: { ignore: true },
      feeGrowthOutside0X128: { ignore: true },
      feeGrowthOutside1X128: { ignore: true },
    },
  },
  PositionSnapshot: { ignore: true },
  Position: { ignore: true },
};

std.mergePatch(import 'base.libsonnet', base)
