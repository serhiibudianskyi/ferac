package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/ferac module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	// ErrTransferLimitExceeded is returned when a creator/team account tries to
	// move more than 25% of its remaining restricted balance within a period.
	ErrTransferLimitExceeded = errors.Register(ModuleName, 1101, "restricted transfer limit exceeded for the current period")
	// ErrExitNoticeRequired is returned when a validator self-undelegates without
	// having served the mandatory exit notice period.
	ErrExitNoticeRequired = errors.Register(ModuleName, 1102, "validator exit notice period has not elapsed")
	// ErrStakeConcentration is returned when a delegation would push a validator
	// beyond the maximum allowed share of total bonded stake.
	ErrStakeConcentration = errors.Register(ModuleName, 1103, "delegation exceeds the maximum validator stake share")
	// ErrSybilValidator is returned when an operator would control more
	// validators than allowed.
	ErrSybilValidator = errors.Register(ModuleName, 1104, "operator already controls the maximum number of validators")
	// ErrInvalidCommission is returned when a validator sets a commission rate
	// different from the protocol rate.
	ErrInvalidCommission = errors.Register(ModuleName, 1105, "validator commission rate does not match the protocol rate")
	// ErrInsufficientSelfDelegation is returned when a validator's self bond is
	// below the protocol minimum.
	ErrInsufficientSelfDelegation = errors.Register(ModuleName, 1106, "self delegation is below the protocol minimum")
	// ErrMaxSupplyExceeded is returned when an operation would push the total
	// supply above the immutable hard cap.
	ErrMaxSupplyExceeded = errors.Register(ModuleName, 1107, "operation would exceed the immutable max supply")
	// ErrUnknownRestrictedAccount is returned when querying a non restricted account.
	ErrUnknownRestrictedAccount = errors.Register(ModuleName, 1108, "account is not subject to transfer restrictions")
)
