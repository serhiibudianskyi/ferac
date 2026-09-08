import { DeliverTxResponse, StdFee } from "@cosmjs/stargate";
import { EncodeObject, GeneratedType, OfflineSigner, Registry } from "@cosmjs/proto-signing";
import { IgniteClient } from "../client";
import { Api } from "./rest";
import { MsgUpdateParams } from "./types/ferac/ferac/v1/tx";
import { MsgAnnounceValidatorExit } from "./types/ferac/ferac/v1/tx";
import { MsgSetValidatorQuality } from "./types/ferac/ferac/v1/tx";
export { MsgUpdateParams, MsgAnnounceValidatorExit, MsgSetValidatorQuality };
type sendMsgUpdateParamsParams = {
    value: MsgUpdateParams;
    fee?: StdFee;
    memo?: string;
};
type sendMsgAnnounceValidatorExitParams = {
    value: MsgAnnounceValidatorExit;
    fee?: StdFee;
    memo?: string;
};
type sendMsgSetValidatorQualityParams = {
    value: MsgSetValidatorQuality;
    fee?: StdFee;
    memo?: string;
};
type msgUpdateParamsParams = {
    value: MsgUpdateParams;
};
type msgAnnounceValidatorExitParams = {
    value: MsgAnnounceValidatorExit;
};
type msgSetValidatorQualityParams = {
    value: MsgSetValidatorQuality;
};
export declare const registry: Registry;
interface TxClientOptions {
    addr: string;
    prefix: string;
    signer?: OfflineSigner;
}
export declare const txClient: ({ signer, prefix, addr }?: TxClientOptions) => {
    sendMsgUpdateParams({ value, fee, memo }: sendMsgUpdateParamsParams): Promise<DeliverTxResponse>;
    sendMsgAnnounceValidatorExit({ value, fee, memo }: sendMsgAnnounceValidatorExitParams): Promise<DeliverTxResponse>;
    sendMsgSetValidatorQuality({ value, fee, memo }: sendMsgSetValidatorQualityParams): Promise<DeliverTxResponse>;
    msgUpdateParams({ value }: msgUpdateParamsParams): EncodeObject;
    msgAnnounceValidatorExit({ value }: msgAnnounceValidatorExitParams): EncodeObject;
    msgSetValidatorQuality({ value }: msgSetValidatorQualityParams): EncodeObject;
};
interface QueryClientOptions {
    addr: string;
}
export declare const queryClient: ({ addr: addr }?: QueryClientOptions) => Api<unknown>;
declare class SDKModule {
    query: ReturnType<typeof queryClient>;
    tx: ReturnType<typeof txClient>;
    structure: Record<string, unknown>;
    registry: Array<[string, GeneratedType]>;
    constructor(client: IgniteClient);
    updateTX(client: IgniteClient): void;
}
declare const IgntModule: (test: IgniteClient) => {
    module: {
        FeracFeracV_1: SDKModule;
    };
    registry: [string, GeneratedType][];
};
export default IgntModule;
