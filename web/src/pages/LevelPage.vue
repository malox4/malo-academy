<script setup>
import { computed, onMounted } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { useAuth } from "../stores/auth";
import { useProgress } from "../stores/progress";
import PaywallCard from "../components/PaywallCard.vue";
import { extrasForLevel } from "../content/path";
import { topicById, labIdForTopic } from "../content/topics";
import { LABS, drillKind } from "../content/practice";

const route = useRoute();
const auth = useAuth();
const progress = useProgress();

onMounted(() => progress.load());

const found = computed(() => {
  const cat = progress.catalog;
  if (!cat) return null;
  for (const g of cat.grades) {
    const level = g.levels.find((l) => l.id === route.params.id);
    if (level) return { grade: g, level };
  }
  return null;
});

const locked = computed(() => found.value && !auth.canGrade(found.value.grade.id));
const lessons = computed(() => progress.catalog?.lessons || []);

function meta(id) {
  return lessons.value.find((l) => l.id === id);
}

const extras = computed(() => {
  const fromCat = found.value?.level?.materialIds;
  const ids = Array.isArray(fromCat) && fromCat.length ? fromCat : extrasForLevel(route.params.id);
  return ids
    .map((id) => {
      const topic = topicById(id);
      if (!topic) return null;
      const labId = labIdForTopic(id);
      const lab = LABS.find((l) => l.id === labId);
      return { topic, lab, labId };
    })
    .filter(Boolean);
});
</script>

<template>
  <div v-if="!found">Уровень не найден.</div>
  <div v-else>
    <RouterLink to="/hall" class="text-sm text-mute hover:text-accent">← Зал</RouterLink>
    <p class="mt-6 font-mono text-[11px] uppercase tracking-[0.2em] text-mute">{{ found.grade.title }} · {{ found.level.hours }}</p>
    <h1 class="font-display mt-2 text-4xl md:text-6xl">{{ found.level.title }}</h1>
    <p class="mt-3 max-w-xl text-[17px] leading-7 text-mute">{{ found.level.subtitle }}</p>

    <PaywallCard v-if="locked" class="mt-8" :title="found.grade.title + ' закрыт на бесплатном тарифе'" />

    <div v-if="!locked" class="mt-8 space-y-3">
      <RouterLink
        v-for="(id, i) in found.level.moduleIds"
        :key="id"
        :to="`/lesson/${id}`"
        class="card card-hover flex items-center gap-4 rounded-[24px] p-4"
      >
        <div
          class="grid h-12 w-12 place-items-center rounded-2xl font-display text-lg"
          :class="progress.isDone(id) ? 'bg-ok/15 text-ok' : 'bg-ink text-white'"
        >
          {{ i + 1 }}
        </div>
        <div class="min-w-0 flex-1">
          <div class="font-medium">{{ meta(id)?.title || id }}</div>
          <p class="truncate text-sm text-mute">{{ meta(id)?.teaser }}</p>
        </div>
        <span class="hidden font-mono text-[11px] text-mute sm:block">{{ meta(id)?.minutes }} мин</span>
      </RouterLink>
    </div>

    <section v-if="extras.length" class="mt-12">
      <p class="kicker">материалы пути</p>
      <h2 class="font-display mt-2 text-3xl">Статья и лаба этого этажа</h2>
      <p class="mt-2 max-w-xl text-sm leading-6 text-mute">
        Не отдельный каталог: эти темы входят в intern→senior. 48 уроков на месте, рядом — закрепление.
      </p>
      <div class="mt-6 grid gap-3 md:grid-cols-2">
        <article v-for="row in extras" :key="row.topic.id" class="card rounded-[24px] p-5">
          <div class="flex flex-wrap items-center gap-2 font-mono text-[11px] text-mute">
            <span>{{ row.topic.grade }}</span>
            <span v-if="row.lab">{{ drillKind(row.lab) }}</span>
            <span v-if="row.topic.free" class="stamp stamp-ok">intern</span>
          </div>
          <div class="font-display mt-2 text-2xl">{{ row.topic.title }}</div>
          <p class="mt-2 text-sm leading-6 text-mute">{{ row.topic.teaser }}</p>
          <div class="mt-4 flex flex-wrap gap-3 text-sm">
            <RouterLink :to="'/materials/' + row.topic.id" class="font-medium text-accent">Статья →</RouterLink>
            <RouterLink v-if="row.lab" :to="'/practice?lab=' + row.lab.id" class="text-mute hover:text-accent">Лаба →</RouterLink>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>
