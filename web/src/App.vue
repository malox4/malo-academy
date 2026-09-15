<script setup>
import { computed } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";
import { useAuth } from "./stores/auth";
import BrandMark from "./components/BrandMark.vue";

const route = useRoute();
const auth = useAuth();
const publicPage = computed(() => route.meta.public && !auth.user);
const shortName = computed(() => {
  const n = String(auth.user?.name || "").trim();
  if (!n) return "Я";
  const parts = n.split(/\s+/);
  return parts[0];
});
</script>

<template>
  <div class="paper-wash" aria-hidden="true" />
  <div v-if="publicPage" class="shell min-h-dvh">
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </div>
  <div v-else class="shell min-h-dvh has-dock">
    <header class="site-bar sticky top-0 z-30">
      <div class="mx-auto flex h-12 max-w-6xl items-center gap-2 px-3 sm:h-14 sm:gap-3 sm:px-4 md:px-6">
        <BrandMark to="/hall" compact />
        <nav class="ml-2 hidden items-center gap-1 text-[13px] font-medium md:flex">
          <RouterLink to="/hall" class="nav-link">Зал</RouterLink>
          <RouterLink to="/boards" class="nav-link">Доски</RouterLink>
          <RouterLink to="/pet" class="nav-link">Пет</RouterLink>
          <RouterLink to="/practice" class="nav-link">Практика</RouterLink>
          <RouterLink to="/materials" class="nav-link">Материалы</RouterLink>
          <RouterLink to="/interview" class="nav-link">Собес</RouterLink>
          <RouterLink to="/live" class="nav-link">Live собес</RouterLink>
          <RouterLink to="/pricing" class="nav-link">Сравнение</RouterLink>
          <RouterLink v-if="auth.user?.role === 'admin'" to="/admin" class="nav-link text-accent">Журнал</RouterLink>
        </nav>
        <div class="ml-auto flex min-w-0 items-center gap-2 text-[12px]">
          <span v-if="auth.trialActive" class="stamp hidden sm:inline">Пробный PRO · {{ auth.trialLabel }}</span>
          <span v-else-if="auth.isPro" class="hidden font-mono text-[11px] text-mute sm:inline">PRO</span>
          <RouterLink to="/profile" class="btn btn-ghost max-w-[7.5rem] truncate py-1 px-2.5 text-[12px] sm:max-w-[12rem] sm:px-3">{{ shortName }}</RouterLink>
        </div>
      </div>
    </header>
    <main class="mx-auto max-w-6xl px-3 py-5 sm:px-4 sm:py-8 md:px-6 md:py-10">
      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <nav class="dock md:hidden" aria-label="Разделы зала">
      <RouterLink to="/hall" exact-active-class="is-on">Зал</RouterLink>
      <RouterLink to="/pet" active-class="is-on">Пет</RouterLink>
      <RouterLink to="/practice" active-class="is-on">Практика</RouterLink>
      <RouterLink to="/interview" active-class="is-on">Собес</RouterLink>
      <RouterLink to="/profile" active-class="is-on">Я</RouterLink>
    </nav>
  </div>
</template>
