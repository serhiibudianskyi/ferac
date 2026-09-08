<template>
  <header class="site-header flex flex-wrap items-center gap-3 p-4 sm:p-5">
    <nav class="header-nav flex min-w-0 flex-1 items-center justify-between gap-2">
      <ul class="header-links flex min-w-0 items-center">
        <li
          class="text-3 px-2 sm:px-4 font-normal"
          v-for="item in navItems"
          :key="item.label"
        >
          <IgntLink :item="item"></IgntLink>
        </li>
      </ul>
      <div class="shrink-0"><IgntAcc /></div>
    </nav>
  </header>
</template>
<script setup lang="ts">
import { IgntLink } from "@ignt/vue-library";
import { type PropType, ref } from "vue";

import IgntAcc from "@/components/IgntAcc.vue";

type MenuItem = {
  label: string;
  to?: string;
  href?: string;
};
const props = defineProps({
  navItems: {
    type: Array as PropType<Array<MenuItem>>,
    required: true,
  },
});
const navItems = ref(props.navItems);
</script>
