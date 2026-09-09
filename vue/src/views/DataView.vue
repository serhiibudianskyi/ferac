<template>
  <div class="container mx-auto px-4 py-10 sm:px-6">
    <header class="mb-8">
      <h1 class="text-4xl font-semibold">{{ label("networkData") }}</h1>
      <p class="mt-2 text-sm text-gray-600">
        {{ label("liveTokenomics") }}
      </p>
    </header>

    <div class="grid gap-6 lg:grid-cols-2">
      <section class="rounded-xl bg-gray-50 p-6">
        <h2 class="mb-5 text-xl font-semibold">{{ label("tokenomics") }}</h2>
        <dl class="grid grid-cols-1 gap-y-4 text-sm sm:grid-cols-2 sm:gap-x-6">
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("denom") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ params?.denom ?? "-" }}</dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("maxSupply") }}</dt>
            <dd class="mt-1 break-words font-medium">
              {{ formatDenomAmount(params?.max_supply ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("networkFee") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ params?.network_fee_rate ?? "-" }}</dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("validatorShare") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ params?.network_fee_validator_share ?? "-" }}</dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("genesisValidators") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ params?.genesis_validator_count ?? "-" }}</dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("finalityThreshold") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ params?.genesis_finality_threshold ?? "-" }}</dd>
          </div>
        </dl>
      </section>

      <section class="rounded-xl bg-gray-50 p-6">
        <h2 class="mb-5 text-xl font-semibold">{{ label("validatorReserve") }}</h2>
        <dl class="grid grid-cols-1 gap-y-4 text-sm sm:grid-cols-2 sm:gap-x-6">
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("remaining") }}</dt>
            <dd class="mt-1 break-words font-medium">
              {{ formatDenomAmount(reserve?.remaining ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("nextRelease") }}</dt>
            <dd class="mt-1 break-words font-medium">
              {{ formatDenomAmount(reserveData?.next_release ?? "0", "uferac") }} FERAC
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("epoch") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ reserve?.epoch ?? "-" }}</dd>
          </div>
          <div class="min-w-0">
            <dt class="text-gray-500">{{ label("lastEpoch") }}</dt>
            <dd class="mt-1 break-words font-medium">{{ reserve?.last_epoch_time ?? "-" }}</dd>
          </div>
        </dl>
      </section>
    </div>

    <section class="mt-6 overflow-hidden rounded-xl bg-gray-50 p-6">
      <h2 class="mb-5 text-xl font-semibold">{{ label("restrictedAccounts") }}</h2>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[680px] text-left text-sm">
          <thead class="border-b border-gray-200 text-gray-500">
            <tr>
              <th class="pb-3 pr-4 font-normal">{{ label("address") }}</th>
              <th class="pb-3 pr-4 font-normal">{{ label("category") }}</th>
              <th class="pb-3 pr-4 text-right font-normal">{{ label("remaining") }}</th>
              <th class="pb-3 text-right font-normal">{{ label("released") }}</th>
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
                {{ label("noRestrictedAccounts") }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="mt-6 overflow-hidden rounded-xl bg-gray-50 p-6">
      <h2 class="mb-5 text-xl font-semibold">{{ label("validatorCommissions") }}</h2>
      <div class="overflow-x-auto">
        <table class="w-full min-w-[760px] text-left text-sm">
          <thead class="border-b border-gray-200 text-gray-500">
            <tr>
              <th class="pb-3 pr-4 font-normal">{{ label("validator") }}</th>
              <th class="pb-3 pr-4 font-normal">{{ label("operatorAddress") }}</th>
              <th class="pb-3 pr-4 text-right font-normal">{{ label("commission") }}</th>
              <th class="pb-3 pr-4 text-right font-normal">{{ label("outstandingRewards") }}</th>
              <th class="pb-3 text-right font-normal">{{ label("totalVisibleRewards") }}</th>
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
                {{ label("noValidatorCommissions") }}
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
import { useLanguage } from "@/def-composables/useLanguage";

const { QueryParams, QueryValidatorReserve, QueryRestrictedAccounts } =
  useFeracFeracV1();
const { t } = useLanguage();
const label = (key: Parameters<typeof t>[0]) => t(key).value;
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
