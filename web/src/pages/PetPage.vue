<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getJSON } from "../lib/http";
import { useAuth } from "../stores/auth";
import { useProgress } from "../stores/progress";
import { PET_INTRO, PET_TABS, LESSON_TABS, lessonFor, lessonCases, consoleCases } from "../content/petLab";

const auth = useAuth();
const progress = useProgress();
const route = useRoute();
const router = useRouter();

const tab = computed(() => {
  const t = String(route.query.tab || "console");
  if (t === "requests" || t === "api") return "api";
  if (t === "sql") return "sql";
  if (LESSON_TABS.includes(t)) return t;
  return "console";
});
const isLesson = computed(() => LESSON_TABS.includes(tab.value));
const activeLesson = computed(() => (isLesson.value ? lessonFor(tab.value) : null));

function setTab(id) {
  router.replace({ path: "/pet", query: { tab: id } });
}

const trial = ref(null);
const wallets = ref([]);
const holds = ref([]);
const journals = ref([]);
const last = ref(null);
const err = ref("");
const running = ref(false);

const idemKey = ref("p2p-" + Date.now().toString(36));
const idemKey2 = ref("p2p2-" + Date.now().toString(36));

const opsCases = computed(() => consoleCases(idemKey.value));
const lessonCaseList = computed(() => lessonCases(tab.value, idemKey.value, idemKey2.value));
const caseList = computed(() => (isLesson.value ? lessonCaseList.value : opsCases.value));
const activeCaseId = ref("wallets");
const activeCase = computed(() => caseList.value.find((c) => c.id === activeCaseId.value) || caseList.value[0]);

const workMethod = ref("GET");
const workPath = ref("/api/v1/wallets");
const workHeaders = ref("");
const workBody = ref("");

function loadCase(c) {
  if (!c) return;
  activeCaseId.value = c.id;
  workMethod.value = c.method;
  workPath.value = c.path;
  workHeaders.value = c.headers
    ? Object.entries(c.headers)
        .map(([k, v]) => k + ": " + v)
        .join("\n")
    : "";
  workBody.value = c.body ? JSON.stringify(c.body, null, 2) : "";
  err.value = "";
}

watch(
  () => tab.value,
  (id) => {
    if (LESSON_TABS.includes(id)) loadCase(lessonCases(id, idemKey.value, idemKey2.value)[0]);
    else if (id === "console") loadCase(opsCases.value[0]);
  },
);

function parseHeaders(text) {
  const headers = {};
  for (const line of String(text || "").split("\n")) {
    const i = line.indexOf(":");
    if (i > 0) headers[line.slice(0, i).trim()] = line.slice(i + 1).trim();
  }
  return headers;
}

const requestPreview = computed(() => {
  const lines = [`${workMethod.value} ${workPath.value}`];
  const headers = parseHeaders(workHeaders.value);
  for (const [k, v] of Object.entries(headers)) lines.push(`${k}: ${v}`);
  if (workMethod.value === "GET") lines.push("", "(тела нет)");
  else if (workBody.value.trim()) lines.push("", workBody.value.trim());
  return lines.join("\n");
});

async function refresh() {
  const [t, w, h, j] = await Promise.all([
    getJSON("/api/v1/ledger/trial-balance"),
    getJSON("/api/v1/wallets"),
    getJSON("/api/v1/holds"),
    getJSON("/api/v1/ledger/journals"),
  ]);
  trial.value = t.data;
  wallets.value = w.data.items || [];
  holds.value = h.data.items || [];
  journals.value = (j.data.items || []).slice(-10).reverse();
}

async function runCase() {
  err.value = "";
  running.value = true;
  const method = workMethod.value;
  const path = workPath.value;
  const headers = parseHeaders(workHeaders.value);
  let body;
  if (method !== "GET") {
    const raw = workBody.value.trim();
    if (raw) {
      try {
        body = JSON.parse(raw);
      } catch {
        err.value = "Тело запроса не JSON.";
        running.value = false;
        return;
      }
    } else {
      body = {};
    }
  }
  const { ok, status, data } = await getJSON(path, {
    method,
    headers,
    body: method === "GET" ? undefined : JSON.stringify(body ?? {}),
  });
  last.value = {
    label: activeCase.value?.title || method,
    status,
    ok,
    data,
    path,
    method,
    headers,
    body: method === "GET" ? undefined : body,
  };
  if (!ok) {
    const msg = data.error?.message || data.error || String(status);
    err.value = data.error?.code === "PLAN" ? "Запись закрыта тарифом Intern (код PLAN). Метод, путь и тело всё равно видны — это граница контракта, не скрытая кнопка." : msg;
  }
  await refresh();
  running.value = false;
}

function newKeys() {
  idemKey.value = "p2p-" + Date.now().toString(36);
  idemKey2.value = "p2p2-" + Date.now().toString(36);
  const list = tab.value === "p2p" ? lessonCases("p2p", idemKey.value, idemKey2.value) : caseList.value;
  const c = list.find((x) => x.id === activeCaseId.value) || list[0];
  if (c) loadCase(c);
}

const sqlTasks = ref([]);
const sqlSchema = ref([]);
const sqlText = ref("SELECT id, name, currency, available FROM intern_wallets");
const sqlOut = ref(null);
const sqlErr = ref("");
const activeTask = ref(null);
const internSql = ref(true);
const sqlTask = computed(() => sqlTasks.value.find((t) => t.id === activeTask.value));

const explorer = ref([]);
const group = ref("Счета");
const openRoute = ref(null);
const bodyText = ref("");
const headerKey = ref("");
const runOut = ref(null);
const activeMethod = ref("");

const groups = computed(() => {
  const g = [];
  for (const row of explorer.value) {
    if (!g.includes(row.group)) g.push(row.group);
  }
  return g;
});

const grouped = computed(() => explorer.value.filter((r) => r.group === group.value));
const activeDoc = computed(
  () =>
    explorer.value.find((r) => r.path === openRoute.value && r.method === (activeMethod.value || r.method)) ||
    grouped.value[0],
);

watch(grouped, (rows) => {
  if (rows[0] && !rows.some((r) => r.path === openRoute.value)) pickDoc(rows[0]);
});

function pickDoc(row) {
  openRoute.value = row.path;
  activeMethod.value = row.method;
  bodyText.value = row.example ? JSON.stringify(row.example, null, 2) : "";
  headerKey.value = row.headers?.["Idempotency-Key"] || "";
  runOut.value = null;
}

function hintToSql(t) {
  const h = String(t.hint || t.gold || "").trim();
  if (/^select\b/i.test(h)) return h.replace(/;+\s*$/, "");
  return internSql.value
    ? "SELECT id, name, currency, available FROM intern_wallets"
    : "SELECT id, name, currency, ledger, hold, (ledger - hold) AS available FROM wallets";
}

async function loadSql() {
  const [t, sch] = await Promise.all([getJSON("/api/sql/tasks"), getJSON("/api/sql/schema")]);
  sqlTasks.value = t.data.items || [];
  internSql.value = Boolean(t.data.intern);
  sqlSchema.value = sch.data.tables || [];
  const first = internSql.value
    ? sqlTasks.value.find((x) => x.internSafe)
    : sqlTasks.value.find((x) => x.grade === "junior") || sqlTasks.value[0];
  if (first) useTask(first);
}

async function loadExplorer() {
  const { data } = await getJSON("/api/explorer");
  explorer.value = data.items || [];
  if (explorer.value[0]) pickDoc(explorer.value[0]);
}

async function runSql(check) {
  sqlErr.value = "";
  sqlOut.value = null;
  const path = check ? "/api/sql/check" : "/api/sql/run";
  const payload = { sql: sqlText.value };
  if (check && activeTask.value) payload.taskId = activeTask.value;
  const { ok, data } = await getJSON(path, { method: "POST", body: JSON.stringify(payload) });
  sqlOut.value = data;
  if (!ok) sqlErr.value = typeof data.error === "string" ? data.error : data.error?.message || "Запрос не прошёл.";
  if (check && data.ok) await loadSql();
}

function useTask(t) {
  activeTask.value = t.id;
  sqlText.value = hintToSql(t);
}

async function runDoc() {
  const doc = activeDoc.value;
  if (!doc) return;
  let body;
  try {
    body = bodyText.value.trim() ? JSON.parse(bodyText.value) : {};
  } catch {
    runOut.value = { error: "Тело не JSON." };
    return;
  }
  const headers = {};
  if (headerKey.value) headers["Idempotency-Key"] = headerKey.value;
  const { status, data } = await getJSON(doc.path, {
    method: doc.method,
    headers,
    body: doc.method === "GET" ? undefined : JSON.stringify(body),
  });
  runOut.value = { method: doc.method, path: doc.path, status, data };
  await refresh();
}

onMounted(async () => {
  await progress.load();
  await refresh();
  await loadSql();
  await loadExplorer();
  if (isLesson.value) loadCase(lessonCaseList.value[0]);
  else loadCase(opsCases.value[0]);
});

const anna = computed(() => wallets.value.find((w) => w.id === "wal_anna"));
const boris = computed(() => wallets.value.find((w) => w.id === "wal_boris"));
const openHolds = computed(() => holds.value.filter((h) => String(h.status || "").toUpperCase() === "OPEN"));
const sqlColumns = computed(() => {
  const o = sqlOut.value;
  return o?.columns || o?.result?.columns || Object.keys((o?.rows || o?.result?.rows || [])[0] || {});
});
const sqlRows = computed(() => sqlOut.value?.rows || sqlOut.value?.result?.rows || []);
const lastP2p = computed(() => {
  const d = last.value?.data;
  if (!d) return null;
  if (d.journalId && d.fromWalletId) return d;
  return null;
});

function pretty(v) {
  return JSON.stringify(v, null, 2);
}
</script>

<template>
  <div>
    <p class="font-mono text-[11px] uppercase tracking-[0.22em] text-mute">Analyst Hall · учебный банк</p>
    <h1 class="font-display mt-2 text-4xl md:text-6xl">{{ PET_INTRO.title }}</h1>
    <p class="mt-3 max-w-3xl text-[17px] leading-7 text-mute">{{ PET_INTRO.lead }}</p>

    <div v-if="tab === 'console'" class="mt-6 grid gap-3 md:grid-cols-3">
      <section class="card rounded-[24px] p-5">
        <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Что это</p>
        <p class="mt-2 text-sm leading-6">{{ PET_INTRO.what }}</p>
      </section>
      <section class="card rounded-[24px] p-5">
        <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Зачем</p>
        <p class="mt-2 text-sm leading-6">{{ PET_INTRO.why }}</p>
      </section>
      <section class="card rounded-[24px] p-5">
        <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Как работать</p>
        <ol class="mt-2 space-y-2 text-sm leading-6">
          <li v-for="s in PET_INTRO.how" :key="s.n"><span class="font-medium">{{ s.n }}. {{ s.title }}. </span>{{ s.text }}</li>
        </ol>
      </section>
    </div>

    <div class="mt-6 flex flex-wrap gap-2">
      <button
        v-for="t in PET_TABS"
        :key="t.id"
        type="button"
        class="tab-plate"
        :class="tab === t.id ? 'is-on' : ''"
        :title="t.about"
        @click="setTab(t.id)"
      >
        {{ t.label }}
      </button>
    </div>
    <p class="mt-2 text-sm text-mute">{{ PET_TABS.find((t) => t.id === tab)?.about }}</p>

    <div v-if="tab === 'console' || isLesson" class="mt-8">
      <div v-if="isLesson && activeLesson" class="mb-6 space-y-4">
        <section class="card rounded-[28px] p-5 md:p-6">
          <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-accent">Термин</p>
          <h2 class="font-display mt-2 text-3xl">{{ activeLesson.name }} <span class="text-mute">/ {{ activeLesson.en }}</span></h2>
          <p class="mt-3 text-[16px] leading-7">{{ activeLesson.term }}</p>
          <p class="mt-5 font-mono text-[10px] uppercase tracking-widest text-accent">Объяснение</p>
          <p v-for="(p, i) in activeLesson.about" :key="i" class="mt-2 text-[16px] leading-7">{{ p }}</p>
          <p class="mt-5 font-mono text-[10px] uppercase tracking-widest text-accent">Примеры</p>
          <ul class="mt-2 list-disc space-y-1 pl-5 text-[15px] leading-7">
            <li v-for="(p, i) in activeLesson.example" :key="i">{{ p }}</li>
          </ul>
          <p class="mt-5 font-mono text-[10px] uppercase tracking-widest text-accent">Для чего</p>
          <p class="mt-2 text-[16px] leading-7">{{ activeLesson.purpose }}</p>
          <p class="mt-4 rounded-2xl bg-paper px-4 py-3 text-sm leading-6">{{ activeLesson.watch }}</p>
        </section>
        <div v-if="lastP2p" class="grid gap-3 sm:grid-cols-2">
          <div class="card rounded-[24px] p-5">
            <p class="font-mono text-[10px] uppercase tracking-widest text-mute">Дебет · отправитель</p>
            <p class="font-display mt-2 text-2xl">DR Анна {{ lastP2p.amount }}</p>
            <p class="mt-1 text-sm text-mute">Уменьшение пассива клиента. journalId {{ lastP2p.journalId }}</p>
          </div>
          <div class="card rounded-[24px] p-5">
            <p class="font-mono text-[10px] uppercase tracking-widest text-mute">Кредит · получатель</p>
            <p class="font-display mt-2 text-2xl">CR Борис {{ lastP2p.amount }}</p>
            <p class="mt-1 text-sm text-mute">Увеличение пассива клиента. Одна сумма, один journalId.</p>
          </div>
        </div>
      </div>

      <div class="grid gap-3 sm:grid-cols-3">
        <div class="card rounded-[24px] p-5">
          <div class="font-mono text-[11px] text-mute">Анна · available</div>
          <div class="font-display mt-2 text-4xl">{{ anna?.available ?? "—" }}</div>
          <div class="mt-1 font-mono text-[11px] text-mute">{{ anna?.currency }} · холд {{ anna?.hold ?? 0 }} · ledger {{ anna?.ledger ?? "—" }}</div>
        </div>
        <div class="card rounded-[24px] p-5">
          <div class="font-mono text-[11px] text-mute">Борис · available</div>
          <div class="font-display mt-2 text-4xl">{{ boris?.available ?? "—" }}</div>
          <div class="mt-1 font-mono text-[11px] text-mute">{{ boris?.currency }} · холд {{ boris?.hold ?? 0 }} · ledger {{ boris?.ledger ?? "—" }}</div>
        </div>
        <div class="card rounded-[24px] p-5">
          <div class="font-mono text-[11px] text-mute">Пробный баланс UZS</div>
          <div class="font-display mt-2 text-4xl" :class="trial?.balanced ? 'text-ok' : 'text-bad'">
            {{ trial?.balanced ? "сходится" : trial?.diff ?? "—" }}
          </div>
          <div class="mt-1 font-mono text-[11px] text-mute">DR {{ trial?.debit }} / CR {{ trial?.credit }}</div>
        </div>
      </div>

      <div class="mt-6 grid items-start gap-5 lg:grid-cols-[minmax(260px,0.85fr)_minmax(0,1.2fr)]">
        <aside class="card max-h-[78vh] space-y-1 overflow-auto rounded-[28px] p-3">
          <p class="px-2 pt-2 font-mono text-[10px] uppercase tracking-widest text-mute">
            {{ isLesson ? "Кейсы по порядку" : "Что сделать" }}
          </p>
          <button
            v-for="c in caseList"
            :key="c.id"
            type="button"
            class="w-full rounded-2xl px-3 py-3 text-left text-sm"
            :class="activeCaseId === c.id ? 'bg-ink text-white' : 'text-mute hover:bg-paper'"
            @click="loadCase(c)"
          >
            <div class="font-mono text-[10px] uppercase tracking-widest" :class="activeCaseId === c.id ? 'text-white/50' : 'text-accent'">
              {{ c.n ? "кейс " + c.n : c.method }} {{ c.write ? "· запись" : "· чтение" }}
            </div>
            <div class="mt-1 leading-5">{{ c.title }}</div>
          </button>
          <button v-if="tab === 'p2p'" type="button" class="mt-2 w-full rounded-2xl px-3 py-2 text-left text-[12px] text-mute hover:bg-paper" @click="newKeys">
            Новый Idempotency-Key
          </button>
        </aside>

        <div class="space-y-4">
          <section v-if="activeCase" class="card rounded-[28px] p-5 md:p-6">
            <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-accent">Термин</p>
            <h2 class="font-display mt-2 text-3xl">{{ activeCase.term }}</h2>
            <p class="mt-3 text-[16px] leading-7">{{ activeCase.about }}</p>
            <p class="mt-4 text-sm leading-6"><span class="font-medium">Для чего. </span>{{ activeCase.purpose }}</p>
            <p class="mt-3 text-sm leading-6"><span class="font-medium">Что должно получиться. </span>{{ activeCase.expect }}</p>
          </section>

          <section class="card rounded-[28px] p-5 md:p-6">
            <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-accent">Запрос к ядру</p>
            <h3 class="font-display mt-2 text-2xl">{{ workMethod }} {{ workPath }}</h3>
            <pre class="mt-3 max-h-48 overflow-auto rounded-2xl bg-paper px-4 py-3 font-mono text-[12px] leading-5">{{ requestPreview }}</pre>
            <label v-if="workMethod !== 'GET'" class="mt-4 block text-sm font-medium">
              Тело (можно править)
              <textarea v-model="workBody" rows="7" class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 font-mono text-[12px] leading-5 outline-none focus:border-accent" />
            </label>
            <p v-else class="mt-4 text-sm text-mute">GET без тела. Менять нечего — выполните запрос и читайте ответ.</p>
            <label v-if="workMethod !== 'GET'" class="mt-3 block text-sm font-medium">
              Заголовки
              <textarea v-model="workHeaders" rows="2" class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 font-mono text-[12px] outline-none focus:border-accent" placeholder="Idempotency-Key: …" />
            </label>
            <div class="mt-4 flex flex-wrap gap-2">
              <button type="button" class="btn btn-accent btn-sm" :disabled="running" @click="runCase">Выполнить</button>
              <button type="button" class="btn btn-ghost btn-sm" @click="refresh">Обновить состояние</button>
            </div>
            <p v-if="!auth.isPro && activeCase?.write" class="mt-3 text-sm text-mute">
              Intern: чтение проходит. POST вернёт PLAN — вы всё равно видите контракт вызова.
            </p>
            <p v-if="err" class="mt-3 text-sm text-bad">{{ err }}</p>
          </section>

          <section v-if="last" class="card rounded-[28px] p-5">
            <p class="font-mono text-[11px] uppercase tracking-widest text-mute">Ответ</p>
            <p class="mt-2 font-mono text-[13px]">{{ last.method }} {{ last.path }} → HTTP {{ last.status }}</p>
            <pre class="mt-3 max-h-72 overflow-auto font-mono text-[12px] leading-5 text-mute">{{ pretty(last.data) }}</pre>
          </section>
        </div>
      </div>

      <div class="mt-8 grid gap-4 md:grid-cols-2">
        <section class="card rounded-[24px] p-5">
          <h2 class="font-medium">Кошельки</h2>
          <ul class="mt-3 space-y-2 text-sm">
            <li v-for="w in wallets" :key="w.id" class="flex justify-between gap-3 border-b border-line py-2">
              <span>{{ w.name }} <span class="font-mono text-[11px] text-mute">{{ w.id }}</span></span>
              <span class="font-mono">{{ w.available }} {{ w.currency }}</span>
            </li>
          </ul>
        </section>
        <section class="card rounded-[24px] p-5">
          <h2 class="font-medium">{{ isLesson && tab !== "hold" ? "Последние журналы" : "Открытые холды" }}</h2>
          <template v-if="isLesson && tab !== 'hold'">
            <p v-if="!journals.length" class="mt-3 text-sm text-mute">Журнал пуст или сид ещё грузится.</p>
            <ul v-else class="mt-3 space-y-2 text-sm">
              <li v-for="j in journals" :key="j.id" class="flex justify-between gap-3 border-b border-line py-2 font-mono text-[12px]">
                <span>{{ j.kind }} · {{ j.id }}</span>
                <span>{{ j.debit }}/{{ j.credit }}</span>
              </li>
            </ul>
          </template>
          <template v-else>
            <p v-if="!openHolds.length" class="mt-3 text-sm text-mute">Открытых холдов нет — либо сид ещё грузится.</p>
            <ul v-else class="mt-3 space-y-2 text-sm">
              <li v-for="h in openHolds" :key="h.id" class="flex justify-between border-b border-line py-2 font-mono text-[12px]">
                <span>{{ h.id }} · {{ h.status }}</span>
                <span>{{ h.amount }}</span>
              </li>
            </ul>
          </template>
        </section>
      </div>
    </div>

    <div v-else-if="tab === 'sql'" class="mt-8">
      <section class="card mb-5 rounded-[28px] p-5">
        <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Зачем SQL здесь</p>
        <p class="mt-2 text-[16px] leading-7">
          Тикет приходит словами. Доказательство — выборка. Intern читает витрины intern_*. PRO видит полную схему: wallets, holds, journals, transfers.
          DROP и ALTER запрещены. Эталон открывается после успешной проверки.
        </p>
      </section>
      <div class="grid gap-5 lg:grid-cols-[minmax(260px,0.9fr)_minmax(0,1.2fr)]">
        <aside class="card max-h-[74vh] space-y-2 overflow-auto rounded-[28px] p-3">
          <p class="px-2 pt-2 font-mono text-[10px] uppercase tracking-widest text-mute">Задачи</p>
          <button
            v-for="t in sqlTasks"
            :key="t.id"
            type="button"
            class="w-full rounded-2xl px-3 py-3 text-left text-sm"
            :class="activeTask === t.id ? 'bg-ink text-white' : 'text-mute hover:bg-paper'"
            @click="useTask(t)"
          >
            <div class="font-mono text-[10px] uppercase tracking-widest" :class="activeTask === t.id ? 'text-white/50' : 'text-accent'">
              {{ t.grade }} <span v-if="t.locked">· PRO</span> <span v-if="t.gold">· эталон открыт</span>
            </div>
            <div class="mt-1 leading-5">{{ t.title }}</div>
          </button>
        </aside>
        <div class="space-y-4">
          <section v-if="sqlTask" class="card rounded-[28px] p-5">
            <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Ситуация</p>
            <p class="mt-2 text-[16px] leading-7">{{ sqlTask.story }}</p>
            <p class="mt-4 font-medium">{{ sqlTask.question }}</p>
            <p class="mt-3 text-sm text-mute">Подсказка: {{ sqlTask.hint }}</p>
            <p class="mt-2 text-sm leading-6"><span class="font-medium">Зачем. </span>{{ sqlTask.why }}</p>
          </section>
          <section class="card rounded-[28px] p-5">
            <div class="flex flex-wrap items-center justify-between gap-2">
              <h2 class="font-display text-2xl">Запрос</h2>
              <span class="font-mono text-[11px] text-mute">{{ internSql ? "Intern: intern_wallets / intern_holds / intern_journals / intern_tickets" : "PRO: полная схема ядра" }}</span>
            </div>
            <textarea v-model="sqlText" rows="8" class="mt-3 w-full rounded-2xl border border-line bg-white/90 px-3 py-2 font-mono text-[13px] leading-6 outline-none focus:border-accent" />
            <div class="mt-3 flex flex-wrap gap-2">
              <button type="button" class="btn btn-accent btn-sm" @click="runSql(false)">Выполнить</button>
              <button type="button" class="btn btn-sm" :disabled="!activeTask || sqlTask?.locked" @click="runSql(true)">Проверить задачу</button>
            </div>
            <p class="mt-2 text-[12px] text-mute">Один SELECT за раз. Точка с запятой в конце допустима. DROP / TRUNCATE / ALTER запрещены.</p>
            <p v-if="sqlErr" class="mt-2 text-sm text-bad">{{ sqlErr }}</p>
            <p v-if="sqlOut?.why" class="mt-2 text-[15px] leading-7">{{ sqlOut.why }}</p>
            <p v-if="sqlOut?.gold" class="mt-2 font-mono text-[12px] leading-5 text-mute">Эталон: {{ sqlOut.gold }}</p>
            <div v-if="sqlRows.length" class="mt-4 overflow-auto">
              <table class="w-full text-left font-mono text-[12px]">
                <thead>
                  <tr>
                    <th v-for="c in sqlColumns" :key="c" class="border-b border-line py-1 pr-3">{{ c }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(row, i) in sqlRows" :key="i">
                    <td v-for="c in sqlColumns" :key="c" class="border-b border-line py-1 pr-3">{{ row[c] }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
          <section class="card rounded-[24px] p-5">
            <h3 class="font-medium">Таблицы лабы</h3>
            <ul class="mt-3 space-y-2 text-sm text-mute">
              <li v-for="t in sqlSchema" :key="t.name"><span class="font-mono text-ink">{{ t.name }}</span> · {{ t.cols }}</li>
            </ul>
          </section>
        </div>
      </div>
    </div>

    <div v-else class="mt-8">
      <section class="card mb-5 rounded-[28px] p-5">
        <p class="font-mono text-[10px] uppercase tracking-widest text-accent">Справочник, не Swagger «попробуйте»</p>
        <p class="mt-2 text-[16px] leading-7">
          Выберите метод слева. Справа — когда вызывать, что меняется в ядре, типичные ошибки и живой запрос.
          GET читается на Intern. POST и PATCH — PRO или триал.
        </p>
      </section>
      <div class="grid gap-5 lg:grid-cols-[minmax(240px,0.85fr)_minmax(0,1.25fr)]">
        <aside class="space-y-3">
          <div class="flex flex-wrap gap-1">
            <button v-for="g in groups" :key="g" type="button" class="tab-plate" :class="group === g ? 'is-on' : ''" @click="group = g">{{ g }}</button>
          </div>
          <div class="card max-h-[68vh] space-y-1 overflow-auto rounded-[28px] p-3">
            <button
              v-for="row in grouped"
              :key="row.method + row.path"
              type="button"
              class="w-full rounded-2xl px-3 py-3 text-left text-sm"
              :class="activeDoc?.path === row.path && activeDoc?.method === row.method ? 'bg-ink text-white' : 'hover:bg-paper'"
              @click="pickDoc(row)"
            >
              <div class="font-mono text-[10px] uppercase tracking-widest opacity-60">{{ row.method }}</div>
              <div class="mt-1 leading-5">{{ row.title }}</div>
            </button>
          </div>
        </aside>
        <div v-if="activeDoc" class="space-y-4">
          <section class="card rounded-[28px] p-5 md:p-6">
            <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-accent">{{ activeDoc.group }} · {{ activeDoc.method }} {{ activeDoc.path }}</p>
            <h2 class="font-display mt-2 text-3xl">{{ activeDoc.title }}</h2>
            <p class="mt-3 text-[16px] leading-7"><span class="font-medium">Когда вызывать. </span>{{ activeDoc.when }}</p>
            <p class="mt-4 text-sm leading-6 text-mute"><span class="font-medium text-ink">Доступ. </span>{{ activeDoc.auth }}</p>
            <p class="mt-3 text-sm leading-6"><span class="font-medium">Что делает ядро. </span>{{ activeDoc.db }}</p>
            <p v-if="activeDoc.errors" class="mt-3 text-sm leading-6"><span class="font-medium">Типичные ошибки. </span>{{ activeDoc.errors }}</p>
            <p v-if="activeDoc.bug" class="mt-3 rounded-2xl bg-paper px-4 py-3 text-sm leading-6">{{ activeDoc.bug }}</p>
            <pre class="mt-5 overflow-auto rounded-2xl bg-paper px-4 py-3 font-mono text-[12px] leading-5">{{ activeDoc.method }} {{ activeDoc.path }}{{ headerKey ? "\nIdempotency-Key: " + headerKey : "" }}{{ activeDoc.method === "GET" ? "\n\n(тела нет)" : bodyText ? "\n\n" + bodyText : "" }}</pre>
            <label v-if="activeDoc.method !== 'GET'" class="mt-5 block text-sm font-medium">
              Тело
              <textarea v-model="bodyText" rows="7" class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 font-mono text-[12px] leading-5 outline-none focus:border-accent" />
            </label>
            <p v-else class="mt-5 text-sm text-mute">GET без тела. Ответ появится ниже после выполнения.</p>
            <label v-if="activeDoc.method !== 'GET'" class="mt-3 block text-sm font-medium">
              Idempotency-Key
              <input v-model="headerKey" class="mt-2 w-full rounded-2xl border border-line bg-white px-3 py-2 font-mono text-[13px] outline-none focus:border-accent" />
            </label>
            <button type="button" class="mt-4 rounded-full bg-accent px-5 py-2.5 text-sm font-semibold text-white" @click="runDoc">Выполнить</button>
            <pre v-if="runOut" class="mt-4 max-h-64 overflow-auto font-mono text-[12px] leading-5 text-mute">{{ pretty(runOut) }}</pre>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>
