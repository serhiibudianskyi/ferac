package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// NetworkFee returns the 0.01% protocol fee for a transferred amount and its
// split into the validator part and the treasury part (0.005% / 0.005%).
func (k Keeper) NetworkFee(ctx context.Context, amount math.Int) (validatorPart, treasuryPart math.Int, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}

	fee := params.NetworkFeeRate.MulInt(amount).TruncateInt()
	if !fee.IsPositive() {
		return math.ZeroInt(), math.ZeroInt(), nil
	}

	validatorPart = params.NetworkFeeValidatorShare.MulInt(fee).TruncateInt()
	treasuryPart = fee.Sub(validatorPart)

	return validatorPart, treasuryPart, nil
}

// ChargeNetworkFee collects the 0.01% network fee from the payer and routes half
// to the validators (through the fee collector, which the distribution module
// pays out) and half to the treasury.
func (k Keeper) ChargeNetworkFee(ctx context.Context, payer sdk.AccAddress, amount math.Int) error {
	validatorPart, treasuryPart, err := k.NetworkFee(ctx, amount)
	if err != nil {
		return err
	}
	if !validatorPart.IsPositive() && !treasuryPart.IsPositive() {
		return nil
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	if validatorPart.IsPositive() {
		// Reserve 5% of the validator pool for the proposer of this block.
		// The SDK proposer-reward parameters are deprecated in v0.53, so this
		// split is applied explicitly before the remaining pool is distributed.
		proposerPart := validatorPart.MulRaw(5).QuoRaw(100)
		poolPart := validatorPart.Sub(proposerPart)

		if poolPart.IsPositive() {
			poolCoins := sdk.NewCoins(sdk.NewCoin(params.Denom, poolPart))
			if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payer, authtypes.FeeCollectorName, poolCoins); err != nil {
				return err
			}
		}

		if proposerPart.IsPositive() {
			proposer, err := k.stakingKeeper.GetValidatorByConsAddr(
				ctx,
				sdk.ConsAddress(sdk.UnwrapSDKContext(ctx).BlockHeader().ProposerAddress),
			)
			if err == nil {
				proposerCoins := sdk.NewCoins(sdk.NewCoin(params.Denom, proposerPart))
				if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payer, authtypes.FeeCollectorName, proposerCoins); err != nil {
					return err
				}
				if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, distrtypes.ModuleName, proposerCoins); err != nil {
					return err
				}
				if err := k.distrKeeper.AllocateTokensToValidator(ctx, proposer, sdk.NewDecCoinsFromCoins(proposerCoins...)); err != nil {
					return err
				}
			} else {
				poolCoins := sdk.NewCoins(sdk.NewCoin(params.Denom, proposerPart))
				if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payer, authtypes.FeeCollectorName, poolCoins); err != nil {
					return err
				}
			}
		}
	}

	if treasuryPart.IsPositive() {
		if err := k.sendToTreasury(ctx, params, payer, treasuryPart); err != nil {
			return err
		}
	}

	return nil
}

// ChargeDexFee collects the 0.02% DEX swap fee and routes all of it to the
// treasury. It is exported for use by the DEX execution path.
func (k Keeper) ChargeDexFee(ctx context.Context, trader sdk.AccAddress, amount math.Int) (math.Int, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	fee := params.DexFeeRate.MulInt(amount).TruncateInt()
	if !fee.IsPositive() {
		return math.ZeroInt(), nil
	}

	return fee, k.sendToTreasury(ctx, params, trader, fee)
}

func (k Keeper) sendToTreasury(ctx context.Context, params types.Params, from sdk.AccAddress, amount math.Int) error {
	treasury, err := k.addressCodec.StringToBytes(params.TreasuryAddress)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(params.Denom, amount))

	return k.bankKeeper.SendCoins(ctx, from, sdk.AccAddress(treasury), coins)
}
