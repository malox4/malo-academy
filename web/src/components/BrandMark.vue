<script setup>
import { RouterLink } from "vue-router";
import { BRAND } from "../brand";

defineProps({
  to: { type: String, default: "/" },
  compact: { type: Boolean, default: false },
});
</script>

<template>
  <RouterLink :to="to" class="brand-mark" :class="{ 'brand-mark-compact': compact }" :aria-label="BRAND.name">
    <span class="brand-mark-analyst">{{ BRAND.name.split(" ")[0] }}</span>
    <span class="brand-mark-hall">{{ BRAND.name.split(" ")[1] }}</span>
    <span class="brand-mark-dot" aria-hidden="true" />
  </RouterLink>
</template>
