<script setup>
import { RouterLink } from "vue-router";
import BrandMark from "./BrandMark.vue";

defineProps({
  solidCta: { type: Boolean, default: false },
});
</script>

<template>
  <header class="site-bar sticky top-0 z-30">
    <div class="mx-auto flex h-14 max-w-6xl items-center gap-3 px-4 md:px-6">
      <BrandMark to="/" />
      <nav class="ml-4 hidden items-center gap-1 text-[13px] font-medium md:flex">
        <RouterLink to="/" class="nav-link">Главная</RouterLink>
        <RouterLink to="/pet" class="nav-link">Пет</RouterLink>
        <RouterLink to="/practice" class="nav-link">Практика</RouterLink>
        <RouterLink to="/materials" class="nav-link">Материалы</RouterLink>
        <RouterLink to="/interview" class="nav-link">Собес</RouterLink>
        <RouterLink to="/live" class="nav-link">Live собес</RouterLink>
        <RouterLink to="/pricing" class="nav-link">Сравнение</RouterLink>
      </nav>
      <div class="ml-auto flex items-center gap-2">
        <RouterLink to="/login" class="px-2 py-1.5 text-[13px] font-medium text-mute hover:text-ink">Войти</RouterLink>
        <RouterLink
          to="/register"
          class="btn text-[13px]"
          :class="solidCta ? '' : 'btn-accent'"
        >
          Начать бесплатно
        </RouterLink>
      </div>
    </div>
  </header>
</template>
