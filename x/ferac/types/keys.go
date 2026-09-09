package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "ferac"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// ReserveName is the module account holding the 4% validator reserve.
	ReserveName = "ferac_validator_reserve"

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("p_ferac")
	// RestrictedAccountsKey is the prefix for creator/team restricted accounts.
	RestrictedAccountsKey = collections.NewPrefix("ra_ferac")
	// ValidatorReserveKey is the prefix for the validator reserve state.
	ValidatorReserveKey = collections.NewPrefix("vr_ferac")
	// PerformancesKey is the prefix for per-validator uptime/quality counters.
	PerformancesKey = collections.NewPrefix("vp_ferac")
	// ValidatorExitsKey is the prefix for validator exit announcements.
	ValidatorExitsKey = collections.NewPrefix("ve_ferac")
)
