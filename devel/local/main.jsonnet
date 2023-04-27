local base = {
  Pool: {
    fields: {
      totalValueLockedETHUntracked: { ignore: true },
    },
  },
  PoolDayData: {
    fields: {
      totalValueLockedUSD: { ignore: true },
    },
  },
  PoolHourData: {
    fields: {
      totalValueLockedUSD: { ignore: true },
    },
  },
  TokenDayData: {
    fields: {
      tokenPrice: { ignore: true },
      volumeUSDUntracked: { rename: 'untrackedVolumeUSD' },
    },
  },
  TokenHourData: {
    fields: {
      tokenPrice: { ignore: true },
      volumeUSDUntracked: { rename: 'untrackedVolumeUSD' },
    },
  },
  UniswapDayData: {
    fields: {
      totalValueLockedUSD: { ignore: true },
    },
  },
};


local to_fix = {
};

std.mergePatch(import 'base.libsonnet', std.mergePatch(base, to_fix))
