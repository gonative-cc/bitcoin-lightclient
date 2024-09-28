package btclightclient

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "bitcoin-lightclient/api/bitcoinlightclient/btclightclient"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "GetHeader",
					Use:            "get-header [height]",
					Short:          "Query get-header",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "height"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "InsertHeaders",
					Use:            "insert-headers [header]",
					Short:          "Send a insert-headers tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "header"}, {ProtoField: "other"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
