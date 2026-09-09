<template>
  <div>
    <div class="container mx-auto px-4 sm:px-6">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6">
        <IgntAssets
          :key="address"
          class="min-w-0 mb-4 md:mb-10"
          :display-limit="3"
        />
        <IgntTransfer class="w-full max-w-xl mx-auto" />
        <IgntTransactions :key="address" class="min-w-0" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAddress } from "@/def-composables/useAddress";

import IgntAssets from "../components/IgntAssets.vue";
import IgntTransactions from "../components/IgntTransactions.vue";
import IgntTransfer from "../components/IgntTransfer.vue";
const { address } = useAddress();
</script>
