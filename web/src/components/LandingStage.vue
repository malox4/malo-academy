<script setup>
import { computed, onMounted, onUnmounted, ref } from "vue";

const scenes = [
  {
    id: "read",
    n: "01",
    stamp: "Intern",
    write: false,
    title: "Чтение кошелька",
    method: "GET",
    path: "/api/v1/wallets",
    head: "",
    body: "",
    res: "200  [{ name: Борис, ledger, hold, available }]",
    note: "Тела нет. Intern это видит всегда.",
    wallets: [
      { name: "Анна", ledger: 92000, hold: 0, avail: 92000 },
      { name: "Борис", ledger: 100000, hold: 0, avail: 100000 },
    ],
  },
  {
    id: "hold",
    n: "02",
    stamp: "PRO",
    write: true,
    title: "Холд: резерв, не проводка",
    method: "POST",
    path: "/api/v1/cards/authorize",
    head: "",
    body: '{ "walletId": "wal_boris", "cardId": "crd_boris", "amount": 12000, "currency": "UZS" }',
    res: "201  { hold: { status: OPEN } }",
    note: "Касса SUCCESS. Журнал ещё 100 000. Available 88 000.",
    wallets: [
      { name: "Анна", ledger: 92000, hold: 0, avail: 92000 },
      { name: "Борис", ledger: 100000, hold: 12000, avail: 88000 },
    ],
  },
  {
    id: "p2p",
    n: "03",
    stamp: "PRO",
    write: true,
    title: "P2P и ключ попытки",
    method: "POST",
    path: "/api/v1/transfers",
    head: "Idempotency-Key: try-1",
    body: '{ "fromWalletId": "wal_anna", "toWalletId": "wal_boris", "amount": 5000, "currency": "UZS" }',
    res: "201  { journalId, amount: 5000 }",
    note: "Повтор с тем же ключом не списывает второй раз. Без ключа учебный банк проведёт снова.",
    wallets: [
      { name: "Анна", ledger: 87000, hold: 0, avail: 87000 },
      { name: "Борис", ledger: 105000, hold: 0, avail: 105000 },
    ],
  },
  {
    id: "iban",
    n: "04",
    stamp: "PRO",
    write: true,
    title: "Входящий IBAN",
    method: "POST",
    path: "/api/v1/bank/incoming",
    head: "",
    body: '{ "iban": "UZ12MALO000000000001", "amount": 25000, "currency": "UZS" }',
    res: "201  кредит клиента · дебет nostro",
    note: "Это не P2P: другой путь, другие счета. Чужой IBAN — невыясненные, не «похожий» клиент.",
    wallets: [
      { name: "Анна", ledger: 117000, hold: 0, avail: 117000 },
      { name: "Борис", ledger: 100000, hold: 0, avail: 100000 },
    ],
  },
];

const idx = ref(0);
const pinned = ref(false);
const scene = computed(() => scenes[idx.value]);
let timer = 0;

function money(n) {
  return Number(n).toLocaleString("ru-RU");
}

function go(i) {
  idx.value = i;
  pinned.value = true;
}

function tick() {
  if (pinned.value) return;
  idx.value = (idx.value + 1) % scenes.length;
}

onMounted(() => {
  const reduce = window.matchMedia?.("(prefers-reduced-motion: reduce)")?.matches;
  if (!reduce) timer = window.setInterval(tick, 3800);
});

onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});
</script>

<template>
  <div class="stage" @mouseenter="pinned = true" @mouseleave="pinned = false">
    <div class="stage-bar">
      <button
        v-for="(s, i) in scenes"
        :key="s.id"
        type="button"
        class="stage-tab"
        :aria-pressed="i === idx"
        @click="go(i)"
      >
        <span class="plate-n">{{ s.n }}</span>
        <span class="stage-tab-t">{{ s.title }}</span>
        <span class="stamp" :class="s.write ? 'stamp-solid' : ''">{{ s.stamp }}</span>
      </button>
    </div>

    <div class="stage-wallets">
      <article v-for="w in scene.wallets" :key="w.name + scene.id" class="stage-wallet rise">
        <p class="stamp">{{ w.name }}</p>
        <dl>
          <div>
            <dt>ledger</dt>
            <dd>{{ money(w.ledger) }}</dd>
          </div>
          <div :class="w.hold ? 'is-hold' : ''">
            <dt>hold</dt>
            <dd>{{ money(w.hold) }}</dd>
          </div>
          <div class="is-avail">
            <dt>available</dt>
            <dd>{{ money(w.avail) }}</dd>
          </div>
        </dl>
      </article>
    </div>

    <div class="stage-http">
      <div class="stage-http-meta">
        <span class="live-dot stage-dot" aria-hidden="true" />
        <span>{{ scene.write ? "запись · PRO / триал" : "чтение · Intern" }}</span>
      </div>
      <p class="stage-line">
        <span class="text-accent">{{ scene.method }}</span>
        {{ scene.path }}
      </p>
      <p v-if="scene.head" class="stage-line mute">{{ scene.head }}</p>
      <p v-if="scene.body" class="stage-line">{{ scene.body }}</p>
      <p v-else class="stage-line mute">(тела нет)</p>
      <p class="stage-res">← {{ scene.res }}</p>
    </div>

    <p class="stage-note">{{ scene.note }}</p>
    <div class="stage-progress" aria-hidden="true">
      <i :key="scene.id" :class="pinned ? 'is-paused' : ''" />
    </div>
  </div>
</template>
