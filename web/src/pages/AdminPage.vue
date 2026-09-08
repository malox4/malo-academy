<script setup>
import { onMounted, ref } from "vue";
import { useAuth } from "../stores/auth";

const auth = useAuth();
const users = ref([]);

async function load() {
  const res = await fetch("/api/admin/users", { credentials: "include" });
  const data = await res.json();
  users.value = data.items || [];
}

async function setPlan(id, plan) {
  await fetch(`/api/admin/users/${id}`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ plan }),
  });
  await load();
}

onMounted(load);
</script>

<template>
  <div v-if="auth.user?.role !== 'admin'">Нет доступа.</div>
  <div v-else>
    <h1 class="font-display text-5xl">Касса</h1>
    <p class="mt-3 text-mute">PRO открывает Junior+ и запись в учебный банк. Free оставляет Intern. Триал хранится на пользователе.</p>
    <ul class="mt-6 divide-y divide-line">
      <li v-for="u in users" :key="u.id" class="flex flex-wrap items-center justify-between gap-3 py-3">
        <div>
          <div class="font-medium">{{ u.name }}</div>
          <div class="text-sm text-mute">{{ u.email }} · {{ u.plan }}<span v-if="u.trial_until"> · trial {{ u.trial_until }}</span></div>
        </div>
        <div class="flex gap-2">
          <button type="button" class="rounded-full border border-line px-3 py-1 text-sm" @click="setPlan(u.id, 'free')">free</button>
          <button type="button" class="rounded-full bg-accent px-3 py-1 text-sm text-white" @click="setPlan(u.id, 'pro')">pro</button>
        </div>
      </li>
    </ul>
  </div>
</template>
