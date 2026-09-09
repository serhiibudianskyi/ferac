<template>
  <div class="transfer-panel min-w-0">
  <IgntTabs
    :key="locale"
    :tabHeaderClasses="[
      'text-3xl',
      'font-semibold',
      'p-0',
      'm-0',
      'mb-2.5',
      'flex-1',
    ]"
    :tabLinkClasses="['pr-4']"
    :inactiveLinkClasses="['text-gray-400']"
    :activeLinkClasses="['text-black']"
  >
    <div class="" :tabTitle="label('send')">
      <IgntSend v-if="address" />
    </div>
    <div class="" :tabTitle="label('receive')">
      <IgntCard v-if="address">
        <template #header>
          <div
            class="flex bg-gray-100 align-center items-center justify-center w-full py-10"
          >
            <IgntQRCode :value="address" color="#000" :width="112" />
          </div>
        </template>
        <template #default>
          <div class="p-5 break-all">
            {{ address }}
          </div>
          <div class="p-5 pt-0 text-right">
            <IgntClipboard :text="address" />
          </div>
        </template>
      </IgntCard>
    </div>
  </IgntTabs>
  </div>
</template>
<script setup lang="ts">
import { IgntTabs } from "@ignt/vue-library";
import { IgntQRCode } from "@ignt/vue-library";
import { IgntCard } from "@ignt/vue-library";
import { IgntClipboard } from "@ignt/vue-library";

import { useAddress } from "@/def-composables/useAddress";
import { useLanguage } from "@/def-composables/useLanguage";

import IgntSend from "./IgntSend.vue";

const { address } = useAddress();
const { locale, t } = useLanguage();
const label = (key: Parameters<typeof t>[0]) => t(key).value;
</script>
