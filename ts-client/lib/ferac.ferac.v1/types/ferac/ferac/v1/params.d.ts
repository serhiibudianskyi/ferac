import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
export declare const protobufPackage = "ferac.ferac.v1";
/** Params defines the FERAC tokenomics and consensus policy parameters. */
export interface Params {
    /** denom is the base (smallest) denomination of FERAC. */
    denom: string;
    /**
     * max_supply is the immutable hard cap expressed in the base denom.
     * 88,888,888 FERAC with 6 decimals = 88888888000000 uferac.
     */
    maxSupply: string;
    /**
     * restricted_release_rate is the share of the *remaining* balance that a
     * creator/team account may move within one restriction period (0.25).
     */
    restrictedReleaseRate: string;
    /** restricted_period_seconds is the length of one creator/team period (12 months). */
    restrictedPeriodSeconds: number;
    /**
     * reserve_release_rate is the share of the *remaining* validator reserve
     * released for rewards each reserve period (0.05).
     */
    reserveReleaseRate: string;
    /** reserve_period_seconds is the length of one validator reserve period (12 months). */
    reservePeriodSeconds: number;
    /** stake_weight is the weight of bonded stake in the reward formula (0.40). */
    stakeWeight: string;
    /** uptime_weight is the weight of uptime in the reward formula (0.40). */
    uptimeWeight: string;
    /** quality_weight is the weight of the quality score in the reward formula (0.20). */
    qualityWeight: string;
    /** network_fee_rate is the protocol fee taken from a transfer (0.0001 = 0.01%). */
    networkFeeRate: string;
    /**
     * network_fee_validator_share is the part of the network fee routed to
     * validators; the remainder goes to the treasury (0.5 => 0.005% / 0.005%).
     */
    networkFeeValidatorShare: string;
    /** dex_fee_rate is the DEX swap fee routed entirely to the treasury (0.0002 = 0.02%). */
    dexFeeRate: string;
    /**
     * max_validator_stake_share caps the share of total bonded stake a single
     * validator may hold, limiting stake concentration.
     */
    maxValidatorStakeShare: string;
    /**
     * min_self_delegation is the minimum self-bonded amount required to run a
     * validator. It raises the cost of Sybil validator creation.
     */
    minSelfDelegation: string;
    /**
     * validator_unbonding_seconds is the mandatory exit notice for a validator
     * self-undelegation (14 days). Delegators use the staking unbonding time.
     */
    validatorUnbondingSeconds: number;
    /** genesis_validator_count is the number of independent genesis validators (4). */
    genesisValidatorCount: number;
    /**
     * genesis_finality_threshold is the number of genesis validators required to
     * finalize a block while the set is still at genesis size (3 of 4).
     */
    genesisFinalityThreshold: number;
    /**
     * supermajority_threshold is the voting power fraction required once the
     * validator set grows beyond the genesis size (~2/3+).
     */
    supermajorityThreshold: string;
    /** validator_commission_rate is the fixed validator commission (0.05). */
    validatorCommissionRate: string;
    /** treasury_address receives the treasury part of network and DEX fees. */
    treasuryAddress: string;
    /**
     * max_validators_per_operator limits how many validators a single funding
     * identity may operate, mitigating Sybil validator farms.
     */
    maxValidatorsPerOperator: number;
}
export declare const Params: MessageFns<Params>;
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
