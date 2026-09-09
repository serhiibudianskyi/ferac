package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// NewRestrictedAccount creates a restricted account whose first period limit is
// the release rate applied to the full allocation.
func NewRestrictedAccount(address string, category AllocationCategory, amount math.Int, rate math.LegacyDec, periodStart int64) RestrictedAccount {
	return RestrictedAccount{
		Address:     address,
		Category:    category,
		Remaining:   amount,
		PeriodLimit: rate.MulInt(amount).TruncateInt(),
		Released:    math.ZeroInt(),
		PeriodStart: periodStart,
	}
}

// Validate performs stateless validation of a restricted account.
func (a RestrictedAccount) Validate() error {
	if a.Address == "" {
		return fmt.Errorf("restricted account address cannot be empty")
	}
	if a.Category != ALLOCATION_CATEGORY_CREATOR && a.Category != ALLOCATION_CATEGORY_TEAM {
		return fmt.Errorf("only creator and team allocations can be restricted, got %s", a.Category)
	}
	if a.Remaining.IsNil() || a.Remaining.IsNegative() {
		return fmt.Errorf("remaining of %s must be non-negative", a.Address)
	}
	if a.PeriodLimit.IsNil() || a.PeriodLimit.IsNegative() {
		return fmt.Errorf("period limit of %s must be non-negative", a.Address)
	}
	if a.Released.IsNil() || a.Released.IsNegative() {
		return fmt.Errorf("released of %s must be non-negative", a.Address)
	}
	if a.Released.GT(a.PeriodLimit) {
		return fmt.Errorf("released of %s exceeds its period limit", a.Address)
	}
	return nil
}

// Available returns how much may still leave the account in the current period.
func (a RestrictedAccount) Available() math.Int {
	available := a.PeriodLimit.Sub(a.Released)
	if available.IsNegative() {
		return math.ZeroInt()
	}
	if available.GT(a.Remaining) {
		return a.Remaining
	}
	return available
}

// Rollover advances the account to the period containing blockTime. The limit of
// a new period is always the release rate applied to the *remaining* balance, so
// an unused limit never carries over.
func (a RestrictedAccount) Rollover(blockTime, periodSeconds int64, rate math.LegacyDec) RestrictedAccount {
	if periodSeconds <= 0 || blockTime < a.PeriodStart+periodSeconds {
		return a
	}

	elapsed := (blockTime - a.PeriodStart) / periodSeconds
	a.PeriodStart += elapsed * periodSeconds
	a.PeriodLimit = rate.MulInt(a.Remaining).TruncateInt()
	a.Released = math.ZeroInt()
	return a
}

// Uptime returns the share of blocks the validator signed during the current
// epoch. A validator with no recorded blocks is treated as fully available so
// that a freshly bonded validator is not punished by the reward formula.
func (p ValidatorPerformance) Uptime() math.LegacyDec {
	if p.TotalBlocks == 0 {
		return math.LegacyOneDec()
	}
	return math.LegacyNewDec(int64(p.SignedBlocks)).Quo(math.LegacyNewDec(int64(p.TotalBlocks)))
}
