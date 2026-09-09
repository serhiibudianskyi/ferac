import { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";
import { QueryRestrictedAccountsRequest } from "./types/ferac/ferac/v1/query";
import type { SnakeCasedPropertiesDeep } from 'type-fest';
export type QueryParamsType = Record<string | number, any>;
export type FlattenObject<TValue> = CollapseEntries<CreateObjectEntries<TValue, TValue>>;
type Entry = {
    key: string;
    value: unknown;
};
type EmptyEntry<TValue> = {
    key: '';
    value: TValue;
};
type ExcludedTypes = Date | Set<unknown> | Map<unknown, unknown>;
type ArrayEncoder = `[${bigint}]`;
type EscapeArrayKey<TKey extends string> = TKey extends `${infer TKeyBefore}.${ArrayEncoder}${infer TKeyAfter}` ? EscapeArrayKey<`${TKeyBefore}${ArrayEncoder}${TKeyAfter}`> : TKey;
type CollapseEntries<TEntry extends Entry> = {
    [E in TEntry as EscapeArrayKey<E['key']>]: E['value'];
};
type CreateArrayEntry<TValue, TValueInitial> = OmitItself<TValue extends unknown[] ? {
    [k: ArrayEncoder]: TValue[number];
} : TValue, TValueInitial>;
type OmitItself<TValue, TValueInitial> = TValue extends TValueInitial ? EmptyEntry<TValue> : OmitExcludedTypes<TValue, TValueInitial>;
type OmitExcludedTypes<TValue, TValueInitial> = TValue extends ExcludedTypes ? EmptyEntry<TValue> : CreateObjectEntries<TValue, TValueInitial>;
type CreateObjectEntries<TValue, TValueInitial> = TValue extends object ? {
    [TKey in keyof TValue]-?: TKey extends string ? CreateArrayEntry<TValue[TKey], TValueInitial> extends infer TNestedValue ? TNestedValue extends Entry ? TNestedValue['key'] extends '' ? {
        key: TKey;
        value: TNestedValue['value'];
    } : {
        key: `${TKey}.${TNestedValue['key']}`;
        value: TNestedValue['value'];
    } | {
        key: TKey;
        value: TValue[TKey];
    } : never : never : never;
}[keyof TValue] : EmptyEntry<TValue>;
export type ChangeProtoToJSPrimitives<T extends object> = {
    [key in keyof T]: T[key] extends Uint8Array | Date ? string : T[key] extends object ? ChangeProtoToJSPrimitives<T[key]> : T[key];
};
export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
    /** set parameter to `true` for call `securityWorker` for this request */
    secure?: boolean;
    /** request path */
    path: string;
    /** content type of request body */
    type?: ContentType;
    /** query params */
    query?: QueryParamsType;
    /** format of response (i.e. response.json() -> format: "json") */
    format?: ResponseType;
    /** request body */
    body?: unknown;
}
export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;
export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
    securityWorker?: (securityData: SecurityDataType | null) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
    secure?: boolean;
    format?: ResponseType;
}
export declare enum ContentType {
    Json = "application/json",
    FormData = "multipart/form-data",
    UrlEncoded = "application/x-www-form-urlencoded"
}
export declare class HttpClient<SecurityDataType = unknown> {
    instance: AxiosInstance;
    private securityData;
    private securityWorker?;
    private secure?;
    private format?;
    constructor({ securityWorker, secure, format, ...axiosConfig }?: ApiConfig<SecurityDataType>);
    setSecurityData: (data: SecurityDataType | null) => void;
    private mergeRequestParams;
    private createFormData;
    request: <T = any>({ secure, path, type, query, format, body, ...params }: FullRequestParams) => Promise<AxiosResponse<T>>;
}
/**
 * @title ferac.ferac.v1
 */
export declare class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
    /**
     * QueryParams
     *
     * @tags Query
     * @name queryParams
     * @request GET:/serhiibudianskyi/ferac/ferac/v1/params
     */
    queryParams: (query?: Record<string, any>, params?: RequestParams) => Promise<AxiosResponse<{
        params: {
            denom: string;
            max_supply: string;
            restricted_release_rate: string;
            restricted_period_seconds: number;
            reserve_release_rate: string;
            reserve_period_seconds: number;
            stake_weight: string;
            uptime_weight: string;
            quality_weight: string;
            network_fee_rate: string;
            network_fee_validator_share: string;
            dex_fee_rate: string;
            max_validator_stake_share: string;
            min_self_delegation: string;
            validator_unbonding_seconds: number;
            genesis_validator_count: number;
            genesis_finality_threshold: number;
            supermajority_threshold: string;
            validator_commission_rate: string;
            treasury_address: string;
            max_validators_per_operator: number;
        };
    }, any>>;
    /**
     * QueryValidatorReserve
     *
     * @tags Query
     * @name queryValidatorReserve
     * @request GET:/serhiibudianskyi/ferac/ferac/v1/validator_reserve
     */
    queryValidatorReserve: (query?: Record<string, any>, params?: RequestParams) => Promise<AxiosResponse<{
        reserve: {
            remaining: string;
            last_epoch_time: number;
            epoch: number;
        };
        next_release: string;
    }, any>>;
    /**
     * QueryRestrictedAccount
     *
     * @tags Query
     * @name queryRestrictedAccount
     * @request GET:/serhiibudianskyi/ferac/ferac/v1/restricted_accounts/{address}
     */
    queryRestrictedAccount: (address: string, query?: Record<string, any>, params?: RequestParams) => Promise<AxiosResponse<{
        account: {
            address: string;
            category: import("./types/ferac/ferac/v1/tokenomics").AllocationCategory;
            remaining: string;
            period_limit: string;
            released: string;
            period_start: number;
        };
        available: string;
    }, any>>;
    /**
     * QueryRestrictedAccounts
     *
     * @tags Query
     * @name queryRestrictedAccounts
     * @request GET:/serhiibudianskyi/ferac/ferac/v1/restricted_accounts
     */
    queryRestrictedAccounts: (query?: Omit<FlattenObject<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QueryRestrictedAccountsRequest>>>, "">, params?: RequestParams) => Promise<AxiosResponse<{
        accounts: {
            address: string;
            category: import("./types/ferac/ferac/v1/tokenomics").AllocationCategory;
            remaining: string;
            period_limit: string;
            released: string;
            period_start: number;
        }[];
        pagination: {
            next_key: string;
            total: number;
        };
    }, any>>;
    /**
     * QueryValidatorPerformance
     *
     * @tags Query
     * @name queryValidatorPerformance
     * @request GET:/serhiibudianskyi/ferac/ferac/v1/performances/{validator_address}
     */
    queryValidatorPerformance: (validator_address: string, query?: Record<string, any>, params?: RequestParams) => Promise<AxiosResponse<{
        performance: {
            validator_address: string;
            signed_blocks: number;
            total_blocks: number;
            quality_score: string;
        };
    }, any>>;
}
export {};
