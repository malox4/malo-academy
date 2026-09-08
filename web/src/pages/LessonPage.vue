<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import BoardView from "../components/BoardView.vue";
import DeskView from "../components/DeskView.vue";
import BriefingPanel from "../components/BriefingPanel.vue";
import PrimerArticle from "../components/PrimerArticle.vue";
import { useProgress } from "../stores/progress";
import { applyPlay } from "../lib/playFx";
import PaywallCard from "../components/PaywallCard.vue";
import { topicForLesson } from "../content/topics";

const route = useRoute();
const router = useRouter();
const progress = useProgress();
const lesson = ref(null);
const paywall = ref(false);
const brief = ref("");
const saved = ref(false);
const fx = ref({});

async function load() {
  paywall.value = false;
  lesson.value = null;
  brief.value = "";
  saved.value = false;
  fx.value = {};
  const res = await fetch(`/api/lessons/${route.params.id}`, { credentials: "include" });
  const data = await res.json();
  if (res.status === 403 && data.paywall) {
    paywall.value = true;
    lesson.value = data.lesson;
    return;
  }
  lesson.value = data;
}

onMounted(() => {
  progress.load();
  load();
});
watch(() => route.params.id, load);

const primer = computed(() => {
  if (lesson.value?.primer) return lesson.value.primer;
  return topicForLesson(lesson.value?.id)?.sections || null;
});

const nextTo = computed(() => {
  if (lesson.value?.next) return `/lesson/${lesson.value.next}`;
  if (lesson.value?.levelId) return `/level/${lesson.value.levelId}`;
  return "/hall";
});

async function goNext() {
  if (lesson.value?.id) await progress.complete(lesson.value.id);
  router.push(nextTo.value);
}

function onPlay(ev) {
  const board = lesson.value?.board;
  if (!board) return;
  fx.value = applyPlay(board, ev);
}

async function saveBrief() {
  if (!lesson.value?.brief) return;
  saved.value = false;
  await fetch("/api/practice", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      kind: "module",
      targetId: lesson.value.id,
      title: lesson.value.brief.title,
      body: brief.value,
    }),
  });
  saved.value = true;
}
</script>

<template>
  <PaywallCard v-if="paywall" :title="lesson?.title || 'Урок закрыт'" />
  <div v-else-if="!lesson" class="font-mono text-sm text-mute">Собираем доску…</div>
  <div v-else>
    <div class="flex flex-wrap items-center justify-between gap-3">
      <RouterLink :to="`/level/${lesson.levelId}`" class="text-sm text-mute hover:text-accent">← {{ lesson.levelId }}</RouterLink>
      <button type="button" class="text-sm font-semibold text-accent" @click="goNext">Дальше →</button>
    </div>
    <p class="mt-6 font-mono text-[11px] uppercase tracking-[0.22em] text-mute">
      {{ lesson.gradeId }} · {{ lesson.minutes }} мин · контур
    </p>
    <h1 class="font-display mt-2 text-4xl leading-[0.95] md:text-6xl">{{ lesson.title }}</h1>
    <p class="mt-4 max-w-2xl text-[17px] leading-7 text-mute">{{ lesson.teaser }}</p>

    <PrimerArticle v-if="primer" class="mt-8" :primer="primer" part="head" />

    <div class="mt-8 grid items-start gap-5 lg:grid-cols-[minmax(0,1.35fr)_minmax(280px,0.9fr)]">
      <div class="space-y-5">
        <BoardView v-if="lesson.board" :board="lesson.board" :fx="fx" />
        <DeskView v-if="lesson.desk" :desk="lesson.desk" @play="onPlay" />
        <section v-if="lesson.brief" class="card rounded-[28px] p-5 md:p-6">
          <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-ok">Бриф · дверь не держит</div>
          <h3 class="font-display mt-1 text-2xl">{{ lesson.brief.title }}</h3>
          <p class="mt-2 text-[15px] leading-7">{{ lesson.brief.prompt }}</p>
          <textarea
            v-model="brief"
            :placeholder="lesson.brief.placeholder"
            rows="5"
            class="mt-4 w-full rounded-2xl border border-line bg-white p-4 text-[16px] leading-7 outline-none focus:border-accent"
          />
          <button type="button" class="mt-3 btn btn-sm" @click="saveBrief">
            {{ saved ? "В зале" : "Сдать в зал" }}
          </button>
        </section>
      </div>
      <BriefingPanel :why="lesson.why" :facts="lesson.facts" :scene="lesson.scene" :trap="lesson.trap" />
    </div>

    <PrimerArticle v-if="primer" class="mt-10" :primer="primer" part="rest">
      <template #drill>
        <p class="text-sm text-mute">Упражнение на доске выше. Неверный ход показывает причину. Переход дальше не блокируется.</p>
      </template>
    </PrimerArticle>

    <div class="sticky bottom-20 z-10 mt-10 flex items-center justify-between gap-3 rounded-2xl bg-ink p-4 text-white md:bottom-6">
      <p class="text-sm text-white/60">Упражнение на доске. К следующему уроку можно перейти сразу.</p>
      <button type="button" class="btn btn-accent btn-sm" @click="goNext">
        {{ lesson.next ? "Следующий контур" : "К уровню" }}
      </button>
    </div>
  </div>
</template>
