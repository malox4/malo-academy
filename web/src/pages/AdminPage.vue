<script setup>
import { computed, onMounted, ref } from "vue";
import { getJSON } from "../lib/http";
import { useAuth } from "../stores/auth";

const auth = useAuth();
const tab = ref("desk");
const filter = ref("all");
const q = ref("");
const err = ref("");
const busy = ref("");
const storeMode = ref("");
const users = ref([]);
const practice = ref([]);
const openId = ref("");
const notes = ref({});
const scores = ref({});
const stats = ref({ all: 0, online: 0, silent: 0, walking: 0 });

const filtered = computed(() => {
  const needle = q.value.trim().toLowerCase();
  return users.value.filter((u) => {
    if (filter.value === "online" && !u.online) return false;
    if (filter.value === "path" && !(u.lessonsDone > 0)) return false;
    if (filter.value === "quiet" && u.lessonsDone > 0) return false;
    if (filter.value === "pro" && !u.isPro) return false;
    if (!needle) return true;
    return [u.name, u.email, u.here, u.floorTitle, u.assessLevel, u.lastLesson?.title]
      .join(" ")
      .toLowerCase()
      .includes(needle);
  });
});

const practiceView = computed(() => {
  const items = [...practice.value].sort((a, b) => String(b.created_at || "").localeCompare(String(a.created_at || "")));
  if (!q.value.trim()) return items;
  const needle = q.value.trim().toLowerCase();
  return items.filter((p) => [p.title, p.email, p.author, p.kind, p.status].join(" ").toLowerCase().includes(needle));
});

function accessLabel(u) {
  if (u.role === "admin") return "хозяин";
  if (u.plan === "pro") return "PRO";
  if (u.trialActive) return "триал";
  return "intern";
}

function seenLabel(u) {
  if (!u.lastSeen) return "ещё не заходил";
  if (u.online) return "в зале";
  const t = new Date(u.lastSeen).getTime();
  if (!t) return "ещё не заходил";
  const min = Math.max(1, Math.round((Date.now() - t) / 60000));
  if (min < 60) return `${min} мин назад`;
  const h = Math.round(min / 60);
  if (h < 48) return `${h} ч назад`;
  return `${Math.round(h / 24)} дн назад`;
}

function fail(data, fallback) {
  return data?.error?.message || fallback;
}

async function loadUsers() {
  const { ok, data } = await getJSON("/api/admin/users");
  if (!ok) {
    err.value = fail(data, "Журнал не открылся.");
    return;
  }
  users.value = data.items || [];
  storeMode.value = data.store || "";
  stats.value = data.stats || stats.value;
}

async function loadPractice() {
  const { ok, data } = await getJSON("/api/admin/practice");
  if (!ok) {
    err.value = fail(data, "Очередь практики не открылась.");
    return;
  }
  practice.value = data.items || [];
  const n = { ...notes.value };
  const s = { ...scores.value };
  for (const item of practice.value) {
    if (n[item.id] == null) n[item.id] = item.reviewer_note || "";
    if (s[item.id] == null) s[item.id] = item.score ?? "";
  }
  notes.value = n;
  scores.value = s;
}

async function load() {
  err.value = "";
  await Promise.all([loadUsers(), loadPractice()]);
}

async function patchUser(id, body) {
  busy.value = id;
  err.value = "";
  const { ok, data } = await getJSON(`/api/admin/users/${id}`, {
    method: "PATCH",
    body: JSON.stringify(body),
  });
  busy.value = "";
  if (!ok) {
    err.value = fail(data, "Доступ не выдан.");
    return;
  }
  await loadUsers();
}

async function review(item, status) {
  busy.value = item.id;
  err.value = "";
  const scoreRaw = scores.value[item.id];
  const payload = { status, note: notes.value[item.id] || "" };
  if (scoreRaw !== "" && scoreRaw != null) payload.score = Number(scoreRaw);
  const { ok, data } = await getJSON(`/api/admin/practice/${item.id}`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  });
  busy.value = "";
  if (!ok) {
    err.value = fail(data, "Разбор не сохранился.");
    return;
  }
  const next = data.item || item;
  practice.value = practice.value.map((p) => (p.id === item.id ? { ...p, ...next } : p));
}

function toggle(id) {
  openId.value = openId.value === id ? "" : id;
}

onMounted(load);
</script>

<template>
  <div>
    <p class="kicker">журнал</p>
    <h1 class="font-display mt-3 text-[1.85rem] leading-[0.95] sm:text-4xl md:text-6xl">Кто в зале</h1>
    <span class="section-rule" aria-hidden="true" />
    <p class="mt-4 max-w-2xl text-[17px] leading-7 text-mute">
      Где человек сейчас, какой этаж, сколько уроков закрыл, живой ли собес. Доступ — внутри карточки, не вместо следа.
    </p>

    <div class="mt-6 flex flex-wrap gap-3">
      <span class="stamp stamp-ink">всего {{ stats.all }}</span>
      <span class="stamp stamp-ok">в зале {{ stats.online }}</span>
      <span class="stamp">путь {{ stats.walking }}</span>
      <span class="stamp stamp-case">тихие {{ stats.silent }}</span>
    </div>

    <div class="mt-8 flex flex-wrap gap-2">
      <button type="button" class="tab-plate" :class="{ 'is-on': tab === 'desk' }" @click="tab = 'desk'">След</button>
      <button type="button" class="tab-plate" :class="{ 'is-on': tab === 'practice' }" @click="tab = 'practice'">Практика</button>
    </div>

    <div v-if="tab === 'desk'" class="mt-4 flex flex-wrap gap-2">
      <button type="button" class="tab-plate" :class="{ 'is-on': filter === 'all' }" @click="filter = 'all'">Все</button>
      <button type="button" class="tab-plate" :class="{ 'is-on': filter === 'online' }" @click="filter = 'online'">В зале</button>
      <button type="button" class="tab-plate" :class="{ 'is-on': filter === 'path' }" @click="filter = 'path'">Идут</button>
      <button type="button" class="tab-plate" :class="{ 'is-on': filter === 'quiet' }" @click="filter = 'quiet'">Тихие</button>
      <button type="button" class="tab-plate" :class="{ 'is-on': filter === 'pro' }" @click="filter = 'pro'">PRO</button>
    </div>

    <label class="mt-6 block max-w-md text-sm font-medium">
      Поиск
      <input
        v-model="q"
        class="mt-1 w-full border-b-2 border-ink bg-transparent py-2 outline-none"
        :placeholder="tab === 'desk' ? 'имя, экран, этаж' : 'работа, автор'"
      />
    </label>

    <p v-if="err" class="mt-4 text-sm text-bad">{{ err }}</p>
    <p v-if="storeMode === 'file'" class="mt-4 max-w-2xl text-sm text-bad">
      Учётки в файле контейнера. На Railway след пропадает после рестарта — нужен Postgres.
    </p>

    <ol v-if="tab === 'desk'" class="desk-list mt-8">
      <li v-for="(u, i) in filtered" :key="u.id" class="desk-item">
        <button type="button" class="desk-head" @click="toggle(u.id)">
          <span class="plate-n">{{ String(i + 1).padStart(2, "0") }}</span>
          <span class="min-w-0 flex-1 text-left">
            <span class="flex flex-wrap items-center gap-2">
              <span class="font-display text-2xl leading-none">{{ u.name }}</span>
              <span class="stamp" :class="u.online ? 'stamp-ok' : 'stamp-ink'">{{ seenLabel(u) }}</span>
              <span class="stamp" :class="u.role === 'admin' ? 'stamp-solid' : u.isPro ? 'stamp-ok' : 'stamp-ink'">{{ accessLabel(u) }}</span>
            </span>
            <span class="mt-2 block text-sm text-mute">
              {{ u.here || "ещё нигде" }}
              <span v-if="u.lastLesson?.title"> · последний урок {{ u.lastLesson.title }}</span>
            </span>
          </span>
          <span class="desk-pct">
            <span class="font-mono text-[12px]">{{ u.pct || 0 }}%</span>
            <span class="desk-meter" aria-hidden="true"><i :style="{ width: (u.pct || 0) + '%' }" /></span>
            <span class="mt-1 block font-mono text-[11px] text-mute">{{ u.lessonsDone || 0 }}/{{ u.lessonsTotal || 0 }} · {{ u.floorTitle || "intern" }}</span>
          </span>
        </button>
        <div v-if="openId === u.id" class="desk-body">
          <p class="text-sm text-mute">{{ u.email }}</p>
          <div class="mt-5 grid gap-4 md:grid-cols-4">
            <div v-for="f in u.floors || []" :key="f.id">
              <div class="flex justify-between font-mono text-[11px] text-mute">
                <span>{{ f.title }}</span>
                <span>{{ f.done }}/{{ f.total }}</span>
              </div>
              <span class="desk-meter mt-2" aria-hidden="true"><i :style="{ width: (f.pct || 0) + '%' }" /></span>
            </div>
          </div>
          <div class="mt-6 grid gap-6 md:grid-cols-2">
            <div>
              <p class="stamp stamp-ink">последние уроки</p>
              <ul v-if="u.recent?.length" class="mt-3 space-y-1 text-sm">
                <li v-for="r in u.recent" :key="r.id">{{ r.title }}</li>
              </ul>
              <p v-else class="mt-3 text-sm text-mute">Уроков ещё нет.</p>
            </div>
            <div>
              <p class="stamp stamp-ink">собес / практика</p>
              <p class="mt-3 text-sm text-mute">
                Разобрано карточек {{ u.interviews || 0 }}.
                <span v-if="u.assessLevel"> Live: {{ u.assessLevel }}.</span>
                Практик {{ u.practiceTotal || 0 }}<span v-if="u.practiceOpen"> · в очереди {{ u.practiceOpen }}</span>.
              </p>
            </div>
          </div>
          <div class="mt-6 flex flex-wrap gap-2">
            <button type="button" class="btn btn-accent btn-sm" :disabled="busy === u.id || u.plan === 'pro'" @click="patchUser(u.id, { plan: 'pro' })">Выдать PRO</button>
            <button type="button" class="btn btn-ghost btn-sm" :disabled="busy === u.id || u.role === 'admin'" @click="patchUser(u.id, { trialDays: 3 })">Триал 3 дня</button>
            <button type="button" class="btn btn-ghost btn-sm" :disabled="busy === u.id || u.role === 'admin' || (u.plan === 'free' && !u.trialActive)" @click="patchUser(u.id, { plan: 'free', trialDays: 0 })">Intern</button>
            <button v-if="u.role !== 'admin'" type="button" class="btn btn-ghost btn-sm" :disabled="busy === u.id" @click="patchUser(u.id, { role: 'admin' })">Хозяин</button>
            <button v-else-if="u.id !== auth.user?.id" type="button" class="btn btn-ghost btn-sm" :disabled="busy === u.id" @click="patchUser(u.id, { role: 'student' })">Снять хозяина</button>
          </div>
        </div>
      </li>
      <li v-if="!filtered.length" class="py-6 text-sm text-mute">Никого по этому следу.</li>
    </ol>

    <ul v-else class="mt-8 space-y-4">
      <li v-for="item in practiceView" :key="item.id" class="card px-5 py-5">
        <div class="flex flex-wrap items-center gap-2">
          <span class="stamp" :class="item.status === 'reviewed' ? 'stamp-ok' : 'stamp-ink'">{{ item.status }}</span>
          <span class="stamp stamp-ink">{{ item.kind }}</span>
        </div>
        <h2 class="font-display mt-3 text-2xl">{{ item.title || "Без названия" }}</h2>
        <p class="mt-1 text-sm text-mute">{{ item.author }} · {{ item.email }}</p>
        <p class="mt-3 whitespace-pre-wrap text-[15px] leading-7">{{ item.body }}</p>
        <div class="mt-4 flex flex-wrap items-end gap-3">
          <label class="text-sm font-medium">
            Балл
            <input v-model="scores[item.id]" type="number" min="0" max="100" class="mt-1 w-20 border-b-2 border-ink bg-transparent py-1 outline-none" />
          </label>
          <label class="min-w-[12rem] flex-1 text-sm font-medium">
            Заметка
            <input v-model="notes[item.id]" class="mt-1 w-full border-b-2 border-ink bg-transparent py-1 outline-none" />
          </label>
          <button type="button" class="btn btn-accent btn-sm" :disabled="busy === item.id" @click="review(item, 'reviewed')">Принять</button>
          <button type="button" class="btn btn-ghost btn-sm" :disabled="busy === item.id" @click="review(item, 'submitted')">Вернуть</button>
        </div>
      </li>
      <li v-if="!practiceView.length" class="text-sm text-mute">Очередь пуста.</li>
    </ul>
  </div>
</template>
