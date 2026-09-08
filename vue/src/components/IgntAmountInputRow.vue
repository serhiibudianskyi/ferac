<template>
  <div>
    <IgntDenom :denom="amount.denom" modifier="avatar" class="z-10" />
    <div class="flex flex-col justify-between ml-4 z-10">
      <div class="font-semibold">
        <IgntDenom :denom="amount.denom" />
      </div>

      <div
        class="text-xs"
        :class="{
          error: !hasEnoughBalance,
        }"
      >
        {{ formatDenomAmount(balance?.amount ?? "0", props.amount.denom) }}
        available
      </div>
    </div>

    <div class="flex-1 w-full h-full">
      <IgntAmountInput
        :max-decimals="amount.denom === 'uferac' ? 6 : 0"
        class="absolute w-full left-0 text-right h-full top-0 outline-0 focus:bg-gray-100 text-3xl font-medium rounded-lg px-4"
        @update="handleChange"
      />

      <div class="focus-background"></div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { IgntAmountInput } from "@ignt/vue-library";
import BigNumber from "bignumber.js";
import { computed, type PropType, ref } from "vue";

import { useAsset } from "../def-composables/useAsset";
import { formatDenomAmount } from "../def-composables/useDenom";
import type { Amount } from "../utils/interfaces";
import IgntDenom from "./IgntDenom.vue";

const props = defineProps({
  amount: {
    type: Object as PropType<Amount>,
    required: true,
  },
});
const { balance } = useAsset(props.amount.denom);

const emit = defineEmits(["change"]);
const value = ref(
  new BigNumber(
    props.amount.amount != ""
      ? props.amount.denom === "uferac"
        ? new BigNumber(props.amount.amount).dividedBy(1_000_000)
        : props.amount.amount
      : 0
  )
);
const hasEnoughBalance = computed(() => {
  const balanceBN = new BigNumber(balance.value?.amount ?? 0);
  const requestedAmount =
    props.amount.denom === "uferac"
      ? value.value.multipliedBy(1_000_000)
      : value.value;
  if (requestedAmount.isPositive()) {
    return balanceBN.gte(requestedAmount);
  } else {
    return true;
  }
});
const handleChange = (amount: BigNumber) => {
  value.value = amount;
  if (hasEnoughBalance.value) {
    const rawAmount =
      props.amount.denom === "uferac"
        ? amount.multipliedBy(1_000_000).toFixed(0)
        : amount.toString();
    emit("change", rawAmount);
  }
};
</script>
