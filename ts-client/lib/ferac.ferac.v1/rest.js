import axios from "axios";
export var ContentType;
(function (ContentType) {
    ContentType["Json"] = "application/json";
    ContentType["FormData"] = "multipart/form-data";
    ContentType["UrlEncoded"] = "application/x-www-form-urlencoded";
})(ContentType || (ContentType = {}));
export class HttpClient {
    constructor({ securityWorker, secure, format, ...axiosConfig } = {}) {
        this.securityData = null;
        this.setSecurityData = (data) => {
            this.securityData = data;
        };
        this.request = async ({ secure, path, type, query, format, body, ...params }) => {
            const secureParams = ((typeof secure === "boolean" ? secure : this.secure) &&
                this.securityWorker &&
                (await this.securityWorker(this.securityData))) ||
                {};
            const requestParams = this.mergeRequestParams(params, secureParams);
            const responseFormat = (format && this.format) || void 0;
            if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
                requestParams.headers.common = { Accept: "*/*" };
                requestParams.headers.post = {};
                requestParams.headers.put = {};
                body = this.createFormData(body);
            }
            return this.instance.request({
                ...requestParams,
                headers: {
                    ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
                    ...(requestParams.headers || {}),
                },
                params: query,
                responseType: responseFormat,
                data: body,
                url: path,
            });
        };
        this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
        this.secure = secure;
        this.format = format;
        this.securityWorker = securityWorker;
    }
    mergeRequestParams(params1, params2) {
        return {
            ...this.instance.defaults,
            ...params1,
            ...(params2 || {}),
            headers: {
                ...(this.instance.defaults.headers),
                ...(params1.headers || {}),
                ...((params2 && params2.headers) || {}),
            },
        };
    }
    createFormData(input) {
        return Object.keys(input || {}).reduce((formData, key) => {
            const property = input[key];
            formData.append(key, property instanceof Blob
                ? property
                : typeof property === "object" && property !== null
                    ? JSON.stringify(property)
                    : `${property}`);
            return formData;
        }, new FormData());
    }
}
/**
 * @title ferac.ferac.v1
 */
export class Api extends HttpClient {
    constructor() {
        super(...arguments);
        /**
         * QueryParams
         *
         * @tags Query
         * @name queryParams
         * @request GET:/serhiibudianskyi/ferac/ferac/v1/params
         */
        this.queryParams = (query, params = {}) => this.request({
            path: `/serhiibudianskyi/ferac/ferac/v1/params`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * QueryValidatorReserve
         *
         * @tags Query
         * @name queryValidatorReserve
         * @request GET:/serhiibudianskyi/ferac/ferac/v1/validator_reserve
         */
        this.queryValidatorReserve = (query, params = {}) => this.request({
            path: `/serhiibudianskyi/ferac/ferac/v1/validator_reserve`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * QueryRestrictedAccount
         *
         * @tags Query
         * @name queryRestrictedAccount
         * @request GET:/serhiibudianskyi/ferac/ferac/v1/restricted_accounts/{address}
         */
        this.queryRestrictedAccount = (address, query, params = {}) => this.request({
            path: `/serhiibudianskyi/ferac/ferac/v1/restricted_accounts/${address}`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * QueryRestrictedAccounts
         *
         * @tags Query
         * @name queryRestrictedAccounts
         * @request GET:/serhiibudianskyi/ferac/ferac/v1/restricted_accounts
         */
        this.queryRestrictedAccounts = (query, params = {}) => this.request({
            path: `/serhiibudianskyi/ferac/ferac/v1/restricted_accounts`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * QueryValidatorPerformance
         *
         * @tags Query
         * @name queryValidatorPerformance
         * @request GET:/serhiibudianskyi/ferac/ferac/v1/performances/{validator_address}
         */
        this.queryValidatorPerformance = (validator_address, query, params = {}) => this.request({
            path: `/serhiibudianskyi/ferac/ferac/v1/performances/${validator_address}`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
    }
}
