package keeper

import (
	"context"

	"cosmossdk.io/core/comet"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// TrackUptime records whether each bonded validator signed the previous block.
// The counters feed the uptime component of the reward formula and are reset at
// the end of every reserve epoch.
func (k Keeper) TrackUptime(ctx context.Context) error {
	info := sdk.UnwrapSDKContext(ctx).CometInfo()
	if info == nil {
		return nil
	}

	commit := info.GetLastCommit()
	if commit == nil {
		return nil
	}

	votes := commit.Votes()
	valCodec := k.stakingKeeper.ValidatorAddressCodec()

	for i := 0; i < votes.Len(); i++ {
		vote := votes.Get(i)

		validator, err := k.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(vote.Validator().Address()))
		if err != nil {
			continue
		}

		operator := validator.GetOperator()
		if _, err := valCodec.StringToBytes(operator); err != nil {
			continue
		}

		perf, err := k.Performances.Get(ctx, operator)
		if err != nil {
			perf = types.ValidatorPerformance{
				ValidatorAddress: operator,
				QualityScore:     math.LegacyOneDec(),
			}
		}

		perf.TotalBlocks++
		if vote.GetBlockIDFlag() == comet.BlockIDFlagCommit {
			perf.SignedBlocks++
		}

		if err := k.Performances.Set(ctx, operator, perf); err != nil {
			return err
		}
	}

	return nil
}

// MaybeReleaseReserve releases 5% of the *remaining* validator reserve once per
// reserve period and distributes it with the 40% stake / 40% uptime / 20%
// quality formula. The reserve therefore decays geometrically and is never
// emptied.
func (k Keeper) MaybeReleaseReserve(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	reserve, err := k.ValidatorReserve.Get(ctx)
	if err != nil {
		return err
	}

	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	if reserve.LastEpochTime == 0 {
		reserve.LastEpochTime = blockTime
		return k.ValidatorReserve.Set(ctx, reserve)
	}
	if blockTime < reserve.LastEpochTime+params.ReservePeriodSeconds {
		return nil
	}

	release := params.ReserveReleaseRate.MulInt(reserve.Remaining).TruncateInt()
	if release.IsPositive() {
		if err := k.distributeRewards(ctx, params, release); err != nil {
			return err
		}
		reserve.Remaining = reserve.Remaining.Sub(release)
	}

	reserve.LastEpochTime += params.ReservePeriodSeconds
	reserve.Epoch++

	if err := k.ValidatorReserve.Set(ctx, reserve); err != nil {
		return err
	}

	return k.resetPerformances(ctx)
}

// NextReserveRelease returns the amount that the next reserve epoch will pay out.
func (k Keeper) NextReserveRelease(ctx context.Context) (math.Int, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	reserve, err := k.ValidatorReserve.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	return params.ReserveReleaseRate.MulInt(reserve.Remaining).TruncateInt(), nil
}

// distributeRewards splits `amount` across the bonded validators proportionally
// to 0.4*stake + 0.4*uptime + 0.2*quality, each normalised across the set.
func (k Keeper) distributeRewards(ctx context.Context, params types.Params, amount math.Int) error {
	validators, err := k.stakingKeeper.GetLastValidators(ctx)
	if err != nil {
		return err
	}
	if len(validators) == 0 {
		return nil
	}

	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	scores := make([]math.LegacyDec, len(validators))
	totalScore := math.LegacyZeroDec()
	totalUptime := math.LegacyZeroDec()
	totalQuality := math.LegacyZeroDec()

	uptimes := make([]math.LegacyDec, len(validators))
	qualities := make([]math.LegacyDec, len(validators))

	for i, val := range validators {
		perf, err := k.Performances.Get(ctx, val.GetOperator())
		if err != nil {
			perf = types.ValidatorPerformance{QualityScore: math.LegacyOneDec()}
		}
		if perf.QualityScore.IsNil() {
			perf.QualityScore = math.LegacyOneDec()
		}

		uptimes[i] = perf.Uptime()
		qualities[i] = perf.QualityScore
		totalUptime = totalUptime.Add(uptimes[i])
		totalQuality = totalQuality.Add(qualities[i])
	}

	for i, val := range validators {
		stakeShare := math.LegacyZeroDec()
		if totalBonded.IsPositive() {
			stakeShare = math.LegacyNewDecFromInt(val.GetBondedTokens()).QuoInt(totalBonded)
		}

		uptimeShare := math.LegacyZeroDec()
		if totalUptime.IsPositive() {
			uptimeShare = uptimes[i].Quo(totalUptime)
		}

		qualityShare := math.LegacyZeroDec()
		if totalQuality.IsPositive() {
			qualityShare = qualities[i].Quo(totalQuality)
		}

		scores[i] = params.StakeWeight.Mul(stakeShare).
			Add(params.UptimeWeight.Mul(uptimeShare)).
			Add(params.QualityWeight.Mul(qualityShare))
		totalScore = totalScore.Add(scores[i])
	}

	if !totalScore.IsPositive() {
		return nil
	}

	distributed := math.ZeroInt()
	payouts := make([]math.Int, len(validators))
	for i := range validators {
		payouts[i] = scores[i].Quo(totalScore).MulInt(amount).TruncateInt()
		distributed = distributed.Add(payouts[i])
	}
	// Truncation dust goes to the highest scoring validator.
	if remainder := amount.Sub(distributed); remainder.IsPositive() {
		payouts[0] = payouts[0].Add(remainder)
	}

	for i, val := range validators {
		if !payouts[i].IsPositive() {
			continue
		}

		coins := sdk.NewCoins(sdk.NewCoin(params.Denom, payouts[i]))
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ReserveName, distrtypes.ModuleName, coins); err != nil {
			return err
		}

		validator := val
		if err := k.distrKeeper.AllocateTokensToValidator(ctx, stakingtypes.ValidatorI(validator), sdk.NewDecCoinsFromCoins(coins...)); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) resetPerformances(ctx context.Context) error {
	return k.Performances.Walk(ctx, nil, func(key string, perf types.ValidatorPerformance) (bool, error) {
		perf.SignedBlocks = 0
		perf.TotalBlocks = 0
		if err := k.Performances.Set(ctx, key, perf); err != nil {
			return true, err
		}
		return false, nil
	})
}
