package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// SelfTransferDecorator rejects direct bank transfers whose sender and
// recipient resolve to the same account.
type SelfTransferDecorator struct{}

func NewSelfTransferDecorator() SelfTransferDecorator {
	return SelfTransferDecorator{}
}

func (SelfTransferDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if send, ok := msg.(*banktypes.MsgSend); ok {
			from, err := sdk.AccAddressFromBech32(send.FromAddress)
			if err != nil {
				return ctx, err
			}
			to, err := sdk.AccAddressFromBech32(send.ToAddress)
			if err != nil {
				return ctx, err
			}
			if from.Equals(to) {
				return ctx, types.ErrSelfTransfer
			}
		}
	}

	return next(ctx, tx, simulate)
}

var _ sdk.AnteDecorator = SelfTransferDecorator{}