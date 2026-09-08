import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
export declare const protobufPackage = "ferac.ferac.v1";
/** AllocationCategory enumerates the genesis distribution buckets. */
export declare enum AllocationCategory {
    ALLOCATION_CATEGORY_UNSPECIFIED = 0,
    ALLOCATION_CATEGORY_CREATOR = 1,
    ALLOCATION_CATEGORY_TEAM = 2,
    ALLOCATION_CATEGORY_VALIDATOR_RESERVE = 3,
    ALLOCATION_CATEGORY_COMMUNITY = 4,
    ALLOCATION_CATEGORY_DEX_LIQUIDITY = 5,
    ALLOCATION_CATEGORY_TREASURY = 6,
    ALLOCATION_CATEGORY_ECOSYSTEM = 7,
    UNRECOGNIZED = -1
}
export declare function allocationCategoryFromJSON(object: any): AllocationCategory;
export declare function allocationCategoryToJSON(object: AllocationCategory): string;
/**
 * RestrictedAccount tracks the 25%-of-remaining / 12-month transfer limit that
 * applies to the creator and team allocations.
 */
export interface RestrictedAccount {
    address: string;
    category: AllocationCategory;
    /** remaining is the amount still under restriction. */
    remaining: string;
    /** period_limit is 25% of `remaining` measured at the start of the period. */
    periodLimit: string;
    /**
     * released is how much of period_limit has already been transferred out in
     * the current period. Unused limit never carries over.
     */
    released: string;
    /** period_start is the unix time the current 12-month period began. */
    periodStart: number;
}
/**
 * ValidatorReserve holds the state of the 4% validator reserve that pays out
 * 5% of its *remaining* balance every 12 months.
 */
export interface ValidatorReserve {
    remaining: string;
    /** last_epoch_time is the unix time of the last reserve release. */
    lastEpochTime: number;
    /** epoch counts completed reserve releases. */
    epoch: number;
}
/**
 * ValidatorPerformance accumulates the uptime and quality inputs of the
 * 40/40/20 reward formula.
 */
export interface ValidatorPerformance {
    validatorAddress: string;
    /** signed_blocks and total_blocks are counted over the current reserve epoch. */
    signedBlocks: number;
    totalBlocks: number;
    /** quality_score is a governance-set score in [0,1]. */
    qualityScore: string;
}
/**
 * ValidatorExit records a validator exit announcement. A validator may only
 * self-undelegate after the notice period has elapsed.
 */
export interface ValidatorExit {
    validatorAddress: string;
    announcedAt: number;
    allowedAt: number;
}
export declare const RestrictedAccount: MessageFns<RestrictedAccount>;
export declare const ValidatorReserve: MessageFns<ValidatorReserve>;
export declare const ValidatorPerformance: MessageFns<ValidatorPerformance>;
export declare const ValidatorExit: MessageFns<ValidatorExit>;
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
