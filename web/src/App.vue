<script setup>
import { computed } from "vue";
import { RouterLink, RouterView, useRoute } from "vue-router";
import { useAuth } from "./stores/auth";
import BrandMark from "./components/BrandMark.vue";

const route = useRoute();
const auth = useAuth();
const publicPage = computed(() => route.meta.public && !auth.user);
</script>

<template>
  <div class="paper-wash" aria-hidden="true" />
  <div v-if="publicPage" class="shell min-h-screen">
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </div>
  <div v-else class="shell min-h-screen pb-24 md:pb-0">
    <header class="site-bar sticky top-0 z-30">
      <div class="mx-auto flex h-14 max-w-6xl items-center gap-3 px-4 md:px-6">
        <BrandMark to="/hall" />
        <nav class="ml-4 hidden items-center gap-1 text-[13px] font-medium md:flex">
          <RouterLink to="/hall" class="nav-link">Зал</RouterLink>
          <RouterLink to="/boards" class="nav-link">Доски</RouterLink>
          <RouterLink to="/pet" class="nav-link">Пет</RouterLink>
          <RouterLink to="/practice" class="nav-link">Практика</RouterLink>
          <RouterLink to="/materials" class="nav-link">Материалы</RouterLink>
          <RouterLink to="/interview" class="nav-link">Собес</RouterLink>
          <RouterLink to="/live" class="nav-link">Live собес</RouterLink>
          <RouterLink to="/pricing" class="nav-link">Сравнение</RouterLink>
          <RouterLink v-if="auth.user?.role === 'admin'" to="/admin" class="nav-link text-accent">Касса</RouterLink>
        </nav>
        <div class="ml-auto flex items-center gap-3 text-[12px]">
          <span
            v-if="auth.trialActive"
            class="stamp hidden sm:inline"
          >
            Пробный PRO · {{ auth.trialLabel }}
          </span>
          <span v-else-if="auth.isPro" class="hidden font-mono text-[11px] text-mute sm:inline">PRO</span>
          <span v-else class="hidden font-mono text-[11px] text-mute sm:inline">FREE · intern</span>
          <RouterLink to="/profile" class="btn-ghost btn py-1 px-3 text-[12px]">{{ auth.user?.name }}</RouterLink>
        </div>
      </div>
    </header>
    <main class="mx-auto max-w-6xl px-4 py-8 md:px-6 md:py-10">
      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <nav class="fixed inset-x-0 bottom-0 z-20 flex justify-around border-t-2 border-ink bg-white py-2 text-[11px] font-medium md:hidden">
      <RouterLink to="/hall" class="px-1.5 py-1" active-class="" exact-active-class="text-accent">Зал</RouterLink>
      <RouterLink to="/pet" class="px-1.5 py-1">Пет</RouterLink>
      <RouterLink to="/practice" class="px-1.5 py-1">Практика</RouterLink>
      <RouterLink to="/interview" class="px-1.5 py-1">Собес</RouterLink>
      <RouterLink to="/live" class="px-1.5 py-1">Live</RouterLink>
      <RouterLink to="/profile" class="px-1.5 py-1">Я</RouterLink>
    </nav>
  </div>
</template>
