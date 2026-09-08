<script setup>
import { RouterLink } from "vue-router";

defineProps({
  linked: { type: Boolean, default: true },
});

const rows = [
  { method: "GET", path: "/api/v1/wallets", note: "ledger, hold, available", tab: "console", who: "Intern" },
  { method: "GET", path: "/api/v1/ledger/journals", note: "проводки книги", tab: "hole", who: "Intern" },
  { method: "POST", path: "/api/v1/transfers", note: "P2P, заголовок Idempotency-Key", tab: "p2p", who: "PRO" },
  { method: "POST", path: "/api/v1/cards/authorize", note: "холд: резерв, журнала ещё нет", tab: "hold", who: "PRO" },
  { method: "POST", path: "/api/v1/bank/incoming", note: "входящий платёж на IBAN", tab: "iban", who: "PRO" },
];
</script>

<template>
  <aside class="card overflow-hidden">
    <div class="border-b border-line px-5 py-4">
      <p class="kicker">контракт ядра</p>
      <h3 class="font-display mt-2 text-2xl">Куда уходит запрос</h3>
      <p class="mt-2 text-sm leading-6 text-mute">Метод и полный путь. Не строка «POST /api/v1».</p>
    </div>
    <ul>
      <li v-for="r in rows" :key="r.path" class="border-b border-line last:border-b-0">
        <RouterLink
          v-if="linked"
          :to="`/pet?tab=${r.tab}`"
          class="block px-5 py-3 transition-colors hover:bg-paper"
        >
          <div class="flex flex-wrap items-baseline justify-between gap-2">
            <span class="font-mono text-[12px] leading-5">
              <span class="text-accent">{{ r.method }}</span>
              {{ r.path }}
            </span>
            <span class="stamp">{{ r.who }}</span>
          </div>
          <p class="mt-1 text-[13px] leading-5 text-mute">{{ r.note }}</p>
        </RouterLink>
        <div v-else class="px-5 py-3">
          <div class="flex flex-wrap items-baseline justify-between gap-2">
            <span class="font-mono text-[12px] leading-5">
              <span class="text-accent">{{ r.method }}</span>
              {{ r.path }}
            </span>
            <span class="stamp">{{ r.who }}</span>
          </div>
          <p class="mt-1 text-[13px] leading-5 text-mute">{{ r.note }}</p>
        </div>
      </li>
    </ul>
  </aside>
</template>
