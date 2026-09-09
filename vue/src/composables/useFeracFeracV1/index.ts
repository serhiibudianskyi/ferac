/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions, type InfiniteData  } from "@tanstack/vue-query";
import { useClient } from '../useClient';

export default function useFeracFeracV_1() {
  const client = useClient();

  type QueryParamsMethod = typeof client.FeracFeracV_1.query.queryParams;
  type QueryParamsData = Awaited<ReturnType<QueryParamsMethod>>["data"];
  const QueryParams = ( options: Partial<UseQueryOptions<QueryParamsData>>) => {
    const key = { type: 'QueryParams',  };    
    return useQuery<QueryParamsData>({ queryKey: [key], queryFn: async () => {
      const res = await client.FeracFeracV_1.query.queryParams();
        return res.data;
    }, ...options});
  }
  

  type QueryValidatorReserveMethod = typeof client.FeracFeracV_1.query.queryValidatorReserve;
  type QueryValidatorReserveData = Awaited<ReturnType<QueryValidatorReserveMethod>>["data"];
  const QueryValidatorReserve = ( options: Partial<UseQueryOptions<QueryValidatorReserveData>>) => {
    const key = { type: 'QueryValidatorReserve',  };    
    return useQuery<QueryValidatorReserveData>({ queryKey: [key], queryFn: async () => {
      const res = await client.FeracFeracV_1.query.queryValidatorReserve();
        return res.data;
    }, ...options});
  }
  

  type QueryRestrictedAccountMethod = typeof client.FeracFeracV_1.query.queryRestrictedAccount;
  type QueryRestrictedAccountData = Awaited<ReturnType<QueryRestrictedAccountMethod>>["data"];
  const QueryRestrictedAccount = (address: string,  options: Partial<UseQueryOptions<QueryRestrictedAccountData>>) => {
    const key = { type: 'QueryRestrictedAccount',  address };    
    return useQuery<QueryRestrictedAccountData>({ queryKey: [key], queryFn: async () => {
      const { address } = key
      const res = await client.FeracFeracV_1.query.queryRestrictedAccount(address);
        return res.data;
    }, ...options});
  }
  
  type QueryRestrictedAccountsMethod = typeof client.FeracFeracV_1.query.queryRestrictedAccounts;
  type QueryRestrictedAccountsData = Awaited<ReturnType<QueryRestrictedAccountsMethod>>["data"] & { pageParam: number };
  const QueryRestrictedAccounts = (query:  NonNullable<Parameters<QueryRestrictedAccountsMethod>[0]>, options:  Partial<UseInfiniteQueryOptions<QueryRestrictedAccountsData, unknown, InfiniteData<QueryRestrictedAccountsData,number>, Array<string | unknown>, number>> , perPage: number) => {
    const key = { type: 'QueryRestrictedAccounts', query };    
    return useInfiniteQuery<QueryRestrictedAccountsData, unknown, InfiniteData<QueryRestrictedAccountsData,number>, Array<string | unknown>, number>({ queryKey: [key], queryFn: async (context: {pageParam?: number}) => {
      const { pageParam=1 } = context;
      const {query } = key

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      const res = await client.FeracFeracV_1.query.queryRestrictedAccounts(query ?? undefined);
        return { ...res.data, pageParam }; 
    }, ...options,
      initialPageParam: 1,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  

  type QueryValidatorPerformanceMethod = typeof client.FeracFeracV_1.query.queryValidatorPerformance;
  type QueryValidatorPerformanceData = Awaited<ReturnType<QueryValidatorPerformanceMethod>>["data"];
  const QueryValidatorPerformance = (validator_address: string,  options: Partial<UseQueryOptions<QueryValidatorPerformanceData>>) => {
    const key = { type: 'QueryValidatorPerformance',  validator_address };    
    return useQuery<QueryValidatorPerformanceData>({ queryKey: [key], queryFn: async () => {
      const { validator_address } = key
      const res = await client.FeracFeracV_1.query.queryValidatorPerformance(validator_address);
        return res.data;
    }, ...options});
  }
  
  return {QueryParams,QueryValidatorReserve,QueryRestrictedAccount,QueryRestrictedAccounts,QueryValidatorPerformance,
  }
}
