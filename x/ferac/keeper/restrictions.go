package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// SendRestrictionFn enforces the creator/team transfer limit: within any
// 12 month period such an account may move at most 25% of the balance that was
// still restricted when the period started. Unused limit is never carried over.
//
// It is registered on the bank keeper so the rule applies to every transfer
// path, not only to direct MsgSend.
func (k Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.AccAddress, error) {
	from, err := k.addressCodec.BytesToString(fromAddr)
	if err != nil {
		return toAddr, err
	}

	account, err := k.RestrictedAccounts.Get(ctx, from)
	if err != nil {
		// Not a restricted account: nothing to enforce.
		return toAddr, nil
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return toAddr, err
	}

	amount := amt.AmountOf(params.Denom)
	if !amount.IsPositive() {
		return toAddr, nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	account = account.Rollover(sdkCtx.BlockTime().Unix(), params.RestrictedPeriodSeconds, params.RestrictedReleaseRate)

	if available := account.Available(); amount.GT(available) {
		periodEnd := account.PeriodStart + params.RestrictedPeriodSeconds
		return toAddr, errorsmod.Wrapf(
			types.ErrTransferLimitExceeded,
			"%s may transfer at most %s%s before unix time %d, requested %s%s",
			from, available, params.Denom, periodEnd, amount, params.Denom,
		)
	}

	account.Released = account.Released.Add(amount)
	account.Remaining = account.Remaining.Sub(amount)
	if account.Remaining.IsNegative() {
		account.Remaining = math.ZeroInt()
	}

	if err := k.RestrictedAccounts.Set(ctx, from, account); err != nil {
		return toAddr, err
	}

	return toAddr, nil
}

// SetRestrictedAccount registers an account under the creator/team restriction.
func (k Keeper) SetRestrictedAccount(ctx context.Context, address string, category types.AllocationCategory, amount math.Int) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	account := types.NewRestrictedAccount(
		address, category, amount,
		params.RestrictedReleaseRate,
		sdk.UnwrapSDKContext(ctx).BlockTime().Unix(),
	)
	if err := account.Validate(); err != nil {
		return err
	}

	return k.RestrictedAccounts.Set(ctx, address, account)
}

// AssertSupplyCap fails if the circulating supply of the FERAC denom is above
// the immutable hard cap. Minting is disabled, so this is a safety net that
// makes an accidental inflation source impossible to miss.
func (k Keeper) AssertSupplyCap(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	supply := k.bankKeeper.GetSupply(ctx, params.Denom)
	if supply.Amount.GT(params.MaxSupply) {
		return errorsmod.Wrap(
			types.ErrMaxSupplyExceeded,
			fmt.Sprintf("supply %s exceeds max supply %s", supply.Amount, params.MaxSupply),
		)
	}

	return nil
}
