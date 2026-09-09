<template>
  <section>
    <header class="flex items-center justify-between">
      <h2 class="text-3xl text-black font-semibold p-0 m-0 mb-2.5 flex-1">
        {{ label("transactions") }}
      </h2>
    </header>
    <table class="transaction-table table-auto w-full" v-if="txs.length > 0">
      <tr v-for="(tx, i) in txs" :key="tx.txhash + '_' + i">
        <td class="transaction-row flex min-w-0 text-xs py-1">
          <div
            :class="{
              'rotate-180 text-green-500': tx.type == 'received',
              'text-error': tx.type == 'sent',
            }"
            class="text-2xl w-10 h-10 rounded-sm bg-gray-200 flex items-center justify-center mr-2"
          >
            <IgntTxArrowIcon />
          </div>
          <div class="transaction-meta flex min-w-0 flex-col justify-between flex-1">
            <div class="font-medium">
              {{ shortenHash(tx.txhash ?? "") }}
              <span class="font-bold text-warning">{{
                tx.ibc ? "IBC" : ""
              }}</span>
            </div>
              <div
                :class="tx.code === 0 ? 'text-green-700' : 'text-red-600'"
                class="font-medium"
              >
                {{ tx.code === 0 ? label("txSuccess") : `${label("txFailed")} (code ${tx.code})` }}
              </div>
            <div class="opacity-60">
              {{ formatTimestamp(tx.timestamp) }}
            </div>
          </div>
          <div class="transaction-details flex min-w-0 max-w-[48%] flex-col justify-between items-end">
            <div class="font-medium text-right text-xs text-gray-600 inline">
              <span
                v-for="(amount, index) in tx.amount"
                :key="tx.txhash + '_' + i + '_' + index"
                :class="{
                  'bg-green-200': tx.type == 'received',
                  'bg-red-200': tx.type == 'sent',
                }"
                class="p-1 rounded-md"
              >
                {{ tx.type == "received" ? "+" : "-" }}
                {{ formatDenomAmount(amount.amount, amount.denom) }}
                <IgntDenom :denom="amount.denom ?? ''" />
              </span>
            </div>
            <div class="opacity-60">
              <template v-if="tx.type == 'received'">
                {{label("from")}}: {{ tx.sender }}
              </template>
              <template v-else> {{label("to")}}: {{ tx.receiver }}</template>
            </div>
            <div v-if="tx.code !== 0" class="max-w-full text-right text-red-600">
              {{ tx.rawLog }}
            </div>
          </div>
        </td>
      </tr>
    </table>
    <div
      v-else
      class="text-left text-black opacity-75 text-md font-normal py-8"
    >
      {{ label("transactionEmpty") }}
    </div>
    <div
      v-if="hasMoreReceived || hasMoreSent"
      class="shadow-std flex items-center justify-center w-40 rounded-full text-sm font-medium mx-auto inset-x-0 py-2"
      @click="onShowMore"
    >
      {{ label("showMore") }}

      <IgntArrowIcon class="ml-2" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { IgntTxArrowIcon } from "@ignt/vue-library";
import { IgntArrowIcon } from "@ignt/vue-library";
import { computed } from "vue";

import { formatDenomAmount } from "@/def-composables/useDenom";
import { useTransactions } from "@/def-composables/useTransactions";
import { useLanguage } from "@/def-composables/useLanguage";

import IgntDenom from "./IgntDenom.vue";

const { transferTxs, hasMoreReceived, hasMoreSent, fetchReceived, fetchSent } =
  useTransactions();
const { locale, t } = useLanguage();
const label = (key: Parameters<typeof t>[0]) => t(key).value;
const formatTimestamp = (timestamp: string) => {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) {
    return timestamp;
  }

  const locales = {
    en: "en-US",
    ru: "ru-RU",
    uk: "uk-UA",
  } as const;

  return new Intl.DateTimeFormat(locales[locale.value], {
    dateStyle: "long",
    timeStyle: "short",
  }).format(date);
};
const txs = computed(() => {
  return transferTxs.value?.map((x) => normalizeTX(x)) ?? [];
});
const onShowMore = () => {
  if (hasMoreSent?.value) {
    fetchSent();
  }
  if (hasMoreReceived?.value) {
    fetchReceived();
  }
};
const shortenHash = (hash: string) => {
  return hash.slice(0, 6) + "..." + hash.slice(-6);
};
const normalizeTX = (tx: any) => {
  const normalized = {
    ibc: false,
    sender: "",
    receiver: "",
    txhash: "",
    timestamp: "",
    code: 0,
    rawLog: "",
    type: "",
    amount: [
      {
        amount: "",
        denom: "",
      },
    ],
  };
  const isIBC = (tx.tx.body.messages[0]["@type"] as string).includes(
    "ibc.applications.transfer.v1.MsgTransfer"
  );
  const isBankTransfer = (tx.tx.body.messages[0]["@type"] as string).includes(
    "cosmos.bank.v1beta1.MsgSend"
  );
  if (isIBC) {
    normalized.ibc = true;
    normalized.sender = tx.tx.body.messages[0].sender;
    normalized.receiver = tx.tx.body.messages[0].receiver;
    normalized.amount = [tx.tx.body.messages[0].token];
  } else if (isBankTransfer) {
    normalized.sender = tx.tx.body.messages[0].from_address;
    normalized.receiver = tx.tx.body.messages[0].to_address;
    normalized.amount = tx.tx.body.messages[0].amount;
  }
  normalized.txhash = tx.txhash;
  normalized.timestamp = tx.timestamp;
  normalized.code = Number(tx.code ?? 0);
  normalized.rawLog = tx.raw_log ?? "";
  normalized.type = tx.type;

  return normalized;
};
</script>
