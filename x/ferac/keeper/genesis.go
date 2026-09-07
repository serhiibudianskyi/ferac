package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	for _, account := range genState.RestrictedAccounts {
		if account.PeriodStart == 0 {
			account.PeriodStart = sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
		}
		if account.PeriodLimit.IsNil() || account.PeriodLimit.IsZero() {
			account.PeriodLimit = genState.Params.RestrictedReleaseRate.MulInt(account.Remaining).TruncateInt()
		}
		if account.Released.IsNil() {
			account.Released = math.ZeroInt()
		}

		if err := k.RestrictedAccounts.Set(ctx, account.Address, account); err != nil {
			return err
		}
	}

	reserve := genState.ValidatorReserve
	if reserve.LastEpochTime == 0 {
		reserve.LastEpochTime = sdk.UnwrapSDKContext(ctx).BlockTime().Unix()
	}
	if err := k.ValidatorReserve.Set(ctx, reserve); err != nil {
		return err
	}

	// Move the 4% reserve allocation into the module account it is paid out
	// from. Minting is disabled, so this is the only way the reserve is funded.
	if genState.ReserveFunder != "" && reserve.Remaining.IsPositive() {
		funder, err := k.addressCodec.StringToBytes(genState.ReserveFunder)
		if err != nil {
			return err
		}

		coins := sdk.NewCoins(sdk.NewCoin(genState.Params.Denom, reserve.Remaining))
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, funder, types.ReserveName, coins); err != nil {
			return err
		}
	}

	for _, perf := range genState.Performances {
		if err := k.Performances.Set(ctx, perf.ValidatorAddress, perf); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	genesis.ValidatorReserve, err = k.ValidatorReserve.Get(ctx)
	if err != nil {
		return nil, err
	}

	genesis.RestrictedAccounts = []types.RestrictedAccount{}
	if err := k.RestrictedAccounts.Walk(ctx, nil, func(_ string, account types.RestrictedAccount) (bool, error) {
		genesis.RestrictedAccounts = append(genesis.RestrictedAccounts, account)
		return false, nil
	}); err != nil {
		return nil, err
	}

	genesis.Performances = []types.ValidatorPerformance{}
	if err := k.Performances.Walk(ctx, nil, func(_ string, perf types.ValidatorPerformance) (bool, error) {
		genesis.Performances = append(genesis.Performances, perf)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
