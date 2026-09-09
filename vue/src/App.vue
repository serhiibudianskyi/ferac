<script setup lang="ts">
import { RouterView } from "vue-router";
import { computed } from "vue";

import IgntHeader from "./components/IgntHeader.vue";
import { type Locale, useLanguage } from "./def-composables/useLanguage";

const { locale, setLocale, t } = useLanguage();
const changeLocale = (event: Event) => {
  setLocale((event.target as HTMLSelectElement).value as Locale);
};
const items = computed(() => [
  {
    label: t("portfolio").value,
    to: "/",
  },
  {
    label: t("data").value,
    to: "/data",
  },
]);

const currentYear = new Date().getFullYear();
</script>

<template>
  <div class="app-container flex min-h-screen flex-col">
    <IgntHeader :nav-items="items"></IgntHeader>
    <main class="flex-1">
      <RouterView />
    </main>
    <footer class="site-header flex flex-wrap items-center justify-center gap-3 p-4 sm:p-5">
      <label class="sr-only" for="language-select">Language</label>
      <select
        id="language-select"
        :value="locale"
        class="rounded-md border border-gray-200 bg-white px-2 py-1 text-xs font-semibold"
        aria-label="Language"
        @change="changeLocale"
      >
        <option value="en">EN</option>
        <option value="ru">RU</option>
        <option value="uk">UK</option>
      </select>
      <span>© {{ currentYear }} Budianskyi S. {{ t("allRightsReserved").value }}</span>
    </footer>
    </div>
</template>
