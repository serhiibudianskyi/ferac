package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

func (k msgServer) AnnounceValidatorExit(ctx context.Context, req *types.MsgAnnounceValidatorExit) (*types.MsgAnnounceValidatorExitResponse, error) {
	sender, err := k.addressCodec.StringToBytes(req.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(req.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	// Only the validator operator itself may announce its exit.
	if !bytes.Equal(sender, valAddr) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"only the operator of %s can announce its exit", req.ValidatorAddress,
		)
	}

	exit, err := k.AnnounceExit(ctx, req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgAnnounceValidatorExitResponse{AllowedAt: exit.AllowedAt}, nil
}

func (k msgServer) SetValidatorQuality(ctx context.Context, req *types.MsgSetValidatorQuality) (*types.MsgSetValidatorQualityResponse, error) {
	authority, err := k.addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	if _, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(req.ValidatorAddress); err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}

	if req.QualityScore.IsNil() || req.QualityScore.IsNegative() || req.QualityScore.GT(math.LegacyOneDec()) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "quality score must be within [0,1], got %s", req.QualityScore)
	}

	if err := k.Keeper.SetValidatorQuality(ctx, req.ValidatorAddress, req.QualityScore); err != nil {
		return nil, err
	}

	return &types.MsgSetValidatorQualityResponse{}, nil
}
