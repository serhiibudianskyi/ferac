package ante

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/keeper"
)

// NetworkFeeDecorator charges the 0.01% protocol fee on FERAC transfers and
// splits it 0.005% to the validators and 0.005% to the treasury.
//
// It runs as a post handler so that the fee is only charged for transfers that
// actually executed, and it is a separate charge from the gas fee handled by
// the standard fee decorators.
type NetworkFeeDecorator struct {
	keeper keeper.Keeper
}

func NewNetworkFeeDecorator(k keeper.Keeper) NetworkFeeDecorator {
	return NetworkFeeDecorator{keeper: k}
}

func (d NetworkFeeDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate, success bool, next sdk.PostHandler) (sdk.Context, error) {
	if !success {
		return next(ctx, tx, simulate, success)
	}

	params, err := d.keeper.Params.Get(ctx)
	if err != nil {
		return ctx, err
	}
	if params.TreasuryAddress == "" || !params.NetworkFeeRate.IsPositive() {
		return next(ctx, tx, simulate, success)
	}

	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *banktypes.MsgSend:
			if err := d.charge(ctx, m.FromAddress, m.Amount.AmountOf(params.Denom)); err != nil {
				return ctx, err
			}
		case *banktypes.MsgMultiSend:
			for _, in := range m.Inputs {
				if err := d.charge(ctx, in.Address, in.Coins.AmountOf(params.Denom)); err != nil {
					return ctx, err
				}
			}
		}
	}

	return next(ctx, tx, simulate, success)
}

func (d NetworkFeeDecorator) charge(ctx sdk.Context, payer string, amount math.Int) error {
	if !amount.IsPositive() {
		return nil
	}

	addr, err := d.keeper.AddressCodec().StringToBytes(payer)
	if err != nil {
		return err
	}

	return d.keeper.ChargeNetworkFee(ctx, sdk.AccAddress(addr), amount)
}

var _ sdk.PostDecorator = NetworkFeeDecorator{}
