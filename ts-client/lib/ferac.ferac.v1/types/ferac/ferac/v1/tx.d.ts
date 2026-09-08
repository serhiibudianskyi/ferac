import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { Params } from "./params";
export declare const protobufPackage = "ferac.ferac.v1";
/** MsgUpdateParams is the Msg/UpdateParams request type. */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /**
     * params defines the module parameters to update.
     *
     * NOTE: All parameters must be supplied.
     */
    params: Params | undefined;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 */
export interface MsgUpdateParamsResponse {
}
/** MsgAnnounceValidatorExit is the Msg/AnnounceValidatorExit request type. */
export interface MsgAnnounceValidatorExit {
    /** sender must be the account address of the validator operator. */
    sender: string;
    /** validator_address is the operator address of the exiting validator. */
    validatorAddress: string;
}
/** MsgAnnounceValidatorExitResponse defines the response structure. */
export interface MsgAnnounceValidatorExitResponse {
    /** allowed_at is the unix time from which self-undelegation is permitted. */
    allowedAt: number;
}
/** MsgSetValidatorQuality is the Msg/SetValidatorQuality request type. */
export interface MsgSetValidatorQuality {
    authority: string;
    validatorAddress: string;
    /** quality_score must be in [0,1]. */
    qualityScore: string;
}
/** MsgSetValidatorQualityResponse defines the response structure. */
export interface MsgSetValidatorQualityResponse {
}
export declare const MsgUpdateParams: MessageFns<MsgUpdateParams>;
export declare const MsgUpdateParamsResponse: MessageFns<MsgUpdateParamsResponse>;
export declare const MsgAnnounceValidatorExit: MessageFns<MsgAnnounceValidatorExit>;
export declare const MsgAnnounceValidatorExitResponse: MessageFns<MsgAnnounceValidatorExitResponse>;
export declare const MsgSetValidatorQuality: MessageFns<MsgSetValidatorQuality>;
export declare const MsgSetValidatorQualityResponse: MessageFns<MsgSetValidatorQualityResponse>;
/** Msg defines the Msg service. */
export interface Msg {
    /**
     * UpdateParams defines a (governance) operation for updating the module
     * parameters. The authority defaults to the x/gov module account.
     */
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
    /**
     * AnnounceValidatorExit starts the validator exit notice period. A validator
     * cannot reduce its self-delegation before the notice period has elapsed.
     */
    AnnounceValidatorExit(request: MsgAnnounceValidatorExit): Promise<MsgAnnounceValidatorExitResponse>;
    /**
     * SetValidatorQuality sets the governance-controlled quality score used by
     * the 40/40/20 validator reward formula.
     */
    SetValidatorQuality(request: MsgSetValidatorQuality): Promise<MsgSetValidatorQualityResponse>;
}
export declare const MsgServiceName = "ferac.ferac.v1.Msg";
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    private readonly service;
    constructor(rpc: Rpc, opts?: {
        service?: string;
    });
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
    AnnounceValidatorExit(request: MsgAnnounceValidatorExit): Promise<MsgAnnounceValidatorExitResponse>;
    SetValidatorQuality(request: MsgSetValidatorQuality): Promise<MsgSetValidatorQualityResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P : P & {
    [K in keyof P]: Exact<P[K], I[K]>;
} & {
    [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
};
export interface MessageFns<T> {
    encode(message: T, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): T;
    fromJSON(object: any): T;
    toJSON(message: T): unknown;
    create<I extends Exact<DeepPartial<T>, I>>(base?: I): T;
    fromPartial<I extends Exact<DeepPartial<T>, I>>(object: I): T;
}
export {};
