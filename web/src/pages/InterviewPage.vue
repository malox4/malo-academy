<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import BoardView from "../components/BoardView.vue";
import PaywallCard from "../components/PaywallCard.vue";
import { INTERVIEW, TOPICS, scoreAnswer, inArcOrder } from "../content/interview";
import { applyPlay } from "../lib/playFx";
import { useAuth } from "../stores/auth";
import { useProgress } from "../stores/progress";

const route = useRoute();
const router = useRouter();
const auth = useAuth();
const progress = useProgress();

const topic = ref("Порядок");
const open = ref(inArcOrder(INTERVIEW)[0]?.id || "");
const pick = ref(null);
const typed = ref("");
const fx = ref({});
const last = ref(null);

const list = computed(() => {
  if (topic.value === "Порядок") return inArcOrder(INTERVIEW);
  if (topic.value === "Все") return INTERVIEW;
  return INTERVIEW.filter((i) => i.topic === topic.value);
});

const active = computed(() => INTERVIEW.find((i) => i.id === open.value) || list.value[0]);
const locked = computed(() => active.value?.grade === "senior" && !auth.isPro);

watch(
  () => route.query.q,
  (q) => {
    if (q && INTERVIEW.some((i) => i.id === q)) open.value = String(q);
  },
  { immediate: true },
);

onMounted(() => {
  progress.load();
  const q = String(route.query.q || "");
  if (!q || !INTERVIEW.some((i) => i.id === q)) {
    const id = open.value;
    if (id) router.replace({ query: { q: id } });
  }
});

watch(active, () => {
  pick.value = null;
  typed.value = "";
  fx.value = {};
  last.value = null;
});

function selectItem(id) {
  open.value = id;
  progress.markInterview(id);
  router.replace({ query: { q: id } });
}

function react(ev) {
  const board = active.value?.board;
  if (!board) return;
  last.value = ev;
  fx.value = applyPlay(board, ev);
}

function say(reply, index) {
  pick.value = index;
  react({
    good: reply.good,
    why: reply.why,
    hit: reply.hit,
    lights: reply.lights,
    broken: reply.broken,
    tickets: reply.tickets,
  });
  progress.markInterview(active.value.id);
}

function speak() {
  const item = active.value;
  const scored = scoreAnswer(item, typed.value);
  const board = item.board;
  const hit = board.nodes?.[1]?.id || board.lines?.[0]?.id || board.debit?.[0]?.id;
  pick.value = -1;
  react({
    good: scored.good,
    why: scored.why,
    hit,
    lights: scored.good ? { [hit]: "ok" } : undefined,
    broken: scored.good ? [] : [hit],
  });
  progress.markInterview(item.id);
}

const internCount = INTERVIEW.filter((i) => i.grade === "intern").length;
const seniorCount = INTERVIEW.filter((i) => i.grade === "senior").length;
</script>

<template>
  <div>
    <p class="font-mono text-[11px] uppercase tracking-[0.22em] text-mute">Собес · BA / SA</p>
    <h1 class="font-display mt-2 text-4xl md:text-6xl">Собеседование</h1>
    <p class="mt-3 max-w-2xl text-[17px] leading-7 text-mute">
      Вопрос со стола → что скажете → доска или тикет реагирует → почему. Не тест с галочкой.
      Intern ({{ internCount }}) открыт. Расширенный банк ({{ seniorCount }}) — в PRO и в триале.
      Живая сессия на оценку — отдельный пункт Live собес.
    </p>

    <div class="mt-6 flex flex-wrap gap-2">
      <button
        v-for="t in TOPICS"
        :key="t"
        type="button"
        class="tab-plate"
        :class="topic === t ? 'is-on' : ''"
        @click="topic = t"
      >
        {{ t }}
      </button>
    </div>

    <div class="mt-8 grid gap-5 lg:grid-cols-[minmax(240px,0.85fr)_minmax(0,1.35fr)]">
      <aside class="card max-h-[70vh] space-y-1 overflow-auto rounded-[28px] p-3">
        <button
          v-for="i in list"
          :key="i.id"
          type="button"
          class="w-full rounded-2xl px-3 py-3 text-left text-sm"
          :class="active?.id === i.id ? 'bg-ink text-white' : 'text-mute hover:bg-paper'"
          @click="selectItem(i.id)"
        >
          <div class="flex items-center gap-2 font-mono text-[10px] uppercase tracking-widest" :class="active?.id === i.id ? 'text-white/55' : 'text-accent'">
            <span>{{ i.topic }}</span>
            <span v-if="i.grade === 'senior'" class="stamp stamp-solid text-[9px]">PRO</span>
            <span v-if="progress.sawInterview(i.id)" class="ml-auto text-[9px] opacity-70">разобран</span>
          </div>
          <div class="mt-1 leading-5">{{ i.question }}</div>
        </button>
      </aside>

      <div v-if="locked">
        <PaywallCard title="Senior-трек собеса — в PRO" />
        <p class="mt-4 text-sm text-mute">Intern-вопросы слева открыты. Полный банк — на триале и в платном PRO.</p>
      </div>

      <div v-else-if="active" class="space-y-4">
        <BoardView :board="active.board" :fx="fx" />

        <section class="card rounded-[28px] p-5 md:p-6">
          <div class="rounded-2xl rounded-tl-sm bg-ink px-4 py-3 text-white">
            <div class="font-mono text-[11px] uppercase tracking-widest text-white/45">Интервьюер · {{ active.topic }}</div>
            <p class="mt-1 text-[17px] leading-7">{{ active.question }}</p>
          </div>

          <p class="mt-4 font-mono text-[10px] uppercase tracking-[0.2em] text-accent">Что скажете за столом</p>
          <div class="mt-3 grid gap-2">
            <button
              v-for="(r, i) in active.replies"
              :key="r.id"
              type="button"
              class="rounded-2xl border px-4 py-3 text-left text-[16px] leading-6 transition"
              :class="
                pick === i
                  ? r.good
                    ? 'border-ok bg-ok/10'
                    : 'border-bad bg-bad/10'
                  : 'border-line bg-white hover:border-accent'
              "
              @click="say(r, i)"
            >
              {{ r.text }}
            </button>
          </div>

          <label class="mt-5 block text-sm font-medium">
            Свой ответ
            <textarea
              v-model="typed"
              rows="3"
              class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 text-[15px] leading-6 outline-none focus:border-accent"
              placeholder="Ноги журнала, ключ, Need — не лозунг."
            />
          </label>
          <button type="button" class="btn btn-accent btn-sm mt-3" @click="speak">
            Сказать своими словами
          </button>
          <p class="mt-2 text-[12px] text-mute">Можно сменить ход — доска перерисуется. Зелёной галочки как зачёта нет.</p>
        </section>

        <section v-if="last" class="card rounded-[28px] p-5">
          <div class="font-mono text-[10px] uppercase tracking-[0.2em]" :class="last.good ? 'text-ok' : 'text-bad'">
            {{ last.good ? "Шов держится" : "Доска среагировала" }}
          </div>
          <p class="mt-2 text-[16px] leading-7">{{ last.why }}</p>
          <div class="mt-4 text-[11px] uppercase tracking-widest text-mute">Сильный ответ</div>
          <p class="mt-1 text-[15px] leading-7">{{ active.strong }}</p>
        </section>

        <section class="rounded-[28px] border border-bad/20 bg-bad/8 p-5">
          <div class="font-mono text-[10px] uppercase tracking-[0.2em] text-bad">Ловушки</div>
          <ul class="mt-3 space-y-2 text-sm leading-6">
            <li v-for="t in active.traps" :key="t">— {{ t }}</li>
          </ul>
          <p v-if="last" class="mt-4 text-[15px] leading-7">
            <span class="font-mono text-[11px] text-mute">Follow-up · </span>
            А если вам скажут: {{ active.follow }}
          </p>
        </section>
      </div>
    </div>
  </div>
</template>
