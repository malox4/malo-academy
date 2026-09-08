<script setup>
import { computed, onMounted } from "vue";
import { RouterLink } from "vue-router";
import { useAuth } from "../stores/auth";
import { useProgress } from "../stores/progress";
import { extrasForLevel } from "../content/path";

const auth = useAuth();
const progress = useProgress();

onMounted(() => progress.load());

const grades = computed(() => progress.catalog?.grades || []);
const lessons = computed(() => progress.catalog?.lessons || []);
const internOpen = computed(() => lessons.value.find((l) => l.id === "intern-1-profession"));

function doneOn(level) {
  const ids = level.moduleIds || [];
  if (!ids.length) return 0;
  return Math.round((ids.filter((id) => progress.isDone(id)).length / ids.length) * 100);
}

const firstLesson = computed(() => (internOpen.value ? `/lesson/${internOpen.value.id}` : "/level/intern-1"));
</script>

<template>
  <div>
    <p class="kicker">зал аналитика</p>
    <h1 class="font-display mt-3 text-4xl leading-[0.95] md:text-6xl">Этажи intern → senior</h1>
    <span class="section-rule" aria-hidden="true" />
    <p class="mt-4 max-w-2xl text-[17px] leading-7 text-mute">
      48 уроков. На этаже — материалы пути. Intern открыт. Запись в учебный банк — PRO или триал.
    </p>
    <div class="mt-6 flex flex-wrap gap-3">
      <RouterLink :to="firstLesson" class="btn btn-accent">Первый урок</RouterLink>
      <RouterLink to="/pet" class="btn btn-ghost">Учебный банк</RouterLink>
      <RouterLink to="/practice" class="btn btn-ghost">Практика</RouterLink>
      <RouterLink to="/interview" class="btn btn-ghost">Собес</RouterLink>
    </div>

    <div class="mt-14 space-y-12">
      <section v-for="g in grades" :key="g.id">
        <div class="flex items-end justify-between gap-4">
          <div>
            <div class="stamp stamp-ink">{{ g.roman }}</div>
            <h2 class="font-display mt-2 text-3xl md:text-4xl">{{ g.title }}</h2>
            <span class="section-rule" aria-hidden="true" />
            <p class="mt-2 max-w-xl text-sm text-mute">{{ g.tagline }}</p>
          </div>
          <span v-if="!auth.canGrade(g.id)" class="stamp stamp-solid">PRO</span>
        </div>
        <div class="catalog-grid mt-6 md:grid-cols-3">
          <RouterLink
            v-for="lv in g.levels"
            :key="lv.id"
            :to="`/level/${lv.id}`"
            class="card card-hover p-5"
          >
            <div class="flex items-center justify-between">
              <span class="stamp stamp-ink">{{ lv.hours }}</span>
              <span class="plate-n">{{ doneOn(lv) }}%</span>
            </div>
            <div class="mt-3 h-1 overflow-hidden bg-line">
              <div class="h-full bg-accent" :style="{ width: doneOn(lv) + '%' }" />
            </div>
            <div class="mt-4 font-display text-2xl leading-none">{{ lv.title }}</div>
            <p class="mt-2 text-sm leading-6 text-mute">{{ lv.subtitle }}</p>
            <p class="mt-3 font-mono text-[11px] text-accent">
              {{ (lv.moduleIds || []).length }} уроков
              <span v-if="(lv.materialIds || extrasForLevel(lv.id)).length">
                · {{ (lv.materialIds || extrasForLevel(lv.id)).length }} материалов пути
              </span>
            </p>
          </RouterLink>
        </div>
      </section>
    </div>
  </div>
</template>
