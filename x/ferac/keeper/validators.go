package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// AnnounceExit starts the validator exit notice period.
func (k Keeper) AnnounceExit(ctx context.Context, operator string) (types.ValidatorExit, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return types.ValidatorExit{}, err
	}

	now := sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	exit := types.ValidatorExit{
		ValidatorAddress: operator,
		AnnouncedAt:      now,
		AllowedAt:        now + params.ValidatorUnbondingSeconds,
	}

	return exit, k.ValidatorExits.Set(ctx, operator, exit)
}

// AssertExitNoticeServed fails unless the validator announced its exit at least
// `validator_unbonding_seconds` ago. Delegators are not affected: they only wait
// for the staking unbonding time.
func (k Keeper) AssertExitNoticeServed(ctx context.Context, operator string) error {
	exit, err := k.ValidatorExits.Get(ctx, operator)
	if err != nil {
		return errorsmod.Wrapf(
			types.ErrExitNoticeRequired,
			"validator %s must announce its exit before reducing its self delegation", operator,
		)
	}

	if now := sdk.UnwrapSDKContext(ctx).BlockTime().Unix(); now < exit.AllowedAt {
		return errorsmod.Wrapf(
			types.ErrExitNoticeRequired,
			"validator %s may self-undelegate from unix time %d, current time is %d",
			operator, exit.AllowedAt, now,
		)
	}

	return nil
}

// AssertStakeConcentration fails if adding `additional` tokens to the validator
// would push it above the maximum allowed share of the total bonded stake.
// This is what stops a single actor from buying majority voting power.
func (k Keeper) AssertStakeConcentration(ctx context.Context, operator string, additional math.Int) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(operator)
	if err != nil {
		return err
	}

	totalBonded, err := k.stakingKeeper.TotalBondedTokens(ctx)
	if err != nil {
		return err
	}

	current := math.ZeroInt()
	if validator, err := k.stakingKeeper.GetValidator(ctx, sdk.ValAddress(valAddr)); err == nil {
		current = validator.GetBondedTokens()
	}

	newStake := current.Add(additional)
	newTotal := totalBonded.Add(additional)
	if !newTotal.IsPositive() {
		return nil
	}

	// The genesis set is intentionally exempt: with 4 validators every one of
	// them necessarily holds more than the long term concentration cap.
	validators, err := k.stakingKeeper.GetLastValidators(ctx)
	if err != nil {
		return err
	}
	if uint32(len(validators)) <= params.GenesisValidatorCount {
		return nil
	}

	share := math.LegacyNewDecFromInt(newStake).QuoInt(newTotal)
	if share.GT(params.MaxValidatorStakeShare) {
		return errorsmod.Wrapf(
			types.ErrStakeConcentration,
			"validator %s would hold %s of the bonded stake, maximum is %s",
			operator, share, params.MaxValidatorStakeShare,
		)
	}

	return nil
}

// AssertValidatorAdmission enforces the entry rules for a new validator:
// a minimum self bond, the protocol commission rate and a cap on how many
// validators a single operator account may run.
func (k Keeper) AssertValidatorAdmission(ctx context.Context, operator string, selfDelegation math.Int, commissionRate math.LegacyDec) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if selfDelegation.LT(params.MinSelfDelegation) {
		return errorsmod.Wrapf(
			types.ErrInsufficientSelfDelegation,
			"self delegation %s is below the minimum %s", selfDelegation, params.MinSelfDelegation,
		)
	}

	if !commissionRate.Equal(params.ValidatorCommissionRate) {
		return errorsmod.Wrapf(
			types.ErrInvalidCommission,
			"commission rate must be %s, got %s", params.ValidatorCommissionRate, commissionRate,
		)
	}

	valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(operator)
	if err != nil {
		return err
	}

	// One operator account may only back `max_validators_per_operator`
	// validators. Combined with the minimum self bond this makes running a farm
	// of Sybil validators expensive rather than free.
	owned, err := k.countValidatorsBackedBy(ctx, sdk.AccAddress(valAddr), params.MinSelfDelegation)
	if err != nil {
		return err
	}
	if owned >= params.MaxValidatorsPerOperator {
		return errorsmod.Wrapf(
			types.ErrSybilValidator,
			"account already backs %d validators, maximum is %d", owned, params.MaxValidatorsPerOperator,
		)
	}

	return k.AssertStakeConcentration(ctx, operator, selfDelegation)
}

// countValidatorsBackedBy returns how many validators the account funds with a
// stake large enough to qualify as a self bond.
func (k Keeper) countValidatorsBackedBy(ctx context.Context, account sdk.AccAddress, threshold math.Int) (uint32, error) {
	validators, err := k.stakingKeeper.GetLastValidators(ctx)
	if err != nil {
		return 0, err
	}

	var count uint32
	for _, val := range validators {
		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			continue
		}

		delegation, err := k.stakingKeeper.GetDelegation(ctx, account, sdk.ValAddress(valAddr))
		if err != nil {
			continue
		}
		if val.TokensFromShares(delegation.Shares).TruncateInt().GTE(threshold) {
			count++
		}
	}

	return count, nil
}

// SetValidatorQuality stores the governance-set quality score of a validator.
func (k Keeper) SetValidatorQuality(ctx context.Context, operator string, score math.LegacyDec) error {
	perf, err := k.Performances.Get(ctx, operator)
	if err != nil {
		perf = types.ValidatorPerformance{ValidatorAddress: operator}
	}
	perf.QualityScore = score

	return k.Performances.Set(ctx, operator, perf)
}
