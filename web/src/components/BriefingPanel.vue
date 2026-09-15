<script setup>
defineProps({
  why: String,
  facts: { type: Array, default: () => [] },
  scene: Object,
  trap: Object,
});
</script>

<template>
  <aside class="space-y-3">
    <section v-if="why" class="card rise rounded-[24px] p-5">
      <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-accent">Для чего это нужно</div>
      <p class="mt-2 text-[15px] leading-6">{{ why }}</p>
    </section>
    <section v-if="facts?.length" class="card rise rounded-[24px] p-5" style="animation-delay: 80ms">
      <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-mute">Что запомнить</div>
      <ol class="mt-3 space-y-2.5">
        <li v-for="(f, i) in facts" :key="i" class="flex gap-3 text-[14px] leading-6">
          <span class="font-mono text-[11px] text-accent">{{ String(i + 1).padStart(2, "0") }}</span>
          <span>{{ f }}</span>
        </li>
      </ol>
    </section>
    <section v-if="scene" class="card rise overflow-hidden rounded-[24px]" style="animation-delay: 140ms">
      <div class="bg-ink px-5 py-4 text-white">
        <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-white/50">Пример</div>
        <h3 class="font-display mt-1 text-xl">{{ scene.title }}</h3>
        <p class="mt-2 text-[14px] leading-6 text-white/80">{{ scene.body }}</p>
        <div v-if="scene.chips?.length" class="mt-3 flex flex-wrap gap-2">
          <span v-for="c in scene.chips" :key="c" class="rounded-full bg-white/10 px-2.5 py-1 font-mono text-[11px]">{{ c }}</span>
        </div>
      </div>
    </section>
    <section v-if="trap" class="card rise rounded-[24px] border-bad/20 p-5" style="animation-delay: 200ms">
      <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-bad">{{ trap.title || "Типичная ошибка" }}</div>
      <p class="mt-2 text-[14px] leading-6">{{ trap.body }}</p>
    </section>
  </aside>
</template>
