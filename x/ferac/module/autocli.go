package ferac

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/serhiibudianskyi/ferac/x/ferac/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ValidatorReserve",
					Use:       "validator-reserve",
					Short:     "Shows the remaining validator reserve and the next yearly release",
				},
				{
					RpcMethod:      "RestrictedAccount",
					Use:            "restricted-account [address]",
					Short:          "Shows the creator/team transfer limit state of an account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "RestrictedAccounts",
					Use:       "restricted-accounts",
					Short:     "Lists all accounts subject to transfer restrictions",
				},
				{
					RpcMethod:      "ValidatorPerformance",
					Use:            "validator-performance [validator-address]",
					Short:          "Shows the uptime and quality counters of a validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validator_address"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "SetValidatorQuality",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "AnnounceValidatorExit",
					Use:            "announce-validator-exit [validator-address]",
					Short:          "Starts the 14 day validator exit notice period",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "validator_address"}},
				},
			},
		},
	}
}
