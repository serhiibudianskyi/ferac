import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { PageRequest, PageResponse } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params } from "./params";
import { RestrictedAccount, ValidatorPerformance, ValidatorReserve } from "./tokenomics";
export declare const protobufPackage = "ferac.ferac.v1";
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** params holds all the parameters of this module. */
    params: Params | undefined;
}
/** QueryValidatorReserveRequest is the request type for Query/ValidatorReserve. */
export interface QueryValidatorReserveRequest {
}
/** QueryValidatorReserveResponse is the response type for Query/ValidatorReserve. */
export interface QueryValidatorReserveResponse {
    reserve: ValidatorReserve | undefined;
    /** next_release is the amount that will be released at the next epoch. */
    nextRelease: string;
}
/** QueryRestrictedAccountRequest is the request type for Query/RestrictedAccount. */
export interface QueryRestrictedAccountRequest {
    address: string;
}
/** QueryRestrictedAccountResponse is the response type for Query/RestrictedAccount. */
export interface QueryRestrictedAccountResponse {
    account: RestrictedAccount | undefined;
    /** available is how much may still be transferred in the current period. */
    available: string;
}
/** QueryRestrictedAccountsRequest is the request type for Query/RestrictedAccounts. */
export interface QueryRestrictedAccountsRequest {
    pagination: PageRequest | undefined;
}
/** QueryRestrictedAccountsResponse is the response type for Query/RestrictedAccounts. */
export interface QueryRestrictedAccountsResponse {
    accounts: RestrictedAccount[];
    pagination: PageResponse | undefined;
}
/** QueryValidatorPerformanceRequest is the request type for Query/ValidatorPerformance. */
export interface QueryValidatorPerformanceRequest {
    validatorAddress: string;
}
/** QueryValidatorPerformanceResponse is the response type for Query/ValidatorPerformance. */
export interface QueryValidatorPerformanceResponse {
    performance: ValidatorPerformance | undefined;
}
export declare const QueryParamsRequest: MessageFns<QueryParamsRequest>;
export declare const QueryParamsResponse: MessageFns<QueryParamsResponse>;
export declare const QueryValidatorReserveRequest: MessageFns<QueryValidatorReserveRequest>;
export declare const QueryValidatorReserveResponse: MessageFns<QueryValidatorReserveResponse>;
export declare const QueryRestrictedAccountRequest: MessageFns<QueryRestrictedAccountRequest>;
export declare const QueryRestrictedAccountResponse: MessageFns<QueryRestrictedAccountResponse>;
export declare const QueryRestrictedAccountsRequest: MessageFns<QueryRestrictedAccountsRequest>;
export declare const QueryRestrictedAccountsResponse: MessageFns<QueryRestrictedAccountsResponse>;
export declare const QueryValidatorPerformanceRequest: MessageFns<QueryValidatorPerformanceRequest>;
export declare const QueryValidatorPerformanceResponse: MessageFns<QueryValidatorPerformanceResponse>;
/** Query defines the gRPC querier service. */
export interface Query {
    /** Parameters queries the parameters of the module. */
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** ValidatorReserve queries the remaining validator reserve and epoch state. */
    ValidatorReserve(request: QueryValidatorReserveRequest): Promise<QueryValidatorReserveResponse>;
    /** RestrictedAccount queries the creator/team transfer limit state of an account. */
    RestrictedAccount(request: QueryRestrictedAccountRequest): Promise<QueryRestrictedAccountResponse>;
    /** RestrictedAccounts lists all restricted accounts. */
    RestrictedAccounts(request: QueryRestrictedAccountsRequest): Promise<QueryRestrictedAccountsResponse>;
    /** ValidatorPerformance queries the uptime/quality counters of a validator. */
    ValidatorPerformance(request: QueryValidatorPerformanceRequest): Promise<QueryValidatorPerformanceResponse>;
}
export declare const QueryServiceName = "ferac.ferac.v1.Query";
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    private readonly service;
    constructor(rpc: Rpc, opts?: {
        service?: string;
    });
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    ValidatorReserve(request: QueryValidatorReserveRequest): Promise<QueryValidatorReserveResponse>;
    RestrictedAccount(request: QueryRestrictedAccountRequest): Promise<QueryRestrictedAccountResponse>;
    RestrictedAccounts(request: QueryRestrictedAccountsRequest): Promise<QueryRestrictedAccountsResponse>;
    ValidatorPerformance(request: QueryValidatorPerformanceRequest): Promise<QueryValidatorPerformanceResponse>;
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
