package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:             DefaultParams(),
		RestrictedAccounts: []RestrictedAccount{},
		ValidatorReserve: ValidatorReserve{
			Remaining:     ValidatorReserveAllocation,
			LastEpochTime: 0,
			Epoch:         0,
		},
		Performances: []ValidatorPerformance{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(gs.RestrictedAccounts))
	for _, acc := range gs.RestrictedAccounts {
		if _, ok := seen[acc.Address]; ok {
			return fmt.Errorf("duplicate restricted account %s", acc.Address)
		}
		seen[acc.Address] = struct{}{}

		if err := acc.Validate(); err != nil {
			return err
		}
	}

	if gs.ValidatorReserve.Remaining.IsNil() || gs.ValidatorReserve.Remaining.IsNegative() {
		return fmt.Errorf("validator reserve remaining must be non-negative")
	}
	if gs.ValidatorReserve.Remaining.GT(ValidatorReserveAllocation) {
		return fmt.Errorf(
			"validator reserve remaining %s exceeds the 4%% allocation %s",
			gs.ValidatorReserve.Remaining, ValidatorReserveAllocation,
		)
	}

	for _, perf := range gs.Performances {
		if perf.QualityScore.IsNil() || perf.QualityScore.IsNegative() || perf.QualityScore.GT(math.LegacyOneDec()) {
			return fmt.Errorf("quality score of %s must be within [0,1]", perf.ValidatorAddress)
		}
	}

	return nil
}
