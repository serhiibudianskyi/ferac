package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/keeper"
	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// ValidatorPolicyDecorator enforces the validator decentralization rules:
// admission requirements for new validators, the stake concentration cap and
// the 14 day exit notice a validator must serve before reducing its self bond.
type ValidatorPolicyDecorator struct {
	keeper keeper.Keeper
}

func NewValidatorPolicyDecorator(k keeper.Keeper) ValidatorPolicyDecorator {
	return ValidatorPolicyDecorator{keeper: k}
}

func (d ValidatorPolicyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	params, err := d.keeper.Params.Get(ctx)
	if err != nil {
		return ctx, err
	}

	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *stakingtypes.MsgCreateValidator:
			if err := d.keeper.AssertValidatorAdmission(
				ctx, m.ValidatorAddress, m.Value.Amount, m.Commission.Rate,
			); err != nil {
				return ctx, err
			}

		case *stakingtypes.MsgEditValidator:
			if m.CommissionRate != nil && !m.CommissionRate.Equal(params.ValidatorCommissionRate) {
				return ctx, errorsmod.Wrapf(
					types.ErrInvalidCommission,
					"commission rate must remain %s", params.ValidatorCommissionRate,
				)
			}

		case *stakingtypes.MsgDelegate:
			if err := d.keeper.AssertStakeConcentration(ctx, m.ValidatorAddress, m.Amount.Amount); err != nil {
				return ctx, err
			}

		case *stakingtypes.MsgBeginRedelegate:
			if err := d.keeper.AssertStakeConcentration(ctx, m.ValidatorDstAddress, m.Amount.Amount); err != nil {
				return ctx, err
			}

		case *stakingtypes.MsgUndelegate:
			if err := d.assertSelfUndelegationAllowed(ctx, m.DelegatorAddress, m.ValidatorAddress); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

// assertSelfUndelegationAllowed only constrains a validator unbonding its own
// stake. Ordinary delegators keep the short staking unbonding time.
func (d ValidatorPolicyDecorator) assertSelfUndelegationAllowed(ctx sdk.Context, delegator, validator string) error {
	delAddr, err := d.keeper.AddressCodec().StringToBytes(delegator)
	if err != nil {
		return err
	}

	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	if !sdk.AccAddress(valAddr).Equals(sdk.AccAddress(delAddr)) {
		return nil
	}

	return d.keeper.AssertExitNoticeServed(ctx, validator)
}

var _ sdk.AnteDecorator = ValidatorPolicyDecorator{}
