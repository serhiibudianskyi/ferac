<template>
  <div class="container mx-auto px-4 py-10 sm:px-6">
    <header class="mb-8">
      <h1 class="text-4xl font-semibold">Network data</h1>
      <p class="mt-2 text-sm text-gray-600">
        Live tokenomics and allocation state from the Ferac chain.
      </p>
    </header>

    <div class="grid gap-6 lg:grid-cols-2">
      <section class="rounded-xl bg-gray-50 p-6">
        <h2 class="mb-5 text-xl font-semibold">Tokenomics</h2>
        <dl class="grid grid-cols-2 gap-x-6 gap-y-4 text-sm">
          <div>
            <dt class="text-gray-500">Denom</dt>
            <dd class="mt-1 font-medium">{{ params?.denom ?? "-" }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Max supply</dt>
            <dd class="mt-1 font-medium">
              {{ formatDenomAmount(params?.max_supply ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div>
            <dt class="text-gray-500">Network fee</dt>
            <dd class="mt-1 font-medium">{{ params?.network_fee_rate ?? "-" }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Validator share</dt>
            <dd class="mt-1 font-medium">{{ params?.network_fee_validator_share ?? "-" }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Genesis validators</dt>
            <dd class="mt-1 font-medium">{{ params?.genesis_validator_count ?? "-" }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Finality threshold</dt>
            <dd class="mt-1 font-medium">{{ params?.genesis_finality_threshold ?? "-" }}</dd>
          </div>
        </dl>
      </section>

      <section class="rounded-xl bg-gray-50 p-6">
        <h2 class="mb-5 text-xl font-semibold">Validator reserve</h2>
        <dl class="grid grid-cols-2 gap-x-6 gap-y-4 text-sm">
          <div>
            <dt class="text-gray-500">Remaining</dt>
            <dd class="mt-1 font-medium">
              {{ formatDenomAmount(reserve?.remaining ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div>
            <dt class="text-gray-500">Next release</dt>
            <dd class="mt-1 font-medium">
              {{ formatDenomAmount(reserveData?.next_release ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div>
            <dt class="text-gray-500">Epoch</dt>
            <dd class="mt-1 font-medium">{{ reserve?.epoch ?? "-" }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Last epoch</dt>
            <dd class="mt-1 font-medium">{{ reserve?.last_epoch_time ?? "-" }}</dd>
          </div>
        </dl>
      </section>
    </div>

    <section class="mt-6 overflow-hidden rounded-xl bg-gray-50 p-6">
      <h2 class="mb-5 text-xl font-semibold">Restricted accounts</h2>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[680px] text-left text-sm">
          <thead class="border-b border-gray-200 text-gray-500">
            <tr>
              <th class="pb-3 pr-4 font-normal">Address</th>
              <th class="pb-3 pr-4 font-normal">Category</th>
              <th class="pb-3 pr-4 text-right font-normal">Remaining</th>
              <th class="pb-3 text-right font-normal">Released</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="account in accounts"
              :key="account.address"
              class="border-b border-gray-200 last:border-0"
            >
              <td class="py-4 pr-4 font-medium">{{ account.address }}</td>
              <td class="py-4 pr-4">{{ account.category }}</td>
              <td class="py-4 pr-4 text-right">
                {{ formatDenomAmount(account.remaining, "uferac") }} FERAC
              </td>
              <td class="py-4 text-right">
                {{ formatDenomAmount(account.released, "uferac") }} FERAC
              </td>
            </tr>
            <tr v-if="!accounts.length">
              <td colspan="4" class="py-8 text-center text-gray-500">
                No restricted accounts
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="mt-6 overflow-hidden rounded-xl bg-gray-50 p-6">
      <h2 class="mb-5 text-xl font-semibold">Validator commissions</h2>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[760px] text-left text-sm">
          <thead class="border-b border-gray-200 text-gray-500">
            <tr>
              <th class="pb-3 pr-4 font-normal">Validator</th>
              <th class="pb-3 pr-4 font-normal">Operator address</th>
              <th class="pb-3 pr-4 text-right font-normal">Commission</th>
              <th class="pb-3 pr-4 text-right font-normal">Outstanding rewards</th>
              <th class="pb-3 text-right font-normal">Total visible rewards</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="validator in validatorCommissions"
              :key="validator.operatorAddress"
              class="border-b border-gray-200 last:border-0"
            >
              <td class="py-4 pr-4 font-medium">{{ validator.moniker }}</td>
              <td class="py-4 pr-4">{{ validator.operatorAddress }}</td>
              <td class="py-4 pr-4 text-right">
                {{ formatDenomAmount(validator.commission, "uferac") }} FERAC
              </td>
              <td class="py-4 pr-4 text-right">
                {{ formatDenomAmount(validator.outstanding, "uferac") }} FERAC
              </td>
              <td class="py-4 text-right font-medium">
                {{ formatDenomAmount(validator.total, "uferac") }} FERAC
              </td>
            </tr>
            <tr v-if="!validatorCommissions.length">
              <td colspan="5" class="py-8 text-center text-gray-500">
                No validator commissions
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { useQuery } from "@tanstack/vue-query";

import { useClient } from "@/composables/useClient";
import useFeracFeracV1 from "@/composables/useFeracFeracV1";
import { formatDenomAmount } from "@/def-composables/useDenom";

const { QueryParams, QueryValidatorReserve, QueryRestrictedAccounts } =
  useFeracFeracV1();
const paramsQuery = QueryParams({});
const reserveQuery = QueryValidatorReserve({});
const accountsQuery = QueryRestrictedAccounts({} as never, {}, 100);

const params = computed(() => paramsQuery.data.value?.params);
const reserveData = computed(() => reserveQuery.data.value);
const reserve = computed(() => reserveData.value?.reserve);
const accounts = computed(
  () => accountsQuery.data.value?.pages.flatMap((page) => page.accounts) ?? []
);
const client = useClient();
const commissionsQuery = useQuery({
  queryKey: ["validator-commissions"],
  refetchOnWindowFocus: true,
  queryFn: async () => {
    const validators = (await client.CosmosStakingV_1Beta_1.query.queryValidators()).data.validators ?? [];

    return Promise.all(validators.map(async (validator) => {
      const commission = (await client.CosmosDistributionV_1Beta_1.query.queryValidatorCommission(
        validator.operator_address ?? ""
      )).data.commission?.commission?.[0];
      const outstanding = (await client.CosmosDistributionV_1Beta_1.query.queryValidatorOutstandingRewards(
        validator.operator_address ?? ""
      )).data.rewards?.rewards?.[0];
      const commissionAmount = commission?.amount ?? "0";
      const outstandingAmount = outstanding?.amount ?? "0";

      return {
        moniker: validator.description?.moniker ?? validator.operator_address ?? "Unknown",
        operatorAddress: validator.operator_address ?? "",
        commission: commissionAmount,
        outstanding: outstandingAmount,
        total: (Number(commissionAmount) + Number(outstandingAmount)).toFixed(18),
      };
    }));
  },
});
const validatorCommissions = computed(() => commissionsQuery.data.value ?? []);
</script>
