<template>
  <div>
    <div class="pt-4 md:pt-8">
      <div class="text-xs text-gray-600">{{ label("sendTo") }}</div>

      <div>
        <input
          v-model="state.tx.receiver"
          class="mt-1 py-2 px-4 h-12 bg-gray-100 border-xs text-base leading-tight w-full rounded-xl outline-0"
          :class="{
            'border border-red-400':
              state.tx.receiver.length > 0 && !validReceiver,
          }"
          :placeholder="label('recipientAddress')"
        />
        <div
          v-if="state.tx.receiver.length > 0 && !validReceiver"
          class="text-xs text-red-400 mt-1"
        >
          Invalid address
        </div>
      </div>
    </div>
    <div v-if="hasAnyBalance">
      <IgntAmountSelect
        class="token-selector--main"
        :selected="state.tx.amounts"
        :balances="(balances.assets as Amount[])"
        @update="handleTxAmountUpdate"
      />
      <div class="px-4 text-xs text-gray-500">
        {{ label("estimatedFee") }}: {{ formatDenomAmount(estimatedFee, "uferac") }} FERAC
      </div>
    </div>

    <div
      class="flex text-xs font-semibold items-center mt-4 md:mt-8"
      :class="[
        {
          'cursor-pointer': hasAnyBalance,
          'text-gray-400': !hasAnyBalance,
        },
      ]"
      @click="
        (evt: MouseEvent) => {
          toggleAdvanced();
					return evt;
        }
      "
    >
      {{ label("advanced") }}
      <template v-if="hasAnyBalance">
        <IgntChevronDownIcon
          :class="{ 'rotate-180': state.advancedOpen }"
          class="ml-1"
        />
      </template>
    </div>

    <div
      v-if="state.advancedOpen && hasAnyBalance"
      class="h-3 w-full md:h-6"
    />

    <div v-if="state.advancedOpen && hasAnyBalance" class="advanced">
      <div class="text-xs pb-2">{{ label("fees") }}</div>

      <IgntAmountSelect
        class="token-selector"
        :selected="state.tx.fees"
        :balances="(balances.assets as Amount[])"
        @update="handleTxFeesUpdate"
      />

      <div class="text-xs mt-8 text-gray-600">{{ label("reference") }}</div>

      <div class="mb-4">
        <input
          v-model="state.tx.memo"
          class="mt-1 py-2 px-4 h-12 bg-gray-100 border-xs text-base leading-tight w-full rounded-xl outline-0"
          :placeholder="label('enterReference')"
        />
      </div>

      <div class="text-xs text-gray-600">{{ label("channel") }}</div>

      <div class="input-wrapper">
        <input
          v-model="state.tx.ch"
          class="mt-1 py-2 px-4 h-12 bg-gray-100 border-xs text-base leading-tight w-full rounded-xl outline-0"
          :placeholder="label('enterChannel')"
        />
      </div>
    </div>

    <div class="h-3 w-full md:h-6" />

    <div>
      <IgntButton
        style="width: 100%"
        :disabled="!ableToTx"
        @click="sendTx"
        :busy="isTxOngoing"
        >{{ label("send") }}</IgntButton
      >
      <div
        v-if="isTxError"
        class="flex items-center justify-center text-xs text-red-500 italic mt-2"
      >
        {{ state.errorMessage || "Error submitting Tx" }}
      </div>

      <div
        v-if="isTxSuccess"
        class="flex items-center justify-center text-xs text-green-500 italic mt-2"
      >
        {{ label("success") }}
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { fromBech32 } from "@cosmjs/encoding";
import { IgntButton } from "@ignt/vue-library";
import { IgntChevronDownIcon } from "@ignt/vue-library";
import { useQueryClient } from "@tanstack/vue-query";
import BigNumber from "bignumber.js";
import Long from "long";
import { computed, reactive, watch } from "vue";

import { useClient } from "@/composables/useClient";
import { useAddress } from "@/def-composables/useAddress";
import { useAssets } from "@/def-composables/useAssets";
import { formatDenomAmount } from "@/def-composables/useDenom";
import { useLanguage } from "@/def-composables/useLanguage";
import type { Amount } from "@/utils/interfaces";
import { env } from "@/env";

import IgntAmountSelect from "./IgntAmountSelect.vue";
interface TxData {
  receiver: string;
  ch: string;
  amounts: Array<Amount>;
  memo: string;
  fees: Array<Amount>;
}

enum UI_STATE {
  FRESH = 1,

  BOOTSTRAPED = 2,

  WALLET_LOCKED = 3,

  SEND = 100,
  SEND_ADD_TOKEN = 101,

  TX_SIGNING = 300,
  TX_SUCCESS = 301,
  TX_ERROR = 302,
}

interface State {
  tx: TxData;
  currentUIState: UI_STATE;
  advancedOpen: boolean;
  errorMessage: string;
}

const initialState: State = {
  tx: {
    receiver: "",
    ch: "",
    amounts: [],
    memo: "",
    fees: [],
  },
  currentUIState: UI_STATE.SEND,
  advancedOpen: false,
  errorMessage: "",
};
const state = reactive(initialState);
const client = useClient();
const queryClient = useQueryClient();
const { address } = useAddress();
const { balances } = useAssets(100);
const { t } = useLanguage();
const label = (key: Parameters<typeof t>[0]) => t(key).value;

onMounted(() => {
  const recipient = new URLSearchParams(window.location.search).get("recipient");
  if (recipient) {
    state.tx.receiver = recipient;
  }
});

const getBlockHeight = async (): Promise<number> => {
  const response = await fetch(`${env.rpcURL}status`);
  const data = (await response.json()) as {
    result?: { sync_info?: { latest_block_height?: string } };
  };
  return Number(data.result?.sync_info?.latest_block_height ?? 0);
};

const waitForNextBlock = async (height: number): Promise<void> => {
  for (let attempt = 0; attempt < 50; attempt += 1) {
    if ((await getBlockHeight()) > height) {
      return;
    }
    await new Promise((resolve) => setTimeout(resolve, 200));
  }
};

const resetTx = (): void => {
  state.tx.amounts = [];
  state.tx.receiver = "";
  state.tx.memo = "";
  state.tx.ch = "";
  state.tx.fees = [];
  state.errorMessage = "";

  state.currentUIState = UI_STATE.SEND;
};
const sendTx = async (): Promise<void> => {
  state.currentUIState = UI_STATE.TX_SIGNING;

  const fee: Array<Amount> = state.tx.fees.map((x) => ({
    denom: x.denom,
    amount: x.amount == "" ? "0" : x.amount,
  }));
  const transactionFee =
    fee.length > 0 ? fee : [{ amount: "20", denom: "uferac" }];

  const amount: Array<Amount> = state.tx.amounts.map((x) => ({
    denom: x.denom,
    amount: x.amount == "" ? "0" : x.amount,
  }));

  const memo = state.tx.memo;

  const isIBC = state.tx.ch !== "";

  let send;

  let payload: any = {
    amount,
    toAddress: state.tx.receiver,
    fromAddress: address.value,
  };

  try {
    if (isIBC) {
      payload = {
        ...payload,
        sourcePort: "transfer",
        sourceChannel: state.tx.ch,
        sender: address.value,
        receiver: state.tx.receiver,
        timeoutHeight: 0,
        timeoutTimestamp: Long.fromNumber(
          new Date().getTime() + 60000
        ).multiply(1000000),
        token: state.tx.amounts[0],
      };

      send = () =>
        client.IbcApplicationsTransferV_1.tx.sendMsgTransfer({
          value: payload,
          fee: { amount: transactionFee as Readonly<Amount>[], gas: "300000" },
          memo,
        });
    } else {
      send = () =>
        client.CosmosBankV_1Beta_1.tx.sendMsgSend({
          value: payload,
          fee: { amount: transactionFee as Readonly<Amount[]>, gas: "300000" },
          memo,
        });
    }

    const txResult = await send();

    if (txResult.code) {
      throw new Error(txResult.rawLog || `Transaction failed with code ${txResult.code}`);
    }
    await queryClient.invalidateQueries({
      predicate: ({ queryKey }) => {
        const queryType =
          typeof queryKey[0] === "string"
            ? queryKey[0]
            : (queryKey[0] as { type?: string })?.type;
        return [
          "ServiceGetTxsEvent",
          "QueryAllBalances",
          "QueryBalance",
          "QuerySpendableBalances",
          "QuerySpendableBalanceByDenom",
          "validator-commissions",
        ].includes(queryType ?? "");
      },
    });
    await queryClient.refetchQueries({
      queryKey: ["validator-commissions"],
      type: "active",
    });
    await waitForNextBlock(Number(txResult.height));
    await queryClient.refetchQueries({
      queryKey: ["validator-commissions"],
      type: "all",
    });
    resetTx();
    state.currentUIState = UI_STATE.TX_SUCCESS;
    setTimeout(() => {
      resetTx();
    }, 2500);
  } catch (e) {
    console.error(e);
    state.errorMessage = e instanceof Error ? e.message : "Error submitting Tx";
    state.currentUIState = UI_STATE.TX_ERROR;
  }
};
const toggleAdvanced = () => {
  if (hasAnyBalance.value) {
    state.advancedOpen = !state.advancedOpen;
  }
};
const handleTxAmountUpdate = (selected: Amount[]) => {
  state.tx.amounts = selected;
};
const handleTxFeesUpdate = (selected: Amount[]) => {
  state.tx.fees = selected;
};
const parseAmount = (amount: string): BigNumber => {
  return amount == "" ? new BigNumber(0) : new BigNumber(amount);
};
const hasAnyBalance = computed<boolean>(
  () =>
    balances.value.assets.length > 0 &&
    balances.value.assets.some((x) => parseAmount(x.amount ?? "0").isPositive())
);
const estimatedFee = computed(() => {
  const safeAmount = (amount: string | undefined) => {
    const parsedAmount = new BigNumber(amount || 0);
    return parsedAmount.isFinite() ? parsedAmount : new BigNumber(0);
  };
  const transferAmount = state.tx.amounts.reduce(
    (total, amount) =>
      amount.denom === "uferac"
        ? total.plus(safeAmount(amount.amount))
        : total,
    new BigNumber(0)
  );
  const gasFee = state.tx.fees
    .filter((fee) => fee.denom === "uferac")
    .reduce(
      (total, fee) => total.plus(safeAmount(fee.amount)),
      new BigNumber(0)
    );
  const baseGasFee = state.tx.fees.length > 0 ? gasFee : new BigNumber(20);

  return baseGasFee.plus(transferAmount.multipliedBy(0.0001)).toFixed(0);
});
const isTxOngoing = computed<boolean>(() => {
  return state.currentUIState === UI_STATE.TX_SIGNING;
});
const isTxSuccess = computed<boolean>(() => {
  return state.currentUIState === UI_STATE.TX_SUCCESS;
});
const isTxError = computed<boolean>(() => {
  return state.currentUIState === UI_STATE.TX_ERROR;
});
const validTxFees = computed<boolean>(() =>
  state.tx.fees.every((x) => {
    const parsedAmount = parseAmount(x.amount);

    return !parsedAmount.isNaN() && parsedAmount.isPositive();
  })
);
const validTxAmount = computed<boolean>(() => {
  return (
    state.tx.amounts.length > 0 &&
    state.tx.amounts.every((x) => {
      const parsedAmount = parseAmount(x.amount);

      return (
        !parsedAmount.isNaN() &&
        parsedAmount.isPositive() &&
        !parsedAmount.isZero()
      );
    })
  );
});
const validReceiver = computed<boolean>(() => {
  let valid: boolean;

  try {
    valid = !!fromBech32(state.tx.receiver);
  } catch {
    valid = false;
  }

  return valid;
});
const ableToTx = computed<boolean>(
  () =>
    validTxAmount.value &&
    validReceiver.value &&
    validTxFees.value &&
    !!address.value
);
watch(
  () => balances.value.assets,
  (assets) => {
    if (state.tx.amounts.length === 0 && assets.length > 0) {
      const feracBalance = assets.find((asset) => asset.denom === "uferac");
      const defaultBalance = feracBalance ?? assets[0];
      state.tx.amounts = [{ ...defaultBalance, amount: "" }];
    }
  },
  { immediate: true }
);
</script>
