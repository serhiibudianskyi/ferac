<script setup lang="ts">
import { RouterView } from "vue-router";

import IgntHeader from "./components/IgntHeader.vue";

const items = [
  {
    label: "Portfolio",
    to: "/",
  },
  {
    label: "Data",
    to: "/data",
  },
];

const currentYear = new Date().getFullYear();
</script>

<template>
  <div class="app-container flex min-h-screen flex-col">
    <IgntHeader :nav-items="items"></IgntHeader>
    <main class="flex-1">
      <RouterView />
    </main>
    <footer class="site-header flex flex-wrap items-center gap-3 p-4 sm:p-5 justify-center ">
      © {{ currentYear }} Budianskyi S. All rights reserved.
    </footer>
    </div>
</template>
