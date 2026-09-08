<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import SiteNav from "../components/SiteNav.vue";
import { getJSON } from "../lib/http";
import { useAuth } from "../stores/auth";
import { LIVE_TRACKS, liveTrackById } from "../content/live";

const auth = useAuth();
const route = useRoute();
const router = useRouter();

const active = ref(LIVE_TRACKS[0].id);
const session = ref(null);
const question = ref(null);
const done = ref(false);
const assessment = ref(null);
const err = ref("");
const typed = ref("");
const pick = ref("");
const busy = ref(false);
const view = ref("catalog");

const track = computed(() => liveTrackById(active.value));

function applyPayload(data) {
  session.value = data.session || null;
  question.value = data.question || null;
  done.value = Boolean(data.done);
  if (data.assessment) assessment.value = data.assessment;
  if (data.session?.status === "running" && data.question) view.value = "session";
  if (data.done && data.assessment) view.value = "result";
}

async function refresh() {
  const { ok, data } = await getJSON("/api/interview/live");
  if (!ok) return;
  applyPayload(data);
}

onMounted(async () => {
  if (!auth.user) return;
  await refresh();
  if (route.query.start === "1" && !session.value) await start();
});

async function start() {
  if (!auth.user) {
    router.push({ path: "/register", query: { next: "/live" } });
    return;
  }
  err.value = "";
  busy.value = true;
  const { ok, data } = await getJSON("/api/interview/live/start", {
    method: "POST",
    body: JSON.stringify({ mode: "live" }),
  });
  busy.value = false;
  if (!ok) {
    err.value = data.error?.message || "Не удалось начать сессию.";
    return;
  }
  typed.value = "";
  pick.value = "";
  assessment.value = null;
  applyPayload(data);
  view.value = "session";
  router.replace({ path: "/live" });
}

async function answer(advance) {
  err.value = "";
  busy.value = true;
  const { ok, data } = await getJSON("/api/interview/live/answer", {
    method: "POST",
    body: JSON.stringify({
      pickId: pick.value,
      text: typed.value,
      advance: Boolean(advance),
    }),
  });
  busy.value = false;
  if (!ok) {
    err.value = data.error?.message || "Ответ не принят.";
    return;
  }
  typed.value = "";
  pick.value = "";
  applyPayload(data);
  if (data.done) await finish();
}

async function finish() {
  busy.value = true;
  const { ok, data } = await getJSON("/api/interview/live/finish", { method: "POST", body: "{}" });
  busy.value = false;
  if (!ok) return;
  session.value = data.session || session.value;
  assessment.value = data.assessment;
  view.value = "result";
  done.value = true;
}

function backCatalog() {
  view.value = "catalog";
}

const progressLabel = computed(() => {
  if (!session.value) return "";
  return `Вопрос ${Math.min(session.value.index + 1, session.value.total)} из ${session.value.total}`;
});
</script>

<template>
  <div>
    <SiteNav v-if="!auth.user" />
    <div :class="auth.user ? '' : 'mx-auto max-w-6xl px-4 pb-20 pt-8 md:px-6'">

    <template v-if="view === 'catalog'">
      <p class="font-mono text-[11px] uppercase tracking-[0.22em] text-mute">Самообучение · без потока</p>
      <h1 class="font-display mt-2 text-4xl md:text-6xl">Live собес</h1>
      <p class="mt-3 max-w-2xl text-[17px] leading-7 text-mute">
        Витрина блоков — карта сессии: разминка, SQL, кейс, софты. Не таймер и не замена банка вопросов.
        Сессия — дуга с холдом и журналом; уровень intern…senior пишется в профиль. Разбор одного вопроса — в пункте Собес.
      </p>

      <article class="card mt-8 rounded-[28px] p-6 md:p-8">
        <p class="font-mono text-[11px] text-accent">Сессия зала</p>
        <h2 class="font-display mt-1 text-3xl md:text-4xl">Дуга сессии</h2>
        <p class="mt-3 max-w-xl text-[15px] leading-7 text-mute">
          Роль → уточнения → SDLC → SQL → окно → стейкхолдеры → холд → 3DS → журнал → конфликт → закрытие.
          Ответ словами или ходом. Оценка хардов и софтов — после финиша, не во время каждого клика.
        </p>
        <div class="mt-5 flex flex-wrap gap-3">
          <button
            type="button"
            class="btn"
            :disabled="busy"
            @click="start"
          >
            {{ session?.status === "running" ? "Продолжить сессию" : "Начать сессию" }}
          </button>
          <RouterLink to="/interview" class="btn btn-ghost">
            Банк вопросов
          </RouterLink>
        </div>
        <p v-if="err" class="mt-3 text-sm text-bad">{{ err }}</p>
      </article>

      <div class="mt-10 flex flex-wrap gap-2">
        <button
          v-for="t in LIVE_TRACKS"
          :key="t.id"
          type="button"
          class="tab-plate"
          :class="active === t.id ? 'is-on' : ''"
          @click="active = t.id"
        >
          {{ t.title }}
        </button>
      </div>

      <section class="mt-8">
        <p class="font-mono text-[11px] text-accent">{{ track.kicker }}</p>
        <h2 class="font-display mt-1 text-3xl">{{ track.title }}</h2>
        <p class="mt-2 max-w-xl text-[15px] leading-7 text-mute">{{ track.blurb }}</p>
        <div class="mt-6 grid gap-3 md:grid-cols-2 lg:grid-cols-3">
          <RouterLink
            v-for="c in track.cards"
            :key="c.id"
            :to="{ path: '/interview', query: { q: c.id } }"
            class="card card-hover rounded-[24px] p-5"
          >
            <div class="flex items-center justify-between font-mono text-[11px] text-mute">
              <span>{{ c.minutes }} мин</span>
              <span v-if="c.free === false" class="rounded-full bg-ink px-2 py-0.5 text-white">PRO</span>
              <span v-else class="text-ok">разобрать</span>
            </div>
            <h3 class="mt-2 font-display text-2xl leading-none">{{ c.title }}</h3>
            <p class="mt-2 text-sm leading-6 text-mute">{{ c.teaser }}</p>
          </RouterLink>
        </div>
      </section>
    </template>

    <template v-else-if="view === 'session' && question">
      <button type="button" class="text-sm text-mute hover:text-ink" @click="backCatalog">← К витрине</button>
      <p class="mt-4 font-mono text-[11px] uppercase tracking-[0.22em] text-mute">{{ progressLabel }} · {{ question.stageTitle }}</p>
      <h1 class="font-display mt-2 text-3xl md:text-5xl">{{ question.topic }}</h1>

      <section class="card mt-6 rounded-[28px] p-5 md:p-6">
        <div class="rounded-2xl rounded-tl-sm bg-ink px-4 py-3 text-white">
          <div class="font-mono text-[11px] uppercase tracking-widest text-white/45">Интервьюер</div>
          <p class="mt-1 text-[17px] leading-7">{{ question.question }}</p>
        </div>
        <p v-if="session?.followOpen" class="mt-4 text-[15px] leading-7 text-mute">
          Follow-up: {{ question.follow }}
        </p>
        <p class="mt-4 font-mono text-[10px] uppercase tracking-[0.2em] text-accent">Что скажете</p>
        <div class="mt-3 grid gap-2">
          <button
            v-for="r in question.replies"
            :key="r.id"
            type="button"
            class="rounded-2xl border px-4 py-3 text-left text-[16px] leading-6"
            :class="pick === r.id ? 'border-ink bg-paper' : 'border-line bg-white hover:border-accent'"
            @click="pick = r.id"
          >
            {{ r.text }}
          </button>
        </div>
        <label class="mt-5 block text-sm font-medium">
          Свой текст
          <textarea
            v-model="typed"
            rows="3"
            class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 text-[15px] leading-6 outline-none focus:border-accent"
            placeholder="Ноги журнала, ключ, Need — не лозунг."
          />
        </label>
        <button
          type="button"
          class="btn btn-accent btn-sm mt-3"
          :disabled="busy"
          @click="answer(false)"
        >
          Ответить
        </button>
        <p v-if="session?.lastWhy" class="mt-4 text-[15px] leading-7" :class="session.lastGood ? 'text-ok' : 'text-mute'">
          {{ session.lastWhy }}
        </p>
        <p v-if="err" class="mt-3 text-sm text-bad">{{ err }}</p>
      </section>
    </template>

    <template v-else-if="view === 'result'">
      <button type="button" class="text-sm text-mute hover:text-ink" @click="backCatalog">← К витрине</button>
      <p class="mt-4 font-mono text-[11px] uppercase tracking-[0.22em] text-mute">Сессия закрыта</p>
      <h1 class="font-display mt-2 text-4xl md:text-6xl">{{ assessment?.level || "intern" }}</h1>
      <p v-if="assessment" class="mt-3 max-w-xl text-[16px] leading-7 text-mute">
        Балл {{ assessment.score }}/{{ assessment.max }} · уверенность {{ Math.round((assessment.confidence || 0) * 100) }}%.
        Оценка лежит в профиле. Это ориентир зала, не сертификат потока.
      </p>
      <div class="mt-6 flex flex-wrap gap-3">
        <RouterLink to="/profile" class="btn">Профиль</RouterLink>
        <RouterLink to="/interview" class="btn btn-ghost">Банк вопросов</RouterLink>
        <button type="button" class="btn btn-ghost" @click="start">Ещё раз</button>
      </div>
    </template>
    </div>
  </div>
</template>
