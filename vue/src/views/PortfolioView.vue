<template>
  <div>
    <div class="container mx-auto px-4 sm:px-6">
      <div class="grid grid-cols-1 gap-10 md:grid-cols-2 md:gap-6">
        <div class="min-w-0">
          <IgntAssets :key="address" class="mb-10" :display-limit="3" />
          <IgntTransactions :key="address" />
        </div>
        <IgntTransfer class="w-full max-w-xl mx-auto" />
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
