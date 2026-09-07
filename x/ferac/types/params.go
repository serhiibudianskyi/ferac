package types

import (
	"fmt"

	"cosmossdk.io/math"
)

const (
	// BaseDenom is the smallest FERAC denomination.
	BaseDenom = "uferac"
	// DisplayDenom is the human readable FERAC denomination.
	DisplayDenom = "FERAC"
	// Decimals is the number of decimals of the display denom.
	Decimals = 6

	// Year is the length of the 12 month creator/team and validator reserve periods.
	Year = int64(365 * 24 * 60 * 60)
	// ValidatorUnbondingSeconds is the 14 day validator exit notice period.
	ValidatorUnbondingSeconds = int64(14 * 24 * 60 * 60)
)

// Genesis distribution of the 88,888,888 FERAC fixed supply, in uferac.
var (
	// MaxSupply is the immutable hard cap: 88,888,888 FERAC.
	MaxSupply = math.NewInt(88_888_888_000_000)

	// CreatorAllocation is 10% of the max supply.
	CreatorAllocation = math.NewInt(8_888_888_800_000)
	// TeamAllocation is 3% of the max supply.
	TeamAllocation = math.NewInt(2_666_666_640_000)
	// ValidatorReserveAllocation is 4% of the max supply.
	ValidatorReserveAllocation = math.NewInt(3_555_555_520_000)
	// CommunityAllocation is 35% of the max supply.
	CommunityAllocation = math.NewInt(31_111_110_800_000)
	// DexLiquidityAllocation is 20% of the max supply.
	DexLiquidityAllocation = math.NewInt(17_777_777_600_000)
	// TreasuryAllocation is 15% of the max supply.
	TreasuryAllocation = math.NewInt(13_333_333_200_000)
	// EcosystemAllocation is 13% of the max supply.
	EcosystemAllocation = math.NewInt(11_555_555_440_000)
)

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{
		Denom:                     BaseDenom,
		MaxSupply:                 MaxSupply,
		RestrictedReleaseRate:     math.LegacyNewDecWithPrec(25, 2),
		RestrictedPeriodSeconds:   Year,
		ReserveReleaseRate:        math.LegacyNewDecWithPrec(5, 2),
		ReservePeriodSeconds:      Year,
		StakeWeight:               math.LegacyNewDecWithPrec(40, 2),
		UptimeWeight:              math.LegacyNewDecWithPrec(40, 2),
		QualityWeight:             math.LegacyNewDecWithPrec(20, 2),
		NetworkFeeRate:            math.LegacyNewDecWithPrec(1, 4),
		NetworkFeeValidatorShare:  math.LegacyNewDecWithPrec(50, 2),
		DexFeeRate:                math.LegacyNewDecWithPrec(2, 4),
		MaxValidatorStakeShare:    math.LegacyNewDecWithPrec(20, 2),
		MinSelfDelegation:         math.NewInt(1_000_000_000), // 1,000 FERAC
		ValidatorUnbondingSeconds: ValidatorUnbondingSeconds,
		GenesisValidatorCount:     4,
		GenesisFinalityThreshold:  3,
		SupermajorityThreshold:    math.LegacyMustNewDecFromStr("0.666666666666666667"),
		ValidatorCommissionRate:   math.LegacyNewDecWithPrec(5, 2),
		TreasuryAddress:           "",
		MaxValidatorsPerOperator:  1,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	// The max supply is a hard cap: it can never be raised, not even by governance.
	if p.MaxSupply.IsNil() || !p.MaxSupply.Equal(MaxSupply) {
		return fmt.Errorf("max supply is immutable: expected %s, got %s", MaxSupply, p.MaxSupply)
	}

	rates := map[string]math.LegacyDec{
		"restricted_release_rate":     p.RestrictedReleaseRate,
		"reserve_release_rate":        p.ReserveReleaseRate,
		"network_fee_rate":            p.NetworkFeeRate,
		"network_fee_validator_share": p.NetworkFeeValidatorShare,
		"dex_fee_rate":                p.DexFeeRate,
		"max_validator_stake_share":   p.MaxValidatorStakeShare,
		"validator_commission_rate":   p.ValidatorCommissionRate,
		"supermajority_threshold":     p.SupermajorityThreshold,
		"stake_weight":                p.StakeWeight,
		"uptime_weight":               p.UptimeWeight,
		"quality_weight":              p.QualityWeight,
	}
	for name, rate := range rates {
		if rate.IsNil() || rate.IsNegative() || rate.GT(math.LegacyOneDec()) {
			return fmt.Errorf("%s must be within [0,1], got %s", name, rate)
		}
	}

	if weights := p.StakeWeight.Add(p.UptimeWeight).Add(p.QualityWeight); !weights.Equal(math.LegacyOneDec()) {
		return fmt.Errorf("reward weights must sum to 1, got %s", weights)
	}

	if p.MaxValidatorStakeShare.IsZero() {
		return fmt.Errorf("max_validator_stake_share must be positive")
	}

	// A BFT chain cannot be safe with a quorum below 2/3 of the voting power.
	if p.SupermajorityThreshold.LT(math.LegacyMustNewDecFromStr("0.666666666666666666")) {
		return fmt.Errorf("supermajority_threshold must be at least 2/3, got %s", p.SupermajorityThreshold)
	}

	if p.RestrictedPeriodSeconds <= 0 {
		return fmt.Errorf("restricted_period_seconds must be positive")
	}
	if p.ReservePeriodSeconds <= 0 {
		return fmt.Errorf("reserve_period_seconds must be positive")
	}
	if p.ValidatorUnbondingSeconds <= 0 {
		return fmt.Errorf("validator_unbonding_seconds must be positive")
	}
	if p.MinSelfDelegation.IsNil() || !p.MinSelfDelegation.IsPositive() {
		return fmt.Errorf("min_self_delegation must be positive")
	}
	if p.GenesisValidatorCount == 0 {
		return fmt.Errorf("genesis_validator_count must be positive")
	}
	if p.GenesisFinalityThreshold == 0 || p.GenesisFinalityThreshold > p.GenesisValidatorCount {
		return fmt.Errorf("genesis_finality_threshold must be within (0, genesis_validator_count]")
	}
	// 3-of-4 at genesis: the threshold must still exceed 2/3 of the genesis set.
	minGenesisThreshold := math.LegacyNewDec(int64(p.GenesisValidatorCount)).Mul(p.SupermajorityThreshold)
	if math.LegacyNewDec(int64(p.GenesisFinalityThreshold)).LTE(minGenesisThreshold) {
		return fmt.Errorf(
			"genesis_finality_threshold %d does not reach the %s supermajority of %d validators",
			p.GenesisFinalityThreshold, p.SupermajorityThreshold, p.GenesisValidatorCount,
		)
	}
	if p.MaxValidatorsPerOperator == 0 {
		return fmt.Errorf("max_validators_per_operator must be positive")
	}

	return nil
}

// QuorumFor returns the number of votes required to finalize a block for a
// validator set of the given size. While the set is at its genesis size the
// fixed 3-of-4 rule applies; beyond it the ~2/3+ supermajority rule is used.
func (p Params) QuorumFor(validators uint32) uint32 {
	if validators == 0 {
		return 0
	}
	if validators <= p.GenesisValidatorCount {
		if p.GenesisFinalityThreshold > validators {
			return validators
		}
		return p.GenesisFinalityThreshold
	}

	required := math.LegacyNewDec(int64(validators)).Mul(p.SupermajorityThreshold)
	quorum := uint32(required.TruncateInt64())
	if math.LegacyNewDec(int64(quorum)).LTE(required) {
		quorum++
	}
	return quorum
}
