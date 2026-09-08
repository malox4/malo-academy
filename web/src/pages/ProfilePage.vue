<script setup>
import { computed, onMounted, ref } from "vue";
import { useAuth } from "../stores/auth";
import { useRouter, RouterLink } from "vue-router";
import PlanCompare from "../components/PlanCompare.vue";
import { getJSON } from "../lib/http";

const auth = useAuth();
const router = useRouter();
const assessment = ref(null);

onMounted(async () => {
  const { data } = await getJSON("/api/me/assessment");
  assessment.value = data.assessment;
});

const status = computed(() => {
  if (auth.trialActive) return `Пробный PRO · осталось ${auth.trialLabel}`;
  if (auth.isPro) return "PRO · все грейды и запись в учебный банк";
  if (auth.trialEnded) return "Триал закончился · Intern открыт";
  return "Слушатель · Intern открыт, банк на чтение";
});

const hards = computed(() => assessment.value?.hards || {});
const softs = computed(() => assessment.value?.softs || {});

async function out() {
  await auth.logout();
  router.push("/");
}
</script>

<template>
  <div>
    <p class="font-mono text-[11px] uppercase tracking-[0.2em] text-mute">Analyst Hall · профиль</p>
    <h1 class="font-display mt-2 text-5xl">{{ auth.user?.name }}</h1>
    <p class="mt-3 text-[17px] leading-7 text-mute">{{ auth.user?.email }} · {{ status }}</p>
    <section v-if="assessment" class="card mt-8 rounded-[28px] p-5">
      <p class="font-mono text-[11px] uppercase tracking-widest text-accent">Собес · уровень</p>
      <h2 class="font-display mt-1 text-3xl">{{ assessment.level }}</h2>
      <p class="mt-2 text-sm text-mute">Уверенность {{ Math.round((assessment.confidence || 0) * 100) }}% · балл {{ assessment.score }}/{{ assessment.max }}</p>
      <div class="mt-4 grid gap-4 md:grid-cols-2">
        <div>
          <div class="text-sm font-medium">Харды</div>
          <ul class="mt-2 space-y-1 text-sm text-mute">
            <li v-for="(v, k) in hards" :key="k">{{ k }} · {{ v }}</li>
          </ul>
        </div>
        <div>
          <div class="text-sm font-medium">Софты</div>
          <ul class="mt-2 space-y-1 text-sm text-mute">
            <li v-for="(v, k) in softs" :key="k">{{ k }} · {{ v }}</li>
          </ul>
        </div>
      </div>
    </section>
    <div class="mt-6 flex flex-wrap gap-3">
      <RouterLink to="/pet?tab=sql" class="btn btn-ghost btn-sm">Пет · SQL</RouterLink>
      <RouterLink to="/pet?tab=api" class="btn btn-ghost btn-sm">Пет · API</RouterLink>
      <RouterLink to="/pet?tab=p2p" class="btn btn-ghost btn-sm">Пет · P2P</RouterLink>
      <RouterLink to="/practice" class="btn btn-ghost btn-sm">Практика</RouterLink>
      <RouterLink to="/interview" class="btn btn-ghost btn-sm">Собес</RouterLink>
      <RouterLink to="/live" class="btn btn-ghost btn-sm">Live собес</RouterLink>
    </div>
    <button type="button" class="btn btn-ghost btn-sm mt-8" @click="out">Выйти</button>
    <PlanCompare class="mt-14" :cta="false" />
  </div>
</template>
