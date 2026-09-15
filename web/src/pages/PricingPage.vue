<script setup>
import { useAuth } from "../stores/auth";
import SiteNav from "../components/SiteNav.vue";
import PlanCompare from "../components/PlanCompare.vue";

const auth = useAuth();
</script>

<template>
  <div>
    <SiteNav v-if="!auth.user" />
    <div :class="auth.user ? '' : 'mx-auto max-w-6xl px-4 py-12 md:px-6'">
      <PlanCompare :cta="!auth.user" />
    </div>
  </div>
</template>
