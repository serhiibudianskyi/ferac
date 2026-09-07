package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	authKeeper    types.AuthKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistrKeeper

	Schema collections.Schema
	Params collections.Item[types.Params]
	// RestrictedAccounts holds the creator/team accounts subject to the
	// 25%-of-remaining per 12 months transfer limit.
	RestrictedAccounts collections.Map[string, types.RestrictedAccount]
	// ValidatorReserve holds the state of the 4% validator reserve.
	ValidatorReserve collections.Item[types.ValidatorReserve]
	// Performances holds the uptime/quality counters keyed by operator address.
	Performances collections.Map[string, types.ValidatorPerformance]
	// ValidatorExits holds validator exit announcements keyed by operator address.
	ValidatorExits collections.Map[string, types.ValidatorExit]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		authority:     authority,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		RestrictedAccounts: collections.NewMap(
			sb, types.RestrictedAccountsKey, "restricted_accounts",
			collections.StringKey, codec.CollValue[types.RestrictedAccount](cdc),
		),
		ValidatorReserve: collections.NewItem(
			sb, types.ValidatorReserveKey, "validator_reserve",
			codec.CollValue[types.ValidatorReserve](cdc),
		),
		Performances: collections.NewMap(
			sb, types.PerformancesKey, "performances",
			collections.StringKey, codec.CollValue[types.ValidatorPerformance](cdc),
		),
		ValidatorExits: collections.NewMap(
			sb, types.ValidatorExitsKey, "validator_exits",
			collections.StringKey, codec.CollValue[types.ValidatorExit](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// AddressCodec returns the account address codec used by the module.
func (k Keeper) AddressCodec() address.Codec {
	return k.addressCodec
}
