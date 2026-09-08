<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink, useRoute } from "vue-router";
import BoardView from "../components/BoardView.vue";
import BriefingPanel from "../components/BriefingPanel.vue";

const route = useRoute();
const boards = ref([]);
const current = ref(null);
const err = ref("");

const id = computed(() => route.params.id);
const activeId = computed(() => id.value || current.value?.id);

async function loadList() {
  const cat = await fetch("/api/catalog", { credentials: "include" }).then((r) => r.json());
  boards.value = cat.boards || [];
}

async function loadBoard(bid) {
  err.value = "";
  const res = await fetch(`/api/boards/${bid}`, { credentials: "include" });
  const data = await res.json();
  if (!res.ok) {
    err.value = data.error?.message || "Доска недоступна на текущем тарифе.";
    current.value = null;
    return;
  }
  current.value = data;
}

onMounted(async () => {
  await loadList();
  const first = id.value || boards.value.find((b) => b.id === "hold-capture")?.id || boards.value[0]?.id;
  if (first) await loadBoard(first);
});

watch(id, (v) => {
  if (v) loadBoard(v);
});
</script>

<template>
  <div>
    <p class="font-mono text-[11px] uppercase tracking-[0.22em] text-mute">Схемы операций</p>
    <h1 class="font-display mt-2 text-4xl md:text-6xl">Доски</h1>
    <p class="mt-4 max-w-2xl text-[17px] leading-7 text-mute">
      Короткие схемы платёжного контура: hold, журнал, IBAN, сверка, учебный API.
      Здесь не полная статья — один термин, шаги операции и типичная ошибка аналитика.
    </p>
    <ul class="mt-5 max-w-2xl space-y-2 text-[15px] leading-7">
      <li><strong>Зачем.</strong> Увидеть, что происходит с деньгами на каждом шаге — не только подпись на слайде.</li>
      <li><strong>Как.</strong> Выберите схему, нажмите блок слева, справа прочитайте правило, пример и ошибку.</li>
      <li><strong>Дальше.</strong> Тот же контур разбирается в материалах и уроках intern.</li>
    </ul>

    <div class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <RouterLink
        v-for="b in boards"
        :key="b.id"
        :to="`/boards/${b.id}`"
        class="card card-hover rounded-[24px] p-5"
        :class="activeId === b.id ? 'ring-2 ring-accent' : ''"
      >
        <div class="flex flex-wrap items-center gap-2">
          <span v-if="b.term" class="stamp">{{ b.term }}</span>
          <span class="font-mono text-[11px] text-mute">{{ b.minutes }} мин · {{ b.grade }}</span>
        </div>
        <div class="font-display mt-3 text-2xl leading-none">{{ b.title }}</div>
        <p v-if="b.learn" class="mt-2 text-[15px] leading-6">{{ b.learn }}</p>
        <p class="mt-2 text-sm leading-6 text-mute">{{ b.teaser }}</p>
      </RouterLink>
    </div>

    <p v-if="err" class="mt-6 text-bad">{{ err }}</p>
    <div v-else-if="current?.board" class="mt-10">
      <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-mute">Схема</p>
      <p class="mt-2 max-w-2xl text-[15px] leading-7 text-mute">
        Слева — шаги. Нажмите блок, чтобы прочитать, что меняется. Справа — для чего схема, что запомнить, пример, типичная ошибка.
      </p>
      <div class="mt-5 grid items-start gap-5 lg:grid-cols-[minmax(0,1.3fr)_minmax(280px,0.9fr)]">
        <BoardView :board="current.board" />
        <BriefingPanel :why="current.why" :facts="current.facts" :scene="current.scene" :trap="current.trap" />
      </div>
    </div>
  </div>
</template>
