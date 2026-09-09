import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { Params } from "./params";
import { RestrictedAccount, ValidatorPerformance, ValidatorReserve } from "./tokenomics";
export declare const protobufPackage = "ferac.ferac.v1";
/** GenesisState defines the ferac module's genesis state. */
export interface GenesisState {
    /** params defines all the parameters of the module. */
    params: Params | undefined;
    /**
     * restricted_accounts are the creator/team accounts subject to the
     * 25%-of-remaining per 12 months transfer limit.
     */
    restrictedAccounts: RestrictedAccount[];
    /** validator_reserve is the state of the 4% validator reserve. */
    validatorReserve: ValidatorReserve | undefined;
    /** performances are the per-validator uptime/quality counters. */
    performances: ValidatorPerformance[];
    /**
     * reserve_funder is the genesis account whose FERAC is moved into the
     * validator reserve module account at chain launch. Minting is disabled, so
     * the reserve can only be funded from the genesis distribution.
     */
    reserveFunder: string;
}
export declare const GenesisState: MessageFns<GenesisState>;
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
