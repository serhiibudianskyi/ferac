package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

func (q queryServer) ValidatorReserve(ctx context.Context, req *types.QueryValidatorReserveRequest) (*types.QueryValidatorReserveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	reserve, err := q.k.ValidatorReserve.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	next, err := q.k.NextReserveRelease(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorReserveResponse{Reserve: reserve, NextRelease: next}, nil
}

func (q queryServer) RestrictedAccount(ctx context.Context, req *types.QueryRestrictedAccountRequest) (*types.QueryRestrictedAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	account, err := q.k.RestrictedAccounts.Get(ctx, req.Address)
	if err != nil {
		return nil, status.Error(codes.NotFound, types.ErrUnknownRestrictedAccount.Error())
	}

	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	account = account.Rollover(sdk.UnwrapSDKContext(ctx).BlockTime().Unix(), params.RestrictedPeriodSeconds, params.RestrictedReleaseRate)

	return &types.QueryRestrictedAccountResponse{Account: account, Available: account.Available()}, nil
}

func (q queryServer) RestrictedAccounts(ctx context.Context, req *types.QueryRestrictedAccountsRequest) (*types.QueryRestrictedAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	accounts, pageRes, err := query.CollectionPaginate(
		ctx, q.k.RestrictedAccounts, req.Pagination,
		func(_ string, account types.RestrictedAccount) (types.RestrictedAccount, error) {
			return account, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRestrictedAccountsResponse{Accounts: accounts, Pagination: pageRes}, nil
}

func (q queryServer) ValidatorPerformance(ctx context.Context, req *types.QueryValidatorPerformanceRequest) (*types.QueryValidatorPerformanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	perf, err := q.k.Performances.Get(ctx, req.ValidatorAddress)
	if err != nil {
		if err == collections.ErrNotFound {
			return nil, status.Error(codes.NotFound, "validator performance not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidatorPerformanceResponse{Performance: perf}, nil
}
