package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	feracante "github.com/serhiibudianskyi/ferac/x/ferac/ante"
)

// setFeracHandlers replaces the default ante handler with one that also
// enforces the validator decentralization rules, and installs a post handler
// that collects the FERAC network fee for transfers that actually executed.
func (app *App) setFeracHandlers() error {
	base, err := ante.NewAnteHandler(ante.HandlerOptions{
		AccountKeeper:   app.AuthKeeper,
		BankKeeper:      app.BankKeeper,
		SignModeHandler: app.txConfig.SignModeHandler(),
		FeegrantKeeper:  app.FeeGrantKeeper,
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	})
	if err != nil {
		return fmt.Errorf("failed to create ante handler: %w", err)
	}

	validatorPolicy := sdk.ChainAnteDecorators(feracante.NewValidatorPolicyDecorator(app.FeracKeeper))

	app.SetAnteHandler(func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		ctx, err := base(ctx, tx, simulate)
		if err != nil {
			return ctx, err
		}

		return validatorPolicy(ctx, tx, simulate)
	})

	app.SetPostHandler(sdk.ChainPostDecorators(feracante.NewNetworkFeeDecorator(app.FeracKeeper)))

	return nil
}
